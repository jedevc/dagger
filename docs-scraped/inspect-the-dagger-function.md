# Inspect the Dagger Function



```go
package main

import "dagger/hello-dagger/internal/dagger"

type HelloDagger struct{}

// Build the application container
func (m *HelloDagger) Build(source *dagger.Directory) *dagger.Container {
	// get the build environment container
	// by calling another Dagger Function
	build := m.BuildEnv(source).
		// build the application
		WithExec([]string{"npm", "run", "build"}).
		// get the build output directory
		Directory("./dist")
	// start from a slim NGINX container
	return dag.Container().From("nginx:1.25-alpine").
		// copy the build output directory to the container
		WithDirectory("/usr/share/nginx/html", build).
		// expose the container port
		WithExposedPort(80)
}
```

```python
import dagger
from dagger import dag, function, object_type


@object_type
class HelloDagger:
    @function
    def build(self, source: dagger.Directory) -> dagger.Container:
        """Build the application container"""
        build = (
            # get the build environment container
            # by calling another Dagger Function
            self.build_env(source)
            # build the application
            .with_exec(["npm", "run", "build"])
            #  get the build output directory
            .directory("./dist")
        )
        return (
            dag.container()
            # start from a slim NGINX container
            .from_("nginx:1.25-alpine")
            # copy the build output directory to the container
            .with_directory("/usr/share/nginx/html", build)
            # expose the container port
            .with_exposed_port(80)
        )
```

```typescript
import { dag, Container, Directory, object, func } from "@dagger.io/dagger"

@object()
class HelloDagger {
  /**
   * Build the application container
   */
  @func()
  build(source: Directory): Container {
    // get the build environment container
    // by calling another Dagger Function
    const build = this.buildEnv(source)
      // build the application
      .withExec(["npm", "run", "build"])
      // get the build output directory
      .directory("./dist")
    return (
      dag
        .container()
        // start from a slim NGINX container
        .from("nginx:1.25-alpine")
        // copy the build output directory to the container
        .withDirectory("/usr/share/nginx/html", build)
        // expose the container port
        .withExposedPort(80)
    )
  }
}
```

