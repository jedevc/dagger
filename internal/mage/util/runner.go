package util

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"sort"
	"text/template"

	"dagger.io/dagger"
	"golang.org/x/exp/maps"
)

const (
	runnerBinName      = "dagger-runner"
	containerdTomlPath = "/etc/dagger/containerd.toml"

	// engineTomlPath     = "/etc/dagger/engine.toml"
	// EngineDefaultStateDir = "/var/lib/dagger"

	runnerEntrypointPath = "/usr/local/bin/runner-entrypoint.sh"
)

const runnerEntrypointTmpl = `#!/bin/sh
set -e

# cgroup v2: enable nesting
# see https://github.com/moby/moby/blob/38805f20f9bcc5e87869d6c79d432b166e1c88b4/hack/dind#L28
if [ -f /sys/fs/cgroup/cgroup.controllers ]; then
	# move the processes from the root group to the /init group,
	# otherwise writing subtree_control fails with EBUSY.
	# An error during moving non-existent process (i.e., "cat") is ignored.
	mkdir -p /sys/fs/cgroup/init
	xargs -rn1 < /sys/fs/cgroup/cgroup.procs > /sys/fs/cgroup/init/cgroup.procs || :
	# enable controllers
	sed -e 's/ / +/g' -e 's/^/+/' < /sys/fs/cgroup/cgroup.controllers \
		> /sys/fs/cgroup/cgroup.subtree_control
fi

{{.ContainerdBin}} --config {{.ContainerdConfig}} --log-level debug &
sleep 5
{{.RunnerBin}}
exit 0
`

const containerdConfig = `
version = 2
`

func getRunnerEntrypoint(opts ...DevEngineOpts) (string, error) {
	mergedOpts := map[string]string{}
	for _, opt := range opts {
		maps.Copy(mergedOpts, opt.EntrypointArgs)
	}
	keys := maps.Keys(mergedOpts)
	sort.Strings(keys)

	var entrypoint string

	type entrypointTmplParams struct {
		ContainerdBin    string
		ContainerdConfig string
		RunnerBin        string
	}
	tmpl := template.Must(template.New("entrypoint").Parse(runnerEntrypointTmpl))
	buf := new(bytes.Buffer)
	err := tmpl.Execute(buf, entrypointTmplParams{
		ContainerdBin:    "/usr/local/bin/containerd",
		ContainerdConfig: containerdTomlPath,
		RunnerBin:        "/usr/local/bin/" + runnerBinName,
	})
	if err != nil {
		panic(err)
	}
	entrypoint = buf.String()

	return entrypoint, nil
}

// DevRunnerContainer returns a container that runs a dev engine
func DevRunnerContainer(c *dagger.Client, arches []string, version string, opts ...DevEngineOpts) []*dagger.Container {
	return devRunnerContainers(c, arches, version, opts...)
}

func devRunnerContainer(c *dagger.Client, arch string, version string, opts ...DevEngineOpts) *dagger.Container {
	engine := devEngineContainer(c, arch, version, opts...)

	dir, err := os.MkdirTemp("", "runner")
	if err != nil {
		panic(err)
	}

	image := "image.tar"

	ctx := context.TODO()
	_, err = engine.Export(ctx, filepath.Join(dir, image))
	if err != nil {
		panic(err)
	}

	runnerEntrypoint, err := getRunnerEntrypoint(opts...)
	if err != nil {
		panic(err)
	}
	return c.Container(dagger.ContainerOpts{Platform: dagger.Platform("linux/" + arch)}).
		From("alpine:"+alpineVersion).
		WithExec([]string{
			"apk", "add",
			// for Buildkit
			"git", "openssh", "pigz", "xz",
			// for CNI
			"iptables", "ip6tables", "dnsmasq",
		}).
		WithFile("/usr/local/bin/runc", runcBin(c, arch), dagger.ContainerWithFileOpts{
			Permissions: 0o700,
		}).
		WithFile("/usr/local/bin/"+shimBinName, shimBin(c, arch)). // XXX: remove this
		WithFile("/usr/local/bin/"+runnerBinName, runnerBin(c, arch)).
		WithFile("/usr/local/bin/"+engineBinName, engineBin(c, arch, version)).
		WithDirectory("/usr/local/bin", qemuBins(c, arch)).
		WithDirectory("/usr/local/bin", containerdBin(c, arch)).
		WithDirectory("/", cniPlugins(c, arch)).
		WithDirectory(EngineDefaultStateDir, c.Directory()).
		WithDirectory("/var/lib/containerd", c.Directory()).
		WithNewFile(containerdTomlPath, dagger.ContainerWithNewFileOpts{
			Contents:    containerdConfig,
			Permissions: 0o600,
		}).
		WithNewFile(runnerEntrypointPath, dagger.ContainerWithNewFileOpts{
			Contents:    runnerEntrypoint,
			Permissions: 0o755,
		}).
		WithFile(image, c.Host().File(filepath.Join(dir, image))).
		WithEntrypoint([]string{filepath.Base(runnerEntrypointPath)})
}

func devRunnerContainers(c *dagger.Client, arches []string, version string, opts ...DevEngineOpts) []*dagger.Container {
	platformVariants := make([]*dagger.Container, 0, len(arches))
	for _, arch := range arches {
		platformVariants = append(platformVariants, devRunnerContainer(c, arch, version, opts...))
	}

	return platformVariants
}

// helper functions for building the dev engine container

func runnerBin(c *dagger.Client, arch string) *dagger.File {
	return goBase(c).
		WithEnvVariable("GOOS", "linux").
		WithEnvVariable("GOARCH", arch).
		WithExec([]string{
			"go", "build",
			"-o", "./bin/" + runnerBinName,
			"-ldflags", "-s -w",
			"/app/cmd/runner",
		}).
		File("./bin/" + runnerBinName)
}
