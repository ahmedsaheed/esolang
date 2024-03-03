package evaluator

import (
	"esolang/lang-esolang/lexer"
	"esolang/lang-esolang/object"
	"esolang/lang-esolang/parser"
	"testing"
)

func TestEvaluateIntegerExpression(t *testing.T) {
	tests := []struct {
		input   string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testIntegerObject(t, evaluated, test.expected)
	}
}

func testIntegerObject(t *testing.T, obj object.Object ,expected int64) bool {
	evaluatedIntegerObject, ok := obj.(*object.Integer)

	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}

	if evaluatedIntegerObject.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", evaluatedIntegerObject.Value, expected)
		return false
	}
	return true


}

func testEval(input string) object.Object {
	lexer := lexer.New(input)
	parser := parser.New(lexer)
	program := parser.ParseProgram()
	return Eval(program)
}