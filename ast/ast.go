/*
Package ast implements the Abstract Syntax Tree for the Monkey programming language.
*/

package ast

import (
	"bytes"
	"esolang/lang-esolang/token"
	"fmt"
	"strings"
)

// Node is the interface that all nodes in the AST implement.
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement is the interface that all statement nodes in the AST implement.
type Statement interface {
	Node
	statementNode()
}

// Expression is the interface that all expression nodes in the AST implement.
type Expression interface {
	Node
	expressionNode()
}

/*
Identifier represents an identifier.
For example, `x`.
*/
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

/*
LetStatement represents a let statement.
For example, `let x = 5;`.
*/
type LetStatement struct {
	Token token.Token // token.LET
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())

	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

/*
ReturnStatement represents a return statement.
For example, `return 5;`.
*/
type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

/*
ExpressionStatement represents an expression statement.
For example, `5 + 5;`.
*/
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

/*
IntegerLiteral represents an integer literal.
For example, `5`.
*/
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// FloatLiteral holds a floating-point number
// For example, `5.0`.
type FloatLiteral struct {
	Token token.Token

	Value float64
}

func (fl *FloatLiteral) expressionNode() {}

func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }

func (fl *FloatLiteral) String() string { return fl.Token.Literal }

/*
PrefixExpression represents a prefix expression.
For example, the `-` in `-5`.
*/
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

/*
InfixExpression represents an infix expression.
For example, the `+` in `5 + 5`.
*/
type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

/*
Boolean represents a boolean value.
`true` or `false`.
*/
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

/*
BlockStatement represents a block statement.

	A block statement is a sequence of statements enclosed in braces.
*/
type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

/*
IfExpression represents an if expression.
For example, `if (x < y) { x } else { y }`.
*/
type IfExpression struct {
	Token       token.Token
	Condition   Expression      // The condition to be evaluated
	Consequence *BlockStatement // The block statement to be executed if the condition is true
	Alternative *BlockStatement // The block statement to be executed if the condition is false
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

/*
FunctionLiteral represents a function literal.

	For example, `fn(x, y) { x + y; }`.
	fn is the actual token literal.
	x and y are the parameters - an array of identifiers.
	x + y; is the body - a block statement.
*/
type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}

	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(")")
	out.WriteString(fl.Body.String())

	return out.String()
}

/*
CallExpression represents a call expression.

	 <expression>(<expression>, <expression>, ...)
	 For example, `add(1, 2)`
	`reduce(2 + 2, 4 * 4 + 1)`
*/
type CallExpression struct {
	Token     token.Token // the `(` token
	Function  Expression  // tokenLiteral || identifier
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, arg := range ce.Arguments {
		args = append(args, arg.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ","))
	out.WriteString(")")
	return out.String()
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

type WhileLoopExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
}

func (fle WhileLoopExpression) expressionNode()      {}
func (fle WhileLoopExpression) TokenLiteral() string { return fle.Token.Literal }
func (fle WhileLoopExpression) String() string {
	var output bytes.Buffer
	output.WriteString("while (")
	output.WriteString(fle.Condition.String())
	output.WriteString(") {")
	output.WriteString(fle.Consequence.String())
	output.WriteString("}")
	return output.String()
}

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}

	for _, element := range al.Elements {
		elements = append(elements, element.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var output bytes.Buffer
	output.WriteString("(")
	output.WriteString(ie.Left.String())
	output.WriteString("[")
	output.WriteString(ie.Index.String())
	output.WriteString("])")

	return output.String()
}

type HashLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode()      {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) String() string {
	var output bytes.Buffer

	pairs := []string{}
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}
	output.WriteString("{")
	output.WriteString(strings.Join(pairs, ", "))
	output.WriteString("}")

	return output.String()
}

type ObjectCallExpression struct {
	Token  token.Token
	Object Expression
	Call   Expression
}

func (oce *ObjectCallExpression) expressionNode()      {}
func (oce *ObjectCallExpression) TokenLiteral() string { return oce.Token.Literal }
func (oce *ObjectCallExpression) String() string {
	var out bytes.Buffer

	out.WriteString(oce.Object.String())
	out.WriteString(".")
	out.WriteString(oce.Call.String())

	return out.String()
}

// ImportExpression represents an `import` expression and holds the name
// of the module being imported.
type ImportExpression struct {
	Token token.Token // The 'import' token
	Name  Expression
}

func (ie *ImportExpression) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (ie *ImportExpression) TokenLiteral() string { return ie.Token.Literal }

// String returns a stringified version of the AST for debugging
func (ie *ImportExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ie.TokenLiteral())
	out.WriteString("(")
	out.WriteString(fmt.Sprintf("\"%s\"", ie.Name))
	out.WriteString(")")

	return out.String()
}

type BindExpression struct {
	Token token.Token // The := token
	Left  Expression
	Value Expression
}

func (be *BindExpression) expressionNode() {}

// TokenLiteral prints the literal value of the token associated with this node
func (be *BindExpression) TokenLiteral() string { return be.Token.Literal }

// String returns a stringified version of the AST for debugging
func (be *BindExpression) String() string {
	var out bytes.Buffer

	out.WriteString(be.Left.String())
	out.WriteString(be.TokenLiteral())
	out.WriteString(be.Value.String())

	return out.String()
}

type FunctionDefineLiteral struct {
	// Token holds the token
	Token token.Token

	// Paremeters holds the function parameters.
	Parameters []*Identifier

	// Defaults holds any default-arguments.
	Defaults map[string]Expression

	// Body holds the set of statements in the functions' body.
	Body *BlockStatement
}

func (fl *FunctionDefineLiteral) expressionNode() {}

// TokenLiteral returns the literal token.
func (fl *FunctionDefineLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

// String returns this object as a string.
func (fl *FunctionDefineLiteral) String() string {
	var out bytes.Buffer
	params := make([]string, 0)
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())
	return out.String()

}

type PostfixExpression struct {
	// Token holds the token we're operating upon
	Token token.Token
	// Operator holds the postfix token, e.g. ++
	Operator string
}

func (pe *PostfixExpression) expressionNode() {}

// TokenLiteral returns the literal token.
func (pe *PostfixExpression) TokenLiteral() string { return pe.Token.Literal }

// String returns this object as a string.
func (pe *PostfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Token.Literal)
	out.WriteString(pe.Operator)
	out.WriteString(")")
	return out.String()
}

type AssignStatement struct {
	Token    token.Token
	Name     *Identifier
	Operator string
	Value    Expression
}

func (as *AssignStatement) expressionNode() {}

// TokenLiteral returns the literal token.
func (as *AssignStatement) TokenLiteral() string { return as.Token.Literal }

// String returns this object as a string.
func (as *AssignStatement) String() string {
	var out bytes.Buffer
	out.WriteString(as.Name.String())
	out.WriteString(as.Operator)
	out.WriteString(as.Value.String())
	return out.String()
}
