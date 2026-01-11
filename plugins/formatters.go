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

type OrExpressionFormatter struct {
	*OrExpression
}

func (of *OrExpressionFormatter) WriteTo(cw *ast.CodeWriter) {
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

func TransformStatements(stmts []ast.Statement) []ast.Statement {
	result := []ast.Statement{}
	for _, stmt := range stmts {
		result = append(result, TransformStatement(stmt))
	}
	return result
}

func TransformExpression(exp ast.Expression) ast.Expression {
	switch v := exp.(type) {
	// wrap custom expressions with a "formatter" and return it
	case *OrExpression:
		return &OrExpressionFormatter{OrExpression: v}
	// traverse the rest of the tree
	case *DeferFunctionExpression:
		v.Body.Statements = TransformStatements(v.Body.Statements)
	case *ast.BinaryExpression:
		v.Left = TransformExpression(v.Left)
		v.Right = TransformExpression(v.Right)
	case *ast.UnaryExpression:
		v.Right = TransformExpression(v.Right)
	case *ast.PostfixExpression:
		v.Left = TransformExpression(v.Left)
	case *ast.GroupedExpression:
		v.Expression = TransformExpression(v.Expression)
	case *ast.CallExpression:
		v.Function = TransformExpression(v.Function)
		args := []ast.Expression{}
		for _, arg := range v.Arguments {
			args = append(args, TransformExpression(arg))
		}
		v.Arguments = args
	case *ast.MemberExpression:
		v.Object = TransformExpression(v.Object)
		v.Property = TransformExpression(v.Property)
	case *ast.AssignmentExpression:
		v.Left = TransformExpression(v.Left)
		v.Value = TransformExpression(v.Value)
	case *ast.CompoundAssignmentExpression:
		v.Left = TransformExpression(v.Left)
		v.Value = TransformExpression(v.Value)
	case *ast.ArrayLiteral:
		elems := []ast.Expression{}
		for _, elem := range v.Elements {
			elems = append(elems, TransformExpression(elem))
		}
		v.Elements = elems
	case *ast.ObjectLiteral:
		props := []ast.ObjectProperty{}
		for _, prop := range v.Properties {
			prop.Key = TransformExpression(prop.Key)
			prop.Value = TransformExpression(prop.Value)
			props = append(props, prop)
		}
		v.Properties = props
	}
	return exp
}

func TransformStatement(node ast.Statement) ast.Statement {
	switch v := node.(type) {
	// wrap custom statements with a "formatter" and return it
	case *DeferStatement:
		return &DeferStatementFormatter{DeferStatement: v}
	// traverse the rest of the tree
	case *ast.Program:
		v.Statements = TransformStatements(v.Statements)
	case *DeferFunctionDeclaration:
		v.Body.Statements = TransformStatements(v.Body.Statements)
	case *ExpressionStatement:
		v.Expression = TransformExpression(v.Expression)
	case *LetStatement:
		v.Value = TransformExpression(v.Value)
	case *ast.IfStatement:
		v.Condition = TransformExpression(v.Condition)
	case *ast.WhileStatement:
		v.Condition = TransformExpression(v.Condition)
	case *ast.ReturnStatement:
		v.ReturnValue = TransformExpression(v.ReturnValue)
	case *ast.ForStatement:
		v.Init = TransformExpression(v.Init)
		v.Condition = TransformExpression(v.Condition)
		v.Update = TransformExpression(v.Update)
	}
	return node
}
