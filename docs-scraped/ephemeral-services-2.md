# Ephemeral Services

See it in action:

![Container to host networking](/img/current_docs/features/service-container.gif)

This also works in the opposite direction: containers in Dagger Functions can communicate with services running on the host. Here's an example of how a pipeline running in a Dagger Function can access and query a MariaDB database service running on the host:

```go
package main

import (
	"context"
	"dagger/my-module/internal/dagger"
)

type MyModule struct{}

func (m *MyModule) UserList(
	ctx context.Context,
	svc *dagger.Service,
) (string, error) {
	return dag.Container().
		From("mariadb:10.11.2").
		WithServiceBinding("db", svc).
		WithExec([]string{"/usr/bin/mysql", "--user=root", "--password=secret", "--host=db", "-e", "SELECT Host, User FROM mysql.user"}).
		Stdout(ctx)
}
```

```python
import dagger
from dagger import dag, function, object_type


@object_type
class MyModule:
    @function
    async def user_list(self, svc: dagger.Service) -> str:
        return await (
            dag.container()
            .from_("mariadb:10.11.2")
            .with_service_binding("db", svc)
            .with_exec(
                [
                    "/usr/bin/mysql",
                    "--user=root",
                    "--password=secret",
                    "--host=db",
                    "-e",
                    "SELECT Host, User FROM mysql.user",
                ]
            )
            .stdout()
        )
```

```typescript
import { dag, object, func, Service } from "@dagger.io/dagger"

@object()
class MyModule {
  @func()
  async userList(svc: Service): Promise<string> {
    return await dag
      .container()
      .from("mariadb:10.11.2")
      .withServiceBinding("db", svc)
      .withExec([
        "/usr/bin/mysql",
        "--user=root",
        "--password=secret",
        "--host=db",
        "-e",
        "SELECT Host, User FROM mysql.user",
      ])
      .stdout()
  }
}
```

