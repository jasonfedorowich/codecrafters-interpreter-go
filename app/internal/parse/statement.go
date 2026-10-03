package parse

import "github.com/codecrafters-io/interpreter-starter-go/app/internal/scanner"

type Statement interface {
	Accept(visitor Visitor) any
}

type ExpressionStatement struct {
	E Expr
}

func (e *ExpressionStatement) Accept(visitor Visitor) any {
	return visitor.VisitExpressionStatement(e)
}

type PrintStatement struct {
	E Expr
}

func (p *PrintStatement) Accept(visitor Visitor) any {
	return visitor.VisitPrintStatement(p)
}

type VarStatement struct {
	Name  *scanner.Token
	Value Expr
}

func (v *VarStatement) Accept(visitor Visitor) any {
	return visitor.VisitVarStatement(v)
}

type BlockStatement struct {
	Statements []Statement
}

func (b *BlockStatement) Accept(visitor Visitor) any {
	return visitor.VisitBlockStatement(b)
}
