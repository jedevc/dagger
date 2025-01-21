# Optional arguments

Function arguments can be marked as optional. In this case, the Dagger CLI will not display an error if the argument is omitted in the function call.

Here's an example of a Dagger Function with an optional argument:

```go
package main

import (
	"context"
	"fmt"
)

type MyModule struct{}

func (m *MyModule) Hello(
	ctx context.Context,
	// +optional
	name string,
) (string, error) {
	if name != "" {
		return fmt.Sprintf("Hello, %s", name), nil
	} else {
		return "Hello, world", nil
	}
}
```

```python
from dagger import function, object_type


@object_type
class MyModule:
    @function
    def hello(self, name: str | None) -> str:
        if name is None:
            name = "world"
        return f"Hello, {name}"
```

```typescript
import { object, func } from "@dagger.io/dagger"

@object()
class MyModule {
  @func()
  hello(name?: string): string {
    if (name) {
      return `Hello, ${name}`
    }
    return "Hello, world"
  }
}
```

