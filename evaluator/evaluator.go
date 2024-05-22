package evaluator

import (
	"esolang/lang-esolang/ast"
	"esolang/lang-esolang/builtins"
	"esolang/lang-esolang/lexer"
	"esolang/lang-esolang/object"
	"esolang/lang-esolang/parser"
	"esolang/lang-esolang/utils"
	"fmt"
	"math"
	"os"
	"strings"
)

var (
	// TRUE and FALSE are the only instances of the Boolean object - add lil' optimization instead of creating new instances
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

/*
	Eval evaluates an AST node and returns a type of `object.Object` from the object system.

The eval typically receives an *ast.Program and  recursively traverses whilst evaluating each AST node
*/
func Eval(node ast.Node, env *object.Environment) object.Object {
	// 1. Check the type of the node
	switch node := node.(type) {

	// 2. If the node is a *ast.Program, evaluate the statements
	case *ast.Program:
		return evalProgram(node, env)

		// 3. If the node is a *ast.ExpressionStatement, recursively evaluate the expression
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

		// Work with expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node, node.Operator, right)
	case *ast.PostfixExpression:
		return evalPostfixExpression(node, env, node.Operator)

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, node, left, right)
	case *ast.WhileLoopExpression:
		return evalWhileLoopExpression(node, env)
	case *ast.ImportExpression:
		return evalImportExpression(node, env)
	case *ast.ObjectCallExpression:
		return evalObjectCallExpression(node, env)
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ast.BindExpression:
		value := Eval(node.Value, env)
		if isError(value) {
			return value
		}

		if ident, ok := node.Left.(*ast.Identifier); ok {
			if obj, ok := value.(object.Copyable); ok {
				env.Set(ident.Value, obj.Copy())
			} else {
				env.Set(ident.Value, value)
			}

			return &object.Null{}
		}
		return newError(node.Token.FileName, node.Token.Line, node.Token.Column, "expected identifier on left got=%T", node.Left)

	case *ast.AssignStatement:
		return evalAssignStatement(node, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Body: body, Env: env}
	case *ast.FunctionDefineLiteral:
		params := node.Parameters
		body := node.Body
		// defaults := node.Defaults
		env.Set(node.TokenLiteral(), &object.Function{Parameters: params, Env: env, Body: body})
		return NULL
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(node, left, index)

	case *ast.HashLiteral:
		return evalHashLiteral(node, env)

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(node, function, args)

	case *ast.ReturnStatement:
		value := Eval(node.ReturnValue, env)
		if isError(value) {
			return value
		}
		return &object.ReturnValue{Value: value}
	}

	return nil
}

func evalIndexExpression(node *ast.IndexExpression, left object.Object, index object.Object) object.Object {
	currLine := node.Token.Line
	currCol := node.Token.Column
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(node, left, index)
	case left.Type() == object.MODULE_TYPE:
		return evalModuleIndexExpression(node, left, index)
	case left.Type() == object.STRING_OBJ:
		return evalStringIndexExpression(left, index)
	default:

		if left.Type() == object.ARRAY_OBJ {
			return newError(node.Token.FileName, currLine, currCol, `index operation on %s only uses ["index"] accessor`, left.Type())
		}

		return newError(node.Token.FileName, currLine, currCol, "index operator not supported: %s", left.Type())
	}

}

func evalArrayIndexExpression(array object.Object, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	maxIdx := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > maxIdx {
		return NULL
	}

	return arrayObject.Elements[idx]
}

func evalStringIndexExpression(str object.Object, index object.Object) object.Object {
	strObj := str.(*object.String)
	idx := index.(*object.Integer).Value

	if idx < 0 || idx >= int64(len(strObj.Value)) {
		return NULL
	}
	return &object.String{Value: string(strObj.Value[idx])}
}

func evalImportExpression(ie *ast.ImportExpression, env *object.Environment) object.Object {
	name := Eval(ie.Name, env)
	if isError(name) {
		return name
	}

	if s, ok := name.(*object.String); ok {
		attrs := Module(ie, s.Value)
		if isError(attrs) {
			return attrs
		}
		return &object.Module{Name: s.Value, Attrs: attrs}
	}
	return newError(ie.Token.FileName, ie.Token.Line, ie.Token.Column, "ImportError: invalid import path '%s'", name)
}

func evalWhileLoopExpression(flExpression *ast.WhileLoopExpression, env *object.Environment) object.Object {
	var result object.Object

	for {
		condition := Eval(flExpression.Condition, env)
		if isError(condition) {
			return condition
		}

		if isTruthy(condition) {
			result = Eval(flExpression.Consequence, env)
		} else {
			break
		}
	}

	if result != nil {
		return result
	}
	return &object.Null{}
}

func evalHashIndexExpression(node *ast.IndexExpression, hash object.Object, index object.Object) object.Object {
	hashObj := hash.(*object.Hash)
	key, ok := index.(object.Hashable)
	if !ok {
		return newError(node.Token.FileName, node.Token.Line, node.Token.Column, "unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObj.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}

func evalModuleIndexExpression(node *ast.IndexExpression, module, index object.Object) object.Object {
	moduleObject := module.(*object.Module)
	return evalHashIndexExpression(node, moduleObject.Attrs, index)
}

func evalExpressions(expression []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object
	for _, expr := range expression {
		evaluated := Eval(expr, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	if builtin, ok := builtins.Builtins[node.Value]; ok {
		return builtin
	}
	return newError(node.Token.FileName, node.Token.Line, node.Token.Column, "cannot find '%s' in scope", node.Value)
}

func evalIfExpression(ifExpressionNode *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ifExpressionNode.Condition, env)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(ifExpressionNode.Consequence, env)
	} else if ifExpressionNode.Alternative != nil {
		return Eval(ifExpressionNode.Alternative, env)
	} else {
		return NULL
	}
}

func evalInfixExpression(operator string, node ast.Expression, leftOperand, rightOperand object.Object) object.Object {

	// check node type is ast.InfixExpression

	switch node := node.(type) {
	case *ast.InfixExpression:
		switch {
		case leftOperand.Type() == object.INTEGER_OBJ && rightOperand.Type() == object.INTEGER_OBJ:
			return evalIntegerInfixExpression(node, operator, leftOperand, rightOperand)
		case leftOperand.Type() == object.FLOAT_OBJ && rightOperand.Type() == object.FLOAT_OBJ:
			return evalFloatInfixExpression(node, operator, leftOperand, rightOperand)
		case leftOperand.Type() == object.FLOAT_OBJ && rightOperand.Type() == object.INTEGER_OBJ:
			return evalFloatIntegerInfixExpression(node, operator, leftOperand, rightOperand)
		case leftOperand.Type() == object.INTEGER_OBJ && rightOperand.Type() == object.FLOAT_OBJ:
			return evalIntegerFloatInfixExpression(node, operator, leftOperand, rightOperand)
		case leftOperand.Type() == object.STRING_OBJ && (rightOperand.Type() == object.STRING_OBJ || rightOperand.Type() == object.INTEGER_OBJ):
			return evalStringInfixExpression(node, leftOperand, rightOperand)
		case leftOperand.Type() == object.ARRAY_OBJ && rightOperand.Type() == object.ARRAY_OBJ && operator == "+":
			return &object.Array{Elements: append(leftOperand.(*object.Array).Elements, rightOperand.(*object.Array).Elements...)}
		case operator == "==":
			return nativeBoolToBooleanObject(leftOperand == rightOperand)
		case operator == "!=":
			return nativeBoolToBooleanObject(leftOperand != rightOperand)
		case operator == "&&":
			return nativeBoolToBooleanObject(objectToNativeBoolean(leftOperand) && objectToNativeBoolean(rightOperand))
		case operator == "-|":
			return nativeBoolToBooleanObject(objectToNativeBoolean(leftOperand) || objectToNativeBoolean(rightOperand))
		case leftOperand.Type() != rightOperand.Type():
			return newError(node.Token.FileName, node.Token.Line, node.Token.Column, "type mismatch: %s %s %s", leftOperand.Type(), operator, rightOperand.Type())
		default:
			return newError(node.Token.FileName, node.Token.Line, node.Token.Column, "unknown operator: %s %s %s", leftOperand.Type(), operator, rightOperand.Type())
		}

	case *ast.AssignStatement:
		switch {
		case leftOperand.Type() == object.INTEGER_OBJ && rightOperand.Type() == object.INTEGER_OBJ:
			return evalIntegerInfixExpression(node, operator, leftOperand, rightOperand)
		case leftOperand.Type() == object.FLOAT_OBJ && rightOperand.Type() == object.FLOAT_OBJ:
			return evalFloatInfixExpression(node, operator, leftOperand, rightOperand)
		case leftOperand.Type() == object.FLOAT_OBJ && rightOperand.Type() == object.INTEGER_OBJ:
			return evalFloatIntegerInfixExpression(node, operator, leftOperand, rightOperand)
		case leftOperand.Type() == object.INTEGER_OBJ && rightOperand.Type() == object.FLOAT_OBJ:
			return evalIntegerFloatInfixExpression(node, operator, leftOperand, rightOperand)
		case leftOperand.Type() == object.STRING_OBJ && (rightOperand.Type() == object.STRING_OBJ || rightOperand.Type() == object.INTEGER_OBJ):
			return evalStringInfixExpression(node, leftOperand, rightOperand)
		case leftOperand.Type() == object.ARRAY_OBJ && rightOperand.Type() == object.ARRAY_OBJ && operator == "+":
			return &object.Array{Elements: append(leftOperand.(*object.Array).Elements, rightOperand.(*object.Array).Elements...)}
		case operator == "==":
			return nativeBoolToBooleanObject(leftOperand == rightOperand)
		case operator == "!=":
			return nativeBoolToBooleanObject(leftOperand != rightOperand)
		case operator == "&&":
			return nativeBoolToBooleanObject(objectToNativeBoolean(leftOperand) && objectToNativeBoolean(rightOperand))
		case operator == "-|":
			return nativeBoolToBooleanObject(objectToNativeBoolean(leftOperand) || objectToNativeBoolean(rightOperand))
		case leftOperand.Type() != rightOperand.Type():
			return newError(node.Token.FileName, node.Token.Line, node.Token.Column, "type mismatch: %s %s %s", leftOperand.Type(), operator, rightOperand.Type())
		default:
			return newError(node.Token.FileName, node.Token.Line, node.Token.Column, "unknown operator: %s %s %s", leftOperand.Type(), operator, rightOperand.Type())
		}
	}

	return NULL
}

// evalStringInfixExpression evaluates the string concatenation
func evalStringInfixExpression(node ast.Expression, leftOperand, rightOperand object.Object) object.Object {

	switch node := node.(type) {

	case *ast.InfixExpression:
		operator := node.Operator
		currLine := node.Token.Line
		currColumn := node.Token.Column
		switch rightOperand.Type() {
		case object.INTEGER_OBJ:
			leftValue := leftOperand.(*object.String).Value
			rightIntValue := rightOperand.(*object.Integer).Value
			if operator == "*" {
				return &object.String{Value: strings.Repeat(leftValue, int(rightIntValue))}
			}
		case object.STRING_OBJ:
			leftValue := leftOperand.(*object.String).Value
			rightValue := rightOperand.(*object.String).Value
			switch operator {
			case "+":
				return &object.String{Value: leftValue + rightValue}
			case "==":
				return nativeBoolToBooleanObject(leftValue == rightValue)
			case "!=":
				return nativeBoolToBooleanObject(leftValue != rightValue)
			default:
				return newError(node.Token.FileName, currLine, currColumn, "unknown operator: %s %s %s", leftOperand.Type(), operator, rightOperand.Type())
			}
		default:
			return newError(node.Token.FileName, currLine, currColumn, "unknown operator: %s %s %s", leftOperand.Type(), operator, rightOperand.Type())
		}
	case *ast.AssignStatement:
		operator := node.Operator
		currLine := node.Token.Line
		currColumn := node.Token.Column
		switch rightOperand.Type() {
		case object.INTEGER_OBJ:
			leftValue := leftOperand.(*object.String).Value
			rightIntValue := rightOperand.(*object.Integer).Value
			if operator == "*" {
				return &object.String{Value: strings.Repeat(leftValue, int(rightIntValue))}
			}
		case object.STRING_OBJ:
			leftValue := leftOperand.(*object.String).Value
			rightValue := rightOperand.(*object.String).Value
			switch operator {
			case "+":
				return &object.String{Value: leftValue + rightValue}
			case "==":
				return nativeBoolToBooleanObject(leftValue == rightValue)
			case "!=":
				return nativeBoolToBooleanObject(leftValue != rightValue)
			default:
				return newError(node.Token.FileName, currLine, currColumn, "unknown operator: %s %s %s", leftOperand.Type(), operator, rightOperand.Type())
			}
		default:
			return newError(node.Token.FileName, currLine, currColumn, "unknown operator: %s %s %s", leftOperand.Type(), operator, rightOperand.Type())
		}

	}

	return NULL
}

// evalIntegerInfixExpression evaluates is where the actual arithmetic operations for + , - , / and * performed
func evalIntegerInfixExpression(node ast.Expression, operator string, leftOperand, rightOperand object.Object) object.Object {
	leftValue := leftOperand.(*object.Integer).Value
	rightValue := rightOperand.(*object.Integer).Value

	switch node := node.(type) {
	case *ast.InfixExpression:
		switch operator {
		case "+":
			return &object.Integer{Value: leftValue + rightValue}
		case "-":
			return &object.Integer{Value: leftValue - rightValue}
		case "*":
			return &object.Integer{Value: leftValue * rightValue}
		case "/":
			if rightValue == 0 {
				return newError(node.Token.FileName, node.Token.Line, node.Token.Column, "Can't divide by zero")
			}
			return &object.Integer{Value: leftValue / rightValue}
		case "%":
			if rightValue == 0 {
				return newError(node.Token.FileName, node.Token.Line, node.Token.Column, "Can't divide by zero")
			}
			return &object.Integer{Value: leftValue % rightValue}
		case "<":
			return nativeBoolToBooleanObject(leftValue < rightValue)
		case ">":
			return nativeBoolToBooleanObject(leftValue > rightValue)
		case "==":
			return nativeBoolToBooleanObject(leftValue == rightValue)
		case "!=":
			return nativeBoolToBooleanObject(leftValue != rightValue)
		case "-=":
			return &object.Integer{Value: leftValue - rightValue}
		case "*=":
			return &object.Integer{Value: leftValue * rightValue}
		case "+=":
			return &object.Integer{Value: leftValue + rightValue}
		case "<=":
			return nativeBoolToBooleanObject(leftValue <= rightValue)
		case ">=":
			return nativeBoolToBooleanObject(leftValue >= rightValue)

		default:
			return newError(node.Token.FileName, node.Token.Line, node.Token.Column, "unknown operator: %s %s %s", leftOperand.Type(), operator, rightOperand.Type())
		}
	case *ast.AssignStatement:
		switch operator {
		case "+":
			return &object.Integer{Value: leftValue + rightValue}
		case "-":
			return &object.Integer{Value: leftValue - rightValue}
		case "*":
			return &object.Integer{Value: leftValue * rightValue}
		case "/":
			if rightValue == 0 {
				return newError(node.Token.FileName, node.Token.Line, node.Token.Column, "Can't divide by zero")
			}
			return &object.Integer{Value: leftValue / rightValue}
		case "%":
			if rightValue == 0 {
				return newError(node.Token.FileName, node.Token.Line, node.Token.Column, "Can't divide by zero")
			}
			return &object.Integer{Value: leftValue % rightValue}
		case "<":
			return nativeBoolToBooleanObject(leftValue < rightValue)
		case ">":
			return nativeBoolToBooleanObject(leftValue > rightValue)
		case "==":
			return nativeBoolToBooleanObject(leftValue == rightValue)
		case "!=":
			return nativeBoolToBooleanObject(leftValue != rightValue)
		case "-=":
			return &object.Integer{Value: leftValue - rightValue}
		case "*=":
			return &object.Integer{Value: leftValue * rightValue}
		case "**":
			return &object.Integer{Value: int64(math.Pow(float64(leftValue), float64(rightValue)))}
		case "+=":
			return &object.Integer{Value: leftValue + rightValue}
		case "<=":
			return nativeBoolToBooleanObject(leftValue <= rightValue)
		case ">=":
			return nativeBoolToBooleanObject(leftValue >= rightValue)

		default:
			return newError(node.Token.FileName, node.Token.Line, node.Token.Column, "unknown operator: %s %s %s", leftOperand.Type(), operator, rightOperand.Type())
		}
	}
	return NULL
}

func evalFloatInfixExpression(node ast.Expression, operator string, leftOperand, rightOperand object.Object) object.Object {
	leftValue := leftOperand.(*object.Float).Value
	rightValue := rightOperand.(*object.Float).Value

	switch node := node.(type) {
	case *ast.InfixExpression:

		switch operator {
		case "+":
			return &object.Float{Value: leftValue + rightValue}
		case "-":
			return &object.Float{Value: leftValue - rightValue}
		case "*":
			return &object.Float{Value: leftValue * rightValue}
		case "/":
			if rightValue == 0 {
				return newError(node.Token.FileName, node.Token.Line, node.Token.Column, "Can't divide by zero")
			}
			return &object.Float{Value: leftValue / rightValue}
		case "<":
			return nativeBoolToBooleanObject(leftValue < rightValue)
		case ">":
			return nativeBoolToBooleanObject(leftValue > rightValue)
		case "==":
			return nativeBoolToBooleanObject(leftValue == rightValue)
		case "!=":
			return nativeBoolToBooleanObject(leftValue != rightValue)
		case "-=":
			return &object.Float{Value: leftValue - rightValue}
		case "*=":
			return &object.Float{Value: leftValue * rightValue}
		case "+=":
			return &object.Float{Value: leftValue + rightValue}
		case "<=":
			return nativeBoolToBooleanObject(leftValue <= rightValue)
		case ">=":
			return nativeBoolToBooleanObject(leftValue >= rightValue)
		default:
			return newError(node.Token.FileName, node.Token.Line, node.Token.Column, "unknown operator: %s %s %s", leftOperand.Type(), operator, rightOperand.Type())
		}
	case *ast.AssignStatement:
		switch operator {
		case "+":
			return &object.Float{Value: leftValue + rightValue}
		case "-":
			return &object.Float{Value: leftValue - rightValue}
		case "*":
			return &object.Float{Value: leftValue * rightValue}
		case "/":
			if rightValue == 0 {
				return newError(node.Token.FileName, node.Token.Line, node.Token.Column, "Can't divide by zero")
			}
			return &object.Float{Value: leftValue / rightValue}
		case "<":
			return nativeBoolToBooleanObject(leftValue < rightValue)
		case ">":
			return nativeBoolToBooleanObject(leftValue > rightValue)
		case "==":
			return nativeBoolToBooleanObject(leftValue == rightValue)
		case "!=":
			return nativeBoolToBooleanObject(leftValue != rightValue)
		case "-=":
			return &object.Float{Value: leftValue - rightValue}
		case "*=":
			return &object.Float{Value: leftValue * rightValue}
		case "+=":
			return &object.Float{Value: leftValue + rightValue}
		case "<=":
			return nativeBoolToBooleanObject(leftValue <= rightValue)
		case ">=":
			return nativeBoolToBooleanObject(leftValue >= rightValue)
		default:
			return newError(node.Token.FileName, node.Token.Line, node.Token.Column, "unknown operator: %s %s %s", leftOperand.Type(), operator, rightOperand.Type())
		}
	}
	return NULL
}

func evalFloatIntegerInfixExpression(node ast.Expression, operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := float64(right.(*object.Integer).Value)

	switch node := node.(type) {
	case *ast.InfixExpression:
		switch operator {
		case "+":
			return &object.Float{Value: leftVal + rightVal}
		case "+=":
			return &object.Float{Value: leftVal + rightVal}
		case "-":
			return &object.Float{Value: leftVal - rightVal}
		case "-=":
			return &object.Float{Value: leftVal - rightVal}
		case "*":
			return &object.Float{Value: leftVal * rightVal}
		case "*=":
			return &object.Float{Value: leftVal * rightVal}
		case "/":
			if rightVal == 0 {
				return newError(node.Token.FileName, node.Token.Line, node.Token.Column, "Can't divide by zero")
			}
			return &object.Float{Value: leftVal / rightVal}
		case "<":
			return nativeBoolToBooleanObject(leftVal < rightVal)
		case "<=":
			return nativeBoolToBooleanObject(leftVal <= rightVal)
		case ">":
			return nativeBoolToBooleanObject(leftVal > rightVal)
		case ">=":
			return nativeBoolToBooleanObject(leftVal >= rightVal)
		case "==":
			return nativeBoolToBooleanObject(leftVal == rightVal)
		case "!=":
			return nativeBoolToBooleanObject(leftVal != rightVal)
		default:
			return newError(node.Token.FileName, node.Token.Line, node.Token.Column, "unknown operator: %s %s %s",
				left.Type(), operator, right.Type())
		}
	case *ast.AssignStatement:
		switch operator {
		case "+":
			return &object.Float{Value: leftVal + rightVal}
		case "+=":
			return &object.Float{Value: leftVal + rightVal}
		case "-":
			return &object.Float{Value: leftVal - rightVal}
		case "-=":
			return &object.Float{Value: leftVal - rightVal}
		case "*":
			return &object.Float{Value: leftVal * rightVal}
		case "*=":
			return &object.Float{Value: leftVal * rightVal}
		case "/":
			if rightVal == 0 {
				return newError(node.Token.FileName, node.Token.Line, node.Token.Column, "Can't divide by zero")
			}
			return &object.Float{Value: leftVal / rightVal}
		case "<":
			return nativeBoolToBooleanObject(leftVal < rightVal)
		case "<=":
			return nativeBoolToBooleanObject(leftVal <= rightVal)
		case ">":
			return nativeBoolToBooleanObject(leftVal > rightVal)
		case ">=":
			return nativeBoolToBooleanObject(leftVal >= rightVal)
		case "==":
			return nativeBoolToBooleanObject(leftVal == rightVal)
		case "!=":
			return nativeBoolToBooleanObject(leftVal != rightVal)
		default:
			return newError(node.Token.FileName, node.Token.Line, node.Token.Column, "unknown operator: %s %s %s",
				left.Type(), operator, right.Type())
		}
	}
	return NULL
}

func evalIntegerFloatInfixExpression(node ast.Expression, operator string, left, right object.Object) object.Object {
	leftVal := float64(left.(*object.Integer).Value)
	rightVal := right.(*object.Float).Value

	switch node := node.(type) {
	case *ast.InfixExpression:
		switch operator {
		case "+":
			return &object.Float{Value: leftVal + rightVal}
		case "-":
			return &object.Float{Value: leftVal - rightVal}
		case "*":
			return &object.Float{Value: leftVal * rightVal}

		case "/":
			if rightVal == 0 {
				return newError(node.Token.FileName, node.Token.Line, node.Token.Column, "divide by zero")
			}
			return &object.Float{Value: leftVal / rightVal}
		case "<":
			return nativeBoolToBooleanObject(leftVal < rightVal)
		case ">":
			return nativeBoolToBooleanObject(leftVal > rightVal)
		case "==":
			return nativeBoolToBooleanObject(leftVal == rightVal)
		case "!=":
			return nativeBoolToBooleanObject(leftVal != rightVal)
		case "-=":
			return &object.Float{Value: leftVal - rightVal}
		case "*=":
			return &object.Float{Value: leftVal * rightVal}
		case "+=":
			return &object.Float{Value: leftVal + rightVal}
		case "<=":
			return nativeBoolToBooleanObject(leftVal <= rightVal)
		case ">=":
			return nativeBoolToBooleanObject(leftVal >= rightVal)
		default:
			return newError(node.Token.FileName, node.Token.Line, node.Token.Column, "unknown operator: %s %s %s",
				left.Type(), operator, right.Type())
		}
	case *ast.AssignStatement:
		switch operator {
		case "+":
			return &object.Float{Value: leftVal + rightVal}
		case "-":
			return &object.Float{Value: leftVal - rightVal}
		case "*":
			return &object.Float{Value: leftVal * rightVal}

		case "/":
			if rightVal == 0 {
				return newError(node.Token.FileName, node.Token.Line, node.Token.Column, "divide by zero")
			}
			return &object.Float{Value: leftVal / rightVal}
		case "<":
			return nativeBoolToBooleanObject(leftVal < rightVal)
		case ">":
			return nativeBoolToBooleanObject(leftVal > rightVal)
		case "==":
			return nativeBoolToBooleanObject(leftVal == rightVal)
		case "!=":
			return nativeBoolToBooleanObject(leftVal != rightVal)
		case "-=":
			return &object.Float{Value: leftVal - rightVal}
		case "*=":
			return &object.Float{Value: leftVal * rightVal}
		case "+=":
			return &object.Float{Value: leftVal + rightVal}
		case "<=":
			return nativeBoolToBooleanObject(leftVal <= rightVal)
		case ">=":
			return nativeBoolToBooleanObject(leftVal >= rightVal)
		default:
			return newError(node.Token.FileName, node.Token.Line, node.Token.Column, "unknown operator: %s %s %s",
				left.Type(), operator, right.Type())
		}
	}
	return NULL
}

func evalPrefixExpression(node *ast.PrefixExpression, operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(node, right)
	default:
		return newError(node.Token.FileName, node.Token.Line, node.Token.Column, "unknown operator: %s%s", operator, right.Type())
	}
}

func evalPostfixExpression(node *ast.PostfixExpression, env *object.Environment, operator string) object.Object {
	switch operator {
	case "++":
		val, ok := env.Get(node.Token.Literal)
		if !ok {
			return newError(node.Token.FileName, node.Token.Line, node.Token.Column, "%s is unknown", node.Token.Literal)
		}

		switch arg := val.(type) {
		case *object.Integer:
			v := arg.Value
			env.Set(node.Token.Literal, &object.Integer{Value: v + 1})
			return arg
		default:
			return newError(node.Token.FileName, node.Token.Line, node.Token.Column, "%s is not an int", node.Token.Literal)

		}
	case "--":
		val, ok := env.Get(node.Token.Literal)
		if !ok {
			return newError(node.Token.FileName, node.Token.Line, node.Token.Column, "%s is unknown", node.Token.Literal)
		}

		switch arg := val.(type) {
		case *object.Integer:
			v := arg.Value
			env.Set(node.Token.Literal, &object.Integer{Value: v - 1})
			return arg
		default:
			return newError(node.Token.FileName, node.Token.Line, node.Token.Column, "%s is not an int", node.Token.Literal)
		}
	default:
		return newError(node.Token.FileName, node.Token.Line, node.Token.Column, "unknown operator: %s", operator)
	}
}

// evalMinusPrefixOperatorExpression evaluates the right object and returns a new object with the value negated
func evalMinusPrefixOperatorExpression(node *ast.PrefixExpression, right object.Object) object.Object {

	if right.Type() != object.INTEGER_OBJ {
		return newError(node.Token.FileName, node.Token.Line, node.Token.Column, "unknown operator: -%s", right.Type())
	}
	// if operand ins't integer - we escape
	if right.Type() != object.INTEGER_OBJ {
		return NULL
	}
	// retrieve & negate the value
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

// evalBangOperatorExpression evaluates the right object and returns a new object with the value negated
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

func evalHashLiteral(
	node *ast.HashLiteral,
	env *object.Environment,
) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}
		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError(node.Token.FileName, node.Token.Line, node.Token.Column, "unusable as hash key: %s", key.Type())
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}
	return &object.Hash{Pairs: pairs}
}

func evalStatements(statements []ast.Statement, env *object.Environment) object.Object {
	var evaluatedResults object.Object
	for _, statement := range statements {
		evaluatedResults = Eval(statement, env)

		// if the statement is a return statement, return the value
		if evaluatedResults, ok := evaluatedResults.(*object.ReturnValue); ok {
			return evaluatedResults.Value
		}
	}
	return evaluatedResults
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		result = Eval(statement, env)
		if result != nil {
			resultType := result.Type()
			if resultType == object.RETURN_VALUE_OBJ || resultType == object.ERROR_OBJ {
				return result
			}
		}
	}
	return result
}

func evalObjectCallExpression(call *ast.ObjectCallExpression, env *object.Environment) object.Object {
	objectValue := Eval(call.Object, env)

	if objectValue == nil {
		return newError(call.Token.FileName, call.Token.Line, call.Token.Column, "object is nil")
	}
	if method, ok := call.Call.(*ast.CallExpression); ok {
		args := evalExpressions(call.Call.(*ast.CallExpression).Arguments, env)
		ret := objectValue.InvokeMethod(method.Function.String(), *env, args...)
		if ret != nil {
			return ret
		}
	}
	// TODO: check if the object has the method implemented in esolang

	return newError(call.Token.FileName, call.Token.Line, call.Token.Column, "value of type `%s` has no member `%s`", objectValue.Type(), call.Call.String())
}

func evalAssignStatement(node *ast.AssignStatement, env *object.Environment) (val object.Object) {
	evaluated := Eval(node.Value, env)
	if isError(evaluated) {
		return evaluated
	}
	switch node.Operator {
	case "+=":
		current, ok := env.Get(node.Name.String())

		if !ok {
			return newError(node.Token.FileName, node.Token.Line, node.Token.Column, "%s is unknown", node.Name.String())
		}

		res := evalInfixExpression("+=", node, current, evaluated)
		if isError(res) {
			return res
		}

		env.Set(node.Name.String(), res)
		return res

	case "-=":

		// Get the current value
		current, ok := env.Get(node.Name.String())
		if !ok {
			return newError(node.Token.FileName, node.Token.Line, node.Token.Column, "%s is unknown", node.Name.String())
		}

		res := evalInfixExpression("-=", node, current, evaluated)
		if isError(res) {
			return res
		}

		env.Set(node.Name.String(), res)
		return res

	case "*=":
		// Get the current value
		current, ok := env.Get(node.Name.String())
		if !ok {
			return newError(node.Token.FileName, node.Token.Line, node.Token.Column, "%s is unknown", node.Name.String())
		}

		res := evalInfixExpression("*=", node, current, evaluated)
		if isError(res) {
			return res
		}

		env.Set(node.Name.String(), res)
		return res

	case "=":
		env.Set(node.Name.String(), evaluated)
	}
	return evaluated
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var evaluatedResult object.Object
	for _, statement := range program.Statements {
		evaluatedResult = Eval(statement, env)
		switch evaluatedResult := evaluatedResult.(type) {
		case *object.ReturnValue:
			return evaluatedResult.Value
		case *object.Error:
			return evaluatedResult
		}
	}
	return evaluatedResult
}

func Module(node *ast.ImportExpression, name string) object.Object {
	error := func(message string) object.Object {
		return &object.Error{
			Message: message,
		}
	}

	parseAndGetModuleHash := func(source, code string) object.Object {
		l := lexer.New(source, string(code))
		p := parser.New(l)

		module := p.ParseProgram()
		if len(p.Errors()) != 0 {
			return error(strings.Join(p.Errors(), "\n"))
		}

		env := object.NewEnvironment()
		objeEval := Eval(module, env)
		switch encounterError := objeEval.(type) {
		case *object.Error:
			return encounterError
		}

		return env.ExportedHash()
	}

	if utils.IsBuiltinModule(name) {
		// TODO: line numbers and column numbers for built-in modules
		// errors should be the node's line and column numbers instead
		moduleName := strings.Split(name, "/")[1]
		moduleCode, err := builtins.GetStdLib(moduleName)
		if err != nil {
			return error(err.Error())
		}
		return parseAndGetModuleHash(name, moduleCode)
	}

	filename := utils.FindModule(name)
	if filename == "" {
		return error(fmt.Sprintf("ImportError: no module named '%s'", name))
	}

	b, err := os.ReadFile(filename)
	if err != nil {
		return error(fmt.Sprintf("IOError: error reading module '%s': %s", name, err))
	}

	sourceName := strings.Split(filename, "/")
	source := sourceName[len(sourceName)-2] + "/" + sourceName[len(sourceName)-1]
	return parseAndGetModuleHash(source, string(b))
}

func applyFunction(node *ast.CallExpression, fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {

	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError(node.Token.FileName, node.Token.Line, node.Token.Column, "not a function: %s", fn.Type())
	}
}

func unwrapReturnValue(evaluated object.Object) object.Object {
	if returnValue, ok := evaluated.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return evaluated
}

func extendFunctionEnv(function *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(function.Env)
	for paramIdx, param := range function.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	return env
}

func newError(fileName string, line, column int, format string, a ...interface{}) *object.Error {
	linesAndCol := fmt.Sprintf("%s:%d:%d:", fileName, line, column)
	format = fmt.Sprintf("%s %s", linesAndCol, format)
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

// nativeBoolToBooleanObject converts a native bool to a boolean object
func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func objectToNativeBoolean(o object.Object) bool {
	if r, ok := o.(*object.ReturnValue); ok {
		o = r.Value
	}
	switch obj := o.(type) {
	case *object.Boolean:
		return obj.Value
	case *object.String:
		return obj.Value != ""
	case *object.Null:
		return false
	case *object.Integer:
		if obj.Value == 0 {
			return false
		}
		return true
	case *object.Array:
		if len(obj.Elements) == 0 {
			return false
		}
		return true
	case *object.Hash:
		if len(obj.Pairs) == 0 {
			return false
		}
		return true
	default:
		return true
	}
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
