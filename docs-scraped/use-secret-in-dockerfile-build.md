# Use secret in Dockerfile build

The following code listing demonstrates how to inject a secret into a Dockerfile build. The secret is automatically mounted in the build container at `/run/secrets/SECRET-ID`.

```go
package main

import (
	"context"

	"main/internal/dagger"
)

type MyModule struct{}

// Build a Container from a Dockerfile
func (m *MyModule) Build(
	ctx context.Context,
	// The source code to build
	source *dagger.Directory,
	// The secret to use in the Dockerfile
	secret *dagger.Secret,
) (*dagger.Container, error) {
	// Ensure the Dagger secret's name matches what the Dockerfile
	// expects as the id for the secret mount.
	secretVal, err := secret.Plaintext(ctx)
	if err != nil {
		return nil, err
	}
	buildSecret := dag.SetSecret("gh-secret", secretVal)

	return source.
		DockerBuild(dagger.DirectoryDockerBuildOpts{
			Secrets: []*dagger.Secret{buildSecret},
		}), nil
}
```

```python
from typing import Annotated

import dagger
from dagger import Doc, dag, function, object_type


@object_type
class MyModule:
    @function
    async def build(
        self,
        source: Annotated[dagger.Directory, Doc("The source code to build")],
        secret: Annotated[dagger.Secret, Doc("The secret to use in the Dockerfile")],
    ) -> dagger.Container:
        """Build a Container from a Dockerfile"""
        # Ensure the Dagger secret's name matches what the Dockerfile
        # expects as the id for the secret mount.
        build_secret = dag.set_secret("gh-secret", await secret.plaintext())

        return source.docker_build(secrets=[build_secret])
```

```typescript
import { dag, object, func, Secret } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Build a Container from a Dockerfile
   */
  @func()
  async build(
    /**
     * The source code to build
     */
    source: Directory,
    /**
     * The secret to use in the Dockerfile
     */
    secret: Secret,
  ): Promise<Container> {
    // Ensure the Dagger secret's name matches what the Dockerfile
    // expects as the id for the secret mount.
    const buildSecret = dag.setSecret("gh-secret", await secret.plaintext())

    return source.dockerBuild({ secrets: [buildSecret] })
  }
}
```

