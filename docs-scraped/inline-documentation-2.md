# Inline Documentation

Here is an example of the result from `dagger functions`:

```
Name         Description
hello        Return a greeting.
loud-hello   Return a loud greeting.
```

Here is an example of the result from `dagger call hello --help`:

```
Return a greeting.

USAGE
  dagger call hello [arguments]

ARGUMENTS
      --greeting string   The greeting to display [required]
      --name string       Who to greet [required]
```


The following code snippet shows how to add documentation for an object and its fields in your Dagger module:

```go
package main

// The struct represents a single user of the system.
type MyModule struct {
	Name string
	Age  int
}

func New(
	// The name of the user.
	name string,
	// The age of the user.
	age int,
) *MyModule {
	return &MyModule{
		Name: name,
		Age:  age,
	}
}
```

```python
from typing import Annotated

from dagger import Doc, object_type


@object_type
class MyModule:
    """The object represents a single user of the system."""

    name: Annotated[str, Doc("The name of the user.")]
    age: Annotated[str, Doc("The age of the user.")]
```

```typescript
import { object } from "@dagger.io/dagger"

/**
 * The object represents a single user of the system.
 */
@object()
class MyModule {
  @func()
  name: string
  @func()
  age: number

  constructor(
    /**
     * The name of the user.
     */
    age: number,
    /**
     * The age of the user.
     */
    name: string,
  ) {
    this.name = name
    this.age = age
  }
}
```

