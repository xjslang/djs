package djs

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/plugin"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

// A plugin enhances the capabilities of the parser, so that it can "understand"
// constructs that are not necessarily part of the JavaScript language.
//
// In this case we are telling the parser how to parse "defer" statements.
func djsPlugin(b *plugin.Builder) {
	deferTyp := token.RegisterType("defer")

	// now the scanner can scan the "defer" keyword
	b.UseScanner(func(sc *scanner.Scanner, next func() (token.Token, error)) (tok token.Token, err error) {
		if tok, err = next(); err != nil {
			return
		}
		if tok.Type == token.IDENT && tok.Literal == "defer" {
			tok.Type = deferTyp
		}
		return
	})

	// now the parser can parse the "defer" syntax
	b.UseStmtParser(func(p *parser.Parser, next func() (ast.Stmt, error)) (_ ast.Stmt, err error) {
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
}
