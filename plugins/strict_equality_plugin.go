package plugins

import (
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/token"
)

// StrictEqualityPlugin transforms '==' to '===' and '!=' to '!=='
func StrictEqualityPlugin(pb *parser.Builder) {
	lb := pb.LexerBuilder
	lb.UseTokenInterceptor(func(l *lexer.Lexer, next func() token.Token) token.Token {
		ret := next()
		if ret.Type == token.EQ && ret.Literal == "==" {
			ret.Literal = "==="
		}
		if ret.Type == token.NOT_EQ && ret.Literal == "!=" {
			ret.Literal = "!=="
		}
		return ret
	})
}
