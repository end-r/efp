package efp

import (
	"strings"
	"testing"
)

func TestLexerBasicOperators(t *testing.T) {
	SingleLexer(t, tknAssign, "=")
	SingleLexer(t, tknOpenBrace, "{")
	SingleLexer(t, tknCloseBrace, "}")
	SingleLexer(t, tknOpenSquare, "[")
	SingleLexer(t, tknCloseSquare, "]")
	SingleLexer(t, tknRequired, "!")
	SingleLexer(t, tknComma, ",")
	SingleLexer(t, tknOpenBracket, "(")
	SingleLexer(t, tknCloseBracket, ")")
}

func SingleLexer(t *testing.T, tkn tokenType, data string) {
	l := lex([]byte(data))
	assert(t, len(l.tokens) == 1, "Produced wrong number of tokens")
	assert(t, l.tokens[0].tkntype == tkn, "Wrong token type")
}

func TestLexerDuplicateOperators(t *testing.T) {
	MultiLexer(t, tknAssign, "=", 3)
	MultiLexer(t, tknOpenBrace, "{", 2)
	MultiLexer(t, tknCloseBrace, "}", 5)
	MultiLexer(t, tknOpenSquare, "[", 3)
	MultiLexer(t, tknCloseSquare, "]", 8)
	MultiLexer(t, tknRequired, "!", 1)
	MultiLexer(t, tknComma, ",", 11)
	MultiLexer(t, tknOpenBracket, "(", 9)
	MultiLexer(t, tknCloseBracket, ")", 9)
}

func MultiLexer(t *testing.T, tkn tokenType, data string, times int) {
	l := lex([]byte(strings.Repeat(data, times)))
	assert(t, len(l.tokens) == times, "Produced wrong number of tokens")
	for i := 0; i < times; i++ {
		assert(t, l.tokens[i].tkntype == tkn, "Wrong token type")
	}
}

func TestLexerValueLength(t *testing.T) {
	l := lex([]byte("hello this is dog"))
	assert(t, len(l.tokens) == 4, "wrong token number")
	expected := []int{5, 4, 2, 3}
	for i, tk := range l.tokens {
		assert(t, tk.end-tk.start == expected[i], "wrong token string")
	}
}

func TestLexerTokenString(t *testing.T) {
	l := lex([]byte("hello this is dog"))
	expected := []string{"hello", "this", "is", "dog"}
	for i, tk := range l.tokens {
		assert(t, l.tokenString(tk) == expected[i], "Wrong token")
	}
}

func TestLexerTokenLengths(t *testing.T) {
	l := lex([]byte("alias x : y = 5"))
	assert(t, len(l.tokens) == 6, "wrong token number")
}
