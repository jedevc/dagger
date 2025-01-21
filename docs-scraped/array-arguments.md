# Array arguments

To pass an array argument to a Dagger Function, add the corresponding flag, followed by a comma-separated list of values.

```go
package main

import (
	"strings"
)

type MyModule struct{}

func (m *MyModule) Hello(names []string) string {
	message := "Hello"
	if len(names) > 0 {
		message += " " + strings.Join(names, ", ")
	}

	return message
}
```

```python
from dagger import function, object_type


@object_type
class MyModule:
    @function
    def hello(self, names: list[str]) -> str:
        message = "Hello"
        for name in names:
            message += f", {name}"
        return message
```

```typescript
import { object, func } from "@dagger.io/dagger"

@object()
class MyModule {
  @func()
  hello(names: string[]): string {
    let message = "Hello"
    for (const name of names) {
      message += ` ${name}`
    }
    return message
  }
}
```

