package djs

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/printer"
)

// A printer tells the printer how to print specific AST nodes.
//
// In this case we are telling it how to FORMAT defer statements.
func formatter(pr *printer.Printer, node ast.Node, next func(*printer.Printer, ast.Node) error) error {
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
}
