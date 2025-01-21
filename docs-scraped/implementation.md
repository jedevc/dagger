# Implementation

Here is an example of a module `Example` that implements the `Fooer` interface:

```go
package main

import (
	"context"
	"fmt"
)

type Example struct{}

func (m *Example) Foo(ctx context.Context, bar int) (string, error) {
	return fmt.Sprintf("number is: %d", bar), nil
}
```

```ts
import { func, object } from "@dagger.io/dagger"

export interface Fooer {
  // You can also declare it as a method signature (e.g., `foo(): Promise<string>`)
  foo: (bar: number) => Promise<string>
}

@object()
export class Example {
  @func()
  async foo(bar: number): Promise<string> {
    return `number is: ${bar}`
  }
}
```

