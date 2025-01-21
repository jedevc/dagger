# Pre-call filtering

Here are a few examples of useful patterns:

```go

// exclude Go tests and test data
+ignore=["**_test.go", "**/testdata/**"]

// exclude binaries
+ignore=["bin"]

// exclude Python dependencies
+ignore=["**/.venv", "**/__pycache__"]

// exclude Node.js dependencies
+ignore=["**/node_modules"]

// exclude Git metadata
+ignore=[".git", "**/.gitignore"]
```

```python

# exclude Go tests and test data
Ignore(["**_test.go", "**/testdata/**"])

# exclude binaries
Ignore(["bin"])

# exclude Python dependencies
Ignore(["**/.venv", "**/__pycache__"])

# exclude Node.js dependencies
Ignore(["**/node_modules"])

# exclude Git metadata
Ignore([".git", "**/.gitignore"])
```

```typescript

// exclude tests and test data
@argument({ ignore: ["**_test.go", "**/testdata/**"] })

// exclude binaries
@argument({ ignore: ["bin"] })

// exclude Python dependencies
@argument({ ignore: ["**/.venv", "**/__pycache__"] })

// exclude Node.js dependencies
@argument({ ignore: ["**/node_modules"] })

// exclude Git metadata
@argument({ ignore: [".git", "**/.gitignore"] })
```

