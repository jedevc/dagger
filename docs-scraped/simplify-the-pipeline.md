# Simplify the pipeline

Next, update the pipeline to use this module.

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
	address, err := m.Build(source).
		Publish(ctx, fmt.Sprintf("ttl.sh/hello-dagger-%.0f", math.Floor(rand.Float64()*10000000))) //#nosec
	if err != nil {
		return "", err
	}
	return address, nil
}

// Build the application container
func (m *HelloDagger) Build(source *dagger.Directory) *dagger.Container {
	build := dag.Node(dagger.NodeOpts{Ctr: m.BuildEnv(source)}).
		Commands().
		Run([]string{"build"}).
		Directory("./dist")
	return dag.Container().From("nginx:1.25-alpine").
		WithDirectory("/usr/share/nginx/html", build).
		WithExposedPort(80)
}

// Return the result of running unit tests
func (m *HelloDagger) Test(ctx context.Context, source *dagger.Directory) (string, error) {
	return dag.Node(dagger.NodeOpts{Ctr: m.BuildEnv(source)}).
		Commands().
		Run([]string{"test:unit", "run"}).
		Stdout(ctx)
}

// Build a ready-to-use development environment
func (m *HelloDagger) BuildEnv(source *dagger.Directory) *dagger.Container {
	return dag.Node(dagger.NodeOpts{Version: "21"}).
		WithNpm().
		WithSource(source).
		Install().
		Container()
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
            f"ttl.sh/hello-dagger-{random.randrange(10**8)}"
        )

    @function
    def build(self, source: dagger.Directory) -> dagger.Container:
        """Build the application container"""
        build = (
            dag.node(ctr=self.build_env(source))
            .commands()
            .run(["build"])
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
            dag.node(ctr=self.build_env(source))
            .commands()
            .run(["test:unit", "run"])
            .stdout()
        )

    @function
    def build_env(self, source: dagger.Directory) -> dagger.Container:
        """Build a ready-to-use development environment"""
        return (
            dag.node(version="21").with_npm().with_source(source).install().container()
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
      "ttl.sh/hello-dagger-" + Math.floor(Math.random() * 10000000),
    )
  }

  /**
   * Build the application container
   */
  @func()
  build(source: Directory): Container {
    const build = dag
      .node({ ctr: this.buildEnv(source) })
      .commands()
      .run(["build"])
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
    return await dag
      .node({ ctr: this.buildEnv(source) })
      .commands()
      .run(["test:unit", "run"])
      .stdout()
  }

  /**
   * Build a ready-to-use development environment
   */
  @func()
  buildEnv(source: Directory): Container {
    return dag
      .node({ version: "21" })
      .withNpm()
      .withSource(source)
      .install()
      .container()
  }
}
```

