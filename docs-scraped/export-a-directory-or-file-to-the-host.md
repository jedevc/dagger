# Export a directory or file to the host

The following Dagger Functions return a just-in-time directory and file. These outputs can be exported to the host with `dagger call ... export ...`.

```go
package main

import "dagger/my-module/internal/dagger"

type MyModule struct{}

// Return a directory
func (m *MyModule) GetDir() *dagger.Directory {
	return m.Base().
		Directory("/src")
}

// Return a file
func (m *MyModule) GetFile() *dagger.File {
	return m.Base().
		File("/src/foo")
}

// Return a base container
func (m *MyModule) Base() *dagger.Container {
	return dag.Container().
		From("alpine:latest").
		WithExec([]string{"mkdir", "/src"}).
		WithExec([]string{"touch", "/src/foo", "/src/bar"})
}
```

```python
import dagger
from dagger import dag, function, object_type


@object_type
class MyModule:
    @function
    def get_dir(self) -> dagger.Directory:
        """Return a directory"""
        return self.base().directory("/src")

    @function
    def get_file(self) -> dagger.File:
        """Return a file"""
        return self.base().file("/src/foo")

    @function
    def base(self) -> dagger.Container:
        """Return a base container"""
        return (
            dag.container()
            .from_("alpine:latest")
            .with_exec(["mkdir", "/src"])
            .with_exec(["touch", "/src/foo", "/src/bar"])
        )
```

```typescript
import {
  dag,
  Directory,
  Container,
  File,
  object,
  func,
} from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Return a directory
   */
  @func()
  getDir(): Directory {
    return this.base().directory("/src")
  }

  /**
   * Return a file
   */
  @func()
  getFile(): File {
    return this.base().file("/src/foo")
  }

  /**
   * Return a base container
   */
  @func()
  base(): Container {
    return dag
      .container()
      .from("alpine:latest")
      .withExec(["mkdir", "/src"])
      .withExec(["touch", "/src/foo", "/src/bar"])
  }
}
```

