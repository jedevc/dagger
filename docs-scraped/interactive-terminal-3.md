# Interactive Terminal

Under the hood, this creates a new container (defaults to `alpine`) and starts a shell, mounting the directory inside. This container can be customized using additional options. Here is a more complex example, which produces the same result as the previous one but this time using an `ubuntu` container image and `bash` shell instead of the default `alpine` container image and `sh` shell:

```go
package main

import (
	"context"
	"dagger/my-module/internal/dagger"
)

type MyModule struct{}

func (m *MyModule) AdvancedDirectory(ctx context.Context) (string, error) {
	return dag.
		Git("https://github.com/dagger/dagger.git").
		Head().
		Tree().
		Terminal(dagger.DirectoryTerminalOpts{
			Container: dag.Container().From("ubuntu"),
			Cmd:       []string{"/bin/bash"},
		}).
		File("README.md").
		Contents(ctx)
}
```

```python
from dagger import dag, function, object_type


@object_type
class MyModule:
    @function
    async def advanced_directory(self) -> str:
        return await (
            dag.git("https://github.com/dagger/dagger.git")
            .head()
            .tree()
            .terminal(
                container=dag.container().from_("ubuntu"),
                cmd=["/bin/bash"],
            )
            .file("README.md")
            .contents()
        )
```

```typescript
import { dag, Container, object, func } from "@dagger.io/dagger"

@object()
class MyModule {
  @func()
  async advancedDirectory(): Promise<string> {
    return await dag
      .git("https://github.com/dagger/dagger.git")
      .head()
      .tree()
      .terminal({
        container: dag.container().from("ubuntu"),
        cmd: ["/bin/bash"],
      })
      .file("README.md")
      .contents()
  }
}
```

