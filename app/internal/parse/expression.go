package parse

import "github.com/codecrafters-io/interpreter-starter-go/app/internal/scanner"

type Expr interface {
	Accept(visitor Visitor) any
}

type BinaryExpr struct {
	Left     Expr
	Right    Expr
	Operator *scanner.Token
}

func (be *BinaryExpr) Accept(visitor Visitor) any {
	return visitor.VisitBinaryExpr(be)
}

type LiteralExpr struct {
	Obj any
}

func (le *LiteralExpr) Accept(visitor Visitor) any {
	return visitor.VisitLiteralExpr(le)
}

type GroupingExpr struct {
	E Expr
}

func (ge *GroupingExpr) Accept(visitor Visitor) any {
	return visitor.VisitGroupExpr(ge)
}

type UnaryExpr struct {
	Right    Expr
	Operator *scanner.Token
}

func (ue *UnaryExpr) Accept(visitor Visitor) any {
	return visitor.VisitUnaryExpr(ue)
}

type IdentifierExpr struct {
	Name *scanner.Token
}

func (ie *IdentifierExpr) Accept(visitor Visitor) any {
	return visitor.VisitIdentifierExpr(ie)
}

type AssignmentExpr struct {
	Name  *scanner.Token
	Value Expr
}

func (ae *AssignmentExpr) Accept(visitor Visitor) any {
	return visitor.VisitAssignmentExpr(ae)
}

type LogicalExpr struct {
	Left     Expr
	Right    Expr
	Operator *scanner.Token
}

func (le *LogicalExpr) Accept(visitor Visitor) any {
	return visitor.VisitLogicalExpr(le)
}

type Visitor interface {
	VisitBinaryExpr(expr *BinaryExpr) any
	VisitLiteralExpr(expr *LiteralExpr) any
	VisitGroupExpr(expr *GroupingExpr) any
	VisitUnaryExpr(expr *UnaryExpr) any
	VisitExpressionStatement(expr *ExpressionStatement) any
	VisitPrintStatement(expr *PrintStatement) any
	VisitIdentifierExpr(ie *IdentifierExpr) any
	VisitVarStatement(v *VarStatement) any
	VisitAssignmentExpr(ae *AssignmentExpr) any
	VisitBlockStatement(b *BlockStatement) any
	VisitIfStatement(i *IfStatement) any
	VisitLogicalExpr(le *LogicalExpr) any
	VisitWhileStatement(i *WhileStatement) any
}
