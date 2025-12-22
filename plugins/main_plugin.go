package plugins

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/token"
)

type MainFunctionDeclaration struct {
	*ast.FunctionDeclaration
}

func (mfd *MainFunctionDeclaration) WriteTo(cw *ast.CodeWriter) {
	mfd.FunctionDeclaration.WriteTo(cw)
	cw.WriteString(";main()")
}

func MainPlugin(pb *parser.Builder) {
	pb.UseStatementInterceptor(func(p *parser.Parser, next func() ast.Statement) ast.Statement {
		if p.CurrentToken.Type != token.FUNCTION {
			return next()
		}

		// Check if we're NOT inside a function (top-level only)
		if p.IsInFunction() {
			return next()
		}

		// Parse the function normally
		funcDecl := p.ParseFunctionStatement()
		if funcDecl.Name != nil && funcDecl.Name.Value == "main" {
			return &MainFunctionDeclaration{
				FunctionDeclaration: funcDecl,
			}
		}

		return funcDecl
	})
}
