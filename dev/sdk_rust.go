package main

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/sync/errgroup"

	"github.com/dagger/dagger/dev/internal/dagger"
	"github.com/dagger/dagger/dev/internal/util"
)

const (
	rustGeneratedAPIPath = "sdk/rust/crates/dagger-sdk/src/gen.rs"
	rustVersionFilePath  = "sdk/rust/crates/dagger-sdk/src/core/mod.rs"

	// https://hub.docker.com/_/rust
	rustDockerStable = "rust:1.77-bookworm"
	cargoChefVersion = "0.1.62"
)

type RustSDK struct {
	Dagger  *DaggerDev // +private
	Version string
}

func (r RustSDK) withVersion(version string) RustSDK {
	r.Version = version
	return r
}

// Lint the Rust SDK
func (r RustSDK) Lint(ctx context.Context) error {
	base := r.rustBase(rustDockerStable)

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		_, err := base.
			WithExec([]string{"cargo", "check", "--all", "--release"}).
			Sync(ctx)
		return err
	})

	eg.Go(func() error {
		_, err := base.
			WithExec([]string{"cargo", "fmt", "--check"}).
			Sync(ctx)
		return err
	})

	eg.Go(func() error {
		before := r.Dagger.Source()
		after, err := r.Generate(ctx)
		if err != nil {
			return err
		}
		return dag.Dirdiff().AssertEqual(ctx, before, after, []string{"sdk/rust"})
	})

	return eg.Wait()
}

// Test the Rust SDK
func (r RustSDK) Test(ctx context.Context) error {
	installer, err := r.Dagger.installer(ctx, "sdk-rust-test")
	if err != nil {
		return err
	}
	_, err = r.rustBase(rustDockerStable).
		With(installer).
		WithExec([]string{"rustc", "--version"}).
		WithExec([]string{"cargo", "test", "--release", "--all"}).
		Sync(ctx)
	return err
}

// Regenerate the Rust SDK API
func (r RustSDK) Generate(ctx context.Context) (*dagger.Directory, error) {
	installer, err := r.Dagger.installer(ctx, "sdk-rust-generate")
	if err != nil {
		return nil, err
	}

	generated := r.rustBase(rustDockerStable).
		With(installer).
		WithExec([]string{"cargo", "run", "-p", "dagger-bootstrap", "generate", "--output", fmt.Sprintf("/%s", rustGeneratedAPIPath)}).
		WithExec([]string{"cargo", "fix", "--all", "--allow-no-vcs"}).
		WithExec([]string{"cargo", "fmt"}).
		File(strings.TrimPrefix(rustGeneratedAPIPath, "sdk/rust/"))

	return dag.Directory().
		WithDirectory("sdk/rust", r.Dagger.Source().Directory("sdk/rust")).
		WithFile(rustGeneratedAPIPath, generated), nil
}

// Publish the Rust SDK
func (r RustSDK) Publish(
	ctx context.Context,
	tag string,

	// +optional
	dryRun bool,

	// +optional
	cargoRegistryToken *dagger.Secret,
) error {
	version := strings.TrimPrefix(tag, "sdk/rust/v")
	versionArg := version
	if dryRun {
		// just pick any version, it's a dry-run
		versionArg = "--bump=rc"
	}

	crate := "dagger-sdk"

	base := r.
		withVersion(version).
		rustBase(rustDockerStable).
		WithExec([]string{
			"cargo", "install", "cargo-edit", "--locked",
		}).
		WithExec([]string{
			"cargo", "set-version", "-p", crate, versionArg,
		})
	args := []string{
		"cargo", "publish", "-p", crate, "-v", "--all-features",
	}

	if dryRun {
		args = append(args, "--dry-run")
		base = base.WithExec(args)
	} else {
		base = base.
			WithSecretVariable("CARGO_REGISTRY_TOKEN", cargoRegistryToken).
			WithExec(args)
	}

	_, err := base.Sync(ctx)
	return err
}

// Bump the Rust SDK's Engine dependency
//
// Deprecated: this is now included in the Publish step
func (r RustSDK) Bump(
	ctx context.Context,
	// +optional
	version string,
) (*dagger.Directory, error) {
	result := bumpSDK("rust", version, r.Dagger.Source())
	return r.Dagger.Source().Diff(result), nil
}

func (r RustSDK) rustBase(image string) *dagger.Container {
	const appDir = "sdk/rust"

	src := bumpSDK("rust", r.Version, r.Dagger.Source()).Directory(appDir)

	mountPath := fmt.Sprintf("/%s", appDir)

	base := dag.Container().
		From(image).
		WithDirectory(mountPath, src, dagger.ContainerWithDirectoryOpts{
			Include: []string{
				"**/Cargo.toml",
				"**/Cargo.lock",
				"**/main.rs",
				"**/lib.rs",
			},
		}).
		WithWorkdir(mountPath).
		WithEnvVariable("CARGO_HOME", "/root/.cargo").
		WithMountedCache("/root/.cargo", dag.CacheVolume("rust-cargo-"+image)).
		// combine into one layer so there's no assumptions on state of cache volume across steps
		With(util.ShellCmds(
			"rustup component add rustfmt",
			"cargo install --locked cargo-chef@"+cargoChefVersion,
			"cargo chef prepare --recipe-path /tmp/recipe.json",
			"cargo chef cook --release --workspace --recipe-path /tmp/recipe.json",
		)).
		WithMountedDirectory(mountPath, src)

	return base
}
