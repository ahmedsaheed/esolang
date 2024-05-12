/*
Package parser implements the core top-down recursive descent parser i.e a Pratt parser
*/
package parser

import (
	"esolang/lang-esolang/ast"
	"esolang/lang-esolang/lexer"
	"esolang/lang-esolang/token"
	"fmt"
	"strconv"
)

const (
	_int = iota
	LOWEST
	EQUALS
	ANDOR
	LESSGREATER
	SUM
	PRODUCT
	MODULUS
	PREFIX
	CALL
	INDEX
	HIGHEST
)

// precedence map to determine the precedence of the operators
var precedence = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.AND:      ANDOR,
	token.OR:       ANDOR,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.MOD:      MODULUS,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
	token.PERIOD:   CALL,
	token.LBRACKET: INDEX,
}

// Parser is the core struct for the parser
type Parser struct {
	L              *lexer.Lexer                      // lexer
	currentToken   token.Token                       // current token
	peekToken      token.Token                       // next token
	errors         []string                          // errors
	prefixParseFns map[token.TokenType]prefixParseFn // prefix parse functions
	infixParseFns  map[token.TokenType]infixParseFn  // infix parse functions
}

type (
	prefixParseFn func() ast.Expression               // prefix parse function
	infixParseFn  func(ast.Expression) ast.Expression // infix parse function
)

// Errors returns the errors encountered during parsing
func (P *Parser) Errors() []string {
	return P.errors
}

/*
	New creates a new parser

* @param L *lexer.Lexer
* @return *Parser
*/
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
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.WHEN, p.parseWhenLoopExpression)
	p.registerPrefix(token.LBRACE, p.parseHashLiteral)
	p.registerPrefix(token.IMPORT, p.parseImportExpression)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.MOD, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	p.registerInfix(token.PERIOD, p.parseMethodCallExpression)

	return p
}

/*
parseGroupedExpression parses a grouped expression

	expression like (5 + 5) or (5 * 5)
*/
func (P *Parser) parseGroupedExpression() ast.Expression {
	P.nextToken()
	exp := P.parseExpression(LOWEST)

	if !P.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

/*
parseIfExpression parses an if expression

	expression like if (x < y) {x} else {y}
*/
func (P *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: P.currentToken}

	if !P.expectPeek(token.LPAREN) {
		return nil
	}

	P.nextToken()
	expression.Condition = P.parseExpression(LOWEST)

	if !P.expectPeek(token.RPAREN) {
		return nil
	}

	if !P.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = P.parseBlockStatement()

	if P.peekTokenLS(token.ELSE) {
		P.nextToken()

		if !P.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = P.parseBlockStatement()
	}

	return expression
}

/*
parseBoolean parses a boolean expression

	expression like true or false
*/
func (P *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: P.currentToken,
		Value: P.currentTokenLS(token.TRUE),
	}
}

/*
parseIdentifier parses an identifier

	expression like x or y
*/
func (P *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: P.currentToken, Value: P.currentToken.Literal}
}

/*
peekError returns an error message when the next token is not as expected
*/
func (P *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("Line %v Column %v - expected next token to be %s got %s", P.currentToken.Line, P.currentToken.Column, t, P.peekToken.Type)
	P.errors = append(P.errors, msg)
}

/*
nextToken moves to the next token
*/
func (P *Parser) nextToken() {
	P.currentToken = P.peekToken
	P.peekToken = P.L.NextToken()
}

/*
ParseProgram parses the program to create an AST

	such as let x = 5; let y = 10; return x + y;
*/
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

/*
parseStatement parses a statement

	parses a statement such as let x = 5; or return x + y;
*/
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

// registerPrefix registers a prefix parse function
func (P *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	P.prefixParseFns[tokenType] = fn
}

// registerInfix registers an infix parse function
func (P *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	P.infixParseFns[tokenType] = fn
}

// parseReturnStatement parses a return statement
func (P *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: P.currentToken}
	P.nextToken()
	stmt.ReturnValue = P.parseExpression(LOWEST)
	if P.peekTokenLS(token.SEMICOLON) {
		P.nextToken()
	}
	return stmt
}

// parseLetStatement parses a let statement
func (P *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: P.currentToken}

	if !P.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: P.currentToken, Value: P.currentToken.Literal}

	if !P.expectPeek(token.ASSIGN) {
		return nil
	}

	P.nextToken()

	stmt.Value = P.parseExpression(LOWEST)

	if P.peekTokenLS(token.SEMICOLON) {
		P.nextToken()
	}

	return stmt
}

// expectPeek checks if the next token is as expected - returns a boolean
func (P *Parser) expectPeek(t token.TokenType) bool {
	if P.peekTokenLS(t) {
		P.nextToken()
		return true
	} else {
		P.peekError(t)
		return false
	}
}

// currentTokenLS checks if the current matches the expected tokenType - returns a boolean
func (P *Parser) currentTokenLS(t token.TokenType) bool { return P.currentToken.Type == t }

// peekTokenLS checks if the next token matches the expected tokenType - returns a boolean
func (P *Parser) peekTokenLS(t token.TokenType) bool { return P.peekToken.Type == t }

// parseExpressionStatement parses an expression statement such as 5 + 5 or x * y;
func (P *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	// defer untrace(trace("parseExpressionStatement"))
	stmt := &ast.ExpressionStatement{Token: P.currentToken}
	stmt.Expression = P.parseExpression(LOWEST)

	for P.peekTokenLS(token.SEMICOLON) {
		P.nextToken()
	}
	return stmt
}

// parseIntegerLiteral parses an integer literal such as 5 or 10
func (P *Parser) parseIntegerLiteral() ast.Expression {
	// defer untrace(trace("parseIntegerLiteral"))
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

// noPrefixParseFnError returns an error message when no prefix parse function is found
func (P *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("Line %v Column %v - no prefix parse function for %s found", P.currentToken.Line, P.currentToken.Column, t)
	P.errors = append(P.errors, msg)
}

// parsePrefixExpression parses a prefix expression such as -5 or !5
func (P *Parser) parsePrefixExpression() ast.Expression {
	// defer untrace(trace("parsePrefixExpression"))
	expression := &ast.PrefixExpression{
		Token:    P.currentToken,
		Operator: P.currentToken.Literal,
	}
	P.nextToken()
	expression.Right = P.parseExpression(PREFIX)
	return expression
}

// parseInfixExpression parses an infix expression such as 5 + 5 or x * y
func (P *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	// defer untrace(trace("parseInfixExpression"))
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

// parseExpression parses an expression such as 5 + 5 or x * y
func (P *Parser) parseExpression(precedence int) ast.Expression {
	// defer untrace(trace("parseExpression"))
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

// peekPrecedence returns the precedence of the next token
func (P *Parser) peekPrecedence() int {
	if p, ok := precedence[P.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// currPrecedence returns the precedence of the current token
func (P *Parser) currPrecedence() int {
	if p, ok := precedence[P.currentToken.Type]; ok {
		return p
	}
	return LOWEST
}

// parseBlockStatement parses a block statement such as {x + y;} or {y != x;}
func (P *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: P.currentToken}
	block.Statements = []ast.Statement{}

	P.nextToken()

	for !P.currentTokenLS(token.RBRACE) && !P.currentTokenLS(token.EOF) {
		stmt := P.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		P.nextToken()
	}

	return block
}

// parseFunctionLiteral parses a function literal such as fn(x, y) {x + y}
func (P *Parser) parseFunctionLiteral() ast.Expression {
	literal := &ast.FunctionLiteral{Token: P.currentToken}
	if !P.expectPeek(token.LPAREN) {
		return nil
	}
	literal.Parameters = P.parseFunctionParameters()

	if !P.expectPeek(token.LBRACE) {
		return nil
	}

	literal.Body = P.parseBlockStatement()

	return literal

}

// parseFunctionParameters parses the parameters of a function
func (P *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if P.peekTokenLS(token.RPAREN) {
		P.nextToken()
		return identifiers
	}
	P.nextToken()

	ident := &ast.Identifier{Token: P.currentToken, Value: P.currentToken.Literal}
	identifiers = append(identifiers, ident)

	for P.peekTokenLS(token.COMMA) {
		P.nextToken()
		P.nextToken()
		ident := &ast.Identifier{
			Token: P.currentToken,
			Value: P.currentToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !P.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

// parseCallExpression parses a call expression such as add(5, 5) or add(5, 10)
func (P *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	// construct the call expression
	exp := &ast.CallExpression{Token: P.currentToken, Function: function}
	// parse the arguments
	exp.Arguments = P.parseExpressionList(token.RPAREN)
	// return the call expression
	return exp
}

/*
parseCallArguments parses the arguments of a call expression
*/
func (P *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	// no arguments
	if P.peekTokenLS(token.RPAREN) {
		P.nextToken()
		return args
	}

	// skip the LPAREN
	P.nextToken()
	args = append(args, P.parseExpression(LOWEST))

	// get args separated by commas
	for P.peekTokenLS(token.COMMA) {
		P.nextToken()
		P.nextToken()
		args = append(args, P.parseExpression(LOWEST))
	}

	// check for RPAREN
	if !P.expectPeek(token.RPAREN) {
		return nil
	}
	// return slice of arguments as expressions
	return args
}

// parseStringLiteral parses a string literal
func (P *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: P.currentToken, Value: P.currentToken.Literal}
}

// parseArrayLiteral parses an array literal
func (P *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: P.currentToken}
	array.Elements = P.parseExpressionList(token.RBRACKET)
	return array
}

// parseExpressionList parses a list of expressions
func (P *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if P.peekTokenLS(end) {
		P.nextToken()
		return list
	}

	P.nextToken()
	list = append(list, P.parseExpression(LOWEST))

	for P.peekTokenLS(token.COMMA) {
		P.nextToken()
		P.nextToken()
		list = append(list, P.parseExpression(LOWEST))
	}

	if !P.expectPeek(end) {
		return nil
	}
	return list
}

func (P *Parser) parseIndexExpression(leftExpression ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: P.currentToken, Left: leftExpression}

	P.nextToken()
	exp.Index = P.parseExpression(LOWEST)

	if !P.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
}

func (P *Parser) parseWhenLoopExpression() ast.Expression {
	expression := &ast.WhileLoopExpression{Token: P.currentToken}
	if !P.expectPeek(token.LPAREN) {
		return nil
	}

	P.nextToken()
	expression.Condition = P.parseExpression(LOWEST)
	if !P.expectPeek(token.RPAREN) {
		return nil
	}
	if !P.expectPeek(token.LBRACE) {
		return nil
	}
	expression.Consequence = P.parseBlockStatement()
	return expression
}

// parseHashLiteral parses a hash literal
// hash literals are of the form {key: value, key: value}
func (P *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Token: P.currentToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for !P.peekTokenLS(token.RBRACE) {
		P.nextToken()

		key := P.parseExpression(LOWEST)

		if !P.expectPeek(token.COLON) {
			return nil
		}

		P.nextToken()
		value := P.parseExpression(LOWEST)

		hash.Pairs[key] = value
		if !P.peekTokenLS(token.RBRACE) && !P.expectPeek(token.COMMA) {
			return nil
		}
	}
	if !P.expectPeek(token.RBRACE) {
		return nil
	}
	return hash
}

func (P *Parser) parseMethodCallExpression(obj ast.Expression) ast.Expression {
	method := &ast.ObjectCallExpression{Token: P.currentToken, Object: obj}
	P.nextToken()
	name := P.parseIdentifier()
	P.nextToken()
	method.Call = P.parseCallExpression(name)
	return method
}

func (P *Parser) parseImportExpression() ast.Expression {
	exp := &ast.ImportExpression{Token: P.currentToken}
	if !P.expectPeek(token.LPAREN) {
		return nil
	}
	P.nextToken()
	exp.Name = P.parseExpression(LOWEST)

	if !P.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}
