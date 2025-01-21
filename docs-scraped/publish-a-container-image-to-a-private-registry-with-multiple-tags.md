# Publish a container image to a private registry with multiple tags

The following Dagger Function tags a just-in-time container image multiple times and publishes it to a private registry.

```go
package main

import (
	"context"
	"fmt"

	"dagger/my-module/internal/dagger"
)

type MyModule struct{}

// Tag a container image multiple times and publish it to a private registry
func (m *MyModule) Publish(
	ctx context.Context,
	// Registry address
	registry string,
	// Registry username
	username string,
	// Registry password
	password *dagger.Secret,
) ([]string, error) {
	tags := [4]string{"latest", "1.0-alpine", "1.0", "1.0.0"}
	addr := []string{}
	ctr := dag.Container().
		From("nginx:1.23-alpine").
		WithNewFile(
			"/usr/share/nginx/html/index.html",
			"Hello from Dagger!",
			dagger.ContainerWithNewFileOpts{Permissions: 0o400},
		).
		WithRegistryAuth(registry, username, password)

	for _, tag := range tags {
		a, err := ctr.Publish(ctx, fmt.Sprintf("%s/%s/my-nginx:%s", registry, username, tag))
		if err != nil {
			return addr, err
		}
		addr = append(addr, a)
	}
	return addr, nil
}
```

```python
from typing import Annotated

import dagger
from dagger import Doc, dag, function, object_type


@object_type
class MyModule:
    @function
    async def publish(
        self,
        registry: Annotated[str, Doc("Registry address")],
        username: Annotated[str, Doc("Registry username")],
        password: Annotated[dagger.Secret, Doc("Registry password")],
    ) -> list[str]:
        """Tag a container image multiple times and publish it to a private registry"""
        tags = ["latest", "1.0-alpine", "1.0", "1.0.0"]
        addr = []
        container = (
            dag.container()
            .from_("nginx:1.23-alpine")
            .with_new_file(
                "/usr/share/nginx/html/index.html",
                "Hello from Dagger!",
                permissions=0o400,
            )
            .with_registry_auth(registry, username, password)
        )
        for tag in tags:
            a = await container.publish(f"{registry}/{username}/my-nginx:{tag}")
            addr.append(a)
        return addr
```

```typescript
import { dag, Secret, object, func } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Tag a container image multiple times and publish it to a private registry
   */
  @func()
  async publish(
    /**
     * Registry address
     */
    registry: string,
    /**
     * Registry username
     */
    username: string,
    /**
     * Registry password
     */
    password: Secret,
  ): Promise<string[]> {
    const tags = ["latest", "1.0-alpine", "1.0", "1.0.0"]

    const addr: string[] = []

    const container = dag
      .container()
      .from("nginx:1.23-alpine")
      .withNewFile("/usr/share/nginx/html/index.html", "Hello from Dagger!", {
        permissions: 0o400,
      })
      .withRegistryAuth(registry, username, password)

    for (const tag in tags) {
      const a = await container.publish(
        `${registry}/${username}/my-nginx:${tags[tag]}`,
      )
      addr.push(a)
    }

    return addr
  }
}
```

