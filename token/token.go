/*
Package token defines the token type and the token struct.
*/
package token

type TokenType string

// Token represents a token
type Token struct {
	Type     TokenType // Type of token
	Literal  string    // Literal value of token
	Line     int       // Line number of token
	Column   int       // Column number of token
	FileName string    // File name where token is found - used for debugging
}

const (
	ILLEGAL     = "ILLEGAL"
	EOF         = "EOF"
	IDENT       = "IDENT"
	INT         = "INT"
	FLOAT       = "FLOAT"
	ASSIGN      = "="
	PLUS        = "+"
	PLUS_PLUS   = "++"
	ASTERISK_EQ = "*="
	PLUS_EQ     = "+="
	MINUS_EQ    = "-="
	MINUS_MINUS = "--"
	COMMA       = ","
	SEMICOLON   = ";"
	LPAREN      = "("
	RPAREN      = ")"
	LBRACE      = "{"
	RBRACE      = "}"
	FUNCTION    = "FUNCTION"
	DEF_FN      = "DEF_FUNTION"
	LET         = "LET"
	BANG        = "!"
	SLASH       = "/"
	ASTERISK    = "*"
	MINUS       = "-"
	LT          = "<"
	GT          = ">"
	LT_EQ       = "<="
	GT_EQ       = ">="
	TRUE        = "TRUE"
	FALSE       = "FALSE"
	IF          = "IF"
	ELSE        = "ELSE"
	RETURN      = "RETURN"
	EQ          = "=="
	NOT_EQ      = "!="
	STRING      = "STRING"
	LBRACKET    = "["
	RBRACKET    = "]"
	COLON       = ":"
	DOUBLECOL   = "::"
	BIND        = ":="
	WHEN        = "WHEN"
	MOD         = "%"
	AND         = "&&"
	OR          = "-|"
	PERIOD      = "."
	IMPORT      = "IMPORT"
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
	"when":   WHEN,
	"import": IMPORT,
	"func":   DEF_FN,
}

// LookupIdent checks if the identifier is a keyword
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
