# Post-call filtering

Here are a few examples of useful patterns:

```go

// exclude all Markdown files
dirOpts := dagger.ContainerWithDirectoryOpts{
  Exclude: "*.md*",
}

// include only the build output directory
dirOpts := dagger.ContainerWithDirectoryOpts{
  Include: "build",
}

// include only ZIP files
dirOpts := dagger.DirectoryWithDirectoryOpts{
  Include: "*.zip",
}

// exclude Git metadata
dirOpts := dagger.DirectoryWithDirectoryOpts{
  Exclude: "*.git",
}
```

```python

# exclude all Markdown files
dir_opts = {"exclude": ["*.md*"]}

# include only the build output directory
dir_opts = {"include": ["build"]}

# include only ZIP files
dir_opts = {"include": ["*.zip"]}

# exclude Git metadata
dir_opts = {"exclude": ["*.git"]}
```

```typescript

// exclude all Markdown files
const dirOpts = { exclude: ["*.md*"] }

// include only the build output directory
const dirOpts = { include: ["build"] }

// include only ZIP files
const dirOpts = { include: ["*.zip"] }

// exclude Git metadata
const dirOpts = { exclude: ["*.git"] }
```

