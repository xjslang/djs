# DJS (Defers for JavaScript) - AI Coding Agent Instructions

DJS (Defers for JavaScript) is a partial implementation of JavaScript that incorporates the `defer` and `go` constructs. For example:

```djs
let sqlite = require('better-sqlite3')

(function main() {
  let db = sqlite('mydata.db') or |err| {
    console.log('Cannot connect to database', err)
    return
  }
  defer db.close()

  // prepare and execute queries
  let stmt = db.prepare('SELECT * FROM users WHERE active = ?');
  let users = stmt.all(1);
  console.log(`Found ${users.length} active users`);
})()
```

DJS is built on top of the [XJS](https://github.com/xjslang/xjs) tool, a JavaScript parser that allows new constructs to be incorporated into the language through the use of plugins. For example:

```djs
input := `console.log('Hello, world!')`
program, _ := parser.NewBuilder(lb).
  // plugins are executed in the same order they have been installed (FIFO)
	Install(plugins.DeferPlugin).
	Install(plugins.OrPlugin).
	Install(plugins.StrictEqualityPlugin).
	Install(plugins.NewPlugin).
	Install(plugins.ThrowPlugin).
  // create the parser and parse the `input` source to generate the AST
  Build(input).
  ParserProgram()

// finally, compile the AST to code
result, _ := compiler.New().Compile(program)
fmt.Println(result.Code)
```

It is important to note that `XJS` is not a complete implementation of JavaScript, and therefore only supports a limited number of features. For example:

- Only `let` is accepted; `const` and `var` are not allowed.
- Only single-line comments `//` are accepted. Multi-line comments `/* .. */` are not allowed.
- Semicolons are not required.
- `==` are transpiled to `===`. And `===` is not allowed.

This is done on purpose, since the goal of `XJS` is to maintain "sufficient and necessary" structures, while allowing the developer to incorporate new and genuine, not necessarily standard, constructs such as `defer` or `or`.

# General Project Instructions

## Language

- Code must be in **English**: variable names, functions, classes, files, commits, etc.
- Documentation must be in **English**: README, code comments, JSDoc/GoDoc, etc.
- Communication with the user can be in **Spanish** if the user writes in Spanish.

## Examples
```go
// Correct
func GetUserByID(id string) (*User, error)

// Incorrect
func ObtenerUsuarioPorID(id string) (*Usuario, error)
```

## Response Behavior

**Always ask before modifying files.** Do not edit, create, or delete files unless the user explicitly confirms they want changes made.

If you think a code change might be helpful, first explain what you would do and ask for confirmation.

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
