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

type IfStatement struct {
	Cond Expr
	Then Statement
	Else Statement
}

func (i *IfStatement) Accept(visitor Visitor) any {
	return visitor.VisitIfStatement(i)
}

type WhileStatement struct {
	Cond Expr
	Body Statement
}

func (i *WhileStatement) Accept(visitor Visitor) any {
	return visitor.VisitWhileStatement(i)
}

type FuncStatement struct {
	Name   *scanner.Token
	Params []*scanner.Token
	Body   []Statement
}

func (f *FuncStatement) Accept(visitor Visitor) any {
	return visitor.VisitFuncStatement(f)
}

type ReturnStmt struct {
	E Expr
}

func (r *ReturnStmt) Accept(visitor Visitor) any {
	return visitor.VisitReturnStmt(r)
}
