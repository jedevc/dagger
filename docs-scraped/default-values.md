# Default values

Function arguments can define a default value if no value is supplied for them.

Here's an example of a Dagger Function with a default value for a string argument:

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
	// +default="world"
	name string,
) (string, error) {
	return fmt.Sprintf("Hello, %s", name), nil
}
```

```python
from dagger import function, object_type


@object_type
class MyModule:
    @function
    def hello(self, name: str = "world") -> str:
        return f"Hello, {name}"
```

```typescript
import { object, func } from "@dagger.io/dagger"

@object()
class MyModule {
  @func()
  hello(name = "world"): string {
    return `Hello, ${name}`
  }
}
```

