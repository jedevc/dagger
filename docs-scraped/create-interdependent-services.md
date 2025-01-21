# Create interdependent services

The following Dagger Function runs two services, service A and service B, that depend on each other. The services are set up with custom hostnames, `svca` and `svcb`, allowing each service to communicate with the other by hostname.

```go
package main

import (
	"context"
	"dagger/my-module/internal/dagger"
)

type MyModule struct{}

// Run two services which are dependent on each other
func (m *MyModule) Services(ctx context.Context) (*dagger.Service, error) {

	svcA := dag.Container().From("nginx").
		WithExposedPort(80).
		AsService(dagger.ContainerAsServiceOpts{Args: []string{"sh", "-c", `nginx & while true; do curl svcb:80 && sleep 1; done`}}).
		WithHostname("svca")

	_, err := svcA.Start(ctx)
	if err != nil {
		return nil, err
	}

	svcB := dag.Container().From("nginx").
		WithExposedPort(80).
		AsService(dagger.ContainerAsServiceOpts{Args: []string{"sh", "-c", `nginx & while true; do curl svca:80 && sleep 1; done`}}).
		WithHostname("svcb")

	svcB, err = svcB.Start(ctx)
	if err != nil {
		return nil, err
	}

	return svcB, nil
}
```

```python
import dagger
from dagger import dag, function, object_type


@object_type
class MyModule:
    @function
    async def services(self) -> dagger.Service:
        """Run two services which are dependent on each other"""
        svc_a = (
            dag.container()
            .from_("nginx")
            .with_exposed_port(80)
            .as_service(
                args=[
                    "sh",
                    "-c",
                    "nginx & while true; do curl svcb:80 && sleep 1; done",
                ]
            )
            .with_hostname("svca")
        )

        await svc_a.start()

        svc_b = (
            dag.container()
            .from_("nginx")
            .with_exposed_port(80)
            .as_service(
                args=[
                    "sh",
                    "-c",
                    "nginx & while true; do curl svca:80 && sleep 1; done",
                ]
            )
            .with_hostname("svcb")
        )

        await svc_b.start()

        return svc_b
```

```typescript
import { dag, object, func, Service } from "@dagger.io/dagger"

@object()
class MyModule {
  // Run two services which are dependent on each other
  @func()
  async services(): Promise<Service> {
    const svcA = dag
      .container()
      .from("nginx")
      .withExposedPort(80)
      .asService({
        args: [
          "sh",
          "-c",
          `nginx & while true; do curl svcb:80 && sleep 1; done`,
        ],
      })
      .withHostname("svca")

    await svcA.start()

    const svcB = dag
      .container()
      .from("nginx")
      .withExposedPort(80)
      .asService({
        args: [
          "sh",
          "-c",
          `nginx & while true; do curl svca:80 && sleep 1; done`,
        ],
      })
      .withHostname("svcb")

    await svcB.start()

    return svcB
  }
}
```

