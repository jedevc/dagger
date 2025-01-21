# Example: MariaDB database service for application tests

The following example demonstrates how services can be used in Dagger Functions, by creating a Dagger Function for application unit/integration testing against a bound MariaDB database service.

The application used in this example is [Drupal](https://www.drupal.org/), a popular open-source PHP CMS. Drupal includes a large number of unit tests, including tests which require an active database connection. All Drupal 10.x tests are written and executed using the [PHPUnit](https://phpunit.de/) testing framework. Read more about [running PHPUnit tests in Drupal](https://phpunit.de/).

```go
package main

import (
	"context"

	"dagger/my-module/internal/dagger"
)

type MyModule struct{}

// Run unit tests against a database service
func (m *MyModule) Test(ctx context.Context) (string, error) {
	mariadb := dag.Container().
		From("mariadb:10.11.2").
		WithEnvVariable("MARIADB_USER", "user").
		WithEnvVariable("MARIADB_PASSWORD", "password").
		WithEnvVariable("MARIADB_DATABASE", "drupal").
		WithEnvVariable("MARIADB_ROOT_PASSWORD", "root").
		WithExposedPort(3306).
		AsService(dagger.ContainerAsServiceOpts{UseEntrypoint: true})

	// get Drupal base image
	// install additional dependencies
	drupal := dag.Container().
		From("drupal:10.0.7-php8.2-fpm").
		WithExec([]string{"composer", "require", "drupal/core-dev", "--dev", "--update-with-all-dependencies"})

	// add service binding for MariaDB
	// run kernel tests using PHPUnit
	return drupal.
		WithServiceBinding("db", mariadb).
		WithEnvVariable("SIMPLETEST_DB", "mysql://user:password@db/drupal").
		WithEnvVariable("SYMFONY_DEPRECATIONS_HELPER", "disabled").
		WithWorkdir("/opt/drupal/web/core").
		WithExec([]string{"../../vendor/bin/phpunit", "-v", "--group", "KernelTests"}).
		Stdout(ctx)
}
```

```python
from dagger import dag, function, object_type


@object_type
class MyModule:
    @function
    async def test(self) -> str:
        """Run unit tests against a database service."""
        # get MariaDB base image
        mariadb = (
            dag.container()
            .from_("mariadb:10.11.2")
            .with_env_variable("MARIADB_USER", "user")
            .with_env_variable("MARIADB_PASSWORD", "password")
            .with_env_variable("MARIADB_DATABASE", "drupal")
            .with_env_variable("MARIADB_ROOT_PASSWORD", "root")
            .with_exposed_port(3306)
            .as_service(use_entrypoint=True)
        )

        # get Drupal base image
        # install additional dependencies
        drupal = (
            dag.container()
            .from_("drupal:10.0.7-php8.2-fpm")
            .with_exec(
                [
                    "composer",
                    "require",
                    "drupal/core-dev",
                    "--dev",
                    "--update-with-all-dependencies",
                ]
            )
        )

        # add service binding for MariaDB
        # run kernel tests using PHPUnit
        return await (
            drupal.with_service_binding("db", mariadb)
            .with_env_variable("SIMPLETEST_DB", "mysql://user:password@db/drupal")
            .with_env_variable("SYMFONY_DEPRECATIONS_HELPER", "disabled")
            .with_workdir("/opt/drupal/web/core")
            .with_exec(["../../vendor/bin/phpunit", "-v", "--group", "KernelTests"])
            .stdout()
        )
```

```typescript
import { dag, object, func } from "@dagger.io/dagger"

@object()
class MyModule {
  /**
   * Run unit tests against a database service
   */
  @func()
  async test(): Promise<string> {
    const mariadb = dag
      .container()
      .from("mariadb:10.11.2")
      .withEnvVariable("MARIADB_USER", "user")
      .withEnvVariable("MARIADB_PASSWORD", "password")
      .withEnvVariable("MARIADB_DATABASE", "drupal")
      .withEnvVariable("MARIADB_ROOT_PASSWORD", "root")
      .withExposedPort(3306)
      .asService({ useEntrypoint: true })

    // get Drupal base image
    // install additional dependencies
    const drupal = dag
      .container()
      .from("drupal:10.0.7-php8.2-fpm")
      .withExec([
        "composer",
        "require",
        "drupal/core-dev",
        "--dev",
        "--update-with-all-dependencies",
      ])

    // add service binding for MariaDB
    // run kernel tests using PHPUnit
    return await drupal
      .withServiceBinding("db", mariadb)
      .withEnvVariable("SIMPLETEST_DB", "mysql://user:password@db/drupal")
      .withEnvVariable("SYMFONY_DEPRECATIONS_HELPER", "disabled")
      .withWorkdir("/opt/drupal/web/core")
      .withExec(["../../vendor/bin/phpunit", "-v", "--group", "KernelTests"])
      .stdout()
  }
}
```

