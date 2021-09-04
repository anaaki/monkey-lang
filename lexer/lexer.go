package lexer

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
