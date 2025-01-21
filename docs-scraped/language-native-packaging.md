# Language-native packaging

The structure of a Dagger module mimics that of each language's conventional packaging mechanisms and tools.

```go

// go.work
go 1.21.7

use (
	./path/to/mymodule
)
```

```shell

# executed by the runtime
uv pip install -r requirements.lock -e ./sdk -e .
```

```shell

# executed by the runtime container
yarn install --production
```

