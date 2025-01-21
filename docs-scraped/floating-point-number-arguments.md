# Floating-point number arguments

To pass a floating-point number as argument to a Dagger function, add the corresponding flag to the `dagger call` command, followed by the value.

Here is an example of a Dagger function that accepts a floating-point number as argument:

```go
package main

type MyModule struct{}

func (m *MyModule) AddFloat(a float64, b float64) float64 {
  return a + b
}
```

```python
from dagger import function, object_type


@object_type
class MyModule:
    @function
    async def add_float(self, a: float, b: float) -> float:
        return a + b
```

```typescript
import type { float } from "@dagger.io/dagger"
import { object, func } from "@dagger.io/dagger"

@object()
export class MyModule {
  @func()
  addFloat(a: float, b: float): float {
    return a + b
  }
}
```

