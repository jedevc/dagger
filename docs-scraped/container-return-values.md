# Container return values

Similar to directories and files, just-in-time containers are produced by calling a Dagger Function that returns the `Container` type. This type provides a complete API for building, running and distributing containers.

Just-in-time containers might be produced by a Dagger Function that:

- Builds a container
- Minifies a container
- Downloads a container image from a running registry
- Exports a container from Docker or other container runtimes
- Snapshots the state of a running container

You can think of a just-in-time container, and the `Container` type that represents it, as a build stage in Dockerfile. Each operation produces a new immutable state, which can be further processed, or exported as an OCI image. Dagger Functions can accept, return and pass containers between themselves, just like regular variables.

Here's an example of a Dagger Function that returns a base `alpine` container image with a list of additional specified packages:

```go
package main

import (
	"context"
	"dagger/my-module/internal/dagger"
)

type MyModule struct{}

func (m *MyModule) AlpineBuilder(ctx context.Context, packages []string) *dagger.Container {
	ctr := dag.Container().
		From("alpine:latest")
	for _, pkg := range packages {
		ctr = ctr.WithExec([]string{"apk", "add", pkg})
	}
	return ctr
}
```

```python
import dagger
from dagger import dag, function, object_type


@object_type
class MyModule:
    @function
    def alpine_builder(self, packages: list[str]) -> dagger.Container:
        ctr = dag.container().from_("alpine:latest")
        for pkg in packages:
            ctr = ctr.with_exec(["apk", "add", pkg])
        return ctr
```

```typescript
import { dag, Container, object, func } from "@dagger.io/dagger"

@object()
class MyModule {
  @func()
  alpineBuilder(packages: string[]): Container {
    let ctr = dag.container().from("alpine:latest")
    for (const pkg in packages) {
      ctr = ctr.withExec(["apk", "add", pkg])
    }
    return ctr
  }
}
```

