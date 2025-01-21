# Debug container builds

The following Dagger Function opens an interactive terminal session at different stages in a Dagger pipeline to debug a container build.

```go
package main

import (
	"dagger/my-module/internal/dagger"
)

type MyModule struct{}

func (m *MyModule) Container() *dagger.Container {
	return dag.Container().
		From("alpine:latest").
		Terminal().
		WithExec([]string{"sh", "-c", "echo hello world > /foo && cat /foo"}).
		Terminal()
}
```

```python
import dagger
from dagger import dag, function, object_type


@object_type
class MyModule:
    @function
    def container(self) -> dagger.Container:
        return (
            dag.container()
            .from_("alpine:latest")
            .terminal()
            .with_exec(["sh", "-c", "echo hello world > /foo && cat /foo"])
            .terminal()
        )
```

```typescript
import { dag, Container, object, func } from "@dagger.io/dagger"

@object()
class MyModule {
  @func()
  container(): Container {
    return dag
      .container()
      .from("alpine:latest")
      .terminal()
      .withExec(["sh", "-c", "echo hello world > /foo && cat /foo"])
      .terminal()
  }
}
```

