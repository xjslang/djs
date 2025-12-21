package plugins

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/token"
)

// NewExpression represents a new expression (e.g., new Constructor(args))
type NewExpression struct {
	Token       token.Token      // The 'new' token
	Constructor ast.Expression   // The constructor function
	Arguments   []ast.Expression // Optional arguments
}

func (ne *NewExpression) WriteTo(cw *ast.CodeWriter) {
	cw.WriteString("new ")
	ne.Constructor.WriteTo(cw)
	cw.WriteRune('(')
	for i, arg := range ne.Arguments {
		if i > 0 {
			cw.WriteRune(',')
		}
		arg.WriteTo(cw)
	}
	cw.WriteRune(')')
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
	pb.RegisterPrefixOperator(newTokenType, func(tok token.Token, right func() ast.Expression) ast.Expression {
		expr := &NewExpression{
			Token:     tok,
			Arguments: []ast.Expression{},
		}

		// Parse the constructor expression
		constructor := right()
		expr.Constructor = constructor

		// Check if there's a call expression (arguments)
		// If the constructor is already a CallExpression, extract its parts
		if callExpr, ok := constructor.(*ast.CallExpression); ok {
			expr.Constructor = callExpr.Function
			expr.Arguments = callExpr.Arguments
		}

		return expr
	})
}
