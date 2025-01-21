# Installation

You can call Dagger Functions from any other Dagger module in your own Dagger module simply by adding it as a module dependency with `dagger install`, as in the following example:

```shell
dagger install github.com/shykes/daggerverse/hello@v0.3.0
```

This module will be added to your `dagger.json`:

```json
...
"dependencies": [
  {
    "name": "hello",
    "source": "github.com/shykes/daggerverse/hello@54d86c6002d954167796e41886a47c47d95a626d"
  }
]
```

When you add a dependency to your module with `dagger install`, the dependent module will be added to the code-generation routines and can be accessed from your own module's code.

The entrypoint to accessing dependent modules from your own module's code is `dag`, the Dagger client, which is pre-initialized. It contains all the core types (like `Container`, `Directory`, etc.), as well as bindings to any dependencies your module has declared.

Here is an example of accessing the installed `hello` module from your own module's code:

```go

func (m *MyModule) Greeting(ctx context.Context) (string, error) {
  return dag.Hello().Hello(ctx)
}
```

```python

@function
async def greeting(self) -> str:
  return await dag.hello().hello()
```

```typescript

@func()
async greeting(): Promise<string> {
  return await dag.hello().hello()
}
```

