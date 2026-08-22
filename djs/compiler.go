package djs

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/xjslang/xjs"
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/printer"
)

func Compile(node ast.Node) (string, error) {
	p := xjs.PrinterBuilder().
		UsePrinter(compiler).
		Build(printer.Compact())
	p.Print(node)
	return p.Output()
}

// A printer tells the printer how to print specific AST nodes.
//
// In this case we are telling it how to COMPILE blocks and defer statements.
func compiler(p *printer.Printer, node ast.Node, next func(*printer.Printer, ast.Node) error) error {
	switch v := node.(type) {
	case *js.BlockStmt:
		// tells the printer how to print blocks that contains "defer" statements
		if isDeferBlock(v) {
			defersName := fmt.Sprintf("__defers_%s__", rndID())

			// share a context with its child nodes
			ctx := p.PushContext()
			defer p.PopContext()
			ctx["defersName"] = defersName

			p.Line().Print(fmt.Sprintf("{ let %s = []; try {", defersName))
			if err := js.PrintBlockStmt(p, v); err != nil {
				return err
			}
			p.Print(fmt.Sprintf("} finally {"+
				"for (let defer of %s) { "+
				"try { defer() } catch (e) { console.error(e) }"+
				"}}}", defersName))
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
			Print(fmt.Sprintf("%s.unshift(() => {", defersName), v.Stmt, "});")
		return nil
	}
	return next(p, node) // delegate in the next printer
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
