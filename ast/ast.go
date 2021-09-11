package ast

import "monkey/token"

type Node interface {
	TokenLiteral() string //デバッグ用にリテラルを返す。
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

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type LetStatement struct {
	Token token.Token // token.Let トークン
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

type ReturnStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *ReturnStatement) statementNode()       {}
func (ls *ReturnStatement) TokenLiteral() string { return ls.Token.Literal }
