package main

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dagger/dagger/dev/internal/dagger"
	"github.com/dagger/dagger/dev/internal/util"
)

const (
	phpSDKPath         = "sdk/php"
	phpSDKGeneratedDir = "generated"
	phpSDKVersionFile  = "src/Connection/version.php"
)

type PHPSDK struct {
	Dagger *DaggerDev // +private
}

// Lint the PHP SDK
func (t PHPSDK) Lint(ctx context.Context) error {
	before := t.Dagger.Source()
	after, err := t.Generate(ctx)
	if err != nil {
		return err
	}
	return dag.Dirdiff().AssertEqual(ctx, before, after, []string{filepath.Join(phpSDKPath, phpSDKGeneratedDir)})
}

// Test the PHP SDK
func (t PHPSDK) Test(ctx context.Context) error {
	return fmt.Errorf("not implemented")
}

// Regenerate the PHP SDK API
func (t PHPSDK) Generate(ctx context.Context) (*dagger.Directory, error) {
	installer, err := t.Dagger.installer(ctx, "sdk-php-generate")
	if err != nil {
		return nil, err
	}

	generated := t.phpBase().
		With(installer).
		With(util.ShellCmds(
			fmt.Sprintf("rm -f %s/*.php", phpSDKGeneratedDir),
			"ls -lha",
			"$_EXPERIMENTAL_DAGGER_CLI_BIN run ./codegen",
		)).
		Directory(".")
	return dag.Directory().WithDirectory(phpSDKPath, generated), nil
}

// Publish the PHP SDK
func (t PHPSDK) Publish(
	ctx context.Context,
	tag string,

	// +optional
	dryRun bool,

	// +optional
	// +default="https://github.com/dagger/dagger-php-sdk.git"
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
		sdk:         "php",
		source:      "https://github.com/dagger/dagger.git",
		sourcePath:  "sdk/php/",
		sourceTag:   tag,
		dest:        gitRepo,
		destTag:     strings.TrimPrefix(tag, "sdk/php/"),
		username:    gitUserName,
		email:       gitUserEmail,
		githubToken: githubToken,
		dryRun:      dryRun,
	})
}

// Bump the PHP SDK's Engine dependency
//
// Deprecated: this is now included in the Publish step
func (t PHPSDK) Bump(
	ctx context.Context,
	// +optional
	version string,
) (*dagger.Directory, error) {
	result := bumpSDK("php", version, t.Dagger.Source())
	return t.Dagger.Source().Diff(result), nil
}

// phpBase returns a PHP container with the PHP SDK source files
// added and dependencies installed.
func (t PHPSDK) phpBase() *dagger.Container {
	src := t.Dagger.Source().Directory(phpSDKPath)
	return dag.Container().
		From("php:8.2-zts-bookworm").
		WithExec([]string{"apt-get", "update"}).
		WithExec([]string{"apt-get", "install", "-y", "git", "unzip"}).
		WithFile("/usr/bin/composer", dag.Container().
			From("composer:2").
			File("/usr/bin/composer"),
		).
		WithMountedCache("/root/.composer", dag.CacheVolume("composer-cache-8.2-zts-bookworm")).
		WithEnvVariable("COMPOSER_HOME", "/root/.composer").
		WithEnvVariable("COMPOSER_ALLOW_SUPERUSER", "1").
		WithWorkdir(fmt.Sprintf("/%s", phpSDKPath)).
		WithFile("composer.json", src.File("composer.json")).
		WithFile("composer.lock", src.File("composer.lock")).
		WithExec([]string{"composer", "install"}).
		WithDirectory(".", src)
}
