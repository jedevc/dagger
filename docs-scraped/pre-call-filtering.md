# Pre-call filtering

Pre-call filtering means that a directory is filtered before it's uploaded to the Dagger Engine container. This is useful for:

- Large monorepos. Typically your Dagger Function only operates on a subset of the monorepo, representing a specific component or feature. Uploading the entire worktree imposes a prohibitive cost.

- Large files, such as audio/video files and other binary content. These files take time to upload. If they're not directly relevant, you'll usually want your Dagger Function to ignore them.

  :::tip
  The `.git` directory is a good example of both these cases. It contains a lot of data, including large binary objects, and for projects with a long version history, it can sometimes be larger than your actual source code.
  :::

- Dependencies. If you're developing locally, you'll typically have your project dependencies installed locally: `node_modules` (Node.js), `.venv` (Python), `vendor` (PHP) and so on. When you call your Dagger Function locally, Dagger will upload all these installed dependencies as well. This is both bad practice and inefficient. Typically, you'll want your Dagger Function to ignore locally-installed dependencies and only operate on the project source code.

:::note
At the time of writing, Dagger [does not read exclusion patterns from existing `.dockerignore`/`.gitignore` files](https://github.com/dagger/dagger/issues/6627). If you already use these files, you'll need to manually implement the same patterns in your Dagger Function.
:::

To implement a pre-call filter in your Dagger Function, add an `ignore` parameter to your `Directory` argument. The `ignore` parameter follows the [`.gitignore` syntax](https://git-scm.com/docs/gitignore). Some important points to keep in mind are:

- The order of arguments is significant: the pattern `"**", "!**"` includes everything but `"!**", "**"` excludes everything.
- Prefixing a path with `!` negates a previous ignore: the pattern `"!foo"` has no effect, since nothing is previously ignored, while the pattern `"**", "!foo"` excludes everything except `foo`.

```go
package main

import (
	"context"
	"dagger/my-module/internal/dagger"
)

type MyModule struct{}

func (m *MyModule) Foo(
	ctx context.Context,
	// +ignore=["*", "!**/*.go", "!go.mod", "!go.sum"]
	source *dagger.Directory,
) (*dagger.Container, error) {
	return dag.
		Container().
		From("alpine:latest").
		WithDirectory("/src", source).
		Sync(ctx)
}
```

```python
from typing import Annotated

import dagger
from dagger import Ignore, dag, function, object_type


@object_type
class MyModule:
    @function
    async def foo(
        self,
        source: Annotated[dagger.Directory, Ignore(["*", "!**/*.py"])],
    ) -> dagger.Container:
        return await (
            dag.container().from_("alpine:latest").with_directory("/src", source).sync()
        )
```

```typescript
import {
  dag,
  object,
  argument,
  func,
  Directory,
  Container,
} from "@dagger.io/dagger"

@object()
class MyModule {
  @func()
  async foo(
    @argument({ ignore: ["*", "!**/*.ts"] }) source: Directory,
  ): Promise<Container> {
    return await dag
      .container()
      .from("alpine:latest")
      .withDirectory("/src", source)
      .sync()
  }
}
```

