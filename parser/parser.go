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
	ASSIGN
	EQUALS
	ANDOR
	LESSGREATER
	SUM
	PRODUCT
	POWER
	MODULUS
	PREFIX
	CALL
	INDEX
	HIGHEST
)

// precedence map to determine the precedence of the operators
var precedence = map[token.TokenType]int{
	token.BIND:        ASSIGN,
	token.EQ:          EQUALS,
	token.NOT_EQ:      EQUALS,
	token.LT:          LESSGREATER,
	token.GT:          LESSGREATER,
	token.LT_EQ:       LESSGREATER,
	token.GT_EQ:       LESSGREATER,
	token.AND:         ANDOR,
	token.OR:          ANDOR,
	token.PLUS:        SUM,
	token.PLUS_EQ:     SUM,
	token.MINUS:       SUM,
	token.MINUS_EQ:    SUM,
	token.MOD:         MODULUS,
	token.SLASH:       PRODUCT,
	token.ASTERISK:    PRODUCT,
	token.ASTERISK_EQ: PRODUCT,
	token.LPAREN:      CALL,
	token.PERIOD:      CALL,
	token.LBRACKET:    INDEX,
	token.DOUBLECOL:   INDEX,
}

// Parser is the core struct for the parser
type Parser struct {
	L               *lexer.Lexer                      // lexer
	currentToken    token.Token                       // current token
	peekToken       token.Token                       // next token
	previousToken   token.Token                       // previous token used for `--` and `++` operators
	errors          []string                          // errors
	prefixParseFns  map[token.TokenType]prefixParseFn // prefix parse functions
	infixParseFns   map[token.TokenType]infixParseFn  // infix parse functions
	postfixParseFns map[token.TokenType]postfixParseFn
}

type (
	prefixParseFn  func() ast.Expression               // prefix parse function
	infixParseFn   func(ast.Expression) ast.Expression // infix parse function
	postfixParseFn func() ast.Expression
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
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.DEF_FN, p.parseFunctionDefinition)
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
	p.registerInfix(token.PLUS_EQ, p.parseAssignExpression)
	p.registerInfix(token.MINUS_EQ, p.parseAssignExpression)
	p.registerInfix(token.ASSIGN, p.parseAssignExpression)
	p.registerInfix(token.ASTERISK_EQ, p.parseAssignExpression)
	p.registerInfix(token.ASSIGN, p.parseAssignExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LT_EQ, p.parseInfixExpression)
	p.registerInfix(token.GT_EQ, p.parseInfixExpression)
	p.registerInfix(token.MOD, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	p.registerInfix(token.PERIOD, p.parseMethodCallExpression)
	p.registerInfix(token.DOUBLECOL, p.parseSelectorExpression)
	p.registerInfix(token.BIND, p.parseBindExpression)

	p.postfixParseFns = make(map[token.TokenType]postfixParseFn)
	p.registerPostfix(token.PLUS_PLUS, p.parsePostfixExpression)
	p.registerPostfix(token.MINUS_MINUS, p.parsePostfixExpression)
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
	msg := fmt.Sprintf("%s: Line %v Column %v - expected next token to be %s got %s", P.currentToken.FileName, P.peekToken.Line, P.peekToken.Column, t, P.peekToken.Type)
	P.errors = append(P.errors, msg)
}

/*
nextToken moves to the next token
*/
func (P *Parser) nextToken() {
	P.previousToken = P.currentToken
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

func (p *Parser) registerPostfix(tokenType token.TokenType, fn postfixParseFn) {
	p.postfixParseFns[tokenType] = fn
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
		msg := fmt.Sprintf("%s Line %v Column %v - could not parse %q as integer", P.currentToken.FileName, P.currentToken.Line, P.currentToken.Column, P.currentToken)
		P.errors = append(P.errors, msg)
		return nil
	}
	lit.Value = value
	return lit

}

// noPrefixParseFnError returns an error message when no prefix parse function is found
func (P *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("%s Line %v Column %v - no prefix parse function for %s found", P.currentToken.FileName, P.currentToken.Line, P.currentToken.Column, t)
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
func (p *Parser) parseExpression(precedence int) ast.Expression {
	// defer untrace(trace("parseExpression"))
	postfix := p.postfixParseFns[p.currentToken.Type]
	if postfix != nil {
		return (postfix())
	}
	prefix := p.prefixParseFns[p.currentToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.currentToken.Type)
		return nil
	}
	leftExp := prefix()
	for !p.peekTokenLS(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
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

func (P *Parser) parseSelectorExpression(exp ast.Expression) ast.Expression {
	P.expectPeek(token.IDENT)
	index := &ast.StringLiteral{Token: P.currentToken, Value: P.currentToken.Literal}
	return &ast.IndexExpression{Left: exp, Index: index, Token: P.currentToken}
}

func (P *Parser) parseBindExpression(exp ast.Expression) ast.Expression {
	switch node := exp.(type) {
	case *ast.Identifier:
	default:
		msg := fmt.Sprintf("expected identifier expression on left but got %T %#v", node, exp)
		P.errors = append(P.errors, msg)
		return nil
	}
	be := &ast.BindExpression{Token: P.currentToken, Left: exp}

	P.nextToken()

	be.Value = P.parseExpression(LOWEST)

	// Correctly bind the function literal to its name so that self-recursive
	// // functions work. This is used by the compiler to emit LoadSelf so a ref
	// // to the current function is available.
	// if fl, ok := be.Value.(*ast.FunctionLiteral); ok {
	// 	ident := be.Left.(*ast.Identifier)
	// 	fl.Name = ident.Value
	// }

	return be
}

func (P *Parser) parseFunctionDefinition() ast.Expression {
	P.nextToken()
	lit := &ast.FunctionDefineLiteral{Token: P.currentToken}
	if !P.expectPeek(token.LPAREN) {
		return nil
	}
	lit.Defaults, lit.Parameters = P.parseFunctionDefParameter()
	if !P.expectPeek(token.LBRACE) {
		return nil
	}
	lit.Body = P.parseBlockStatement()
	return lit
}

// parseFunctionParameters parses the parameters used for a function.
func (P *Parser) parseFunctionDefParameter() (map[string]ast.Expression, []*ast.Identifier) {

	// Any default parameters.
	m := make(map[string]ast.Expression)

	// The argument-definitions.
	identifiers := make([]*ast.Identifier, 0)

	// Is the next parameter ")" ?  If so we're done. No args.
	if P.peekTokenLS(token.RPAREN) {
		P.nextToken()
		return m, identifiers
	}
	P.nextToken()

	// Keep going until we find a ")"
	for !P.currentTokenLS(token.RPAREN) {

		if P.currentTokenLS(token.EOF) {
			P.errors = append(P.errors, "unterminated function parameters")
			return nil, nil
		}

		// Get the identifier.
		ident := &ast.Identifier{Token: P.currentToken, Value: P.currentToken.Literal}
		identifiers = append(identifiers, ident)
		P.nextToken()

		// If there is "=xx" after the name then that's
		// the default parameter.
		if P.currentTokenLS(token.ASSIGN) {
			P.nextToken()
			// Save the default value.
			m[ident.Value] = P.parseExpressionStatement().Expression
			P.nextToken()
		}

		// Skip any comma.
		if P.currentTokenLS(token.COMMA) {
			P.nextToken()
		}
	}

	return m, identifiers
}

func (P *Parser) parseFloatLiteral() ast.Expression {
	flo := &ast.FloatLiteral{Token: P.currentToken}
	value, err := strconv.ParseFloat(P.currentToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as float around line", P.currentToken.Literal)
		P.errors = append(P.errors, msg)
		return nil
	}
	flo.Value = value
	return flo
}

func (P *Parser) parseAssignExpression(name ast.Expression) ast.Expression {
	stmt := &ast.AssignStatement{Token: P.currentToken}
	if n, ok := name.(*ast.Identifier); ok {
		stmt.Name = n
	} else {
		msg := "expected assign token to be IDENT, got null instead"

		if name != nil {
			msg = fmt.Sprintf("%s Line %v Column %v - expected assign token to be IDENT, got %s", P.currentToken.FileName, P.currentToken.Line, P.currentToken.Column, name.TokenLiteral())
		}
		P.errors = append(P.errors, msg)
	}

	oper := P.currentToken
	P.nextToken()
	switch oper.Type {
	case token.PLUS_EQ:
		stmt.Operator = "+="
	case token.MINUS_EQ:
		stmt.Operator = "-="
	case token.ASTERISK_EQ:
		stmt.Operator = "*="
	default:
		stmt.Operator = "="
	}
	stmt.Value = P.parseExpression(LOWEST)
	return stmt
}

func (p *Parser) parsePostfixExpression() ast.Expression {
	expression := &ast.PostfixExpression{
		Token:    p.previousToken,
		Operator: p.currentToken.Literal,
	}
	return expression
}
