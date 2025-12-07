package plugins

import (
	"fmt"
	"strings"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/token"
)

type ExpressionStatement struct {
	*ast.ExpressionStatement
}

type LetStatement struct {
	*ast.LetStatement
}

type OrExpression struct {
	Exprression   ast.Expression
	FallbackBlock *ast.BlockStatement
}

// Override ast.ExpressionStatement.WriteTo
func (es *ExpressionStatement) WriteTo(b *strings.Builder) {
	if stmt, ok := es.Expression.(*OrExpression); ok {
		b.WriteString("try{")
		stmt.Exprression.WriteTo(b)
		b.WriteString("}catch")
		stmt.FallbackBlock.WriteTo(b)
	} else {
		es.ExpressionStatement.WriteTo(b)
	}
}

// Override ast.LetStatement.WriteTo
func (ls *LetStatement) WriteTo(b *strings.Builder) {
	if oe, ok := ls.Value.(*OrExpression); ok {
		b.WriteString("let ")
		ls.Name.WriteTo(b)
		b.WriteString(";try{")
		ls.Name.WriteTo(b)
		b.WriteRune('=')
		ls.Value.WriteTo(b)
		b.WriteString("}catch")
		oe.FallbackBlock.WriteTo(b)
	} else {
		ls.LetStatement.WriteTo(b)
	}
}

func (oe *OrExpression) WriteTo(b *strings.Builder) {
	oe.Exprression.WriteTo(b)
}

func OrPlugin(pb *parser.Builder) {
	lb := pb.LexerBuilder
	orTokenType := lb.RegisterTokenType("or")
	lb.UseTokenInterceptor(func(l *lexer.Lexer, next func() token.Token) token.Token {
		ret := next()
		if ret.Literal == "or" {
			ret.Type = orTokenType
		}
		return ret
	})

	pb.UseStatementInterceptor(func(p *parser.Parser, next func() ast.Statement) ast.Statement {
		ret := next()
		switch stmt := ret.(type) {
		case *ast.ExpressionStatement:
			return &ExpressionStatement{
				ExpressionStatement: stmt,
			}
		case *ast.LetStatement:
			return &LetStatement{
				LetStatement: stmt,
			}
		}
		return ret
	})

	pb.UseExpressionInterceptor(func(p *parser.Parser, next func() ast.Expression) ast.Expression {
		exp := next()
		if p.PeekToken.Type == orTokenType {
			p.NextToken() // consume 'or' and move to {
			if p.PeekToken.Type != token.LBRACE {
				p.AddError(fmt.Sprintf("expected { after or, got %v", p.PeekToken))
				return exp
			}
			p.NextToken() // consume {
			fallbackBlock := p.ParseBlockStatement()
			return &OrExpression{
				Exprression:   exp,
				FallbackBlock: fallbackBlock,
			}
		}
		return exp
	})
}
