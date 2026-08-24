package djs

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/xjslang/xjs"
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

// DeferStmt represents a "defer" statement
type DeferStmt struct {
	ast.BaseStmt
	Layout struct {
		Defer token.Token
	}
	Stmt ast.Stmt
}

func Parse(input string) (*js.Program, error) {
	deferTyp := token.RegisterType("DEFER", "defer")

	sb := xjs.ScannerBuilder()
	// now the scanner can scan the "defer" keyword
	sb.UseScanner(func(sc *scanner.Scanner, next func(*scanner.Scanner) (token.Token, error)) (tok token.Token, err error) {
		if tok, err = next(sc); err != nil {
			return
		}
		if tok.Type == token.IDENT && tok.Literal == "defer" {
			tok.Type = deferTyp
		}
		return
	})

	// now the parser can parse the "defer" syntax
	pb := xjs.ParserBuilder()
	pb.UseStmtParser(func(p *parser.Parser, next func(*parser.Parser) (ast.Stmt, error)) (_ ast.Stmt, err error) {
		if p.CurrentToken.Type == deferTyp {
			node := &DeferStmt{}
			node.Layout.Defer = p.CurrentToken
			p.AdvanceToken()
			if node.Stmt, err = p.ParseStmt(); err != nil {
				return
			}
			return node, nil
		}
		return next(p)
	})

	s := sb.Build(input)
	p := pb.Build(s)
	return js.ParseProgram(p)
}

func Compile(node ast.Node) (string, error) {
	p := xjs.PrinterBuilder().
		UsePrinter(func(p *printer.Printer, node ast.Node, next func(*printer.Printer, ast.Node) error) error {
			switch v := node.(type) {
			case *js.BlockStmt:
				// tells the printer how to print blocks that contains "defer" statements
				if isDeferBlock(v) {
					defersName := "__defers_" + rndID() + "__"

					// share a context with its child nodes
					ctx := p.PushContext()
					defer p.PopContext()
					ctx["defersName"] = defersName

					p.Line().Print("{ let ", defersName, " = []; try {")
					if err := js.PrintBlockStmt(p, v); err != nil {
						return err
					}
					p.Print("} finally {",
						"for (let defer of ", defersName, ") { ",
						"try { defer() } catch (e) { console.error(e) }",
						"}}}")
					return nil
				}
			case *DeferStmt:
				// tells the printer how to print "defer" statements
				ctx := p.Context()
				defersName, ok := ctx["defersName"]
				if !ok {
					// return printer.ErrorAt(v.Layout.Defer.Position, "defer cannot be used outside a block")
					return p.Error("defer cannot be used outside a block")
				}

				p.
					Line(). // ensure a new line is added before printing
					Print(defersName, ".unshift(() => {", v.Stmt, "});")
				return nil
			}
			return next(p, node) // delegate in the next printer
		}).
		Build(printer.Compact())
	p.Print(node)
	return p.Output()
}

func Format(node ast.Node, opts ...printer.Option) (string, error) {
	p := xjs.PrinterBuilder().
		UsePrinter(func(pr *printer.Printer, node ast.Node, next func(*printer.Printer, ast.Node) error) error {
			switch v := node.(type) {
			case *DeferStmt:
				pr.
					Line(). // ensure a new line is added before printing
					Print(v.Layout.Defer)
				pr.
					Space(). // ensure a new space is added before printing
					Print(v.Stmt)
				return nil
			}
			return next(pr, node)
		}).
		Build(opts...)
	p.Print(node)
	return p.Output()
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
	rand.Read(b) //nolint
	return hex.EncodeToString(b)
}
