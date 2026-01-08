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

type OrFormatter struct {
	*OrExpression
}

func (of *OrFormatter) WriteTo(cw *ast.CodeWriter) {
	of.Expression.WriteTo(cw)
	cw.WriteString(" or ")
	if of.ErrorParam != nil {
		cw.WriteRune('|')
		of.ErrorParam.WriteTo(cw)
		cw.WriteRune('|')
		cw.WriteSpace()
	}
	of.FallbackBlock.WriteTo(cw)
}

func Transform(stmts []ast.Statement) []ast.Statement {
	result := []ast.Statement{}
	for _, stmt := range stmts {
		switch v := stmt.(type) {
		// defer
		case *DeferStatement:
			result = append(result, &DeferStatementFormatter{
				DeferStatement: v,
			})
		case *DeferFunctionDeclaration:
			v.Body.Statements = Transform(v.Body.Statements)
			result = append(result, &DeferFunctionFormatter{
				DeferFunctionDeclaration: v,
			})
		case *DeferFunctionExpression:
			v.Body.Statements = Transform(v.Body.Statements)
			result = append(result, &DeferFunctionExpressionFormatter{
				DeferFunctionExpression: v,
			})
		// or
		case *ExpressionStatement:
			if w, ok := v.Expression.(*OrExpression); ok {
				v.Expression = &OrFormatter{
					OrExpression: w,
				}
			}
			result = append(result, v)
		case *LetStatement:
			switch w := v.Value.(type) {
			case *DeferFunctionExpression:
				w.Body.Statements = Transform(w.Body.Statements)
			case *OrExpression:
				v.Value = &OrFormatter{
					OrExpression: w,
				}
			}
			result = append(result, &LetStatementFormatter{
				LetStatement: v,
			})
		default:
			result = append(result, v)
		}
	}
	return result
}
