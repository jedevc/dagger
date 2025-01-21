# Secret arguments

Dagger allows you to utilize confidential information, such as passwords, API keys, SSH keys and so on, in your Dagger [modules](../features/modules.mdx) and Dagger Functions, without exposing those secrets in plaintext logs, writing them into the filesystem of containers you're building, or inserting them into the cache.

Secrets can be passed to Dagger Functions as arguments using the `Secret` core type. Here is an example of a Dagger Function which accepts a GitHub personal access token as a secret, and uses the token to authorize a request to the GitHub API:

```go
package main

import (
	"context"

	"main/internal/dagger"
)

type MyModule struct{}

// Query the GitHub API
func (m *MyModule) GithubApi(
	ctx context.Context,
	// GitHub API token
	token *dagger.Secret,
) (string, error) {
	return dag.Container().
		From("alpine:3.17").
		WithSecretVariable("GITHUB_API_TOKEN", token).
		WithExec([]string{"apk", "add", "curl"}).
		WithExec([]string{"sh", "-c", `curl "https://api.github.com/repos/dagger/dagger/issues" --header "Accept: application/vnd.github+json" --header "Authorization: Bearer $GITHUB_API_TOKEN"`}).
		Stdout(ctx)
}
```

```python
from typing import Annotated

import dagger
from dagger import Doc, dag, function, object_type


@object_type
class MyModule:
    @function
    async def github_api(
        self,
        token: Annotated[dagger.Secret, Doc("GitHub API token")],
    ) -> str:
        """Query the GitHub API"""
        return await (
            dag.container(platform=dagger.Platform("linux/amd64"))
            .from_("alpine:3.17")
            .with_secret_variable("GITHUB_API_TOKEN", token)
            .with_exec(["apk", "add", "curl"])
            .with_exec(
                [
                    "sh",
                    "-c",
                    (
                        'curl "https://api.github.com/repos/dagger/dagger/issues"'
                        ' --header "Authorization: Bearer $GITHUB_API_TOKEN"'
                        ' --header "Accept: application/vnd.github+json"'
                    ),
                ]
            )
            .stdout()
        )
```

```typescript
import { dag, object, func, Secret } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Query the GitHub API
   */
  @func()
  async githubApi(
    /**
     * GitHub API token
     */
    token: Secret,
  ): Promise<string> {
    return await dag
      .container()
      .from("alpine:3.17")
      .withSecretVariable("GITHUB_API_TOKEN", token)
      .withExec(["apk", "add", "curl"])
      .withExec([
        "sh",
        "-c",
        `curl "https://api.github.com/repos/dagger/dagger/issues" --header "Accept: application/vnd.github+json" --header "Authorization: Bearer $GITHUB_API_TOKEN"`,
      ])
      .stdout()
  }
}
```

