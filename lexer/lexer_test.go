package lexer

import (
	"monkey/lang-monkey/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `=+(){},;`
    tests := []struct {
        expectedType    token.TokenType
        expectedLiteral string
    }{
        {token.ASSIGN, "="},
        {token.PLUS, "+"},
        {token.LPAREN, "("},
        {token.RPAREN, ")"},
        {token.LBRACE, "{"},
        {token.RBRACE, "}"},
        {token.COMMA, ","},
        {token.SEMICOLON, ";"},
        {token.EOF, ""},
    }

    l := New(input)

    for i, tokens := range tests {
        tok := l.NextToken()
        
        if tok.Type != tokens.expectedType {
            t.Fatalf("“tests[%d] - tokentype wrong. expected=%q, got=%q",
                i, tokens.expectedType, tok.Type)
        }

        if tok.Literal != tokens.expectedLiteral {
            t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
                i, tokens.expectedLiteral, tok.Literal)
        }
    }

}

