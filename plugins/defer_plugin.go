package plugins

import (
	"github.com/rs/xid"
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/token"
)

// writeFunctionWithDefers writes a function with defer support
func writeFunctionWithDefers(cw *ast.CodeWriter, name *ast.Identifier, parameters []*ast.Identifier, body *ast.BlockStatement, prefix string) {
	cw.WriteString("function")
	if name != nil {
		cw.WriteRune(' ')
		name.WriteTo(cw)
	}
	cw.WriteRune('(')
	for i, param := range parameters {
		if i > 0 {
			cw.WriteRune(',')
		}
		param.WriteTo(cw)
	}

	var hasDefers bool
	for _, stmt := range body.Statements {
		if _, ok := stmt.(*DeferStatement); ok {
			hasDefers = true
			break
		}
	}

	if hasDefers {
		deferName := "defers_" + prefix
		indexName := "i_" + prefix
		errorName := "e_" + prefix
		cw.WriteString(") {let " + deferName + "=[];try")
		body.WriteTo(cw)
		cw.WriteString("finally{" +
			"for(let " + indexName + "=" + deferName + ".length;" + indexName + ">0;" + indexName + "--){" +
			"try{" + deferName + "[" + indexName + "-1]()}catch(" + errorName + "){console.log(" + errorName + ")}}}}",
		)
	} else {
		cw.WriteRune(')')
		body.WriteTo(cw)
	}
}

type DeferFunctionDeclaration struct {
	*ast.FunctionDeclaration
	prefix string
	async  bool
}

func (fd *DeferFunctionDeclaration) WriteTo(cw *ast.CodeWriter) {
	if fd.async {
		cw.WriteString("async ")
	}
	writeFunctionWithDefers(cw, fd.Name, fd.Parameters, fd.Body, fd.prefix)
}

type DeferFunctionExpression struct {
	*ast.FunctionExpression
	prefix string
}

func (fe *DeferFunctionExpression) WriteTo(cw *ast.CodeWriter) {
	writeFunctionWithDefers(cw, fe.Name, fe.Parameters, fe.Body, fe.prefix)
}

type DeferStatement struct {
	Body   *ast.BlockStatement
	prefix string
}

func (ds *DeferStatement) WriteTo(cw *ast.CodeWriter) {
	deferName := "defers_" + ds.prefix
	cw.WriteString(deferName + ".push(() =>")
	ds.Body.WriteTo(cw)
	cw.WriteRune(')')
}

func DeferPlugin(pb *parser.Builder) {
	id := xid.New()
	lb := pb.LexerBuilder
	deferToken := lb.RegisterTokenType("DEFER")
	asyncToken := lb.RegisterTokenType("ASYNC")

	lb.UseTokenInterceptor(func(l *lexer.Lexer, next func() token.Token) token.Token {
		ret := next()
		if ret.Type != token.IDENT {
			return ret
		}
		switch ret.Literal {
		case "defer":
			ret.Type = deferToken
		case "async":
			ret.Type = asyncToken
		}
		return ret
	})

	pb.UseStatementInterceptor(func(p *parser.Parser, next func() ast.Statement) ast.Statement {
		if p.CurrentToken.Type != token.FUNCTION {
			return next()
		}

		return &DeferFunctionDeclaration{
			prefix:              id.String(),
			FunctionDeclaration: p.ParseFunctionStatement(),
		}
	})

	pb.UseStatementInterceptor(func(p *parser.Parser, next func() ast.Statement) ast.Statement {
		if p.CurrentToken.Type != asyncToken {
			return next()
		}

		p.NextToken() // consume 'async'
		return &DeferFunctionDeclaration{
			async:               true,
			prefix:              id.String(),
			FunctionDeclaration: p.ParseFunctionStatement(),
		}
	})

	pb.UseExpressionInterceptor(func(p *parser.Parser, next func() ast.Expression) ast.Expression {
		if p.CurrentToken.Type != token.FUNCTION {
			return next()
		}

		expr := p.ParseFunctionExpression()
		if fe, ok := expr.(*ast.FunctionExpression); ok {
			return &DeferFunctionExpression{
				prefix:             id.String(),
				FunctionExpression: fe,
			}
		}
		return expr
	})

	pb.UseStatementInterceptor(func(p *parser.Parser, next func() ast.Statement) ast.Statement {
		if p.CurrentToken.Type != deferToken {
			return next()
		}

		if !p.IsInFunction() {
			p.AddError("defer statement can only be used inside functions")
			return nil
		}

		stmt := &DeferStatement{prefix: id.String()}
		if p.PeekToken.Type == token.LBRACE {
			p.NextToken() // consume {
			stmt.Body = p.ParseBlockStatement()
		} else {
			p.NextToken() // move to statement
			stmt.Body = &ast.BlockStatement{}
			stmt.Body.Statements = []ast.Statement{p.ParseStatement()}
		}
		return stmt
	})
}
