package lexer

import (
	"monkey/token"
)

type Lexer struct {
	input        string
	position     int // 入力における現在位置(現在の文字)
	readPosition int // これから読み込む位置（現在の文字の次）
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		// 最後まで読んだか終端チェック。終端は0とする 0はASCIIのNULLに対応しているので0を使う。
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF

	default:
		if isLitter(l.ch) { //英数字かどうかを確認して、Tokenになるか判定
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok // readIdentifierでreadChar()を必要分読み切っているため、先にreturn
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) readIdentifier() string {
	// let a = のlからtまでを読み進める。 tの隣は空白,つまりisLetter==Falseにすることで
	// 字句の切れ目を検出している
	position := l.position
	for isLitter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}
func isLitter(ch byte) bool {
	// この企画では英数字だけを対象にする。アンダースコアは例外的にこの言語で許可した範囲とする
	// foo_barみたいな変数が作れるようになる。
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}
