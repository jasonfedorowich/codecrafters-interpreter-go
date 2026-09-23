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

type Visitor interface {
	VisitBinaryExpr(expr *BinaryExpr) any
	VisitLiteralExpr(expr *LiteralExpr) any
}
