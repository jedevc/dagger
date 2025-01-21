# Continue using a container after command execution fails

The following Dagger Function demonstrates how to continue using a container after a command executed within it fails. A common use case for this is to export a report that a test suite tool generates.

:::note
The caveat with this approach is that forcing a zero exit code on a failure caches the failure. This may not be desired depending on the use case.
:::

```go
package main

import (
	"context"
	"fmt"

	"dagger/my-module/internal/dagger"
)

type MyModule struct{}

var script = `#!/bin/sh
echo "Test Suite"
echo "=========="
echo "Test 1: PASS" | tee -a report.txt
echo "Test 2: FAIL" | tee -a report.txt
echo "Test 3: PASS" | tee -a report.txt
exit 1
`

type TestResult struct {
	Report   *dagger.File
	ExitCode int
}

// Handle errors
func (m *MyModule) Test(ctx context.Context) (*TestResult, error) {
	ctr, err := dag.
		Container().
		From("alpine").
		// add script with execution permission to simulate a testing tool
		WithNewFile("/run-tests", script, dagger.ContainerWithNewFileOpts{Permissions: 0o750}).
		// run-tests but allow any return code
		WithExec([]string{"/run-tests"}, dagger.ContainerWithExecOpts{Expect: dagger.Any}).
		// the result of `sync` is the container, which allows continued chaining
		Sync(ctx)
	if err != nil {
		// unexpected error, could be network failure.
		return nil, fmt.Errorf("run tests: %w", err)
	}
	// save report for inspection.
	report := ctr.File("report.txt")

	// use the saved exit code to determine if the tests passed.
	exitCode, err := ctr.ExitCode(ctx)
	if err != nil {
		// exit code not found
		return nil, fmt.Errorf("get exit code: %w", err)
	}

	// Return custom type
	return &TestResult{
		Report:   report,
		ExitCode: exitCode,
	}, nil
}
```

```python
import dagger
from dagger import DaggerError, dag, field, function, object_type

SCRIPT = """#!/bin/sh
echo "Test Suite"
echo "=========="
echo "Test 1: PASS" | tee -a report.txt
echo "Test 2: FAIL" | tee -a report.txt
echo "Test 3: PASS" | tee -a report.txt
exit 1
"""


@object_type
class TestResult:
    report: dagger.File = field()
    exit_code: int = field()


@object_type
class MyModule:
    @function
    async def test(self) -> TestResult:
        """Handle errors"""
        try:
            ctr = await (
                dag.container()
                .from_("alpine")
                # add script with execution permission to simulate a testing tool.
                .with_new_file("/run-tests", SCRIPT, permissions=0o750)
                # run-tests but allow any return code
                .with_exec(["/run-tests"], expect=dagger.ReturnType.ANY)
                # the result of `sync` is the container, which allows continued chaining
                .sync()
            )

            # save report for inspection.
            report = ctr.file("report.txt")

            # use the saved exit code to determine if the tests passed.
            exit_code = await ctr.exit_code()

            return TestResult(report=report, exit_code=exit_code)
        except DaggerError as e:
            # DaggerError is the base class for all errors raised by Dagger
            msg = "Unexpected Dagger error"
            raise RuntimeError(msg) from e


# ruff: noqa: RET505
```

```typescript
import { dag, object, func, File, ReturnType } from "@dagger.io/dagger"

const SCRIPT = `#!/bin/sh
echo "Test Suite"
echo "=========="
echo "Test 1: PASS" | tee -a report.txt
echo "Test 2: FAIL" | tee -a report.txt
echo "Test 3: PASS" | tee -a report.txt
exit 1
`

@object()
class TestResult {
  @func()
  report: File

  @func()
  exitCode: number
}

@object()
class MyModule {
  /**
   * Handle errors
   */
  @func()
  async test(): Promise<TestResult> {
    const ctr = await dag
      .container()
      .from("alpine")
      // add script with execution permission to simulate a testing tool.
      .withNewFile("/run-tests", SCRIPT, { permissions: 0o750 })
      // run-tests but allow any return code
      .withExec(["/run-tests"], { expect: ReturnType.Any })
      // the result of `sync` is the container, which allows continued chaining
      .sync()

    const result = new TestResult()
    // save report for inspection.
    result.report = ctr.file("report.txt")

    // use the saved exit code to determine if the tests passed
    result.exitCode = await ctr.exitCode()

    return result
  }
}
```

