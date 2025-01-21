# String arguments

To pass a string argument to a Dagger Function, add the corresponding flag to the `dagger call` command, followed by the string value.

Here is an example of a Dagger Function that accepts a string argument:

```go
package main

import (
	"context"
	"fmt"
)

type MyModule struct{}

func (m *MyModule) GetUser(ctx context.Context, gender string) (string, error) {
	return dag.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "curl"}).
		WithExec([]string{"apk", "add", "jq"}).
		WithExec([]string{"sh", "-c", fmt.Sprintf("curl https://randomuser.me/api/?gender=%s | jq .results[0].name", gender)}).
		Stdout(ctx)
}
```

```python
from dagger import dag, function, object_type


@object_type
class MyModule:
    @function
    async def get_user(self, gender: str) -> str:
        return await (
            dag.container()
            .from_("alpine:latest")
            .with_exec(["apk", "add", "curl"])
            .with_exec(["apk", "add", "jq"])
            .with_exec(
                [
                    "sh",
                    "-c",
                    (
                        f"curl https://randomuser.me/api/?gender={gender}"
                        " | jq .results[0].name"
                    ),
                ]
            )
            .stdout()
        )
```

```typescript
import { dag, object, func } from "@dagger.io/dagger"

@object()
class MyModule {
  @func()
  async getUser(gender: string): Promise<string> {
    return await dag
      .container()
      .from("alpine:latest")
      .withExec(["apk", "add", "curl"])
      .withExec(["apk", "add", "jq"])
      .withExec([
        "sh",
        "-c",
        `curl https://randomuser.me/api/?gender=${gender} | jq .results[0].name`,
      ])
      .stdout()
  }
}
```

