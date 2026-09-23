package parse

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/app/internal/scanner"
)

type Parser struct {
	tokens []*scanner.Token
	p      int
}

func NewParser(tokens []*scanner.Token) *Parser {
	return &Parser{
		tokens: tokens,
		p:      0,
	}
}

func (p *Parser) Parse() (Expr, error) {
	return p.expression()
}

func (p *Parser) expression() (Expr, error) {
	return p.primary()
}

func (p *Parser) unary() (Expr, error) {
	if p.match(scanner.BANG, scanner.MINUS) {
		pr := p.prev()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &UnaryExpr{
			Right:    right,
			Operator: pr,
		}, nil
	}
	return p.primary()
}

func (p *Parser) primary() (Expr, error) {
	if p.match(scanner.NIL) {
		return &LiteralExpr{
			Obj: nil,
		}, nil
	}
	if p.match(scanner.TRUE) {
		return &LiteralExpr{
			Obj: true,
		}, nil
	}

	if p.match(scanner.FALSE) {
		return &LiteralExpr{
			Obj: false,
		}, nil
	}
	if p.match(scanner.STRING, scanner.NUMBER) {
		return &LiteralExpr{
			Obj: p.prev().Literal,
		}, nil
	}

	if p.match(scanner.LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if !p.match(scanner.RIGHT_PAREN) {
			return nil, fmt.Errorf("unexpected token, %s", p.prev().Literal.String())
		}
		return &GroupingExpr{
			E: expr,
		}, nil
	}
	return nil, fmt.Errorf("unexpected token, %s", p.prev().Literal.String())
}

func (p *Parser) prev() *scanner.Token {
	return p.tokens[p.p-1]
}

func (p *Parser) peek() *scanner.Token {
	return p.tokens[p.p]
}

func (p *Parser) match(tokenType ...scanner.TokenType) bool {
	for _, tt := range tokenType {
		if p.check(tt) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(tokenType scanner.TokenType) bool {
	if len(p.tokens) <= p.p {
		return false
	}
	return p.tokens[p.p].Type == tokenType
}

func (p *Parser) advance() {
	p.p++
}
