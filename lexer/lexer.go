/*
Package lexer takes a source code as an input and outputs tokens that represent the given source code.
It goes through the input and outputs the next recognised token by calling `nextToken()`.
*/
package lexer

import "esolang/lang-esolang/token"

/*
Lexer represents a lexical analyzer.
It processes syntax or source code character by character.
*/
type Lexer struct {
	input        string // The input syntax or source code.
	position     int    // The current position in the input (pointer to the current character).
	readPosition int    // The current reading position in the input (pointer to after the current character).
	char         byte   // The current character under examination.
}

// New creates a new Lexer and returns a pointer to it.
func New(input string) *Lexer {
	L := &Lexer{input: input}
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
	switch L.char {
	default:
		if isLetter(L.char) {
			tok.Literal = L.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(L.char) {
			tok.Type = token.INT
			tok.Literal = L.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, L.char)
		}
	case '=':
		if L.peekChar() == '=' {
			char := L.char
			L.readChar()
			literal := string(char) + string(L.char)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, L.char)
		}
	case '"':
		tok.Type = token.STRING
		tok.Literal = L.readString()
	case ';':
		tok = newToken(token.SEMICOLON, L.char)
	case '(':
		tok = newToken(token.LPAREN, L.char)
	case ')':
		tok = newToken(token.RPAREN, L.char)
	case ',':
		tok = newToken(token.COMMA, L.char)
	case '+':
		tok = newToken(token.PLUS, L.char)

	case '[':
		tok = newToken(token.LBRACKET, L.char)
	case ']':
		tok = newToken(token.RBRACKET, L.char)
	case '{':
		tok = newToken(token.LBRACE, L.char)
	case '}':
		tok = newToken(token.RBRACE, L.char)
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
			tok = newToken(token.BANG, L.char)
		}
	case '/':
		tok = newToken(token.SLASH, L.char)
	case '*':
		tok = newToken(token.ASTERISK, L.char)
	case '-':
		tok = newToken(token.MINUS, L.char)
	case '<':
		tok = newToken(token.LT, L.char)
	case '>':
		tok = newToken(token.GT, L.char)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	}
	L.readChar()
	return tok
}

// newToken creates a new token with the given type and character.
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
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

// readIdentifier reads the next identifier in the input and returns it.
func (L *Lexer) readIdentifier() string {
	position := L.position
	for isLetter(L.char) {
		L.readChar()
	}
	return L.input[position:L.position]

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
