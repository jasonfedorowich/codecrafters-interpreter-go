package parse

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
