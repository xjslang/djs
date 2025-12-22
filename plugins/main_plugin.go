package plugins

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/token"
)

// IIFEFunctionDeclaration wraps a function declaration in an IIFE
type IIFEFunctionDeclaration struct {
	*ast.FunctionDeclaration
}

// WriteTo generates an IIFE: (function name() {...})()
func (iife *IIFEFunctionDeclaration) WriteTo(cw *ast.CodeWriter) {
	cw.WriteString("(function ")
	iife.Name.WriteTo(cw)
	cw.WriteRune('(')
	for i, param := range iife.Parameters {
		if i > 0 {
			cw.WriteRune(',')
		}
		param.WriteTo(cw)
	}
	cw.WriteRune(')')
	iife.Body.WriteTo(cw)
	cw.WriteString(")()")
}

// MainPlugin transforms top-level 'main' functions into IIFEs
func MainPlugin(pb *parser.Builder) {
	pb.UseStatementInterceptor(func(p *parser.Parser, next func() ast.Statement) ast.Statement {
		// Only intercept FUNCTION tokens
		if p.CurrentToken.Type != token.FUNCTION {
			return next()
		}

		// Check if we're NOT inside a function (top-level only)
		if p.IsInFunction() {
			return next()
		}

		// Parse the function normally
		funcDecl := p.ParseFunctionStatement()

		// Check if the function name is "main"
		if funcDecl.Name != nil && funcDecl.Name.Value == "main" {
			return &IIFEFunctionDeclaration{
				FunctionDeclaration: funcDecl,
			}
		}

		// Return the original function if it's not named "main"
		return funcDecl
	})
}
