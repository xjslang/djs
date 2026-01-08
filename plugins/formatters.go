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

func TransformProgram(program *ast.Program) *ast.Program {
	result := &ast.Program{
		Statements: []ast.Statement{},
	}
	// replace DeferFunctionDeclaration with DeferFunctionFormatter
	for _, stmt := range program.Statements {
		switch v := stmt.(type) {
		case *DeferFunctionDeclaration:
			// replace DeferStatement with DeferStatementFormatter
			dfBody := &ast.BlockStatement{
				Statements: []ast.Statement{},
			}
			for _, bodyStmt := range v.Body.Statements {
				if ds, ok := bodyStmt.(*DeferStatement); ok {
					bodyStmt = &DeferStatementFormatter{
						DeferStatement: ds,
					}
				}
				dfBody.Statements = append(dfBody.Statements, bodyStmt)
			}
			v.Body = dfBody

			result.Statements = append(result.Statements, &DeferFunctionFormatter{
				DeferFunctionDeclaration: v,
			})
		case *LetStatement:
			if fe, ok := v.Value.(*DeferFunctionExpression); ok {
				// replace DeferStatement with DeferStatementFormatter
				feBody := &ast.BlockStatement{
					Statements: []ast.Statement{},
				}
				for _, bodyStmt := range fe.Body.Statements {
					if ds, ok := bodyStmt.(*DeferStatement); ok {
						bodyStmt = &DeferStatementFormatter{
							DeferStatement: ds,
						}
					}
					feBody.Statements = append(feBody.Statements, bodyStmt)
				}
				fe.Body = feBody

				result.Statements = append(result.Statements, &LetStatementFormatter{
					LetStatement: v,
				})
			} else {
				result.Statements = append(result.Statements, stmt)
			}
		default:
			result.Statements = append(result.Statements, stmt)
		}
	}
	return result
}
