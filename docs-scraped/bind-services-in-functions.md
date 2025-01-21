# Bind services in functions

A Dagger Function can create and return a service, which can then be used from another Dagger Function or from the calling host. Services in Dagger Functions are returned using the `Service` core type.

Here is an example of a Dagger Function that returns an HTTP service. This service is used by another Dagger Function, which creates a service binding using the alias `www` and then accesses the HTTP service using this alias.

```go
package main

import (
	"context"
	"dagger/my-module/internal/dagger"
)

type MyModule struct{}

// Start and return an HTTP service
func (m *MyModule) HttpService() *dagger.Service {
	return dag.Container().
		From("python").
		WithWorkdir("/srv").
		WithNewFile("index.html", "Hello, world!").
		WithExposedPort(8080).
		AsService(dagger.ContainerAsServiceOpts{Args: []string{"python", "-m", "http.server", "8080"}})
}

// Send a request to an HTTP service and return the response
func (m *MyModule) Get(ctx context.Context) (string, error) {
	return dag.Container().
		From("alpine").
		WithServiceBinding("www", m.HttpService()).
		WithExec([]string{"wget", "-O-", "http://www:8080"}).
		Stdout(ctx)
}
```

```python
import dagger
from dagger import dag, function, object_type


@object_type
class MyModule:
    @function
    def http_service(self) -> dagger.Service:
        """Start and return an HTTP service."""
        return (
            dag.container()
            .from_("python")
            .with_workdir("/srv")
            .with_new_file("index.html", "Hello, world!")
            .with_exposed_port(8080)
            .as_service(args=["python", "-m", "http.server", "8080"])
        )

    @function
    async def get(self) -> str:
        """Send a request to an HTTP service and return the response."""
        return await (
            dag.container()
            .from_("alpine")
            .with_service_binding("www", self.http_service())
            .with_exec(["wget", "-O-", "http://www:8080"])
            .stdout()
        )
```

```typescript
import { dag, object, func, Service } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Start and return an HTTP service
   */
  @func()
  httpService(): Service {
    return dag
      .container()
      .from("python")
      .withWorkdir("/srv")
      .withNewFile("index.html", "Hello, world!")
      .withExposedPort(8080)
      .asService({ args: ["python", "-m", "http.server", "8080"] })
  }

  /**
   * Send a request to an HTTP service and return the response
   */
  @func()
  async get(): Promise<string> {
    return await dag
      .container()
      .from("alpine")
      .withServiceBinding("www", this.httpService())
      .withExec(["wget", "-O-", "http://www:8080"])
      .stdout()
  }
}
```

