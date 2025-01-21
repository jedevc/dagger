# Directories and files

It is possible to automatically load a filesystem path as a `Directory` or `File` object in a Dagger Function, by passing it as a "default path" to the corresponding argument. The `Directory` or `File` loaded in this manner is not merely a string, but it is the actual filesystem state of the specified directory or file, managed by the Dagger Engine and handled in code just like any another variable.

:::important
Default contexts are only available for arguments of type `Directory` and `File`. They are commonly used to load constant filesystem locations, such as an application's source code directory.
:::

When determining how to resolve a default path, Dagger first identifies a "context directory".

- For Git repositories (defined by the presence of a `.git` sub-directory), the context directory is the repository root (for absolute paths), or the directory containing a `dagger.json` file (for relative paths).
- For all other cases, the context directory is the directory containing a `dagger.json` file.

The default path is then resolved starting from the context directory.

:::important
For security reasons, it is not possible to retrieve files or directories outside the context directory.
:::

The best way to understand this is with an example. Consider the following directory structure, representing a project with a Dagger module in `my-module`:

```shell
.
├── README.md
├── my-module
│   ├── dagger.json
│   ├── ...
│   └── src
│       └── main
├── index.html
├── public
│   └── ...
├── src
│   └── ...
```

Here are two Dagger Functions with default paths for their directory and file arguments:

```go
package main

import (
	"context"
	"dagger/my-module/internal/dagger"
)

type MyModule struct{}

func (m *MyModule) ReadDir(
	ctx context.Context,
	// +defaultPath="/"
	source *dagger.Directory,
) ([]string, error) {
	return source.Entries(ctx)
}

func (m *MyModule) ReadFile(
	ctx context.Context,
	// +defaultPath="/README.md"
	source *dagger.File,
) (string, error) {
	return source.Contents(ctx)
}
```

```python
from typing import Annotated

import dagger
from dagger import DefaultPath, function, object_type


@object_type
class MyModule:
    @function
    async def read_dir(
        self,
        source: Annotated[dagger.Directory, DefaultPath("/")],
    ) -> list[str]:
        return await source.entries()

    @function
    async def read_file(
        self,
        source: Annotated[dagger.File, DefaultPath("/README.md")],
    ) -> list[str]:
        return await source.contents()
```

```typescript
import { Directory, File, argument, object, func } from "@dagger.io/dagger"

@object()
class MyModule {
  @func()
  async readDir(
    @argument({ defaultPath: "/" }) source: Directory,
  ): Promise<string[]> {
    return await source.entries()
  }

  @func()
  async readFile(
    @argument({ defaultPath: "/README.md" }) source: File,
  ): Promise<string> {
    return await source.contents()
  }
}
```

