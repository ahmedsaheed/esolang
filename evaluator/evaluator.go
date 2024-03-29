package evaluator

import (
	"esolang/lang-esolang/ast"
	"esolang/lang-esolang/builtins"
	"esolang/lang-esolang/object"
	"fmt"
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

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

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
		return evalInfixExpression(node.Operator, left, right)
	case *ast.WhileLoopExpression:
		return evalWhileLoopExpression(node, env)
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Body: body, Env: env}

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
		return evalIndexExpression(left, index)

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
		return applyFunction(function, args)

	case *ast.ReturnStatement:
		value := Eval(node.ReturnValue, env)
		if isError(value) {
			return value
		}
		return &object.ReturnValue{Value: value}
	}

	return nil
}

func evalIndexExpression(left object.Object, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
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

func evalHashIndexExpression(hash object.Object, index object.Object) object.Object {
	hashObj := hash.(*object.Hash)
	key, ok := index.(object.Hashable)

	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObj.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {

	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
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
	return newError("identifier not found: " + node.Value)
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
	case leftOperand.Type() == object.STRING_OBJ && rightOperand.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, leftOperand, rightOperand)
	case operator == "!=":
		return nativeBoolToBooleanObject(leftOperand != rightOperand)
	case operator == "&&":
		return nativeBoolToBooleanObject(objectToNativeBoolean(leftOperand) && objectToNativeBoolean(rightOperand))
	case operator == "-|":
		return nativeBoolToBooleanObject(objectToNativeBoolean(leftOperand) || objectToNativeBoolean(rightOperand))
	case leftOperand.Type() != rightOperand.Type():
		return newError("type mismatch: %s %s %s", leftOperand.Type(), operator, rightOperand.Type())
	default:
		return newError("unknown operator: %s %s %s", leftOperand.Type(), operator, rightOperand.Type())
	}
}

// evalStringInfixExpression evaluates the string concatenation
func evalStringInfixExpression(operator string, leftOperand, rightOperand object.Object) object.Object {
	if operator != "+" {
		return newError("unknown operator: %s %s %s", leftOperand.Type(), operator, rightOperand.Type())
	}
	leftValue := leftOperand.(*object.String).Value
	rightValue := rightOperand.(*object.String).Value
	return &object.String{Value: leftValue + rightValue}
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
	case "%":
		return &object.Integer{Value: leftValue % rightValue}
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
			return newError("unusable as hash key: %s", key.Type())
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

// nativeBoolToBooleanObject converts a native bool to a boolean object
func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
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

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
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

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
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
