---
description: Create a new DJS parser plugin
agent: agent
tools: ['edit/editFiles', 'search/codebase']
---

# Create a new DJS plugin

Your task is to create a new parser plugin for DJS following the established architecture.

## Plugin structure

Every plugin must:

1. Be created in `plugins/` directory as `<feature>_parser.go`
2. Export a function `func <Feature>Plugin(pb *parser.Builder)`
3. Register custom tokens via `pb.LexerBuilder.RegisterTokenType()`
4. Use interceptors: Token, Statement, and/or Expression as needed
5. Define custom AST nodes with `WriteTo(*strings.Builder)` method

## Reference implementations

Study these existing plugins for patterns:
- [defer_plugin.go](../../plugins/defer_plugin.go) - Token interception, statement transformation, function wrapping
- [or_plugin.go](../../plugins/or_plugin.go) - Expression interception, statement wrapping

## Key patterns

**Token registration and interception:**
```go
tokenType := lb.RegisterTokenType("mytoken")
lb.UseTokenInterceptor(func(l *lexer.Lexer, next func() token.Token) token.Token {
    ret := next()
    if ret.Literal == "mykeyword" {
        ret.Type = tokenType
    }
    return ret
})
```

**Unique variable names (when needed):**
```go
import "github.com/rs/xid"
id := xid.New()
varName := "prefix_" + id.String()
```

**Error handling:**
```go
p.AddError("descriptive error message")
return nil
```

## Output

The `WriteTo` method must generate valid JavaScript. The plugin transforms DJS syntax into equivalent JS at parse time.