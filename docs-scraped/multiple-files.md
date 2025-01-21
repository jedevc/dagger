# Multiple files



```go

// utils/utils.go

import "dagger/<module>/internal/dagger"

func DoThing(client *dagger.Client) *dagger.Directory {
    // we need to pass *dagger.Client in here, since we don't have access to `dag`
	...
}
```

```python

# src/my_module/__init__.py
"""My very own Dagger module"""
from .main import MyModule as MyModule
```

```toml

[build-system]
requires = ["poetry-core>=1.0.0"]
build-backend = "poetry.core.masonry.api"
```

```toml

[build-system]
requires = ["setuptools", "wheel"]
build-backend = "setuptools.build_meta"

[tool.setuptools]
py-modules = ["main"]
```

