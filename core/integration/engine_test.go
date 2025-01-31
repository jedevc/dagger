package core

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	bkconfig "github.com/moby/buildkit/cmd/buildkitd/config"
	"github.com/moby/buildkit/identity"
	bkresolverconfig "github.com/moby/buildkit/util/resolver/config"
	"github.com/pelletier/go-toml"

	"dagger.io/dagger"
	"github.com/dagger/dagger/engine"
	"github.com/dagger/dagger/engine/config"
	"github.com/dagger/dagger/engine/distconsts"
	"github.com/dagger/dagger/internal/testutil"
	"github.com/dagger/dagger/testctx"
	"github.com/stretchr/testify/require"
)

type EngineSuite struct{}

func TestEngine(t *testing.T) {
	testctx.Run(testCtx, t, EngineSuite{}, Middleware()...)
}

func devEngineContainerAsService(ctr *dagger.Container) *dagger.Service {
	return ctr.AsService(dagger.ContainerAsServiceOpts{
		UseEntrypoint:            true,
		InsecureRootCapabilities: true,
	})
}

// devEngineContainer returns a nested dev engine.
func devEngineContainer(c *dagger.Client, withs ...func(*dagger.Container) *dagger.Container) *dagger.Container {
	// This loads the engine.tar file from the host into the container, that
	// was set up by the test caller. This is used to spin up additional dev
	// engines.
	var tarPath string
	if v, ok := os.LookupEnv("_DAGGER_TESTS_ENGINE_TAR"); ok {
		tarPath = v
	} else {
		tarPath = "./bin/engine.tar"
	}
	devEngineTar := c.Host().File(tarPath)

	ctr := c.Container().Import(devEngineTar)
	for _, with := range withs {
		ctr = with(ctr)
	}

	deviceName, cidr := testutil.GetUniqueNestedEngineNetwork()
	return ctr.
		WithMountedCache("/var/lib/dagger", c.CacheVolume("dagger-dev-engine-state-"+identity.NewID())).
		WithExposedPort(1234, dagger.ContainerWithExposedPortOpts{Protocol: dagger.NetworkProtocolTcp}).
		WithDefaultArgs([]string{
			"--addr", "tcp://0.0.0.0:1234",
			"--addr", "unix:///var/run/buildkit/buildkitd.sock",
			// avoid network conflicts with other tests
			"--network-name", deviceName,
			"--network-cidr", cidr,
		})
}

func engineWithConfig(ctx context.Context, t *testctx.T, cfgFns ...func(context.Context, *testctx.T, config.Config) config.Config) func(*dagger.Container) *dagger.Container {
	return func(ctr *dagger.Container) *dagger.Container {
		t.Helper()

		var cfg config.Config

		entries, err := ctr.Directory("/etc/dagger").Entries(ctx)
		require.NoError(t, err)
		if slices.Contains(entries, "engine.json") {
			existingCfgStr, err := ctr.File("/etc/dagger/engine.json").Contents(ctx)
			require.NoError(t, err)
			cfg, err = config.Load(strings.NewReader(existingCfgStr))
			require.NoError(t, err)
		}

		for _, cfgFn := range cfgFns {
			cfg = cfgFn(ctx, t, cfg)
		}

		var buf bytes.Buffer
		require.NoError(t, cfg.Save(&buf))
		return ctr.WithNewFile("/etc/dagger/engine.json", buf.String())
	}
}

func engineWithBkConfig(ctx context.Context, t *testctx.T, cfgFns ...func(context.Context, *testctx.T, bkconfig.Config) bkconfig.Config) func(*dagger.Container) *dagger.Container {
	return func(ctr *dagger.Container) *dagger.Container {
		t.Helper()

		var cfg bkconfig.Config

		entries, err := ctr.Directory("/etc/dagger").Entries(ctx)
		require.NoError(t, err)
		if slices.Contains(entries, "engine.toml") {
			existingCfgStr, err := ctr.File("/etc/dagger/engine.toml").Contents(ctx)
			require.NoError(t, err)

			cfg, err = bkconfig.Load(strings.NewReader(existingCfgStr))
			require.NoError(t, err)
		}

		for _, cfgFn := range cfgFns {
			cfg = cfgFn(ctx, t, cfg)
		}

		newCfgBytes, err := toml.Marshal(cfg)
		require.NoError(t, err)

		return ctr.WithNewFile("/etc/dagger/engine.toml", string(newCfgBytes))
	}
}

func engineClientContainer(ctx context.Context, t *testctx.T, c *dagger.Client, devEngine *dagger.Service) *dagger.Container {
	daggerCli := daggerCliFile(t, c)

	cliBinPath := "/bin/dagger"
	endpoint, err := devEngine.Endpoint(ctx, dagger.ServiceEndpointOpts{Port: 1234, Scheme: "tcp"})
	require.NoError(t, err)
	return c.Container().From(alpineImage).
		WithServiceBinding("dev-engine", devEngine).
		WithMountedFile(cliBinPath, daggerCli).
		WithEnvVariable("_EXPERIMENTAL_DAGGER_CLI_BIN", cliBinPath).
		WithEnvVariable("_EXPERIMENTAL_DAGGER_RUNNER_HOST", endpoint)
}

func (EngineSuite) TestExitsZeroOnSignal(ctx context.Context, t *testctx.T) {
	c := connect(ctx, t)

	// engine should shutdown with exit code 0 when receiving SIGTERM
	ctr := devEngineContainer(c, func(c *dagger.Container) *dagger.Container {
		t.Helper()

		c = c.WithNewFile(
			"/usr/local/bin/dagger-entrypoint.sh",
			`#!/bin/sh
set -ex
/usr/local/bin/dagger-engine --debug &
engine_pid=$!

sleep 5
kill -TERM $engine_pid
wait $engine_pid
exit $?
`,
			dagger.ContainerWithNewFileOpts{Permissions: 0o700},
		)

		// do a sync here so our timeout doesn't include overhead of importing the engine itself
		var err error
		c, err = c.Sync(ctx)
		require.NoError(t, err)
		return c
	})

	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	t = t.WithContext(ctx)
	_, err := ctr.Sync(ctx)
	require.NoError(t, err)
}

func (EngineSuite) TestSetsNameFromEnv(ctx context.Context, t *testctx.T) {
	c := connect(ctx, t)

	engineName := "my-special-engine"
	engineVersion := engine.Version + "-special"
	devEngineSvc := devEngineContainerAsService(devEngineContainer(c, func(c *dagger.Container) *dagger.Container {
		return c.
			WithEnvVariable("_EXPERIMENTAL_DAGGER_ENGINE_NAME", engineName).
			WithEnvVariable("_EXPERIMENTAL_DAGGER_VERSION", engineVersion)
	}))

	clientCtr := engineClientContainer(ctx, t, c, devEngineSvc)

	clientCtr = clientCtr.
		WithNewFile("/query.graphql", `{ version }`).
		WithExec([]string{"dagger", "query", "--doc", "/query.graphql"})
	stdout, err := clientCtr.Stdout(ctx)
	require.NoError(t, err)
	stderr, err := clientCtr.Stderr(ctx)
	require.NoError(t, err)

	require.Contains(t, stderr, engineName)
	require.Contains(t, stderr, engineVersion)

	require.Contains(t, stdout, engineVersion)
}

func (EngineSuite) TestDaggerRun(ctx context.Context, t *testctx.T) {
	c := connect(ctx, t)

	devEngine := devEngineContainerAsService(devEngineContainer(c))

	clientCtr := engineClientContainer(ctx, t, c, devEngine)

	command := fmt.Sprintf(`
		export NO_COLOR=1
		jq -n '{query:"{container{from(address: \"%s\"){file(path: \"/etc/alpine-release\"){contents}}}}"}' | \
		dagger run sh -c 'curl -s \
			-u $DAGGER_SESSION_TOKEN: \
			--max-time 30 \
			-H "content-type:application/json" \
			-d @- \
			http://127.0.0.1:$DAGGER_SESSION_PORT/query'`,
		alpineImage,
	)

	clientCtr = clientCtr.
		WithExec([]string{"apk", "add", "jq", "curl"}).
		WithExec([]string{"sh", "-c", command})

	stdout, err := clientCtr.Stdout(ctx)
	require.NoError(t, err)
	require.Contains(t, stdout, distconsts.AlpineVersion)
	require.JSONEq(t, `{"data": {"container": {"from": {"file": {"contents": "`+distconsts.AlpineVersion+`\n"}}}}}`, stdout)

	stderr, err := clientCtr.Stderr(ctx)
	require.NoError(t, err)
	// verify we got some progress output
	require.Contains(t, stderr, "Container.from")
}

func (EngineSuite) TestVersionCompat(ctx context.Context, t *testctx.T) {
	c := connect(ctx, t)

	tcs := []struct {
		name string

		engineVersion    string
		engineMinVersion string
		clientVersion    string
		clientMinVersion string

		errs []string
	}{
		{
			// v2.0.0 > v1.0.0 for both client and engine
			name:             "compatible",
			engineVersion:    "v2.0.0",
			engineMinVersion: "v1.0.0",
			clientVersion:    "v2.0.0",
			clientMinVersion: "v1.0.0",
		},
		{
			// v2.0.0 > v1.0.0 for both client and engine
			name:             "compatible",
			engineVersion:    "v2.0.0",
			engineMinVersion: "v2.0.0",
			clientVersion:    "v2.0.0",
			clientMinVersion: "v2.0.0",
		},
		{
			// v2.0.0 < v3.0.0 for the client
			name:             "client incompatible",
			engineVersion:    "v2.0.0",
			engineMinVersion: "v1.0.0",
			clientVersion:    "v2.0.0",
			clientMinVersion: "v3.0.0",
			errs: []string{
				"incompatible client version v2.0.0",
			},
		},
		{
			// v2.0.0 < v3.0.0 for the engine
			name:             "engine incompatible",
			engineVersion:    "v2.0.0",
			engineMinVersion: "v3.0.0",
			clientVersion:    "v2.0.0",
			clientMinVersion: "v1.0.0",
			errs: []string{
				"incompatible engine version v2.0.0",
			},
		},
		{
			// v2.0.0 < v3.0.0 for both client and engine
			name:             "client and engine incompatible",
			engineVersion:    "v2.0.0",
			engineMinVersion: "v3.0.0",
			clientVersion:    "v2.0.0",
			clientMinVersion: "v3.0.0",
			errs: []string{
				"incompatible engine version v2.0.0",
			},
		},
		{
			// v2.0.1-foobar > v2.0.0 for both client and engine
			name:             "new dev version",
			engineVersion:    "v2.0.1-foobar",
			engineMinVersion: "v2.0.0",
			clientVersion:    "v2.0.1-foobar",
			clientMinVersion: "v2.0.0",
		},
		{
			// v2.0.1-foobar > v2.0.0 for both client and engine
			name:             "old dev version",
			engineVersion:    "v2.0.0-foobar",
			engineMinVersion: "v2.0.0",
			clientVersion:    "v2.0.0-foobar",
			clientMinVersion: "v2.0.0",
			errs: []string{
				"incompatible engine version v2.0.0-foobar",
			},
		},

		{
			// dev versions for the same version are happily compatible (even
			// if not a perfect match)
			name:             "compatible dev versions",
			engineVersion:    "v2.0.0-dev-123",
			engineMinVersion: "v2.0.0-dev-456",
			clientVersion:    "v2.0.0-dev-456",
			clientMinVersion: "v2.0.0-dev-123",
		},
		{
			// but for different versions, they're incompatible
			name:             "incompatible dev versions",
			engineVersion:    "v2.0.0-dev-123",
			engineMinVersion: "v2.0.1-dev-456",
			clientVersion:    "v2.0.1-dev-456",
			clientMinVersion: "v2.0.0-dev-123",
			errs: []string{
				"incompatible engine version v2.0.0-dev-123",
			},
		},

		{
			// pre-releases match if they're exactly the same
			name:             "compatible prereleases",
			engineVersion:    "v2.0.0-foo-123",
			engineMinVersion: "v2.0.0-foo-123",
			clientVersion:    "v2.0.0-foo-123",
			clientMinVersion: "v2.0.0-foo-123",
		},
		{
			// but can't not be a perfect match (unlike dev versions)
			name:             "incompatible prereleases",
			engineVersion:    "v2.0.0-foo-123",
			engineMinVersion: "v2.0.0-foo-456",
			clientVersion:    "v2.0.0-foo-456",
			clientMinVersion: "v2.0.0-foo-123",
			errs: []string{
				"incompatible engine version v2.0.0-foo-123",
			},
		},

		// empty clients/engines can happen with manual builds
		{
			name:             "compatible empty client",
			engineVersion:    "v2.0.0",
			engineMinVersion: "v1.0.0",
			clientVersion:    "",
			clientMinVersion: "v1.0.0",
		},
		{
			name:             "compatible empty engine",
			engineVersion:    "",
			engineMinVersion: "v1.0.0",
			clientVersion:    "v2.0.0",
			clientMinVersion: "v1.0.0",
		},
	}

	engines := map[string]*dagger.Service{}
	enginesMu := sync.Mutex{}

	for _, tc := range tcs {
		// get a cached engine if possible (saves spinning up more engines than we need to)
		tc := tc
		t.Run(tc.name, func(ctx context.Context, t *testctx.T) {
			devEngineSvcKey := tc.engineVersion + " " + tc.clientMinVersion
			enginesMu.Lock()
			devEngineSvc, ok := engines[devEngineSvcKey]
			if !ok {
				devEngine := devEngineContainer(c, func(c *dagger.Container) *dagger.Container {
					return c.
						WithEnvVariable("_EXPERIMENTAL_DAGGER_VERSION", tc.engineVersion).
						WithEnvVariable("_EXPERIMENTAL_DAGGER_MIN_VERSION", tc.clientMinVersion)
				})
				devEngineSvc = devEngineContainerAsService(devEngine)
				engines[devEngineSvcKey] = devEngineSvc
			}
			enginesMu.Unlock()

			clientCtr := engineClientContainer(ctx, t, c, devEngineSvc)

			clientCtr = clientCtr.
				WithEnvVariable("_EXPERIMENTAL_DAGGER_VERSION", tc.clientVersion).
				WithEnvVariable("_EXPERIMENTAL_DAGGER_MIN_VERSION", tc.engineMinVersion)

			if tc.errs == nil {
				clientCtr = clientCtr.
					WithNewFile("/query.graphql", `{ version }`).
					WithExec([]string{"sh", "-c", "dagger version && dagger query --doc /query.graphql"})
			} else {
				clientCtr = clientCtr.
					WithNewFile("/query.graphql", `{ version }`).
					WithExec([]string{"sh", "-c", "! dagger query --doc /query.graphql"})
			}

			if tc.errs == nil {
				stdout, err := clientCtr.Stdout(ctx)
				require.NoError(t, err)

				// check that both the client and engine versions appear
				// somewhere in the combined output
				require.Contains(t, stdout, tc.clientVersion)
				require.Contains(t, stdout, tc.engineVersion)
			} else {
				stderr, err := clientCtr.Stderr(ctx)
				require.NoError(t, err)

				// check the error is contained
				for _, tcerr := range tc.errs {
					require.Contains(t, stderr, tcerr)
				}
			}
		})
	}
}

func (EngineSuite) TestModuleVersionCompat(ctx context.Context, t *testctx.T) {
	c := connect(ctx, t)

	tcs := []struct {
		name string

		engineVersion    string
		moduleVersion    string
		moduleMinVersion string

		errs []string
	}{
		{
			name:             "compatible equal",
			engineVersion:    "v2.0.0",
			moduleVersion:    "v2.0.0",
			moduleMinVersion: "v1.0.0",
		},
		{
			name:             "compatible less",
			engineVersion:    "v2.0.0",
			moduleVersion:    "v1.0.0",
			moduleMinVersion: "v1.0.0",
		},
		{
			name:             "incompatible too old",
			engineVersion:    "v2.0.0",
			moduleVersion:    "v0.9.0",
			moduleMinVersion: "v1.0.0",
			errs: []string{
				"module requires dagger v0.9.0",
				"support for that version has been removed",
			},
		},
		{
			name:             "incompatible too new",
			engineVersion:    "v2.0.0",
			moduleVersion:    "v2.0.1",
			moduleMinVersion: "v1.0.0",
			errs: []string{
				"module requires dagger v2.0.1, but you have v2.0.0",
			},
		},
		{
			name:             "old style dev version",
			engineVersion:    "v2.0.0",
			moduleVersion:    "badbadbad",
			moduleMinVersion: "v1.0.0",
			errs: []string{
				"module requires dagger v0.11.9", // old-style dev versions are equivalent to v0.11.9
				"support for that version has been removed",
			},
		},
	}

	engines := map[string]*dagger.Service{}
	enginesMu := sync.Mutex{}

	for _, tc := range tcs {
		// get a cached engine if possible (saves spinning up more engines than we need to)
		tc := tc
		t.Run(tc.name, func(ctx context.Context, t *testctx.T) {
			devEngineSvcKey := tc.engineVersion + " " + tc.moduleMinVersion
			enginesMu.Lock()
			devEngineSvc, ok := engines[devEngineSvcKey]
			if !ok {
				devEngine := devEngineContainer(c, func(c *dagger.Container) *dagger.Container {
					return c.
						WithEnvVariable("_EXPERIMENTAL_DAGGER_VERSION", tc.engineVersion).
						WithEnvVariable("_EXPERIMENTAL_DAGGER_MIN_VERSION", tc.moduleMinVersion)
				})
				devEngineSvc = devEngineContainerAsService(devEngine)
				engines[devEngineSvcKey] = devEngineSvc
			}
			enginesMu.Unlock()

			clientCtr := engineClientContainer(ctx, t, c, devEngineSvc)

			clientCtr = clientCtr.
				WithWorkdir("/work").
				// set version to empty, this makes it the latest, we don't want to
				// test client compat (that's the previous tests)
				WithEnvVariable("_EXPERIMENTAL_DAGGER_VERSION", "").
				With(daggerExec("init", "--name=bare", "--sdk=go"))

			clientCtr = clientCtr.
				WithNewFile("/work/dagger.json", `{"name": "bare", "sdk": "go", "engineVersion": "`+tc.moduleVersion+`"}`).
				WithNewFile("/query.graphql", `{bare{containerEcho(stringArg:"hello"){stdout}}}`)

			if tc.errs == nil {
				clientCtr = clientCtr.
					WithExec([]string{"sh", "-c", "dagger query --doc /query.graphql"})
			} else {
				clientCtr = clientCtr.
					WithExec([]string{"sh", "-c", "! dagger query --doc /query.graphql"})
			}

			stderr, err := clientCtr.Stderr(ctx)
			require.NoError(t, err)
			for _, tcerr := range tc.errs {
				require.Contains(t, stderr, tcerr)
			}
		})
	}
}

func (EngineSuite) TestRegistryMirrors(ctx context.Context, t *testctx.T) {
	c := connect(ctx, t)

	imageName := "fakeimagethatisnotreal"
	_, err := c.Container().From("alpine").Publish(ctx, registryRef(imageName))
	require.NoError(t, err)

	f := func(ctx context.Context, t *testctx.T, engine *dagger.Container) {
		engineSvc, err := c.Host().Tunnel(devEngineContainerAsService(engine)).Start(ctx)
		require.NoError(t, err)
		t.Cleanup(func() { engineSvc.Stop(ctx) })

		endpoint, err := engineSvc.Endpoint(ctx, dagger.ServiceEndpointOpts{Scheme: "tcp"})
		require.NoError(t, err)

		c, err := dagger.Connect(ctx, dagger.WithRunnerHost(endpoint), dagger.WithLogOutput(testutil.NewTWriter(t)))
		require.NoError(t, err)
		t.Cleanup(func() { c.Close() })

		result, err := c.Container().From("mymirror.local/fakeimagethatisnotreal").WithExec([]string{"whoami"}).Stdout(ctx)
		require.NoError(t, err)
		require.Equal(t, "root", result)
	}

	// t.Run("engine.json", func(ctx context.Context, t *testctx.T) {
	// 	cfg := engineWithConfig(ctx, t, engineConfigWithMirror("mymirror.local", registryHost))
	// 	f(ctx, t, devEngineContainer(c, cfg))
	// })

	t.Run("engine.toml", func(ctx context.Context, t *testctx.T) {
		cfg := engineWithBkConfig(ctx, t, bkConfigWithMirror("mymirror.local", registryHost))
		f(ctx, t, devEngineContainer(c, cfg))
	})
}

func engineConfigWithMirror(host string, target string) func(context.Context, *testctx.T, config.Config) config.Config {
	return func(ctx context.Context, t *testctx.T, cfg config.Config) config.Config {
		t.Helper()
		cfg.Registries = append(cfg.Registries, config.Registry{Host: host, Mirrors: []string{target}})
		return cfg
	}
}

func bkConfigWithMirror(host string, target string) func(context.Context, *testctx.T, bkconfig.Config) bkconfig.Config {
	return func(ctx context.Context, t *testctx.T, cfg bkconfig.Config) bkconfig.Config {
		t.Helper()
		cfg.Registries[host] = bkresolverconfig.RegistryConfig{
			Mirrors: []string{target},
		}
		return cfg
	}
}
