package djs

import (
	"github.com/xjslang/xjs/ast"
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
