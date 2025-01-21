# Modify a copied directory or remote repository in a container

The following Dagger Function accepts a `Directory` argument, which could reference either a directory from the local filesystem or from a [remote Git repository](../api/arguments.mdx#remote-repositories). It copies the specified directory to the `/src` path in a container, adds a file to it, and returns the modified container.

:::note
Modifications made to a directory's contents after it is written to a container filesystem do not appear on the source. Data flows only one way between Dagger operations, because they are connected in a DAG. To transfer modifications back to the local host, you must explicitly export the directory back to the host filesystem.
:::

:::note
When working with private Git repositories, ensure that [SSH authentication is properly configured](../api/remote-modules.mdx#configuring-ssh-authentication) on your Dagger host.
:::

```go
package main

import (
	"context"

	"dagger/my-module/internal/dagger"
)

type MyModule struct{}

// Return a container with a specified directory and an additional file
func (m *MyModule) CopyAndModifyDirectory(
	ctx context.Context,
	// Source directory
	source *dagger.Directory,
) *dagger.Container {
	return dag.Container().
		From("alpine:latest").
		WithDirectory("/src", source).
		WithExec([]string{"/bin/sh", "-c", `echo foo > /src/foo`})
}
```

```python
from typing import Annotated

import dagger
from dagger import Doc, dag, function, object_type


@object_type
class MyModule:
    @function
    def copy_and_modify_directory(
        self, source: Annotated[dagger.Directory, Doc("Source directory")]
    ) -> dagger.Container:
        """Return a container with a specified directory and an additional file"""
        return (
            dag.container()
            .from_("alpine:latest")
            .with_directory("/src", source)
            .with_exec(["/bin/sh", "-c", "`echo foo > /src/foo`"])
        )
```

```typescript
import { dag, Container, Directory, object, func } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Return a container with a specified directory and an additional file
   */
  @func()
  copyAndModifyDirectory(
    /**
     * Source directory
     */
    source: Directory,
  ): Container {
    return dag
      .container()
      .from("alpine:latest")
      .withDirectory("/src", source)
      .withExec(["/bin/sh", "-c", "`echo foo > /src/foo`"])
  }
}
```

