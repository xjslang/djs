package builder

import (
	"github.com/xjslang/djs/plugins"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
)

func New(lb *lexer.Builder) *parser.Builder {
	return parser.NewBuilder(lb).
		WithSmartSemicolon(true).
		Install(plugins.ThrowPlugin).
		Install(plugins.NewPlugin).
		Install(plugins.StrictEqualityPlugin).
		Install(plugins.OrPlugin).
		Install(plugins.DeferPlugin)
}
