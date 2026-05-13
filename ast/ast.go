// Created and Developed by Farhan Kertadiwangsa
package ast

import (
	"bytes"
	"strings"

	"ky/token"
)

// Node adalah antarmuka dasar untuk semua node AST
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement adalah node yang tidak menghasilkan nilai
type Statement interface {
	Node
	statementNode()
}

// Expression adalah node yang menghasilkan nilai
type Expression interface {
	Node
	expressionNode()
}


type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// Statements
// BuatStatement: buat x = 10
type BuatStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (b *BuatStatement) statementNode()       {}
func (b *BuatStatement) TokenLiteral() string { return b.Token.Literal }
func (b *BuatStatement) String() string {
	return "buat " + b.Name.String() + " = " + b.Value.String()
}

// AssignStatement: x = 20
type AssignStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (a *AssignStatement) statementNode()       {}
func (a *AssignStatement) TokenLiteral() string { return a.Token.Literal }
func (a *AssignStatement) String() string {
	return a.Name.String() + " = " + a.Value.String()
}

// BalikStatement: balik <ekspresi>
type BalikStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (b *BalikStatement) statementNode()       {}
func (b *BalikStatement) TokenLiteral() string { return b.Token.Literal }
func (b *BalikStatement) String() string {
	if b.ReturnValue != nil {
		return "balik " + b.ReturnValue.String()
	}
	return "balik"
}

// ExpressionStatement: ekspresi sebagai pernyataan
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (e *ExpressionStatement) statementNode()       {}
func (e *ExpressionStatement) TokenLiteral() string { return e.Token.Literal }
func (e *ExpressionStatement) String() string {
	if e.Expression != nil {
		return e.Expression.String()
	}
	return ""
}

// BlockStatement: { ... }
type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (b *BlockStatement) statementNode()       {}
func (b *BlockStatement) TokenLiteral() string { return b.Token.Literal }
func (b *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range b.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// HentiStatement: henti (break)
type HentiStatement struct {
	Token token.Token
}

func (h *HentiStatement) statementNode()       {}
func (h *HentiStatement) TokenLiteral() string { return h.Token.Literal }
func (h *HentiStatement) String() string       { return "henti" }

// LanjutStatement: lanjut (continue)
type LanjutStatement struct {
	Token token.Token
}

func (l *LanjutStatement) statementNode()       {}
func (l *LanjutStatement) TokenLiteral() string { return l.Token.Literal }
func (l *LanjutStatement) String() string       { return "lanjut" }

// Expresi

// Identifier: nama variabel atau fungsi
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// IntegerLiteral: 42
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// FloatLiteral: 3.14
type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }

// StringLiteral: "halo"
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return `"` + sl.Value + `"` }

// BooleanLiteral: benar / salah
type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (bl *BooleanLiteral) expressionNode()      {}
func (bl *BooleanLiteral) TokenLiteral() string { return bl.Token.Literal }
func (bl *BooleanLiteral) String() string       { return bl.Token.Literal }

// KosongLiteral: kosong (null)
type KosongLiteral struct {
	Token token.Token
}

func (kl *KosongLiteral) expressionNode()      {}
func (kl *KosongLiteral) TokenLiteral() string { return kl.Token.Literal }
func (kl *KosongLiteral) String() string       { return "kosong" }

// PrefixExpression: -x, !benar, tidak salah
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	return "(" + pe.Operator + pe.Right.String() + ")"
}

// InfixExpression: a + b, x == y, a dan b
type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	return "(" + ie.Left.String() + " " + ie.Operator + " " + ie.Right.String() + ")"
}

// JikaExpression: jika <kondisi> { ... } lain { ... }
type JikaExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (je *JikaExpression) expressionNode()      {}
func (je *JikaExpression) TokenLiteral() string { return je.Token.Literal }
func (je *JikaExpression) String() string {
	var out bytes.Buffer
	out.WriteString("jika ")
	out.WriteString(je.Condition.String())
	out.WriteString(" { ")
	out.WriteString(je.Consequence.String())
	out.WriteString(" }")
	if je.Alternative != nil {
		out.WriteString(" lain { ")
		out.WriteString(je.Alternative.String())
		out.WriteString(" }")
	}
	return out.String()
}

// SelamaExpression: selama <kondisi> { ... }
type SelamaExpression struct {
	Token     token.Token
	Condition Expression
	Body      *BlockStatement
}

func (se *SelamaExpression) expressionNode()      {}
func (se *SelamaExpression) TokenLiteral() string { return se.Token.Literal }
func (se *SelamaExpression) String() string {
	return "selama " + se.Condition.String() + " { " + se.Body.String() + " }"
}

// UlangExpression: ulang buat i = 0; i < 10; i = i + 1 { ... }
type UlangExpression struct {
	Token     token.Token
	Init      Statement
	Condition Expression
	Post      Statement
	Body      *BlockStatement
}

func (ue *UlangExpression) expressionNode()      {}
func (ue *UlangExpression) TokenLiteral() string { return ue.Token.Literal }
func (ue *UlangExpression) String() string {
	return "ulang { " + ue.Body.String() + " }"
}

// FungsiLiteral: fungsi(a, b) { ... }
type FungsiLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
	Name       string
}

func (fl *FungsiLiteral) expressionNode()      {}
func (fl *FungsiLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FungsiLiteral) String() string {
	params := make([]string, len(fl.Parameters))
	for i, p := range fl.Parameters {
		params[i] = p.String()
	}
	name := fl.Name
	if name == "" {
		name = "anonim"
	}
	return "fungsi " + name + "(" + strings.Join(params, ", ") + ") { " + fl.Body.String() + " }"
}

// CallExpression: fungsi(argumen1, argumen2)
type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	args := make([]string, len(ce.Arguments))
	for i, a := range ce.Arguments {
		args[i] = a.String()
	}
	return ce.Function.String() + "(" + strings.Join(args, ", ") + ")"
}

// ArrayLiteral: [1, 2, 3]
type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	elems := make([]string, len(al.Elements))
	for i, e := range al.Elements {
		elems[i] = e.String()
	}
	return "[" + strings.Join(elems, ", ") + "]"
}

// IndexExpression: array[0]
type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	return ie.Left.String() + "[" + ie.Index.String() + "]"
}

// MapLiteral: {"kunci": nilai}
type MapLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (ml *MapLiteral) expressionNode()      {}
func (ml *MapLiteral) TokenLiteral() string { return ml.Token.Literal }
func (ml *MapLiteral) String() string {
	var pairs []string
	for k, v := range ml.Pairs {
		pairs = append(pairs, k.String()+": "+v.String())
	}
	return "{" + strings.Join(pairs, ", ") + "}"
}
