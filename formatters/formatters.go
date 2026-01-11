package formatters

import (
	"github.com/xjslang/djs/plugins"
	"github.com/xjslang/xjs/ast"
)

// DeferFunctionFormatter formats defer function declarations back to DJS syntax
type DeferFunctionFormatter struct {
	*plugins.DeferFunctionDeclaration
}

func (dff *DeferFunctionFormatter) WriteTo(cw *ast.CodeWriter) {
	if dff.AsyncFn {
		cw.WriteString("async ")
	}
	dff.FunctionDeclaration.WriteTo(cw)
}

// DeferFunctionExpressionFormatter formats defer function expressions back to DJS syntax
type DeferFunctionExpressionFormatter struct {
	*plugins.DeferFunctionExpression
}

func (dfef *DeferFunctionExpressionFormatter) WriteTo(cw *ast.CodeWriter) {
	dfef.FunctionExpression.WriteTo(cw)
}

// DeferStatementFormatter formats defer statements back to DJS syntax
type DeferStatementFormatter struct {
	*plugins.DeferStatement
}

func (dsf *DeferStatementFormatter) WriteTo(cw *ast.CodeWriter) {
	cw.WriteString("defer ")
	if len(dsf.Body.Statements) == 1 {
		dsf.Body.Statements[0].WriteTo(cw)
	} else {
		dsf.Body.WriteTo(cw)
	}
}

// LetStatementFormatter formats let statements back to DJS syntax
type LetStatementFormatter struct {
	*plugins.LetStatement
}

func (lsf *LetStatementFormatter) WriteTo(cw *ast.CodeWriter) {
	lsf.LetStatement.WriteTo(cw)
}

// OrExpressionFormatter formats or expressions back to DJS syntax
type OrExpressionFormatter struct {
	*plugins.OrExpression
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

// Prepare wraps the AST nodes with formatters to regenerate DJS source code
func Prepare(node ast.Statement) ast.Statement {
	return prepareStatement(node)
}

func prepareStatements(stmts []ast.Statement) []ast.Statement {
	result := []ast.Statement{}
	for _, stmt := range stmts {
		result = append(result, prepareStatement(stmt))
	}
	return result
}

func prepareExpression(exp ast.Expression) ast.Expression {
	switch v := exp.(type) {
	// wrap custom expressions with a "formatter" and return it
	case *plugins.OrExpression:
		return &OrExpressionFormatter{OrExpression: v}
	// traverse the rest of the tree
	case *plugins.DeferFunctionExpression:
		v.Body.Statements = prepareStatements(v.Body.Statements)
	case *ast.BinaryExpression:
		v.Left = prepareExpression(v.Left)
		v.Right = prepareExpression(v.Right)
	case *ast.UnaryExpression:
		v.Right = prepareExpression(v.Right)
	case *ast.PostfixExpression:
		v.Left = prepareExpression(v.Left)
	case *ast.GroupedExpression:
		v.Expression = prepareExpression(v.Expression)
	case *ast.CallExpression:
		v.Function = prepareExpression(v.Function)
		args := []ast.Expression{}
		for _, arg := range v.Arguments {
			args = append(args, prepareExpression(arg))
		}
		v.Arguments = args
	case *ast.MemberExpression:
		v.Object = prepareExpression(v.Object)
		v.Property = prepareExpression(v.Property)
	case *ast.AssignmentExpression:
		v.Left = prepareExpression(v.Left)
		v.Value = prepareExpression(v.Value)
	case *ast.CompoundAssignmentExpression:
		v.Left = prepareExpression(v.Left)
		v.Value = prepareExpression(v.Value)
	case *ast.ArrayLiteral:
		elems := []ast.Expression{}
		for _, elem := range v.Elements {
			elems = append(elems, prepareExpression(elem))
		}
		v.Elements = elems
	case *ast.ObjectLiteral:
		props := []ast.ObjectProperty{}
		for _, prop := range v.Properties {
			prop.Key = prepareExpression(prop.Key)
			prop.Value = prepareExpression(prop.Value)
			props = append(props, prop)
		}
		v.Properties = props
	}
	return exp
}

func prepareStatement(node ast.Statement) ast.Statement {
	switch v := node.(type) {
	// wrap custom statements with a "formatter" and return it
	case *plugins.DeferStatement:
		return &DeferStatementFormatter{DeferStatement: v}
	// traverse the rest of the tree
	case *ast.Program:
		v.Statements = prepareStatements(v.Statements)
	case *plugins.DeferFunctionDeclaration:
		v.Body.Statements = prepareStatements(v.Body.Statements)
	case *plugins.ExpressionStatement:
		v.Expression = prepareExpression(v.Expression)
	case *plugins.LetStatement:
		v.Value = prepareExpression(v.Value)
	case *ast.IfStatement:
		v.Condition = prepareExpression(v.Condition)
	case *ast.WhileStatement:
		v.Condition = prepareExpression(v.Condition)
	case *ast.ReturnStatement:
		v.ReturnValue = prepareExpression(v.ReturnValue)
	case *ast.ForStatement:
		v.Init = prepareExpression(v.Init)
		v.Condition = prepareExpression(v.Condition)
		v.Update = prepareExpression(v.Update)
	}
	return node
}
