# Set environment variables in a container

The following Dagger Function demonstrates how to set multiple environment variables in a container.

```go
package main

import (
	"context"

	"dagger/my-module/internal/dagger"
)

type MyModule struct{}

type EnvVar struct {
	Name  string
	Value string
}

// Set environment variables in a container
func (m *MyModule) SetEnvVars(ctx context.Context) (string, error) {
	return dag.Container().
		From("alpine").
		With(EnvVariables([]*EnvVar{
			{"ENV_VAR_1", "VALUE 1"},
			{"ENV_VAR_2", "VALUE 2"},
			{"ENV_VAR_3", "VALUE 3"},
		})).
		WithExec([]string{"env"}).
		Stdout(ctx)
}

func EnvVariables(envs []*EnvVar) dagger.WithContainerFunc {
	return func(c *dagger.Container) *dagger.Container {
		for _, e := range envs {
			c = c.WithEnvVariable(e.Name, e.Value)
		}
		return c
	}
}
```

```python
import dagger
from dagger import dag, function, object_type


@object_type
class MyModule:
    @function
    async def set_env_vars(self) -> str:
        """Set environment variables in a container"""
        return await (
            dag.container()
            .from_("alpine")
            .with_(
                self.env_variables(
                    [
                        ("ENV_VAR_1", "VALUE 1"),
                        ("ENV_VAR_2", "VALUE 2"),
                        ("ENV_VAR_3", "VALUE 3"),
                    ]
                )
            )
            .with_exec(["env"])
            .stdout()
        )

    def env_variables(self, envs: list[tuple[str, str]]):
        def env_variables_inner(ctr: dagger.Container):
            for key, value in envs:
                ctr = ctr.with_env_variable(key, value)
            return ctr

        return env_variables_inner
```

```typescript
import { dag, Container, Directory, object, func } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Set environment variables in a container
   */
  @func()
  async setEnvVars(): Promise<string> {
    return await dag
      .container()
      .from("alpine")
      .with(
        envVariables([
          ["ENV_VAR_1", "VALUE 1"],
          ["ENV_VAR_2", "VALUE 2"],
          ["ENV_VAR_3", "VALUE_3"],
        ]),
      )
      .withExec(["env"])
      .stdout()
  }
}

function envVariables(envs: Array<[string, string]>) {
  return (c: Container): Container => {
    for (const [key, value] of envs) {
      c = c.withEnvVariable(key, value)
    }
    return c
  }
}
```

