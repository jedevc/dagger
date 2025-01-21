# Ephemeral Services

Dagger Functions support service containers, enabling users to spin up additional services (as containers) and communicate with those services from their pipelines.

This makes it possible to:
- Instantiate and return services from a Dagger Function, and then:
  - Use those services in other Dagger Functions (container-to-container networking)
  - Use those services from the calling host (container-to-host networking)
- Expose host services for use in a Dagger Function (host-to-container networking).

Services instantiated by a Dagger Function run in service containers, which have the following characteristics:

- Each service container has a canonical, content-addressed hostname and an optional set of exposed ports.
- Service containers are started just-in-time, de-duplicated, and stopped when no longer needed.
- Service containers are health checked prior to running clients.

Some common scenarios for using services with Dagger Functions are:

- Running a database service for local storage or testing
- Running end-to-end integration tests against a service
- Running sidecar services

Here is an example of a Dagger Function that returns an HTTP service, which can then be accessed from the calling host:

```go
package main

import (
	"dagger/my-module/internal/dagger"
)

type MyModule struct{}

func (m *MyModule) HttpService() *dagger.Service {
	return dag.Container().
		From("python").
		WithWorkdir("/srv").
		WithNewFile("index.html", "Hello, world!").
		WithExposedPort(8080).
		AsService(dagger.ContainerAsServiceOpts{Args: []string{"python", "-m", "http.server", "8080"}})
}
```

```python
import dagger
from dagger import dag, function, object_type


@object_type
class MyModule:
    @function
    def http_service(self) -> dagger.Service:
        return (
            dag.container()
            .from_("python")
            .with_workdir("/srv")
            .with_new_file("index.html", "Hello, world!")
            .with_exposed_port(8080)
            .as_service(args=["python", "-m", "http.server", "8080"])
        )
```

```typescript
import { dag, object, func, Service } from "@dagger.io/dagger"

@object()
class MyModule {
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
}
```

