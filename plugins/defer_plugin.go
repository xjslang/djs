package plugins

import (
	"github.com/rs/xid"
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/token"
)

type DeferFunctionDeclaration struct {
	*ast.FunctionDeclaration
	prefix string
}

func (fd *DeferFunctionDeclaration) WriteTo(cw *ast.CodeWriter) {
	cw.WriteString("function ")
	fd.Name.WriteTo(cw)
	cw.WriteRune('(')
	for i, param := range fd.Parameters {
		if i > 0 {
			cw.WriteRune(',')
		}
		param.WriteTo(cw)
	}

	var hasDefers bool
	for _, stmt := range fd.Body.Statements {
		if _, ok := stmt.(*DeferStatement); ok {
			hasDefers = true
			break
		}
	}

	if hasDefers {
		deferName := "defers_" + fd.prefix
		indexName := "i_" + fd.prefix
		errorName := "e_" + fd.prefix
		cw.WriteString(") {let " + deferName + "=[];try")
		fd.Body.WriteTo(cw)
		cw.WriteString("finally{" +
			"for(let " + indexName + "=" + deferName + ".length;" + indexName + ">0;" + indexName + "--){" +
			"try{" + deferName + "[" + indexName + "-1]()}catch(" + errorName + "){console.log(" + errorName + ")}}}}",
		)
	} else {
		cw.WriteRune(')')
		fd.Body.WriteTo(cw)
	}
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
	deferTokenType := lb.RegisterTokenType("DeferStatement")
	lb.UseTokenInterceptor(func(l *lexer.Lexer, next func() token.Token) token.Token {
		ret := next()
		if ret.Type == token.IDENT && ret.Literal == "defer" {
			ret.Type = deferTokenType
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
		if p.CurrentToken.Type != deferTokenType {
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
