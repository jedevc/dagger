# Copy a file to the Dagger module runtime container for custom processing

The following Dagger Function accepts a `File` argument and copies the specified file to the Dagger module runtime container. This makes it possible to add one or more files to the runtime container and manipulate them using custom logic.

```go
package main

import (
	"context"
	"dagger/my-module/internal/dagger"
	"fmt"
	"os"
)

type MyModule struct{}

// Copy a file to the Dagger module runtime container for custom processing
func (m *MyModule) CopyFile(ctx context.Context, source *dagger.File) {
	source.Export(ctx, "foo.txt")
	// your custom logic here
	// for example, read and print the file in the Dagger Engine container
	fmt.Println(os.ReadFile("foo.txt"))
}
```

```python
import anyio

import dagger
from dagger import function, object_type


@object_type
class MyModule:
    @function
    async def copy_file(self, source: dagger.File):
        """Copy a file to the Dagger module runtime container for custom processing"""
        await source.export("foo.txt")
        # your custom logic here
        # for example, read and print the file in the Dagger Engine container
        print(await anyio.Path("foo.txt").read_text())
```

```typescript
import { object, func, File } from "@dagger.io/dagger"
import * as fs from "fs"

@object()
class MyModule {
  // Copy a file to the Dagger module runtime container for custom processing
  @func()
  async copyFile(source: File) {
    await source.export("foo.txt")
    // your custom logic here
    // for example, read and print the file in the Dagger Engine container
    console.log(fs.readFileSync("foo.txt", "utf8"))
  }
}
```

