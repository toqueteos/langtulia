// token defines an extensible token representation.
package token

type Token interface {
	Text() string
}

type token string

// New returns string based Token which should be enough for most usages.
func New(text string) Token {
	return token(text)
}

func (t token) Text() string   { return string(t) }
func (t token) String() string { return string(t) }
