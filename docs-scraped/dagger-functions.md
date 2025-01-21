# Dagger Functions

:::note
While you can use the utilities defined in the automatically-generated code above, you *cannot* edit these files. Even if you edit them locally, any changes will not be persisted when you run the module.
:::

You can now write Dagger Functions using the selected Dagger SDK. Here is an example, which calls a remote API method and returns the result:

```go
package main

import (
	"context"
)

type MyModule struct{}

func (m *MyModule) GetUser(ctx context.Context) (string, error) {
	return dag.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "curl"}).
		WithExec([]string{"apk", "add", "jq"}).
		WithExec([]string{"sh", "-c", "curl https://randomuser.me/api/ | jq .results[0].name"}).
		Stdout(ctx)
}
```

```python
from dagger import dag, function, object_type


@object_type
class MyModule:
    @function
    async def get_user(self) -> str:
        return await (
            dag.container()
            .from_("alpine:latest")
            .with_exec(["apk", "add", "curl"])
            .with_exec(["apk", "add", "jq"])
            .with_exec(
                ["sh", "-c", "curl https://randomuser.me/api/ | jq .results[0].name"]
            )
            .stdout()
        )
```

```typescript
import { dag, object, func } from "@dagger.io/dagger"

@object()
class MyModule {
  @func()
  async getUser(): Promise<string> {
    return await dag
      .container()
      .from("alpine:latest")
      .withExec(["apk", "add", "curl"])
      .withExec(["apk", "add", "jq"])
      .withExec([
        "sh",
        "-c",
        "curl https://randomuser.me/api/ | jq .results[0].name",
      ])
      .stdout()
  }
}
```

