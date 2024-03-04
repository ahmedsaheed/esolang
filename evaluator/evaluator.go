package evaluator

import (
	"esolang/lang-esolang/ast"
	"esolang/lang-esolang/object"
	"fmt"
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
		return evalProgram(node)
	
		// 3. If the node is a *ast.ExpressionStatement, recursively evaluate the expression
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	
		// Work with expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		if isError(right) {return right}
		return evalPrefixExpression(node.Operator, right)
	
	case *ast.InfixExpression:
		left := Eval(node.Left)
		if isError(left) {return left}
		right := Eval(node.Right)
		if isError(right) {return right}
		return evalInfixExpression(node.Operator, left, right)
	
	case *ast.BlockStatement:
		return evalBlockStatement(node)
	
	case *ast.IfExpression:
		return evalIfExpression(node)
	
	case *ast.ReturnStatement:
		value := Eval(node.ReturnValue)
		if isError(value) {return value}
		return &object.ReturnValue{Value: value}
	}

	return nil
}

func evalIfExpression(ifExpressionNode *ast.IfExpression) object.Object {
	condition := Eval(ifExpressionNode.Condition)
	if isError(condition) {return condition}
	if isTruthy(condition) {
		return Eval(ifExpressionNode.Consequence)
	} else if ifExpressionNode.Alternative != nil {
		return Eval(ifExpressionNode.Alternative)
	} else { return NULL }	
}

func isTruthy(condition object.Object) bool {
	switch condition {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func evalInfixExpression(operator string, leftOperand, rightOperand object.Object) object.Object {
	switch {
		case leftOperand.Type() == object.INTEGER_OBJ && rightOperand.Type() == object.INTEGER_OBJ:
			return evalIntegerInfixExpression(operator, leftOperand, rightOperand)
		case operator == "==":
			return nativeBoolToBooleanObject(leftOperand == rightOperand)
		case operator == "!=":
			return nativeBoolToBooleanObject(leftOperand != rightOperand)
		case leftOperand.Type() != rightOperand.Type():
			return newError("type mismatch: %s %s %s", leftOperand.Type(), operator, rightOperand.Type())
		default:
			return newError("unknown operator: %s %s %s", leftOperand.Type(), operator, rightOperand.Type())
	}
}

// evalIntegerInfixExpression evaluates is where the actual arithmetic operations for + , - , / and * performed
func evalIntegerInfixExpression(operator string, leftOperand, rightOperand object.Object) object.Object {
	leftValue := leftOperand.(*object.Integer).Value
	rightValue := rightOperand.(*object.Integer).Value
	
	switch operator {
		case "+":
			return &object.Integer{Value: leftValue + rightValue}
		case "-":
			return &object.Integer{Value: leftValue - rightValue}
		case "*":
			return &object.Integer{Value: leftValue * rightValue}
		case "/":
			return &object.Integer{Value: leftValue / rightValue}
		case "<": 
			return nativeBoolToBooleanObject(leftValue < rightValue)
		case ">":
			return nativeBoolToBooleanObject(leftValue > rightValue)
		case "==":
			return nativeBoolToBooleanObject(leftValue == rightValue)
		case "!=":
			return nativeBoolToBooleanObject(leftValue != rightValue)
		default:
			return newError("unknown operator: %s %s %s", leftOperand.Type(), operator, rightOperand.Type())
	}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type()) 
	}
}
// evalMinusPrefixOperatorExpression evaluates the right object and returns a new object with the value negated
func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}
	// if operand ins't integer - we escape
	if right.Type() != object.INTEGER_OBJ {
		return NULL
	}
	// retrieve & negate the value
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}		 
}
// nativeBoolToBooleanObject converts a native bool to a boolean object
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

		// if the statement is a return statement, return the value
		if evaluatedResults, ok := evaluatedResults.(*object.ReturnValue); ok {
			return evaluatedResults.Value
		}
	}
	return evaluatedResults
}


func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object
	for _, statement := range block.Statements{
		result = Eval(statement)
		if result != nil {
			resultType := result.Type()
			if resultType == object.RETURN_VALUE_OBJ || resultType == object.ERROR_OBJ {
			return result
			}
		}
	}
	return result
}

func evalProgram(program *ast.Program) object.Object {
	var evaluatedResult object.Object
	for _, statement := range program.Statements {
		evaluatedResult = Eval(statement)
		switch evaluatedResult := evaluatedResult.(type) {
			case *object.ReturnValue:
				return evaluatedResult.Value
			case *object.Error:
				return evaluatedResult	
			}
		// if returnValue, ok := evaluatedResult.(*object.ReturnValue); ok {
		// 	return returnValue.Value
		// }
	}
	return evaluatedResult
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}