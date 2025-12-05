---
description: Add or modify commands in magefile.go
agent: agent
tools: ['edit/editFiles', 'search/codebase']
---

# Modify magefile.go

Your task is to add or modify commands in the `magefile.go` file of this project.

## Rules

1. Keep the build tag `//go:build mage` at the top of the file
2. Use the `github.com/magefile/mage/sh` package to run external commands
3. Public functions (PascalCase) become mage commands
4. Add a comment above each function describing what it does (mage uses this as help)
5. Group imports properly
6. Return `error` so mage can report failures

## Common patterns
```go
// Simple command
func NombreComando() error {
    return sh.RunV("comando", "arg1", "arg2")
}

// Command with multiple steps
func NombreComando() error {
    if err := sh.RunV("paso1"); err != nil {
        return err
    }
    return sh.RunV("paso2")
}
```

## Reference of the current file

Check [magefile.go](../../magefile.go) to see the existing structure and maintain consistency.