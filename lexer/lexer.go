/*
Package lexer takes a source code as an input and outputs tokens that represent the given source code.
It goes through the input and outputs the next recognised token by calling `nextToken()`.
*/
package lexer

import (
	"esolang/lang-esolang/token"
	"strings"
	"unicode"
)

/*
Lexer represents a lexical analyzer.
It processes syntax or source code character by character.
*/
type Lexer struct {
	input        string // The input syntax or source code.
	position     int    // The current position in the input (pointer to the current character).
	readPosition int    // The current reading position in the input (pointer to after the current character).
	char         byte   // The current character under examination.
	line         int    // The current line number.
	column       int    // The current column number.
}

// New creates a new Lexer and returns a pointer to it.
func New(input string) *Lexer {
	L := &Lexer{input: input}
	L.line = 1
	L.column = 1
	L.readChar()
	return L
}

// readChar reads the next character in the input and advances the position and readPosition pointers.
func (L *Lexer) readChar() {
	if L.readPosition >= len(L.input) {
		L.char = 0
	} else {
		L.char = L.input[L.readPosition]
	}
	L.position = L.readPosition
	L.readPosition += 1
	L.column++
	if L.char == '\n' { // If the character is a newline, increment the line number.
		L.column = 1
		L.line++
	}
}

// peekChar returns the next character in the input without advancing the position and readPosition pointers.
func (L *Lexer) peekChar() byte {
	if L.readPosition >= len(L.input) {
		return 0
	} else {
		return L.input[L.readPosition]
	}
}

/*
NextToken returns the next token in the input.
It skips white spaces and returns the next token.
It returns a token with the type of ILLEGAL if the character is not recognised.
*/
func (L *Lexer) NextToken() token.Token {
	var tok token.Token
	L.skipWhiteSpaces()

	if L.char == '/' && L.peekChar() == '/' {
		L.skipComment()
		return (L.NextToken())
	}
	switch L.char {
	default:
		if isLetter(L.char) {
			tok.Literal = L.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			tok.Line = L.line
			tok.Column = L.column
			return tok
		} else if isDigit(L.char) {
			tok.Type = token.INT
			tok.Literal = L.readNumber()
			tok.Line = L.line
			tok.Column = L.column
			return tok
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(L.char), Line: L.line, Column: L.column}
		}
	case '=':
		if L.peekChar() == '=' {
			char := L.char
			L.readChar()
			literal := string(char) + string(L.char)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, L.char, L.line, L.column)
		}
	case '.':
		tok = newToken(token.PERIOD, L.char, L.line, L.column)
	case '&':
		if L.peekChar() == '&' {
			char := L.char
			L.readChar()
			literal := string(char) + string(L.char)
			tok = token.Token{Type: token.AND, Literal: literal}
		} else {
			tok = newToken(token.ILLEGAL, L.char, L.line, L.column)
		}
	case ':':
		tok = newToken(token.COLON, L.char, L.line, L.column)
	case '%':
		tok = newToken(token.MOD, L.char, L.line, L.column)
	case '"':
		tok.Type = token.STRING
		tok.Literal = L.readString()
	case ';':
		tok = newToken(token.SEMICOLON, L.char, L.line, L.column)
	case '(':
		tok = newToken(token.LPAREN, L.char, L.line, L.column)
	case ')':
		tok = newToken(token.RPAREN, L.char, L.line, L.column)
	case ',':
		tok = newToken(token.COMMA, L.char, L.line, L.column)
	case '+':
		tok = newToken(token.PLUS, L.char, L.line, L.column)
	case '[':
		tok = newToken(token.LBRACKET, L.char, L.line, L.column)
	case ']':
		tok = newToken(token.RBRACKET, L.char, L.line, L.column)
	case '{':
		tok = newToken(token.LBRACE, L.char, L.line, L.column)
	case '}':
		tok = newToken(token.RBRACE, L.char, L.line, L.column)
	case '!':
		if L.peekChar() == '=' {
			char := L.char
			L.readChar()
			literal := string(char) + string(L.char)
			tok = token.Token{
				Type:    token.NOT_EQ,
				Literal: literal,
			}
		} else {
			tok = newToken(token.BANG, L.char, L.line, L.column)
		}
	case '/':
		tok = newToken(token.SLASH, L.char, L.line, L.column)
	case '*':
		tok = newToken(token.ASTERISK, L.char, L.line, L.column)
	case '-':
		if L.peekChar() == '|' {
			char := L.char
			L.readChar()
			literal := string(char) + string(L.char)
			tok = token.Token{Type: token.OR, Literal: literal}
		} else {
			tok = newToken(token.MINUS, L.char, L.line, L.column)
		}
	case '<':
		tok = newToken(token.LT, L.char, L.line, L.column)
	case '>':
		tok = newToken(token.GT, L.char, L.line, L.column)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	}
	L.readChar()
	return tok
}

// newToken creates a new token with the given type and character.
func newToken(tokenType token.TokenType, ch byte, line, col int) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), Line: line, Column: col}
}

func (L *Lexer) readString() string {
	position := L.position + 1
	for {
		L.readChar()
		if L.char == '"' || L.char == 0 {
			break
		}
	}
	return L.input[position:L.position]

}

func (L *Lexer) skipComment() {
	for L.char != '\n' && L.char != 0 {
		L.readChar()
	}
	L.skipWhiteSpaces()
}

// readIdentifier reads the next identifier in the input and returns it.
func (L *Lexer) readIdentifier() string {

	valid := map[string]bool{
		"directory.glob":     true,
		"math.abs":           true,
		"math.rand":          true,
		"math.sqrt":          true,
		"os.environment":     true,
		"os.getenv":          true,
		"os.setenv":          true,
		"string.interpolate": true,
	}

	types := []string{"string.",
		"array.",
		"integer.",
		"float.",
		"hash.",
		"object."}

	id := ""

	position := L.position
	rposition := L.readPosition

	for isLetter(L.char) {
		id += string(L.char)
		L.readChar()
	}

	if strings.Contains(id, ".") {
		ok := valid[id]

		if !ok {
			for _, t := range types {
				if strings.HasPrefix(id, t) {
					ok = true
				}
			}
		}

		if !ok {
			offset := strings.Index(id, ".")
			id = id[:offset]
			L.position = position
			L.readPosition = rposition
			for offset > 0 {
				L.readChar()
				offset--
			}
		}
	}

	// for isLetter(L.char) {
	// 	L.readChar()
	// }
	// return L.input[position:L.position]

	return id

}

// readNumber reads the next number in the input and returns it.
func (L *Lexer) readNumber() string {
	position := L.position
	for isDigit(L.char) {
		L.readChar()
	}
	return L.input[position:L.position]
}

// isLetter returns true if the given character is a letter.
func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

func isIdentifier(char byte) bool {
	if unicode.IsLetter(rune(char)) || unicode.IsDigit(rune(char)) || rune(char) == '.' || rune(char) == '?' || rune(char) == '$' || rune(char) == '_' {
		return true
	}
	return false
}

// isDigit returns true if the given character is a digit.
func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}

// skipWhiteSpaces skips white spaces in the input.
func (L *Lexer) skipWhiteSpaces() {
	for L.char == ' ' || L.char == '\t' || L.char == '\n' || L.char == '\r' {
		L.readChar()
	}
}
