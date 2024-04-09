/*
Package token defines the token type and the token struct.
*/
package token

type TokenType string

// Token represents a token
type Token struct {
	Type    TokenType // Type of token
	Literal string    // Literal value of token
}

const (
	ILLEGAL   = "ILLEGAL"
	EOF       = "EOF"
	IDENT     = "IDENT"
	INT       = "INT"
	ASSIGN    = "="
	PLUS      = "+"
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	FUNCTION  = "FUNCTION"
	LET       = "LET"
	BANG      = "!"
	SLASH     = "/"
	ASTERISK  = "*"
	MINUS     = "-"
	LT        = "<"
	GT        = ">"
	TRUE      = "TRUE"
	FALSE     = "FALSE"
	IF        = "IF"
	ELSE      = "ELSE"
	RETURN    = "RETURN"
	EQ        = "=="
	NOT_EQ    = "!="
	STRING    = "STRING"
	LBRACKET  = "["
	RBRACKET  = "]"
	COLON     = ":"
	WHILE     = "WHILE"
	MOD       = "%"
	AND       = "&&"
	OR        = "-|"
	PERIOD    = "."
)

// Keywords are reserved words
var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"while":  WHILE,
}

// LookupIdent checks if the identifier is a keyword
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
