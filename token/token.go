package token

type Token interface {
	Text() string
}

type token string

func New(text string) Token {
	return token(text)
}

func (t token) Text() string   { return t }
func (t token) String() string { return t }
