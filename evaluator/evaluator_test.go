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
		{"-15", -15},
		{"-5", -5},
		{"10 + 10 + 10 + 10 - 5", 35},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
		{"100 + 10 / 50 + 20", 120},
		{"40 - 10 + 90 /2", 75 },
		{"(80-20 + 100) / 2", 80},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testIntegerObject(t, evaluated, test.expected)
	}
}

func TestEvaluateBooleanExpression(t *testing.T) {
	tests := []struct{
		input string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testBooleanObject(t, evaluated, test.expected)
	}
}

// evaluating prefix expressions
func TestBangOperator(t *testing.T) {
	tests := []struct {
		input string
		expected bool
	}{
		{"!true", false},
		{"!!true", true},
		{"!false", true},
		{"!!false", false},
		{"!15", false},
		{"!!15", true},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testBooleanObject(t, evaluated, test.expected)
	}
}

func testBooleanObject(t *testing.T, evaluated object.Object, b bool) bool {
	boolean, ok := evaluated.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", evaluated, evaluated)
		return false
	}
	if boolean.Value != b {
		t.Errorf("object value mismatch. got=%t, want=%t", boolean.Value, b)
		return false
	}
	return true
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