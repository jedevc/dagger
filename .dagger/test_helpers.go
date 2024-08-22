package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/remotes/docker"
	"github.com/docker/cli/cli/config/configfile"
	"github.com/moby/buildkit/util/contentutil"
	"github.com/moby/buildkit/version"
	ocispecs "github.com/opencontainers/image-spec/specs-go/v1"

	"github.com/dagger/dagger/.dagger/internal/dagger"
)

func toCredentialsFunc(dt string) func(string) (string, string, error) {
	cfg := configfile.New("config.json")
	err := cfg.LoadFromReader(strings.NewReader(dt))
	if err != nil {
		panic(err)
	}

	return func(host string) (string, string, error) {
		if host == "registry-1.docker.io" {
			host = "https://index.docker.io/v1/"
		}
		ac, err := cfg.GetAuthConfig(host)
		if err != nil {
			return "", "", err
		}
		if ac.IdentityToken != "" {
			return "", ac.IdentityToken, nil
		}
		return ac.Username, ac.Password, nil
	}
}

func (t *Test) setupMirror(ctx context.Context) (*dagger.Service, error) {
	mirror := registryMirror()
	mirror, err := mirror.Start(ctx)
	if err != nil {
		return nil, err
	}
	go mirror.Up(ctx, dagger.ServiceUpOpts{
		Ports: []dagger.PortForward{{
			Backend:  5000,
			Frontend: 5000,
		}},
	})

	var auth docker.Authorizer
	if t.Dagger.DockerCfg != nil {
		hostConfig, err := t.Dagger.DockerCfg.Plaintext(ctx)
		if err != nil {
			return nil, err
		}
		auth = docker.NewDockerAuthorizer(docker.WithAuthCreds(toCredentialsFunc(hostConfig)), docker.WithAuthClient(http.DefaultClient))
	}

	err = copyImagesLocal(auth, "localhost:5000", map[string]string{
		"library/alpine:latest": "docker.io/library/alpine:3.20.1",
		"library/alpine:3":      "docker.io/library/alpine:3.20.1",
		"library/alpine:3.20":   "docker.io/library/alpine:3.20.1",
		"library/alpine:3.20.1": "docker.io/library/alpine:3.20.1",

		// "library/registry:2": "docker.io/library/registry:2",

		// "library/golang:latest":        "docker.io/library/golang:1.22.2-alpine",
		// "library/golang:1.22":          "docker.io/library/golang:1.22.2-alpine",
		// "library/golang:1.22.2":        "docker.io/library/golang:1.22.2-alpine",
		// "library/golang:alpine":        "docker.io/library/golang:1.22.2-alpine",
		// "library/golang:1.22-alpine":   "docker.io/library/golang:1.22.2-alpine",
		// "library/golang:1.22.2-alpine": "docker.io/library/golang:1.22.2-alpine",

		// "library/python:latest":    "docker.io/library/python:3.11-slim",
		// "library/python:3":         "docker.io/library/python:3.11-slim",
		// "library/python:3.11":      "docker.io/library/python:3.11-slim",
		// "library/python:3.11-slim": "docker.io/library/python:3.11-slim",
	})
	if err != nil {
		return nil, err
	}

	return mirror, nil
}

func registryMirror() *dagger.Service {
	return dag.Container().
		From("registry:2").
		WithEnvVariable("LOG_LEVEL", "warn").
		WithEnvVariable("REGISTRY_STORAGE_FILESYSTEM_ROOTDIRECTORY", "/data/registry").
		WithMountedCache("/data/registry", dag.CacheVolume("dagger-registry-mirror-cache")).
		WithExposedPort(5000, dagger.ContainerWithExposedPortOpts{Protocol: dagger.Tcp}).
		WithExec(nil, dagger.ContainerWithExecOpts{
			// without this, we get a context cancelled error, which is unreadable :)
			UseEntrypoint: true,
		}).
		AsService()
}

var localImageCache map[string]map[string]struct{}

func copyImagesLocal(auth docker.Authorizer, host string, images map[string]string) error {
	for to, from := range images {
		if localImageCache == nil {
			localImageCache = map[string]map[string]struct{}{}
		}
		if _, ok := localImageCache[host]; !ok {
			localImageCache[host] = map[string]struct{}{}
		}
		if _, ok := localImageCache[host][to]; ok {
			continue
		}
		localImageCache[host][to] = struct{}{}

		start := time.Now()

		var desc ocispecs.Descriptor
		var provider content.Provider
		var err error
		desc, provider, err = ProviderFromRef(from, auth)
		if err != nil {
			return err
		}

		// already exists check
		_, _, err = docker.NewResolver(docker.ResolverOptions{}).Resolve(context.TODO(), host+"/"+to)
		if err == nil {
			fmt.Printf("copied %s to local mirror %s (skipped)\n", from, host+"/"+to)
			continue
		}

		ingester, err := contentutil.IngesterFromRef(host + "/" + to)
		if err != nil {
			return err
		}
		if err := contentutil.CopyChain(context.TODO(), ingester, provider, desc); err != nil {
			return err
		}
		fmt.Printf("copied %s to local mirror %s in %s\n", from, host+"/"+to, time.Since(start))
	}
	return nil
}

// XXX: I think this is copied from buildkit?
func ProviderFromRef(ref string, auth docker.Authorizer) (ocispecs.Descriptor, content.Provider, error) {
	headers := http.Header{}
	headers.Set("User-Agent", version.UserAgent())
	remote := docker.NewResolver(docker.ResolverOptions{
		Headers:    headers,
		Authorizer: auth,
	})

	name, desc, err := remote.Resolve(context.TODO(), ref)
	if err != nil {
		return ocispecs.Descriptor{}, nil, err
	}

	fetcher, err := remote.Fetcher(context.TODO(), name)
	if err != nil {
		return ocispecs.Descriptor{}, nil, err
	}
	return desc, contentutil.FromFetcher(fetcher), nil
}
