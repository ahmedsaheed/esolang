package parser

import (
	"monkey/lang-monkey/ast"
	"monkey/lang-monkey/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func TestReturnStatement(t *testing.T) {
	input := `
return 5;
return 10;
return 093221;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements got%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral is not 'return', got %q", returnStmt.TokenLiteral())
		}
	}
}

//func TestLetStatements(t *testing.T) {
//	tests := []struct {
//		input              string
//		expectedIdentifier string
//		expectedValue      interface{}
//	}{
//		{"let x = 5;", "x", 5},
//		{"let y = true;", "y", true},
//		{"let foobar = y;", "foobar", "y"},
//	}
//
//	for _, tt := range tests {
//		l := lexer.New(tt.input)
//		p := New(l)
//		program := p.ParseProgram()
//		checkParserErrors(t, p)
//
//		if len(program.Statements) != 1 {
//			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
//				len(program.Statements))
//		}
//
//		stmt := program.Statements[0]
//		if !testLetStatements(t, stmt, tt.expectedIdentifier) {
//			return
//		}
//
//		val := stmt.(*ast.LetStatement).Value
//		if !testLiteralExpression(t, val, tt.expectedValue) {
//			return
//		}
//	}
//}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifierLiteral(t, exp, v)
	case bool:
		testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of expression not handled. got=%T", exp)
	return false
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.tokenLiteral not 'let'. got=%s", s.TokenLiteral())
		return false
	}
	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s", name, letStmt.Name.TokenLiteral())
		return false
	}
	return true
}
func testIntegerLiteral(t *testing.T, expression ast.Expression, val int64) bool     { return false }
func testIdentifierLiteral(t *testing.T, expression ast.Expression, val string) bool { return false }
func testBooleanLiteral(t *testing.T, expression ast.Expression, val bool) bool      { return false }

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
