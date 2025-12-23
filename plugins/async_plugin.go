package plugins

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/token"
)

// AsyncFunctionDeclaration wraps a standard function declaration and emits
// "async function" in the generated JavaScript.
type AsyncFunctionDeclaration struct {
	*ast.FunctionDeclaration
}

func (afd *AsyncFunctionDeclaration) WriteTo(cw *ast.CodeWriter) {
	cw.AddMapping(afd.Token.Line, afd.Token.Column)
	cw.WriteString("async function ")
	afd.Name.WriteTo(cw)
	cw.WriteRune('(')
	for i, param := range afd.Parameters {
		if i > 0 {
			cw.WriteRune(',')
		}
		param.WriteTo(cw)
	}
	cw.WriteRune(')')
	afd.Body.WriteTo(cw)
}

// AsyncFunctionExpression wraps a standard function expression and emits
// "async function" in the generated JavaScript.
type AsyncFunctionExpression struct {
	*ast.FunctionExpression
}

func (afe *AsyncFunctionExpression) WriteTo(cw *ast.CodeWriter) {
	cw.AddMapping(afe.Token.Line, afe.Token.Column)
	cw.WriteString("async function")
	if afe.Name != nil {
		cw.WriteRune(' ')
		afe.Name.WriteTo(cw)
	}
	cw.WriteRune('(')
	for i, param := range afe.Parameters {
		if i > 0 {
			cw.WriteRune(',')
		}
		param.WriteTo(cw)
	}
	cw.WriteRune(')')
	afe.Body.WriteTo(cw)
}

// AwaitExpression represents `await <expr>`.
type AwaitExpression struct {
	Token token.Token
	Right ast.Expression
}

func (ae *AwaitExpression) WriteTo(cw *ast.CodeWriter) {
	cw.AddMapping(ae.Token.Line, ae.Token.Column)
	cw.WriteString("await ")
	if ae.Right != nil {
		ae.Right.WriteTo(cw)
	}
}

// AsyncPlugin adds support for `async function` and `await`.
// Transpilation is direct (no transformation beyond keyword handling).
func AsyncPlugin(pb *parser.Builder) {
	lb := pb.LexerBuilder

	// Register custom token types
	asyncTokenType := lb.RegisterTokenType("ASYNC")
	awaitTokenType := lb.RegisterTokenType("AWAIT")

	// Intercept identifiers and convert to our custom tokens when matching keywords
	lb.UseTokenInterceptor(func(l *lexer.Lexer, next func() token.Token) token.Token {
		ret := next()
		if ret.Type == token.IDENT {
			switch ret.Literal {
			case "async":
				ret.Type = asyncTokenType
			case "await":
				ret.Type = awaitTokenType
			}
		}
		return ret
	})

	// Handle `async function` as a statement (declaration)
	pb.UseStatementInterceptor(func(p *parser.Parser, next func() ast.Statement) ast.Statement {
		if p.CurrentToken.Type != asyncTokenType {
			return next()
		}

		// consume `async`
		p.NextToken()
		if p.CurrentToken.Type != token.FUNCTION {
			p.AddError("async must be followed by function")
			return nil
		}

		fd := p.ParseFunctionStatement()
		if fd == nil {
			return nil
		}
		return &AsyncFunctionDeclaration{FunctionDeclaration: fd}
	})

	// Handle `async function` as an expression and `await <expr>`
	pb.UseExpressionInterceptor(func(p *parser.Parser, next func() ast.Expression) ast.Expression {
		// async function (expression)
		if p.CurrentToken.Type == asyncTokenType {
			p.NextToken()
			if p.CurrentToken.Type != token.FUNCTION {
				p.AddError("async must be followed by function")
				return nil
			}
			expr := p.ParseFunctionExpression()
			if fe, ok := expr.(*ast.FunctionExpression); ok {
				return &AsyncFunctionExpression{FunctionExpression: fe}
			}
			return expr
		}

		// await <expr>
		if p.CurrentToken.Type == awaitTokenType {
			// Only allow inside function context for now
			if !p.IsInFunction() {
				p.AddError("await can only be used inside functions")
				return nil
			}
			tok := p.CurrentToken
			// parse the right-hand side with UNARY precedence
			p.NextToken()
			right := p.ParseExpressionWithPrecedence(parser.UNARY)
			// allow following infix/postfix via remaining expression handling
			return p.ParseRemainingExpression(&AwaitExpression{Token: tok, Right: right})
		}

		return next()
	})
}
