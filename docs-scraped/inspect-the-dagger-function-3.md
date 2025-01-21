# Inspect the Dagger Function



```go
package main

import (
	"context"

	"dagger/hello-dagger/internal/dagger"
)

type HelloDagger struct{}

// Return the result of running unit tests
func (m *HelloDagger) Test(ctx context.Context, source *dagger.Directory) (string, error) {
	// get the build environment container
	// by calling another Dagger Function
	return m.BuildEnv(source).
		// call the test runner
		WithExec([]string{"npm", "run", "test:unit", "run"}).
		// capture and return the command output
		Stdout(ctx)
}
```

```python
import dagger
from dagger import function, object_type


@object_type
class HelloDagger:
    @function
    async def test(self, source: dagger.Directory) -> str:
        """Return the result of running unit tests"""
        return await (
            # get the build environment container
            # by calling another Dagger Function
            self.build_env(source)
            # call the test runner
            .with_exec(["npm", "run", "test:unit", "run"])
            # capture and return the command output
            .stdout()
        )
```

```typescript
import { dag, Container, Directory, object, func } from "@dagger.io/dagger"

@object()
class HelloDagger {
  /**
   * Return the result of running unit tests
   */
  @func()
  async test(source: Directory): Promise<string> {
    // get the build environment container
    // by calling another Dagger Function
    return (
      this.buildEnv(source)
        // call the test runner
        .withExec(["npm", "run", "test:unit", "run"])
        // capture and return the command output
        .stdout()
    )
  }
}
```

