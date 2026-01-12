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
	if exp == nil {
		return nil
	}

	switch v := exp.(type) {
	// updates nodes for appropriate formatting
	case *plugins.OrExpression:
		v.Expression = prepareExpression(v.Expression)
		return &OrExpressionFormatter{OrExpression: v}
	case *ast.BinaryExpression:
		switch v.Operator {
		case "===":
			v.Operator = "=="
		case "!==":
			v.Operator = "!="
		}
		v.Left = prepareExpression(v.Left)
		v.Right = prepareExpression(v.Right)
		return v
	// traverse the rest of the tree
	case *plugins.DeferFunctionExpression:
		v.Body.Statements = prepareStatements(v.Body.Statements)
		return v
	case *ast.FunctionExpression:
		v.Body.Statements = prepareStatements(v.Body.Statements)
		return v
	case *ast.UnaryExpression:
		v.Right = prepareExpression(v.Right)
		return v
	case *ast.PostfixExpression:
		v.Left = prepareExpression(v.Left)
		return v
	case *ast.GroupedExpression:
		v.Expression = prepareExpression(v.Expression)
		return v
	case *ast.CallExpression:
		v.Function = prepareExpression(v.Function)
		args := []ast.Expression{}
		for _, arg := range v.Arguments {
			args = append(args, prepareExpression(arg))
		}
		v.Arguments = args
		return v
	case *ast.MemberExpression:
		v.Object = prepareExpression(v.Object)
		if !v.Computed {
			// Property is an identifier, don't traverse
			return v
		}
		v.Property = prepareExpression(v.Property)
		return v
	case *ast.AssignmentExpression:
		v.Left = prepareExpression(v.Left)
		v.Value = prepareExpression(v.Value)
		return v
	case *ast.CompoundAssignmentExpression:
		v.Left = prepareExpression(v.Left)
		v.Value = prepareExpression(v.Value)
		return v
	case *ast.ArrayLiteral:
		elems := []ast.Expression{}
		for _, elem := range v.Elements {
			elems = append(elems, prepareExpression(elem))
		}
		v.Elements = elems
		return v
	case *ast.ObjectLiteral:
		props := []ast.ObjectProperty{}
		for _, prop := range v.Properties {
			prop.Key = prepareExpression(prop.Key)
			prop.Value = prepareExpression(prop.Value)
			props = append(props, prop)
		}
		v.Properties = props
		return v
	}
	return exp
}

func prepareStatement(node ast.Statement) ast.Statement {
	if node == nil {
		return nil
	}

	switch v := node.(type) {
	// wrap custom statements with a "formatter" and return it
	case *plugins.DeferStatement:
		v.Body.Statements = prepareStatements(v.Body.Statements)
		return &DeferStatementFormatter{DeferStatement: v}
	// traverse the rest of the tree
	case *ast.Program:
		v.Statements = prepareStatements(v.Statements)
		return v
	case *plugins.DeferFunctionDeclaration:
		v.Body.Statements = prepareStatements(v.Body.Statements)
		return v
	case *ast.FunctionDeclaration:
		v.Body.Statements = prepareStatements(v.Body.Statements)
		return v
	case *plugins.ExpressionStatement:
		v.Expression = prepareExpression(v.Expression)
		return v
	case *plugins.LetStatement:
		v.Value = prepareExpression(v.Value)
		return v
	case *ast.IfStatement:
		v.Condition = prepareExpression(v.Condition)
		v.ThenBranch = prepareStatement(v.ThenBranch)
		v.ElseBranch = prepareStatement(v.ElseBranch)
		return v
	case *ast.WhileStatement:
		v.Condition = prepareExpression(v.Condition)
		v.Body = prepareStatement(v.Body)
		return v
	case *ast.ReturnStatement:
		v.ReturnValue = prepareExpression(v.ReturnValue)
		return v
	case *ast.ForStatement:
		v.Init = prepareExpression(v.Init)
		v.Condition = prepareExpression(v.Condition)
		v.Update = prepareExpression(v.Update)
		v.Body = prepareStatement(v.Body)
		return v
	case *ast.BlockStatement:
		v.Statements = prepareStatements(v.Statements)
		return v
	}
	return node
}
