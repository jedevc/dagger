# Boolean arguments

To pass a Boolean argument to a Dagger Function, simply add the corresponding flag, like so:

- To set the argument to true: `--foo=true`, or simply `--foo`
- To set the argument to false: `--foo=false`

Here is an example of a Dagger Function that accepts a Boolean argument:

```go
package main

import (
	"strings"
)

type MyModule struct{}

func (m *MyModule) Hello(shout bool) string {
	message := "Hello, world"
	if shout {
		return strings.ToUpper(message)
	}
	return message
}
```

```python
from dagger import function, object_type


@object_type
class MyModule:
    @function
    def hello(self, shout: bool) -> str:
        message = "Hello, world"
        if shout:
            return message.upper()
        return message
```

```typescript
import { object, func } from "@dagger.io/dagger"

@object()
class MyModule {
  @func()
  hello(shout: boolean): string {
    const message = "Hello, world"
    if (shout) {
      return message.toUpperCase()
    }
    return message
  }
}
```

