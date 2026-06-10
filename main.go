package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/xjslang/xjs"
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/builder"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

type DeferStmt struct {
	ast.BaseStmt
	Layout struct {
		Defer token.Token
	}
	Stmt ast.Stmt
}

func djsPlugin(b *builder.Builder) {
	deferTyp := token.RegisterType("defer")
	b.UseScanner(func(sc *scanner.Scanner, next func() token.Token) token.Token {
		tok := next()
		if tok.Type == token.IDENT && tok.Literal == "defer" {
			tok.Type = deferTyp
		}
		return tok
	})
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

func djsCompiler(p *printer.Printer, node ast.Node, next func(node ast.Node)) {
	defersName := fmt.Sprintf("__defers_%s__", rndID())
	switch v := node.(type) {
	case *js.BlockStmt:
		if isDeferBlock(v) {
			p.SpPrint(fmt.Sprintf("{ let %s = []; try { ", defersName))
			js.PrintBlockStmt(p, v)
			p.SpPrint(fmt.Sprintf("} finally {"+
				"for (let defer of %s) { "+
				"try { defer() } catch (e) { console.error(e) }"+
				"}}}", defersName))
			return
		}
	case *DeferStmt:
		p.LnPrint(fmt.Sprintf("%s.unshift(() =>", defersName))
		p.SpPrint(v.Stmt)
		p.Print(")")
		return
	}
	next(node)
}

func isDeferBlock(node *js.BlockStmt) bool {
	for _, stmt := range node.Stmts {
		if _, ok := stmt.(*DeferStmt); ok {
			return true
		}
	}
	return false
}

func rndID() string {
	b := make([]byte, 4)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func main() {
	input := `function foo() {
		defer console.log('close db')
		defer {console.log('close file')}
		let x = 100
	}`
	b := xjs.NewBuilder().Install(djsPlugin)
	p := b.Build([]byte(input))
	result, err := js.ParseProgram(p)
	if err != nil {
		panic(err)
	}

	pr := xjs.NewPrinter()
	pr.UsePrinter(djsCompiler)
	pr.Print(result)
	fmt.Println(pr.String())
}
