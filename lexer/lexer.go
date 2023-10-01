package lexer

import "monkey/lang-monkey/token"

type Lexer struct {
	input        string // source code
	position     int    // pointer to current char
	readPosition int    // current reading position after currenct char
	char         byte   // currenct char under examination

	/*
	 *The reason for these two “pointers” pointing into our input string
	 *is the fact that we will need to be able to “peek” further into the
	 *input and look after the current character to see what comes up next.
	 *readPosition always points to the “next” character in the input.
	 *position points to the character in the input that corresponds to the ch byte.”
	 */
}

func New(input string) *Lexer {
	L := &Lexer{input: input}
	L.readChar()
	return L
}

func (L *Lexer) readChar() {
	if L.readPosition >= len(L.input) {
		L.char = 0 // “ASCII code for the "NUL” EOF OR NOTHING READ FROM INPUT
	} else {
		L.char = L.input[L.readPosition]
	}
	L.position = L.readPosition
	L.readPosition += 1
}

func (L *Lexer) NextToken() token.Token {
	var tok token.Token
	L.skipWhiteSpaces()
	switch L.char {
	default:
		if isLetter(L.char) {
			tok.Literal = L.readIdentifier()
			tok.Type = token.LookupIdentif(tok.Literal)
			return tok
		} else if isDigit(L.char) {
			tok.Type = token.INT
			tok.Literal = L.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, L.char)
		}
	case '=':
		tok = newToken(token.ASSIGN, L.char)
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
	case '{':
		tok = newToken(token.LBRACE, L.char)
	case '}':
		tok = newToken(token.RBRACE, L.char)
	case '!':
		tok = newToken(token.BANG, L.char)
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

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (L *Lexer) readIdentifier() string {
	position := L.position
	for isLetter(L.char) {
		L.readChar()
	}
	return L.input[position:L.position]

}

func (L *Lexer) readNumber() string {
	position := L.position
	for isDigit(L.char) {
		L.readChar()
	}
	return L.input[position:L.position]
}

func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

func isDigit(char byte) bool {
	return '0' <= char && char <= '9' // whether the byte is a Latin digit between 0 and 9
}

func (L *Lexer) skipWhiteSpaces() {
	for L.char == ' ' || L.char == '\t' || L.char == '\n' || L.char == '\r' {
		L.readChar()
	}
}
