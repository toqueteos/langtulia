package lexer

import (
	"reflect"
	"strings"
	"testing"

	"github.com/toqueteos/langtulia/token"
)

var (
	LPAREN = token.New("(")
	RPAREN = token.New(")")
	A      = token.New("a")
	B      = token.New("b")
	X      = token.New("x")
	Y      = token.New("y")
	ZOO    = token.New("zoo")
	OR     = token.New("|")
	SPACE  = token.New("space")
)

// single function grammar just for testing all public API methods.
func lexText(l *Lexer) StateFn {
	switch r := l.Next(); {
	case r == EOF:
		l.Emit(TokenEOF)
		return nil
	case r == '\n':
		l.Ignore()
	case isSpace(r):
		return lexSpace
	case r == 'a':
		l.Emit(A)
	case r == 'b':
		l.Emit(B)
	case r == 'c':
		return l.Errorf("got c")
	case r == 'd':
		l.Ignore()
	case r == 'x':
		l.AcceptRun("x")
		l.Emit(X)
	case r == 'y':
		l.Accept("yY")
		l.Accept("yY")
		l.Emit(Y)
	case r == 'z':
		if strings.HasPrefix(l.Input(), "oo") {
			l.Next()
			l.Next()
			l.Emit(ZOO)
		} else {
			return l.Errorf("not zoo")
		}
	case r == '|':
		l.Emit(OR)
	case r == '(':
		l.Emit(LPAREN)
	case r == ')':
		l.Emit(RPAREN)
	}
	return lexText
}

func lexSpace(l *Lexer) StateFn {
	for isSpace(l.Peek()) {
		l.Next()
	}
	l.Emit(SPACE)
	return lexText
}

func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

var itemEOF = Token{TokenEOF, ""}

func item(t token.Token) Token { return Token{t, t.Text()} }
func itemErr(err string) Token { return Token{TokenError, err} }

func TestLexer(t *testing.T) {
	var (
		itemLPAREN = item(LPAREN)
		itemRPAREN = item(RPAREN)
		itemA      = item(A)
		itemB      = item(B)
		itemX      = Token{X, "xx"}
		itemY      = Token{Y, "yyy"}
		itemZOO    = item(ZOO)
		itemOR     = item(OR)
	)
	type TestCase struct {
		Input    string
		Expected []Token
	}
	tests := []TestCase{
		{"", []Token{itemEOF}},
		{"aa|b", []Token{itemA, itemA, itemOR, itemB, itemEOF}},
		{"a  a|b", []Token{itemA, Token{SPACE, "  "}, itemA, itemOR, itemB, itemEOF}},
		{"a(a|b)", []Token{itemA, itemLPAREN, itemA, itemOR, itemB, itemRPAREN, itemEOF}},
		{"abc", []Token{itemA, itemB, itemErr("got c")}},
		{"abd", []Token{itemA, itemB, itemEOF}},
		{"axb", []Token{itemA, Token{X, "x"}, itemB, itemEOF}},
		{"axxb", []Token{itemA, itemX, itemB, itemEOF}},
		{"ayyb", []Token{itemA, Token{Y, "yy"}, itemB, itemEOF}},
		{"ayyyb", []Token{itemA, itemY, itemB, itemEOF}},
		{"ayyyyyyb", []Token{itemA, itemY, itemY, itemB, itemEOF}},
		{"zob", []Token{itemErr("not zoo")}},
		{"zoob", []Token{itemZOO, itemB, itemEOF}},
	}

	for idx, tt := range tests {
		_, items := New(tt.Input, lexText)
		var output []Token
		for item := range items {
			output = append(output, item)
		}
		if !reflect.DeepEqual(output, tt.Expected) {
			t.Errorf("%d. Failed!\n\tExpected: %v\n\tGot: %v\n", idx+1, tt.Expected, output)
		}
	}
}

func TestItemString(t *testing.T) {
	type TestCase struct {
		Input    Token
		Expected string
	}
	tests := []TestCase{
		{item(A), `"a"`},
		{item(TokenEOF), "EOF"},
		{Token{TokenError, "some error"}, "some error"},
		{Token{token.New("long"), "test string output for long values"}, `"test strin"...`},
	}

	for idx, tt := range tests {
		output := tt.Input.String()
		if output != tt.Expected {
			t.Errorf("%d. Failed!\n\tExpected: %v\n\tGot: %v\n", idx+1, tt.Expected, output)
		}
	}
}
