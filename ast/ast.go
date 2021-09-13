package ast

import (
	"bytes"
	"monkey/token"
)

type Node interface {
	TokenLiteral() string //デバッグ用にリテラルを返す。
	String() string       //これもデバッグよう。ASTを表示したり、他のASTと比較したりする。
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	// 構文解析器が生成する全てのASTのルートになる
	Statements []Statement
}

type Identifier struct {
	Token token.Token
	Value string
}

// Expressionインターフェイスを満たすように実装。
// 識別子は文(identifire)なのでExpressinでないように字面は見える。
// これは単に実装のわかりやすさのため。Monkeyの他の部分は識別子は値を生成”する”。
// let x = valueProducingIdentifier;
func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

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

type LetStatement struct {
	Token token.Token // token.Let トークン
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
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

type ReturnStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *ReturnStatement) statementNode()       {}
func (ls *ReturnStatement) TokenLiteral() string { return ls.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.Value != nil { // rs.ReturnValue などないので、Valueとした
		out.WriteString(rs.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

type ExpressionStatement struct { // x + 10のようなステートメント
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

type IntegerLiteral struct {
	Token token.Token
	Value int64 // プログラム内では"5"のように実際の数字を入れず、int64とする。
}

func (li *IntegerLiteral) expressionNode()      {}
func (li *IntegerLiteral) TokenLiteral() string { return li.Token.Literal }
func (li *IntegerLiteral) String() string       { return li.Token.Literal }
