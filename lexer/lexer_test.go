package lexer

import (
	"monkey/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `=+(){},;`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBLACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ":"},
		{token.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] token type wrong expexted=%q got=%q",
				i, tt.expectedType, tt.expectedType)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("test[%d] literal wrong expected=%q got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
