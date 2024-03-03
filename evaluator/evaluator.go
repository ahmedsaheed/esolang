package evaluator

import (
	"esolang/lang-esolang/ast"
	"esolang/lang-esolang/object"
)

var (
	// TRUE and FALSE are the only instances of the Boolean object - add lil' optimization instead of creating new instances
	TRUE = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL = &object.Null{}
)

/* Eval evaluates an AST node and returns a type of `object.Object` from the object system.
The eval typically receives an *ast.Program and  recursively traverses whilst evaluating each AST node
*/
func Eval(node ast.Node) object.Object {
	// 1. Check the type of the node
	switch node := node.(type) {
	// 2. If the node is a *ast.Program, evaluate the statements
	case *ast.Program:
		return evalStatements(node.Statements)
	// 3. If the node is a *ast.ExpressionStatement, recursively evaluate the expression
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	// Work with expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	}
	return nil
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalStatements(statements []ast.Statement) object.Object {
	var evaluatedResults object.Object
	for _, statement := range statements {
		evaluatedResults = Eval(statement)
	}
	return evaluatedResults
}