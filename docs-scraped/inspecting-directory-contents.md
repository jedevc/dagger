# Inspecting directory contents

Another way to debug how directories are being filtered is to create a function that receives a `Directory` as input, and returns the same `Directory`:

```go

func (m *MyModule) Debug(
  ctx context.Context,
  // +ignore=["*", "!analytics"]
  source *dagger.Directory,
) *dagger.Directory {
  return source
}
```

```python

@function
async def foo(
    self,
    source: Annotated[
        dagger.Directory, Ignore(["*", "!analytics"])
    ],
) -> dagger.Directory:
    return source
```

```typescript

@func()
debug(
   @argument({ ignore: ["*", "!analytics"] }) source: Directory,
): Directory {
  return source
}
```

