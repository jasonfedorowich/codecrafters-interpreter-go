package parse

import (
	"github.com/codecrafters-io/interpreter-starter-go/app/internal/error_code"
	"github.com/codecrafters-io/interpreter-starter-go/app/internal/scanner"
)

type Parser struct {
	tokens    []*scanner.Token
	p         int
	eh        *error_code.ErrorHandler
	parseExpr bool
}

func NewParser(tokens []*scanner.Token, parseExpr bool) *Parser {
	return &Parser{
		tokens:    tokens,
		p:         0,
		eh:        &error_code.ErrorHandler{},
		parseExpr: parseExpr,
	}
}

func (p *Parser) Parse() ([]Statement, int) {
	var statements []Statement
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	return statements, p.eh.ErrorCode
}

func (p *Parser) declaration() Statement {
	if p.match(scanner.VAR) {
		return p.varDeclaration()
	} else if p.match(scanner.FUN) {
		return p.funcDeclaration()
	}
	return p.statement()
}

func (p *Parser) varDeclaration() Statement {
	var token = p.consume(scanner.IDENTIFIER, "Expect variable name.")
	var expr Expr
	if p.match(scanner.EQUAL) {
		expr = p.expression()
	}
	if !p.match(scanner.SEMICOLON) {
		p.eh.Error(65, "[line %d] Expect ';' after variable declaration.", p.peek().Line)
		return nil
	}
	return &VarStatement{Name: token,
		Value: expr}
}

func (p *Parser) funcDeclaration() Statement {
	name := p.consume(scanner.IDENTIFIER, "Expect function name.")
	p.consume(scanner.LEFT_PAREN, "Expect '(' after function name.")
	var params []*scanner.Token
	for !p.check(scanner.RIGHT_PAREN) && !p.isAtEnd() {
		if len(params) >= 255 {
			p.eh.Error(66, "[line %d] Cannot have more than 255 parameters.", p.peek().Line)
			return nil
		}
		params = append(params, p.consume(scanner.IDENTIFIER, "Expect parameter name."))
		if !p.check(scanner.RIGHT_PAREN) && !p.isAtEnd() {
			p.consume(scanner.COMMA, "Expect ',' between parameters.")
		}
	}

	p.consume(scanner.RIGHT_PAREN, "Expect ')' after parameters.")
	p.consume(scanner.LEFT_BRACE, "Expect '{' before function body.")
	body := p.block()

	return &FuncStatement{Name: name, Params: params, Body: body}
}

func (p *Parser) statement() Statement {
	if p.match(scanner.PRINT) {
		return p.printStatement()
	} else if p.match(scanner.LEFT_BRACE) {
		return p.block()
	} else if p.match(scanner.IF) {
		return p.ifStatement()
	} else if p.match(scanner.WHILE) {
		return p.whileStatement()
	} else if p.match(scanner.FOR) {
		return p.forStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) block() *BlockStatement {
	var statements []Statement
	for !p.check(scanner.RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	p.consume(scanner.RIGHT_BRACE, "Expect '}' after block.")

	return &BlockStatement{Statements: statements}
}

func (p *Parser) expressionStatement() Statement {
	e := p.expression()
	if !p.match(scanner.SEMICOLON) && !p.parseExpr {
		p.eh.Error(65, "[line %d] Expect ';' after expression.", p.peek().Line)
		return nil
	}
	return &ExpressionStatement{
		E: e,
	}
}

func (p *Parser) printStatement() Statement {
	e := p.expression()
	if !p.match(scanner.SEMICOLON) {
		p.eh.Error(65, "[line %d] Expect ';' after value.", p.peek().Line)
		return nil
	}
	return &PrintStatement{
		E: e,
	}
}

func (p *Parser) ifStatement() Statement {
	p.consume(scanner.LEFT_PAREN, "Expect '(' after 'if'.")
	cond := p.expression()
	p.consume(scanner.RIGHT_PAREN, "Expect ')' after if condition.")
	then := p.statement()
	var elseStmt Statement
	if p.match(scanner.ELSE) {
		elseStmt = p.statement()
	}
	return &IfStatement{
		Cond: cond,
		Then: then,
		Else: elseStmt,
	}
}

func (p *Parser) whileStatement() Statement {
	p.consume(scanner.LEFT_PAREN, "Expect '(' after 'if'.")
	cond := p.expression()
	p.consume(scanner.RIGHT_PAREN, "Expect ')' after if condition.")
	body := p.statement()
	return &WhileStatement{
		Cond: cond,
		Body: body,
	}
}

func (p *Parser) forStatement() Statement {
	p.consume(scanner.LEFT_PAREN, "Expect '(' after 'for'.")
	var initializer Statement
	if p.match(scanner.SEMICOLON) {
		initializer = nil
	} else if p.match(scanner.VAR) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}
	var cond Expr
	if !p.check(scanner.SEMICOLON) {
		cond = p.expression()
	}
	p.consume(scanner.SEMICOLON, "Expect ';' after for condition.")
	var inc Expr
	if !p.check(scanner.RIGHT_PAREN) {
		inc = p.expression()
	}
	p.consume(scanner.RIGHT_PAREN, "Expect ')' after for condition.")
	body := p.statement()

	if inc != nil {
		body = &BlockStatement{
			Statements: []Statement{body, &ExpressionStatement{inc}},
		}
	}
	if cond == nil {
		cond = &LiteralExpr{true}
	}
	body = &WhileStatement{
		Cond: cond,
		Body: body,
	}

	if initializer != nil {
		body = &BlockStatement{
			Statements: []Statement{initializer, body},
		}
	}
	return body
}

func (p *Parser) expression() Expr {
	return p.assignment()
}

func (p *Parser) assignment() Expr {
	e := p.or()
	if p.match(scanner.EQUAL) {
		equals := p.prev()
		rval := p.assignment()
		if v, ok := e.(*IdentifierExpr); ok {
			return &AssignmentExpr{
				Name:  v.Name,
				Value: rval,
			}
		}
		p.eh.Error(65, "[line %d] Invalid assignment target.", equals.Line)
	}
	return e
}

func (p *Parser) or() Expr {
	e := p.and()
	if p.eh.InError() {
		return nil
	}
	for p.match(scanner.OR) {
		t := p.prev()
		e = &LogicalExpr{
			Left:     e,
			Right:    p.and(),
			Operator: t,
		}
	}
	return e
}

func (p *Parser) and() Expr {
	e := p.equality()
	if p.eh.InError() {
		return nil
	}
	for p.match(scanner.AND) {
		t := p.prev()
		e = &LogicalExpr{
			Left:     e,
			Right:    p.equality(),
			Operator: t,
		}
	}
	return e
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
	return p.call()
}

func (p *Parser) call() Expr {
	e := p.primary()

	for {
		if p.match(scanner.LEFT_PAREN) {
			e = p.finishCall(e)
		} else {
			break
		}
	}
	return e
}

func (p *Parser) finishCall(e Expr) Expr {
	var args []Expr
	if !p.check(scanner.RIGHT_PAREN) {
		for !p.isAtEnd() {
			if len(args) >= 255 {
				p.eh.Error(65, "[line %d] Can't have more than 255 arguments.", p.peek().Line)
			}
			args = append(args, p.expression())
			if !p.match(scanner.COMMA) {
				break
			}
		}
	}

	t := p.consume(scanner.RIGHT_PAREN, "Expect ')' after arguments.")
	return &CallExpr{
		Callee:    e,
		Arguments: args,
		Right:     t,
	}
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

	if p.match(scanner.IDENTIFIER) {
		return &IdentifierExpr{
			Name: p.prev(),
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

func (p *Parser) advance() *scanner.Token {
	var token = p.tokens[p.p]
	p.p++
	return token
}

func (p *Parser) isAtEnd() bool {
	return len(p.tokens) <= p.p || p.tokens[p.p].Type == scanner.EOF || p.eh.InError()
}

func (p *Parser) consume(tokenType scanner.TokenType, errorMessage string) *scanner.Token {
	if p.check(tokenType) {
		return p.advance()
	}
	p.eh.Error(65, "[line %d] %s", p.peek().Line, errorMessage)
	return nil
}
