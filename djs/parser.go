package djs

import (
	"github.com/xjslang/xjs"
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

func Parse(input []byte) (*js.Program, error) {
	deferTyp := token.RegisterType("DEFER", "defer")

	sb := xjs.ScannerBuilder()
	// now the scanner can scan the "defer" keyword
	sb.UseScanner(func(sc *scanner.Scanner, next func() (token.Token, error)) (tok token.Token, err error) {
		if tok, err = next(); err != nil {
			return
		}
		if tok.Type == token.IDENT && tok.Literal == "defer" {
			tok.Type = deferTyp
		}
		return
	})

	// now the parser can parse the "defer" syntax
	pb := xjs.ParserBuilder()
	pb.UseStmtParser(func(p *parser.Parser, next func() (ast.Stmt, error)) (_ ast.Stmt, err error) {
		if p.CurrentToken.Type == deferTyp {
			node := &DeferStmt{}
			node.Layout.Defer = p.CurrentToken
			p.AdvanceToken()
			if node.Stmt, err = p.ParseStmt(); err != nil {
				return
			}
			return node, nil
		}
		return next()
	})

	s := sb.Build(input)
	p := pb.Build(s)
	return js.ParseProgram(p)
}
