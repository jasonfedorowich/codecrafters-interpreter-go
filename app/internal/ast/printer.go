package ast

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/app/internal/parse"
)

type AstPrinter struct {
}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (ap *AstPrinter) VisitBinaryExpr(expr *parse.BinaryExpr) any {
	return nil
}

func (ap *AstPrinter) VisitLiteralExpr(expr *parse.LiteralExpr) any {
	if expr.Obj == nil {
		return "nil"
	}
	return fmt.Sprint(expr.Obj)
}

func (ap *AstPrinter) VisitGroupExpr(expr *parse.GroupingExpr) any {
	return fmt.Sprintf("(group %s)", expr.E.Accept(ap))
}
