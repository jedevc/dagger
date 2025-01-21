# Programmable Pipelines

Dagger provides a specialized container engine for application delivery pipelines, enabling you to replace your YAML-based CI/CD workflows with code in your favorite programming language. This allows you to build your delivery pipelines in the same language as your application and benefit from your language's existing ecosystem for tooling and best practices.

Dagger Functions are the fundamental unit of computing in Dagger. Dagger Functions let you encapsulate common project operations or workflows, such as "pull a container image", "copy a file", and "forward a TCP port", into portable, reusable program code with clear inputs and outputs.

Here's an example of a Dagger Function that builds and containerizes a Go application:

```go
package main

import (
	"context"
	"dagger/my-module/internal/dagger"
)

type MyModule struct{}

func (m *MyModule) Build(
	ctx context.Context,
	src *dagger.Directory,
	arch string,
	os string,
) *dagger.Container {
	return dag.Container().
		From("golang:1.21").
		WithMountedDirectory("/src", src).
		WithWorkdir("/src").
		WithEnvVariable("GOARCH", arch).
		WithEnvVariable("GOOS", os).
		WithEnvVariable("CGO_ENABLED", "0").
		WithExec([]string{"go", "build", "-o", "build/"})
}
```

```python
import dagger
from dagger import dag, function, object_type


@object_type
class MyModule:
    @function
    def build(self, src: dagger.Directory, arch: str, os: str) -> dagger.Container:
        return (
            dag.container()
            .from_("golang:1.21")
            .with_mounted_directory("/src", src)
            .with_workdir("/src")
            .with_env_variable("GOARCH", arch)
            .with_env_variable("GOOS", os)
            .with_env_variable("CGO_ENABLED", "0")
            .with_exec(["go", "build", "-o", "build/"])
        )
```

```typescript
import { dag, object, Directory, Container, func } from "@dagger.io/dagger"

@object()
class MyModule {
  @func()
  build(src: Directory, arch: string, os: string): Container {
    return dag
      .container()
      .from("golang:1.21")
      .withMountedDirectory("/src", src)
      .withWorkdir("/src")
      .withEnvVariable("GOARCH", arch)
      .withEnvVariable("GOOS", os)
      .withEnvVariable("CGO_ENABLED", "0")
      .withExec(["go", "build", "-o", "build/"])
  }
}
```

