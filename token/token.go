package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// 識別子 + リテラル
	IDENT = "IDENT" // add, foobar, x, y
	INT   = "INT"   // 1343456

	// 演算子
	ASSIGN = "="
	PLUS   = "+"

	//デリミタ
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	LBRACE = "{"
	RPAREN = ")"
	RBRACE = "}"

	// キーワード
	FUNCTION = "fn"
	LET      = "let"
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
