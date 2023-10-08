package parser

import (
	"fmt"
	"monkey/lang-monkey/ast"
	"monkey/lang-monkey/lexer"
	"monkey/lang-monkey/token"
	"strconv"
)

const (
	_int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

type Parser struct {
	L              *lexer.Lexer
	currentToken   token.Token
	peekToken      token.Token
	errors         []string
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

func (P *Parser) Errors() []string {
	return P.errors
}

func New(L *lexer.Lexer) *Parser {
	p := &Parser{
		L:      L,
		errors: []string{},
	}
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	return p
}

func (P *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: P.currentToken, Value: P.currentToken.Literal}
}

func (P *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, P.peekToken.Type)
	P.errors = append(P.errors, msg)
}

func (P *Parser) nextToken() {
	P.currentToken = P.peekToken
	P.peekToken = P.L.NextToken()
}

func (P *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for P.currentToken.Type != token.EOF {
		stmt := P.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		P.nextToken()
	}
	return program
}

func (P *Parser) parseStatement() ast.Statement {
	switch P.currentToken.Type {
	case token.LET:
		return P.parseLetStatement()
	case token.RETURN:
		return P.parseReturnStatement()
	default:
		return P.parseExpressionStatement()
	}
}

func (P *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	P.prefixParseFns[tokenType] = fn
}

func (P *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	P.infixParseFns[tokenType] = fn
}

func (P *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: P.currentToken}
	P.nextToken()
	// TODO: Skipping expression until we encounter a semicolon
	if !P.currentTokenLS(token.SEMICOLON) {
		P.nextToken()
	}
	return stmt
}

func (P *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: P.currentToken}

	if !P.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: P.currentToken, Value: P.currentToken.Literal}

	if !P.expectPeek(token.ASSIGN) {
		return nil
	}

	//TODO: Skip expressions till encounter semicolons

	for !P.currentTokenLS(token.SEMICOLON) {
		P.nextToken()
	}
	return stmt
}

func (P *Parser) expectPeek(t token.TokenType) bool {
	if P.peekTokenLS(t) {
		P.nextToken()
		return true
	} else {
		P.peekError(t)
		return false
	}
}

func (P *Parser) currentTokenLS(t token.TokenType) bool { return P.currentToken.Type == t }
func (P *Parser) peekTokenLS(t token.TokenType) bool    { return P.peekToken.Type == t }

func (P *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: P.currentToken}
	stmt.Expression = P.parseExpression(LOWEST)

	if P.peekTokenLS(token.SEMICOLON) {
		P.nextToken()
	}
	return stmt
}

func (P *Parser) parseExpression(precedence int) ast.Expression {
	prefix := P.prefixParseFns[P.currentToken.Type]
	if prefix == nil {
		return nil
	}
	leftExp := prefix()
	return leftExp
}

func (P *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: P.currentToken}
	value, err := strconv.ParseInt(P.currentToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", P.currentToken)
		P.errors = append(P.errors, msg)
		return nil
	}
	lit.Value = value
	return lit

}
