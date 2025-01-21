# Error Handling

Dagger modules handle errors in the same way as the language they are written in. This allows you to support any kind of error handling that your application requires. You can also use error handling to verify user input.

Here is an example Dagger Function that performs division and throws an error if the denominator is zero:

```go
// A Dagger module for saying hello world!

package main

import (
	"fmt"
)

type MyModule struct {
}

func (*MyModule) Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("cannot divide by zero")
	}
	return a / b, nil
}
```

```python
from dagger import function, object_type


@object_type
class MyModule:
    @function
    def divide(self, a: int, b: int) -> float:
        if b == 0:
            msg = "cannot divide by zero"
            raise ValueError(msg)
        return a / b
```

```typescript
import { object, func } from "@dagger.io/dagger"

@object()
class MyModule {
  @func()
  divide(a: number, b: number): number {
    if (b == 0) {
      throw new Error("cannot divide by zero")
    }

    return a / b
  }
}
```

