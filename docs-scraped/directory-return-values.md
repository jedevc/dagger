# Directory return values

Directory return values might be produced by a Dagger Function that:

- Builds language-specific binaries
- Downloads source code from git or another remote source
- Processes source code (for example, generating documentation or linting code)
- Downloads or processes datasets
- Downloads machine learning models

Here is an example of a Go builder Dagger Function that accepts a remote Git address as a `Directory` argument, builds a Go binary from the source code in that repository, and returns the build directory containing the compiled binary:

```go
package main

import (
	"context"
	"dagger/my-module/internal/dagger"
)

type MyModule struct{}

func (m *MyModule) GoBuilder(ctx context.Context, src *dagger.Directory, arch string, os string) *dagger.Directory {
	return dag.Container().
		From("golang:1.21").
		WithMountedDirectory("/src", src).
		WithWorkdir("/src").
		WithEnvVariable("GOARCH", arch).
		WithEnvVariable("GOOS", os).
		WithEnvVariable("CGO_ENABLED", "0").
		WithExec([]string{"go", "build", "-o", "build/"}).
		Directory("/src/build")
}
```

```python
import dagger
from dagger import dag, function, object_type


@object_type
class MyModule:
    @function
    def go_builder(self, src: dagger.Directory, arch: str, os: str) -> dagger.Directory:
        return (
            dag.container()
            .from_("golang:1.21")
            .with_mounted_directory("/src", src)
            .with_workdir("/src")
            .with_env_variable("GOARCH", arch)
            .with_env_variable("GOOS", os)
            .with_env_variable("CGO_ENABLED", "0")
            .with_exec(["go", "build", "-o", "build/"])
            .directory("/src/build")
        )
```

```typescript
import { dag, Directory, object, func } from "@dagger.io/dagger"

@object()
class MyModule {
  @func()
  goBuilder(src: Directory, arch: string, os: string): Directory {
    return dag
      .container()
      .from("golang:1.21")
      .withMountedDirectory("/src", src)
      .withWorkdir("/src")
      .withEnvVariable("GOARCH", arch)
      .withEnvVariable("GOOS", os)
      .withEnvVariable("CGO_ENABLED", "0")
      .withExec(["go", "build", "-o", "build/"])
      .directory("/src/build")
  }
}
```

