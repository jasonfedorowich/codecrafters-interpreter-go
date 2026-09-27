package parse

import (
	"github.com/codecrafters-io/interpreter-starter-go/app/internal/error_code"
	"github.com/codecrafters-io/interpreter-starter-go/app/internal/scanner"
)

type Parser struct {
	tokens []*scanner.Token
	p      int
	eh     *error_code.ErrorHandler
}

func NewParser(tokens []*scanner.Token) *Parser {
	return &Parser{
		tokens: tokens,
		p:      0,
		eh:     &error_code.ErrorHandler{},
	}
}

func (p *Parser) Parse() (Expr, int) {
	return p.expression(), p.eh.ErrorCode
}

func (p *Parser) expression() Expr {
	return p.equality()
}

func (p *Parser) equality() Expr {
	expr := p.comparison()
	if p.eh.InError() {
		return nil
	}
	for p.match(scanner.BANG_EQUAL, scanner.EQUAL_EQUAL) {
		op := p.prev()
		right := p.comparison()
		if p.eh.InError() {
			return nil
		}
		expr = &BinaryExpr{
			Left:     expr,
			Right:    right,
			Operator: op,
		}
	}
	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()
	if p.eh.InError() {
		return nil
	}
	for p.match(scanner.GREATER, scanner.GREATER_EQUAL, scanner.LESS, scanner.LESS_EQUAL) {
		op := p.prev()
		right := p.term()
		if p.eh.InError() {
			return nil
		}
		expr = &BinaryExpr{
			Left:     expr,
			Right:    right,
			Operator: op,
		}
	}
	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()
	if p.eh.InError() {
		return nil
	}

	for p.match(scanner.PLUS, scanner.MINUS) {
		op := p.prev()
		right := p.factor()
		if p.eh.InError() {
			return nil
		}
		expr = &BinaryExpr{
			Left:     expr,
			Right:    right,
			Operator: op,
		}
	}

	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()
	if p.eh.InError() {
		return nil
	}

	for p.match(scanner.STAR, scanner.SLASH) {
		op := p.prev()
		right := p.unary()
		if p.eh.InError() {
			return nil
		}
		expr = &BinaryExpr{
			Left:     expr,
			Right:    right,
			Operator: op,
		}
	}

	return expr
}

func (p *Parser) unary() Expr {
	if p.match(scanner.BANG, scanner.MINUS) {
		pr := p.prev()
		right := p.unary()
		if p.eh.InError() {
			return nil
		}
		return &UnaryExpr{
			Right:    right,
			Operator: pr,
		}
	}
	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(scanner.NIL) {
		return &LiteralExpr{
			Obj: nil,
		}
	}
	if p.match(scanner.TRUE) {
		return &LiteralExpr{
			Obj: true,
		}
	}

	if p.match(scanner.FALSE) {
		return &LiteralExpr{
			Obj: false,
		}
	}
	if p.match(scanner.STRING, scanner.NUMBER) {
		return &LiteralExpr{
			Obj: p.prev().Literal.S,
		}
	}

	if p.match(scanner.LEFT_PAREN) {
		expr := p.expression()
		if p.eh.InError() {
			return nil
		}
		if !p.match(scanner.RIGHT_PAREN) {
			p.eh.Error(65, "[line %d] Error at ')': Expect expression", p.peek().Line)
			return nil
		}
		return &GroupingExpr{
			E: expr,
		}
	}
	p.eh.Error(65, "[line %d] Unexpected token", p.peek().Line)
	return nil
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
