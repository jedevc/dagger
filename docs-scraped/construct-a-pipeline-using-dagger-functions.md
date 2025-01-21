# Construct a pipeline using Dagger Functions

Dagger Functions are regular code, written in Go, Python or TypeScript using the corresponding Dagger SDK. They consist of a series of method/function calls, such as "pull a container image", "copy a file", "forward a TCP port", and so on, which can be chained together.

:::note
Don't worry about how the Dagger Functions shown below work for the moment - it's explained in detail in  the next sections!
:::

```go
package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"

	"dagger/hello-dagger/internal/dagger"
)

type HelloDagger struct{}

// Publish the application container after building and testing it on-the-fly
func (m *HelloDagger) Publish(ctx context.Context, source *dagger.Directory) (string, error) {
	_, err := m.Test(ctx, source)
	if err != nil {
		return "", err
	}
	return m.Build(source).
		Publish(ctx, fmt.Sprintf("ttl.sh/hello-dagger-%.0f", math.Floor(rand.Float64()*10000000))) //#nosec
}

// Build the application container
func (m *HelloDagger) Build(source *dagger.Directory) *dagger.Container {
	build := m.BuildEnv(source).
		WithExec([]string{"npm", "run", "build"}).
		Directory("./dist")
	return dag.Container().From("nginx:1.25-alpine").
		WithDirectory("/usr/share/nginx/html", build).
		WithExposedPort(80)
}

// Return the result of running unit tests
func (m *HelloDagger) Test(ctx context.Context, source *dagger.Directory) (string, error) {
	return m.BuildEnv(source).
		WithExec([]string{"npm", "run", "test:unit", "run"}).
		Stdout(ctx)
}

// Build a ready-to-use development environment
func (m *HelloDagger) BuildEnv(source *dagger.Directory) *dagger.Container {
	nodeCache := dag.CacheVolume("node")
	return dag.Container().
		From("node:21-slim").
		WithDirectory("/src", source).
		WithMountedCache("/root/.npm", nodeCache).
		WithWorkdir("/src").
		WithExec([]string{"npm", "install"})
}
```

```python
import random

import dagger
from dagger import dag, function, object_type


@object_type
class HelloDagger:
    @function
    async def publish(self, source: dagger.Directory) -> str:
        """Publish the application container after building and testing it on-the-fly"""
        await self.test(source)
        return await self.build(source).publish(
            f"ttl.sh/myapp-{random.randrange(10**8)}"
        )

    @function
    def build(self, source: dagger.Directory) -> dagger.Container:
        """Build the application container"""
        build = (
            self.build_env(source)
            .with_exec(["npm", "run", "build"])
            .directory("./dist")
        )
        return (
            dag.container()
            .from_("nginx:1.25-alpine")
            .with_directory("/usr/share/nginx/html", build)
            .with_exposed_port(80)
        )

    @function
    async def test(self, source: dagger.Directory) -> str:
        """Return the result of running unit tests"""
        return await (
            self.build_env(source)
            .with_exec(["npm", "run", "test:unit", "run"])
            .stdout()
        )

    @function
    def build_env(self, source: dagger.Directory) -> dagger.Container:
        """Build a ready-to-use development environment"""
        node_cache = dag.cache_volume("node")
        return (
            dag.container()
            .from_("node:21-slim")
            .with_directory("/src", source)
            .with_mounted_cache("/root/.npm", node_cache)
            .with_workdir("/src")
            .with_exec(["npm", "install"])
        )
```

```typescript
import { dag, Container, Directory, object, func } from "@dagger.io/dagger"

@object()
class HelloDagger {
  /**
   * Publish the application container after building and testing it on-the-fly
   */
  @func()
  async publish(source: Directory): Promise<string> {
    await this.test(source)
    return await this.build(source).publish(
      "ttl.sh/myapp-" + Math.floor(Math.random() * 10000000),
    )
  }

  /**
   * Build the application container
   */
  @func()
  build(source: Directory): Container {
    const build = this.buildEnv(source)
      .withExec(["npm", "run", "build"])
      .directory("./dist")
    return dag
      .container()
      .from("nginx:1.25-alpine")
      .withDirectory("/usr/share/nginx/html", build)
      .withExposedPort(80)
  }

  /**
   * Return the result of running unit tests
   */
  @func()
  async test(source: Directory): Promise<string> {
    return this.buildEnv(source)
      .withExec(["npm", "run", "test:unit", "run"])
      .stdout()
  }

  /**
   * Build a ready-to-use development environment
   */
  @func()
  buildEnv(source: Directory): Container {
    const nodeCache = dag.cacheVolume("node")
    return dag
      .container()
      .from("node:21-slim")
      .withDirectory("/src", source)
      .withMountedCache("/root/.npm", nodeCache)
      .withWorkdir("/src")
      .withExec(["npm", "install"])
  }
}
```

