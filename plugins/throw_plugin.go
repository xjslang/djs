package plugins

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/token"
)

// ThrowStatement represents a throw statement in the AST
type ThrowStatement struct {
	Token    token.Token    // The 'throw' token
	Argument ast.Expression // The expression to throw
}

func (ts *ThrowStatement) WriteTo(cw *ast.CodeWriter) {
	cw.WriteString("throw ")
	ts.Argument.WriteTo(cw)
	cw.WriteRune(';')
}

// ThrowPlugin adds support for the 'throw' statement
func ThrowPlugin(pb *parser.Builder) {
	lb := pb.LexerBuilder
	throwTokenType := lb.RegisterTokenType("THROW")

	// Intercept 'throw' identifier and convert it to THROW token
	lb.UseTokenInterceptor(func(l *lexer.Lexer, next func() token.Token) token.Token {
		tok := next()
		if tok.Type == token.IDENT && tok.Literal == "throw" {
			tok.Type = throwTokenType
		}
		return tok
	})

	// Statement interceptor for 'throw'
	pb.UseStatementInterceptor(func(p *parser.Parser, next func() ast.Statement) ast.Statement {
		if p.CurrentToken.Type != throwTokenType {
			return next()
		}
		stmt := &ThrowStatement{Token: p.CurrentToken}

		// throw always requires an argument
		if p.PeekToken.Type == token.SEMICOLON || p.PeekToken.Type == token.EOF || p.PeekToken.Type == token.RBRACE {
			p.AddError("throw statement requires an argument")
			return nil
		}

		p.NextToken() // move to the expression
		stmt.Argument = p.ParseExpression()
		p.ExpectSemicolonASI()
		return stmt
	})
}
