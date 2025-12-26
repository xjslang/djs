package builder

import (
	"github.com/xjslang/djs/plugins"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
)

func New(lb *lexer.Builder) *parser.Builder {
	return parser.NewBuilder(lb).
		WithSmartSemicolon(true).
		Install(plugins.DeferPlugin).
		Install(plugins.OrPlugin).
		Install(plugins.StrictEqualityPlugin).
		Install(plugins.NewPlugin).
		Install(plugins.ThrowPlugin)
}
