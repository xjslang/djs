package plugins

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/token"
)

// NewExpression represents a new expression (e.g., new Constructor(args))
type NewExpression struct {
	Token token.Token    // The 'new' token
	Right ast.Expression // The constructor expression (may include arguments)
}

func (ne *NewExpression) WriteTo(cw *ast.CodeWriter) {
	cw.AddMapping(ne.Token.Start)
	cw.WriteString("new ")
	ne.Right.WriteTo(cw)
}

// NewPlugin adds support for the 'new' operator to create instances
func NewPlugin(pb *parser.Builder) {
	lb := pb.LexerBuilder

	// Register NEW token type
	newTokenType := lb.RegisterTokenType("NEW")

	// Intercept 'new' identifier and convert it to NEW token
	lb.UseTokenInterceptor(func(l *lexer.Lexer, next func() token.Token) token.Token {
		ret := next()
		if ret.Type == token.IDENT && ret.Literal == "new" {
			ret.Type = newTokenType
		}
		return ret
	})

	// Register 'new' as a prefix operator
	_ = pb.RegisterPrefixOperator(newTokenType, func(tok token.Token, right func() ast.Expression) ast.Expression {
		return &NewExpression{
			Token: tok,
			Right: right(),
		}
	})
}
