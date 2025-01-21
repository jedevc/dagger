# Call the Dagger Function

Call the Dagger Function:

```shell
dagger call publish --source=.
```

You should see the application being tested, built, and published to the [ttl.sh container registry](https://ttl.sh):

![Publish](/img/current_docs/quickstart/publish.gif)

You can test the published container image by pulling and running it with `docker run`:

![Docker run](/img/current_docs/quickstart/docker.gif)


:::tip FUNCTION CHAINING
[Function chaining](../features/programmable-pipelines.mdx) works the same way, whether you're writing Dagger Function code using a Dagger SDK or invoking a Dagger Function using the Dagger CLI. The following are equivalent:

```go
package main

import (
	"context"

	"dagger/hello-dagger/internal/dagger"
)

type HelloDagger struct{}

// Returns a base container
func (m *HelloDagger) Base() *dagger.Container {
	return dag.Container().From("cgr.dev/chainguard/wolfi-base")
}

// Builds on top of base container and returns a new container
func (m *HelloDagger) Build() *dagger.Container {
	return m.Base().WithExec([]string{"apk", "add", "bash", "git"})
}

// Builds and publishes a container
func (m *HelloDagger) BuildAndPublish(ctx context.Context) (string, error) {
	return m.Build().Publish(ctx, "ttl.sh/bar")
}
```

```python
import dagger
from dagger import dag, function, object_type


@object_type
class HelloDagger:
    @function
    def base(self) -> dagger.Container:
        """Returns a base container"""
        return dag.container().from_("cgr.dev/chainguard/wolfi-base")

    @function
    def build(self) -> dagger.Container:
        """Builds on top of base container and returns a new container"""
        return self.base().with_exec(["apk", "add", "bash", "git"])

    @function
    async def build_and_publish(self) -> str:
        """Builds and publishes a container"""
        return await self.build().publish("ttl.sh/bar")
```

```typescript
import { dag, Container, object, func } from "@dagger.io/dagger"

@object()
class HelloDagger {
  /**
   * Returns a base container
   */
  @func()
  base(): Container {
    return dag.container().from("cgr.dev/chainguard/wolfi-base")
  }

  /**
   * Builds on top of base container and returns a new container
   */
  @func()
  build(): Container {
    return this.base().withExec(["apk", "add", "bash", "git"])
  }

  /**
   * Builds and publishes a container
   */
  @func()
  async buildAndPublish(): Promise<string> {
    return await this.build().publish("ttl.sh/bar")
  }
}
```

```shell

# all equivalent
dagger call base with-exec --args apk,add,bash,git publish --address="ttl.sh/bar"
dagger call build publish --address="ttl.sh/bar"
dagger call build-and-publish
```

