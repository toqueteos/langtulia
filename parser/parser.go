// extensible Pratt parser implementation.
package parser

import (
	"log"

	"github.com/toqueteos/langtulia/expr"
	"github.com/toqueteos/langtulia/token"
)

// InfixParser is associated with a token that appears in the middle of the
// expression it parses.
type InfixParser interface {
	// Parse will be called after the left-hand side has been parsed, and it is
	// responsible for parsing everything that comes after the token.
	//
	// It is also used for postfix expressions, where doesn't consume any more
	// tokens in Parse().
	Parse(*Parser, expr.Expression, token.Token) expr.Expression
	Precedence() int
}

// PrefixParser is associated with a token that appears at the beginning of an
// expression.
type PrefixParser interface {
	// Parse will be called with the consumed leading token. It is **NOT
	// RESPONSIBLE** for parsing anything that comes after that token.
	//
	// It is also used for single-token expressions like variables, where
	// Parse() doesn't consume any more tokens.
	Parse(*Parser, token.Token) expr.Expression
}

// Parser defines an extensible Pratt Parser for parsing any context-free
// grammar.
type Parser struct {
	tokens        chan token.Token
	read          []token.Token
	prefixParsers map[token.Token]PrefixParser
	infixParsers  map[token.Token]InfixParser
}

func NewParser(tokens chan token.Token) *Parser {
	return &Parser{
		tokens:        tokens,
		prefixParsers: make(map[token.Token]PrefixParser),
		infixParsers:  make(map[token.Token]InfixParser),
	}
}

func (p *Parser) RegisterPrefix(t token.Token, pp PrefixParser) {
	p.prefixParsers[t] = pp
}

func (p *Parser) RegisterInfix(t token.Token, ip InfixParser) {
	p.infixParsers[t] = ip
}

func (p *Parser) ParseExpression(precedence int) expr.Expression {
	token := p.consume0()
	prefix := p.prefixParsers[token]

	if prefix == nil {
		log.Fatalf("parser: could not parse %q.", token)
	}

	left := prefix.Parse(p, token)

	for precedence < p.Precedence() {
		token = p.consume0()
		infix := p.infixParsers[token]
		left = infix.Parse(p, left, token)
	}

	return left
}

func (p *Parser) Match(expected token.Token) bool {
	token := p.lookAhead(0)
	if token != expected {
		return false
	}

	p.consume0()
	return true
}

func (p *Parser) Consume(expected token.Token) token.Token {
	token := p.lookAhead(0)
	if token != expected {
		log.Fatalf("parser: expected token %q but found %q", expected, token)
	}
	return p.consume0()
}

func (p *Parser) consume0() token.Token {
	// Make sure we've read the token.
	p.lookAhead(0)

	tok := p.read[0]
	p.read = p.read[1:]
	return tok
}

func (p *Parser) lookAhead(distance int) token.Token {
	// Read in as many as needed.
	for distance >= len(p.read) {
		p.read = append(p.read, <-p.tokens)
	}

	// Get the queued token.
	return p.read[distance]
}

func (p *Parser) Precedence() int {
	if parselet, ok := p.infixParsers[p.lookAhead(0)]; ok {
		return parselet.Precedence()
	}

	return 0
}
