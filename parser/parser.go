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

var precedence = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

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
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)

	return p
}

func (P *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: P.currentToken,
		Value: P.currentTokenLS(token.TRUE),
	}
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
	defer untrace(trace("parseExpressionStatement"))
	stmt := &ast.ExpressionStatement{Token: P.currentToken}
	stmt.Expression = P.parseExpression(LOWEST)

	if P.peekTokenLS(token.SEMICOLON) {
		P.nextToken()
	}
	return stmt
}

func (P *Parser) parseIntegerLiteral() ast.Expression {
	defer untrace(trace("parseIntegerLiteral"))
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

func (P *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	P.errors = append(P.errors, msg)
}

func (P *Parser) parsePrefixExpression() ast.Expression {
	defer untrace(trace("parsePrefixExpression"))
	expression := &ast.PrefixExpression{
		Token:    P.currentToken,
		Operator: P.currentToken.Literal,
	}
	P.nextToken()
	expression.Right = P.parseExpression(PREFIX)
	return expression
}

func (P *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	defer untrace(trace("parseInfixExpression"))
	expression := &ast.InfixExpression{
		Token:    P.currentToken,
		Operator: P.currentToken.Literal,
		Left:     left,
	}

	precedence := P.currPrecedence()
	P.nextToken()
	expression.Right = P.parseExpression(precedence)

	return expression
}

func (P *Parser) parseExpression(precedence int) ast.Expression {
	defer untrace(trace("parseExpression"))
	prefix := P.prefixParseFns[P.currentToken.Type]
	if prefix == nil {
		P.noPrefixParseFnError(P.currentToken.Type)
		return nil
	}
	leftExp := prefix()

	for !P.peekTokenLS(token.SEMICOLON) && precedence < P.peekPrecedence() {
		infix := P.infixParseFns[P.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		P.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (P *Parser) peekPrecedence() int {
	if p, ok := precedence[P.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (P *Parser) currPrecedence() int {
	if p, ok := precedence[P.currentToken.Type]; ok {
		return p
	}
	return LOWEST
}
