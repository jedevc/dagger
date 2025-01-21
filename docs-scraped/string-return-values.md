# String return values

Here is an example of a Dagger Function that returns operating system information for the container as a string:

```go
package main

import (
	"context"

	"main/internal/dagger"
)

type MyModule struct{}

func (m *MyModule) OsInfo(ctx context.Context, ctr *dagger.Container) (string, error) {
	return ctr.
		WithExec([]string{"uname", "-a"}).
		Stdout(ctx)
}
```

```python
import dagger
from dagger import function, object_type


@object_type
class MyModule:
    @function
    async def os_info(self, ctr: dagger.Container) -> str:
        return await ctr.with_exec(["uname", "-a"]).stdout()
```

```typescript
import { Container, object, func } from "@dagger.io/dagger"

@object()
class MyModule {
  @func()
  async osInfo(ctr: Container): Promise<string> {
    return ctr.withExec(["uname", "-a"]).stdout()
  }
}
```

