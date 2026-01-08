package plugins

import (
	"github.com/xjslang/xjs/ast"
)

type DeferFunctionFormatter struct {
	*DeferFunctionDeclaration
}

func (dff *DeferFunctionFormatter) WriteTo(cw *ast.CodeWriter) {
	if dff.asyncFn {
		cw.WriteString("async ")
	}
	dff.FunctionDeclaration.WriteTo(cw)
}

type DeferFunctionExpressionFormatter struct {
	*DeferFunctionExpression
}

func (dfef *DeferFunctionExpressionFormatter) WriteTo(cw *ast.CodeWriter) {
	dfef.FunctionExpression.WriteTo(cw)
}

type DeferStatementFormatter struct {
	*DeferStatement
}

func (dsf *DeferStatementFormatter) WriteTo(cw *ast.CodeWriter) {
	cw.WriteString("defer ")
	if len(dsf.Body.Statements) == 1 {
		dsf.Body.Statements[0].WriteTo(cw)
	} else {
		dsf.Body.WriteTo(cw)
	}
}

type LetStatementFormatter struct {
	*LetStatement
}

func (lsf *LetStatementFormatter) WriteTo(cw *ast.CodeWriter) {
	lsf.LetStatement.WriteTo(cw)
}
