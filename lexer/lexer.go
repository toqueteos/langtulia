// concurrent and extensible lexer.
//
// It is heavily based on text/template lexer:
// http://golang.org/src/text/template/parse/lex.go
//
// More info on Rob Pike's Lexical scanning in Go:
// https://www.youtube.com/watch?v=HxaD_trXwRE
package lexer

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/toqueteos/langtulia/token"
)

var (
	TokenError = token.New("err")
	TokenEOF   = token.New("EOF")
)

// Token contains a token and its value.
type Token struct {
	Token token.Token
	Value string
}

func (t Token) String() string {
	switch t.Token {
	case TokenError:
		return t.Value
	case TokenEOF:
		return "EOF"
	}
	if len(t.Value) > 10 {
		return fmt.Sprintf("%.10q...", t.Value)
	}
	return fmt.Sprintf("%q", t.Value)
}

// StateFn represents the state of the lexer as a function that returns the next
// state.
type StateFn func(*Lexer) StateFn

// EOF denotes there's no next rune
const EOF = -1

type Lexer struct {
	input  string     // the string being scanned.
	start  int        // start position of this token.
	pos    int        // current position in the input.
	width  int        // width of last rune read from input.
	tokens chan Token // channel of scanned tokens.
}

// New creates a new Lexer and starts the scanning process.
func New(input string, start StateFn) (*Lexer, chan Token) {
	l := &Lexer{
		input:  input,
		tokens: make(chan Token),
	}
	go l.run(start)
	return l, l.tokens
}

func (l *Lexer) run(start StateFn) {
	for state := start; state != nil; {
		state = state(l)
	}
	close(l.tokens) // No more tokens will be delivered.
}

// Emit passes a token back to the client.
func (l *Lexer) Emit(t token.Token) {
	l.tokens <- Token{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

// Next returns the next rune in the input.
func (l *Lexer) Next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return EOF
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return
}

// Ignore skips over the pending input before this point.
func (l *Lexer) Ignore() { l.start = l.pos }

// Backup steps back one rune. Can be called only once per call of `Lexer.Next`.
func (l *Lexer) Backup() { l.pos -= l.width }

// Peek returns but does not consume the next rune in the input.
func (l *Lexer) Peek() rune {
	r := l.Next()
	l.Backup()
	return r
}

// Accept consumes the next rune if it's from the valid set.
func (l *Lexer) Accept(valid string) bool {
	if strings.IndexRune(valid, l.Next()) >= 0 {
		return true
	}
	l.Backup()
	return false
}

// AcceptRun consumes a run of runes from the valid set.
func (l *Lexer) AcceptRun(valid string) {
	for strings.IndexRune(valid, l.Next()) >= 0 {
	}
	l.Backup()
}

// Errorf returns an error token and terminates the scan by passing back a nil
// pointer that will be the next state, terminating l.Run.
func (l *Lexer) Errorf(format string, args ...interface{}) StateFn {
	l.tokens <- Token{TokenError, fmt.Sprintf(format, args...)}
	return nil
}

// Input returns the input left to read, if there's any.
func (l *Lexer) Input() string {
	return l.input[l.pos:]
}
