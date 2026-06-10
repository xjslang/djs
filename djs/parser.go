package djs

import (
	"github.com/xjslang/xjs"
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
)

func Parse(input []byte) (ast.Node, error) {
	b := xjs.NewBuilder().Install(djsPlugin)
	p := b.Build(input)
	return js.ParseProgram(p)
}
