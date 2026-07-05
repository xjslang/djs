package djs

import (
	"github.com/xjslang/xjs"
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/printer"
)

func Format(node ast.Node, opts ...printer.Option) (string, error) {
	p := xjs.PrinterBuilder().
		UsePrinter(formatter).
		Build(opts...)
	p.Print(node)
	return p.Output()
}

// A printer tells the printer how to print specific AST nodes.
//
// In this case we are telling it how to FORMAT defer statements.
func formatter(p *printer.Printer, node ast.Node, next func(node ast.Node) error) error {
	switch v := node.(type) {
	case *DeferStmt:
		// LnPrint: ensure a newline is added before printing (equiv: p.EnsureLine(); p.Print(a))
		// SpPrint: ensure a space is added before printing (equiv: p.EnsureSpace(); p.Print(a))
		p.Line().Print(v.Layout.Defer)
		p.Space().Print(v.Stmt)
		return nil
	}
	return next(node)
}
