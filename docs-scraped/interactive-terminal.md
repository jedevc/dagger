# Interactive Terminal

Dagger provides an interactive terminal that can help greatly when trying to debug a pipeline failure.

To use this, set one or more explicit breakpoints in your Dagger pipeline with the `Container.terminal()` method. Dagger then starts an interactive terminal session at each breakpoint. This lets you inspect a `Directory` or a `Container` at any point in your pipeline run, with all the necessary context available to you.

Here is a simple example, which opens an interactive terminal in an `alpine` container:

```shell
dagger core container from --address=alpine terminal
```

Here is an example of a Dagger Function which opens an interactive terminal at two different points in the Dagger pipeline to inspect the built container:

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

