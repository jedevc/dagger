# Integer arguments

To pass an integer argument to a Dagger function, add the corresponding flag to the `dagger call` command, followed by the integer value.

Here is an example of a Dagger function that accepts an integer argument:

```go
package main

type MyModule struct{}

func (m *MyModule) AddInteger(a int, b int) int {
  return a + b
}
```

```python
from dagger import function, object_type


@object_type
class MyModule:
    @function
    async def add_integer(self, a: int, b: int) -> int:
        return a + b
```

```typescript
import { object, func } from "@dagger.io/dagger"

@object()
export class MyModule {
  @func()
  addInteger(a: number, b: number): number {
    return a + b
  }
}
```

