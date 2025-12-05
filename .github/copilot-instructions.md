# DJS (DeferJavaScript) - AI Coding Agent Instructions

# General Project Instructions

## Language

- Code must be in **English**: variable names, functions, classes, files, commits, etc.
- Documentation must be in **English**: README, code comments, JSDoc/GoDoc, etc.
- Communication with the user can be in **Spanish** if the user writes in Spanish.

## Ejemplos
```go
// Correct
func GetUserByID(id string) (*User, error)

// Incorrect
func ObtenerUsuarioPorID(id string) (*Usuario, error)
```

## Project Overview

DJS is a JavaScript language extension that adds Go-style `defer` statements and error-handling `or` blocks. Built on the `xjslang/xjs` parser framework (v0.1.0), it uses a **plugin-based architecture** where each language feature is implemented as a parser plugin that intercepts and transforms the AST during parsing.

## Architecture

### Plugin System (`plugins/`)

Plugins extend the `xjs` parser using three interceptor types:

1. **Token Interceptors** (`UseTokenInterceptor`) - Modify lexer tokens on-the-fly
2. **Statement Interceptors** (`UseStatementInterceptor`) - Transform parsed statements
3. **Expression Interceptors** (`UseExpressionInterceptor`) - Transform parsed expressions

Each plugin defines custom AST nodes that implement `WriteTo(*strings.Builder)` to generate final JavaScript output.

### Key Plugins

**`defer_parser.go`** - Go-style defer statements
- Wraps function bodies with try-finally and defer stack
- Uses `github.com/rs/xid` to generate unique variable names per function (e.g., `defers_<xid>`)
- Pattern: Intercepts `FUNCTION` tokens → wraps in `DeferFunctionDeclaration` → generates `try-finally` with reverse-order execution
- Only transforms functions containing `defer` statements

**`or_parser.go`** - Error handling fallback blocks
- Syntax: `expression or { fallback }` transpiles to `try-catch` blocks
- Wraps standard AST nodes (`ExpressionStatement`, `LetStatement`) to detect and transform `OrExpression`
- Special handling for `let` statements to hoist declaration before try-catch

## Development Workflow

**Linting:**
```bash
mage lint             # Runs golangci-lint
```

The project uses [Mage](https://magefile.org/) for build tasks - see `magefile.go` for available targets. Add new targets incrementally per `.github/prompts/magefile.prompt.md`.

## Conventions

- **No main logic yet** - `main.go` is placeholder Hello World code awaiting actual DJS transpiler implementation
- **AST node embedding** - Custom nodes embed base `xjs` AST nodes (e.g., `*ast.FunctionDeclaration`) and override `WriteTo`
- **Parser state tracking** - Use `p.IsInFunction()` to enforce syntax rules (e.g., defer only inside functions)
- **Unique identifiers** - Generate collision-free variable names using `xid.New().String()` per function/plugin
- **Error handling** - Use `p.AddError()` during parsing; errors are collected not panicked

## Adding New Language Features

1. Create plugin function signature: `func MyPlugin(pb *parser.Builder)`
2. Register custom tokens via `pb.LexerBuilder.RegisterTokenType()`
3. Add interceptors to transform tokens → AST nodes → JavaScript output
4. Define custom AST node structs with `WriteTo(*strings.Builder)` method
5. Hook plugin into parser initialization (currently not shown in main.go - needs implementation)

## Integration Points

- **External dependency**: `github.com/xjslang/xjs` parser framework provides all AST/lexer/parser primitives
- **ID generation**: `github.com/rs/xid` for unique variable naming
- **Output**: Direct JavaScript code generation via `WriteTo` methods (no intermediate IR)

## Critical Notes

- Plugins operate on parser AST, not runtime - transformation happens at parse time
- Multiple statement interceptors chain via `next()` function calls
- Token type registration must happen before lexer interceptor registration
- `WriteTo` methods are responsible for complete JavaScript syntax correctness
