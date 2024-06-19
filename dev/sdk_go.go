package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/moby/buildkit/identity"
	"golang.org/x/sync/errgroup"

	"github.com/dagger/dagger/dev/internal/consts"
	"github.com/dagger/dagger/dev/internal/dagger"
)

type GoSDK struct {
	Dagger *DaggerDev // +private
}

// Lint the Go SDK
func (t GoSDK) Lint(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return dag.
			Go(t.Dagger.Source()).
			Lint(ctx, dagger.GoLintOpts{Packages: []string{"sdk/go"}})
	})
	eg.Go(func() error {
		before := t.Dagger.Source()
		after, err := t.Generate(ctx)
		if err != nil {
			return err
		}
		return dag.Dirdiff().AssertEqual(ctx, before, after, []string{"sdk/go"})
	})
	return eg.Wait()
}

// Test the Go SDK
func (t GoSDK) Test(ctx context.Context) error {
	installer, err := t.Dagger.installer(ctx, "sdk-go-test")
	if err != nil {
		return err
	}

	output, err := t.Dagger.Go().Env().
		With(installer).
		WithWorkdir("sdk/go").
		WithExec([]string{"go", "test", "-v", "-skip=TestProvision", "./..."}).
		Stdout(ctx)
	if err != nil {
		err = fmt.Errorf("test failed: %w\n%s", err, output)
	}
	return err
}

// Regenerate the Go SDK API
func (t GoSDK) Generate(ctx context.Context) (*dagger.Directory, error) {
	installer, err := t.Dagger.installer(ctx, "sdk-go-generate")
	if err != nil {
		return nil, err
	}

	generated := t.Dagger.Go().Env().
		With(installer).
		WithWorkdir("sdk/go").
		WithExec([]string{"go", "generate", "-v", "./..."}).
		WithExec([]string{"go", "mod", "tidy"}).
		Directory(".")
	return dag.Directory().WithDirectory("sdk/go", generated), nil
}

// Publish the Go SDK
func (t GoSDK) Publish(
	ctx context.Context,
	tag string,

	// +optional
	dryRun bool,

	// +optional
	// +default="https://github.com/dagger/dagger-go-sdk.git"
	gitRepo string,
	// +optional
	// +default="dagger-ci"
	gitUserName string,
	// +optional
	// +default="hello@dagger.io"
	gitUserEmail string,

	// +optional
	githubToken *dagger.Secret,
) error {
	return gitPublish(ctx, gitPublishOpts{
		sdk:          "go",
		source:       "https://github.com/dagger/dagger.git",
		sourceTag:    tag,
		sourcePath:   "sdk/go/",
		sourceFilter: "if [ -f go.mod ]; then go mod edit -dropreplace github.com/dagger/dagger; fi",
		sourceEnv:    t.Dagger.Go().Env(),
		dest:         gitRepo,
		destTag:      strings.TrimPrefix(tag, "sdk/go/"),
		username:     gitUserName,
		email:        gitUserEmail,
		githubToken:  githubToken,
		dryRun:       dryRun,
	})
}

// Bump the Go SDK's Engine dependency
//
// Deprecated: this is now included in the Publish step
func (t GoSDK) Bump(
	ctx context.Context,
	// +optional
	version string,
) (*dagger.Directory, error) {
	result := bumpSDK("go", version, t.Dagger.Source())
	return t.Dagger.Source().Diff(result), nil
}

type gitPublishOpts struct {
	sdk string

	source, dest       string
	sourceTag, destTag string
	sourcePath         string
	sourceFilter       string
	sourceEnv          *dagger.Container

	username    string
	email       string
	githubToken *dagger.Secret

	dryRun bool
}

func gitPublish(ctx context.Context, opts gitPublishOpts) error {
	base := opts.sourceEnv
	if base == nil {
		base = dag.Container().
			From(consts.AlpineImage).
			WithExec([]string{"apk", "add", "-U", "--no-cache", "git", "go", "python3"})
	}

	// FIXME: move this into std modules
	git := base.
		WithExec([]string{"git", "config", "--global", "user.name", opts.username}).
		WithExec([]string{"git", "config", "--global", "user.email", opts.email})
	if !opts.dryRun {
		githubTokenRaw, err := opts.githubToken.Plaintext(ctx)
		if err != nil {
			return err
		}
		encodedPAT := base64.URLEncoding.EncodeToString([]byte("pat:" + githubTokenRaw))
		git = git.
			WithEnvVariable("GIT_CONFIG_COUNT", "1").
			WithEnvVariable("GIT_CONFIG_KEY_0", "http.https://github.com/.extraheader").
			WithSecretVariable("GIT_CONFIG_VALUE_0", dag.SetSecret("GITHUB_HEADER", fmt.Sprintf("AUTHORIZATION: Basic %s", encodedPAT)))
	}

	result := git.
		WithEnvVariable("CACHEBUSTER", identity.NewID()).
		WithWorkdir("/src/dagger").
		WithExec([]string{"git", "clone", opts.source, "."}).
		WithExec([]string{"git", "fetch", "origin", "-v", "--update-head-ok", fmt.Sprintf("refs/*%[1]s:refs/*%[1]s", strings.TrimPrefix(opts.sourceTag, "refs/"))}).
		WithMountedFile("/opt/version-update.py", dag.CurrentModule().Source().File("version-update.py")).
		WithExec([]string{"python", "/opt/version-update.py", "--sdk", opts.sdk, "git", "--ref", opts.sourceTag}).
		WithEnvVariable("FILTER_BRANCH_SQUELCH_WARNING", "1").
		WithExec([]string{
			"git", "filter-branch", "-f", "--prune-empty",
			"--subdirectory-filter", opts.sourcePath,
			"--tree-filter", opts.sourceFilter,
			"--", opts.sourceTag,
		})
	if !opts.dryRun {
		result = result.WithExec([]string{
			"git",
			"push",
			// "--force", // NOTE: disabled to avoid accidentally rewriting the history
			opts.dest,
			fmt.Sprintf("%s:%s", opts.sourceTag, opts.destTag),
		})
	} else {
		// on a dry run, just test that the last state of dest is in the current branch (and is a fast-forward)
		history, err := result.
			WithExec([]string{"git", "log", "--oneline", "--no-abbrev-commit"}).
			Stdout(ctx)
		if err != nil {
			return err
		}

		destCommit, err := git.
			WithEnvVariable("CACHEBUSTER", identity.NewID()).
			WithWorkdir("/src/dagger").
			WithExec([]string{"git", "clone", opts.dest, "."}).
			WithExec([]string{"git", "rev-parse", opts.destTag}).
			Stdout(ctx)
		if err != nil {
			if strings.Contains(err.Error(), "invalid reference: "+opts.destTag) {
				// this is a ref that only exists in the source, and not in the
				// dest, so no overwriting will occur
				return nil
			}
			return err
		}
		destCommit = strings.TrimSpace(destCommit)

		if !strings.Contains(history, destCommit) {
			return fmt.Errorf("publish would rewrite history - %s not found\n%s", destCommit, history)
		}
		return nil
	}

	_, err := result.Sync(ctx)
	return err
}

func bumpSDK(sdk string, version string, src *dagger.Directory) *dagger.Directory {
	updater := dag.CurrentModule().Source().File("./version-update.py")
	result := dag.Container().
		From("python:3.12").
		WithFile("/usr/local/bin/changie", dag.Container().From("ghcr.io/miniscruff/changie").File("changie")).
		WithFile("/opt/version-update.py", updater).
		WithMountedDirectory("/src", src).
		WithWorkdir("/src")
	if version == "latest" {
		result = result.WithExec([]string{"sh", "-c",
			fmt.Sprintf("python /opt/version-update.py --sdk=%q generate $(changie latest)", sdk),
		})
	} else {
		result = result.WithExec([]string{"sh", "-c",
			fmt.Sprintf("python /opt/version-update.py --sdk=%q generate %q", sdk, version),
		})
	}
	return result.Directory("/src")
}
