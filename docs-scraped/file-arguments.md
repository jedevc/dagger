# File arguments

File arguments work in the same way as [directory arguments](#directory-arguments). To pass a file to a Dagger Function as an argument, add the corresponding flag, followed by a local filesystem path or a remote Git reference. In both cases, the CLI will convert it to an object referencing that filesystem path or Git repository location, and pass the resulting `File` object as argument to the Dagger Function.

Here's an example of a Dagger Function that accepts a `File` as argument, reads it, and returns its contents:

```go
package main

import (
	"context"
	"dagger/my-module/internal/dagger"
)

type MyModule struct{}

func (m *MyModule) ReadFile(ctx context.Context, source *dagger.File) (string, error) {
	contents, err := dag.Container().
		From("alpine:latest").
		WithFile("/src/myfile", source).
		WithExec([]string{"cat", "/src/myfile"}).
		Stdout(ctx)
	if err != nil {
		return "", err
	}
	return contents, nil
}
```

```python
import dagger
from dagger import dag, function, object_type


@object_type
class MyModule:
    @function
    async def read_file(self, source: dagger.File) -> str:
        return await (
            dag.container()
            .from_("alpine:latest")
            .with_file("/src/myfile", source)
            .with_exec(["cat", "/src/myfile"])
            .stdout()
        )
```

```typescript
import { dag, object, func, File } from "@dagger.io/dagger"

@object()
class MyModule {
  @func()
  async readFile(source: File): Promise<string> {
    return await dag
      .container()
      .from("alpine:latest")
      .withFile("/src/myfile", source)
      .withExec(["cat", "/src/myfile"])
      .stdout()
  }
}
```

