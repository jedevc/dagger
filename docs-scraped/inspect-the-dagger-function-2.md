# Inspect the Dagger Function



```go
package main

import "dagger/hello-dagger/internal/dagger"

type HelloDagger struct{}

// Build a ready-to-use development environment
func (m *HelloDagger) BuildEnv(source *dagger.Directory) *dagger.Container {
	// create a Dagger cache volume for dependencies
	nodeCache := dag.CacheVolume("node")
	return dag.Container().
		// start from a base Node.js container
		From("node:21-slim").
		// add the source code at /src
		WithDirectory("/src", source).
		// mount the cache volume at /root/.npm
		WithMountedCache("/root/.npm", nodeCache).
		// change the working directory to /src
		WithWorkdir("/src").
		// run npm install to install dependencies
		WithExec([]string{"npm", "install"})
}
```

```python
import dagger
from dagger import dag, function, object_type


@object_type
class HelloDagger:
    @function
    def build_env(self, source: dagger.Directory) -> dagger.Container:
        """Build a ready-to-use development environment"""
        # create a Dagger cache volume for dependencies
        node_cache = dag.cache_volume("node")
        return (
            dag.container()
            # start from a base Node.js container
            .from_("node:21-slim")
            # add the source code at /src
            .with_directory("/src", source)
            # mount the cache volume at /root/.npm
            .with_mounted_cache("/root/.npm", node_cache)
            # change the working directory to /src
            .with_workdir("/src")
            # run npm install to install dependencies
            .with_exec(["npm", "install"])
        )
```

```typescript
import { dag, Container, Directory, object, func } from "@dagger.io/dagger"

@object()
class HelloDagger {
  /**
   * Build a ready-to-use development environment
   */
  @func()
  buildEnv(source: Directory): Container {
    // create a Dagger cache volume for dependencies
    const nodeCache = dag.cacheVolume("node")
    return (
      dag
        .container()
        // start from a base Node.js container
        .from("node:21-slim")
        // add the source code at /src
        .withDirectory("/src", source)
        // mount the cache volume at /root/.npm
        .withMountedCache("/root/.npm", nodeCache)
        // change the working directory to /src
        .withWorkdir("/src")
        // run npm install to install dependencies
        .withExec(["npm", "install"])
    )
  }
}
```

