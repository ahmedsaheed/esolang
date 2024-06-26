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
	BACKTICK    = "`"
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
	ELSE_IF     = "ELIF"
	RETURN      = "RETURN"
	EQ          = "=="
	NOT_EQ      = "!="
    STRING_EQ   = "IS"
    STRING_NOT_EQ = "IS_NOT"
	STRING      = "STRING"
	LBRACKET    = "["
	RBRACKET    = "]"
	COLON       = ":"
	DOUBLECOL   = "::"
	BIND        = ":="
	WHEN        = "WHEN"
	MOD         = "%"
	AND         = "&&"
	OR          = "||"
	STRING_OR   = "OR"
	STRING_AND  = "AND"
	PERIOD      = "."
	IMPORT      = "IMPORT"
)

// Keywords are reserved words
var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"and":    AND,
	"or":     OR,
    "is":     EQ,
    "is_not": NOT_EQ,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"elif":   ELSE_IF,
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
