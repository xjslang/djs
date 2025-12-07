---
description: Create a plugin that transforms existing JavaScript syntax
agent: agent
tools: ['edit/createFile', 'edit/editFiles', 'search/codebase']
---

# Create a transformation plugin

Your task is to create a plugin that **modifies existing JavaScript syntax** rather than adding new features.

## Key differences from feature plugins

**Feature plugins** (like `defer` or `or`) add NEW syntax to the language:
- Register custom tokens for new keywords
- Define custom AST nodes with `WriteTo` methods
- Use Statement/Expression interceptors to parse new constructs

**Transformation plugins** (like `strict_equality`) modify EXISTING syntax:
- Use existing tokens (no custom registration needed)
- Only need Token Interceptors (no AST nodes required)
- Transform tokens directly in the lexer

## When to use transformation plugins

Use this approach when you want to:
- Replace operators (`==` → `===`, `var` → `let`, etc.)
- Transform keywords (`function` → `async function`)
- Enforce coding standards (remove semicolons, add semicolons, etc.)
- Modify literals (string quotes, number formats, etc.)

## Structure of a transformation plugin

```go
package plugins

import (
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/token"
)

// MyTransformPlugin describes what transformation it performs
func MyTransformPlugin(pb *parser.Builder) {
	lb := pb.LexerBuilder
	lb.UseTokenInterceptor(func(l *lexer.Lexer, next func() token.Token) token.Token {
		ret := next()
		
		// Check token type and/or literal
		if ret.Type == token.SOME_TYPE && ret.Literal == "original" {
			ret.Literal = "transformed"
		}
		
		return ret
	})
}
```

## Real example: StrictEqualityPlugin

This plugin transforms weak equality into strict equality:

```go
func StrictEqualityPlugin(pb *parser.Builder) {
	lb := pb.LexerBuilder
	lb.UseTokenInterceptor(func(l *lexer.Lexer, next func() token.Token) token.Token {
		ret := next()
		// Transform == to ===
		if ret.Type == token.EQ && ret.Literal == "==" {
			ret.Literal = "==="
		}
		// Transform != to !==
		if ret.Type == token.NOT_EQ && ret.Literal == "!=" {
			ret.Literal = "!=="
		}
		return ret
	})
}
```

**What happens:**
- Input: `1 == 2` → Output: `1 === 2`
- Input: `x != y` → Output: `x !== y`

## Common token types

The `xjs` parser provides these token types (among others):

**Keywords:**
- `token.FUNCTION`, `token.LET`, `token.CONST`, `token.VAR`
- `token.IF`, `token.ELSE`, `token.FOR`, `token.WHILE`
- `token.RETURN`, `token.BREAK`, `token.CONTINUE`

**Operators:**
- `token.EQ` (`==`), `token.NOT_EQ` (`!=`)
- `token.PLUS` (`+`), `token.MINUS` (`-`), `token.ASTERISK` (`*`)
- `token.LT` (`<`), `token.GT` (`>`), `token.LTE` (`<=`), `token.GTE` (`>=`)

**Literals:**
- `token.IDENT` (identifiers like variable names)
- `token.STRING`, `token.INT`, `token.FLOAT`

**Delimiters:**
- `token.LPAREN` (`(`), `token.RPAREN` (`)`), `token.LBRACE` (`{`), `token.RBRACE` (`}`)
- `token.SEMICOLON` (`;`), `token.COMMA` (`,`)

## Step-by-step guide

1. **Identify what to transform**
   - What token type? (`token.EQ`, `token.VAR`, etc.)
   - What literal value? (`"=="`, `"var"`, etc.)

2. **Create the plugin file**
   - File: `plugins/<name>_plugin.go`
   - Function: `func <Name>Plugin(pb *parser.Builder)`

3. **Write the token interceptor**
   ```go
   lb.UseTokenInterceptor(func(l *lexer.Lexer, next func() token.Token) token.Token {
       ret := next()
       if ret.Type == token.TARGET_TYPE && ret.Literal == "target" {
           ret.Literal = "replacement"
       }
       return ret
   })
   ```

4. **Test it**
   - Add to `main.go`: `.Install(plugins.YourPlugin)`
   - Create test input file
   - Run: `go run . test.js`

## Important notes

- **Always call `next()`** - This gets the original token from the lexer
- **Always return a token** - Even if you don't modify it
- **Check both type AND literal** - Some tokens have the same type but different literals
- **Order matters** - Token interceptors are called in the order they're registered
- **No AST needed** - Transformations happen before parsing, so you don't need custom AST nodes

## Common patterns

**Replace operator:**
```go
if ret.Type == token.EQ && ret.Literal == "==" {
    ret.Literal = "==="
}
```

**Replace keyword:**
```go
if ret.Type == token.VAR && ret.Literal == "var" {
    ret.Type = token.LET
    ret.Literal = "let"
}
```

**Conditional transformation:**
```go
if ret.Type == token.FUNCTION && someCondition {
    ret.Literal = "async " + ret.Literal
}
```

## Debugging tips

1. **Print tokens to see what you're working with:**
   ```go
   lb.UseTokenInterceptor(func(l *lexer.Lexer, next func() token.Token) token.Token {
       ret := next()
       fmt.Printf("Token: Type=%v, Literal=%q\n", ret.Type, ret.Literal)
       return ret
   })
   ```

2. **Start simple** - Transform one thing first, then add more

3. **Test with minimal input** - Use small JS files to verify transformations

## Examples of transformation plugins

**NoVarPlugin** - Replace `var` with `let`:
```go
func NoVarPlugin(pb *parser.Builder) {
	lb := pb.LexerBuilder
	lb.UseTokenInterceptor(func(l *lexer.Lexer, next func() token.Token) token.Token {
		ret := next()
		if ret.Type == token.VAR && ret.Literal == "var" {
			ret.Type = token.LET
			ret.Literal = "let"
		}
		return ret
	})
}
```

**DoubleQuotesPlugin** - Force double quotes for strings:
```go
func DoubleQuotesPlugin(pb *parser.Builder) {
	lb := pb.LexerBuilder
	lb.UseTokenInterceptor(func(l *lexer.Lexer, next func() token.Token) token.Token {
		ret := next()
		if ret.Type == token.STRING {
			// Remove existing quotes and add double quotes
			content := strings.Trim(ret.Literal, "'\"")
			ret.Literal = "\"" + content + "\""
		}
		return ret
	})
}
```

## When you need more than token interception

If your transformation requires:
- Modifying AST structure (not just token literals)
- Adding/removing statements
- Changing expression order

Then you need a **feature plugin** with Statement/Expression interceptors. See [plugin.prompt.md](./plugin.prompt.md) for that approach.

## Summary checklist

- [ ] Plugin file created: `plugins/<name>_plugin.go`
- [ ] Function signature: `func <Name>Plugin(pb *parser.Builder)`
- [ ] Token interceptor registered
- [ ] Token type and literal checked
- [ ] Always calls `next()` and returns token
- [ ] Tested with sample input

## Reference implementations

- [strict_equality_plugin.go](../../plugins/strict_equality_plugin.go) - Transforms `==` to `===`
- [defer_plugin.go](../../plugins/defer_plugin.go) - Example of complex feature plugin (for comparison)
- [or_plugin.go](../../plugins/or_plugin.go) - Example of expression transformation
