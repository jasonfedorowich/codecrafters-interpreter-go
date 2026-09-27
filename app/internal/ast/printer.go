package ast

import (
	"fmt"
	"strconv"

	"github.com/codecrafters-io/interpreter-starter-go/app/internal/parse"
)

type AstPrinter struct {
}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (ap *AstPrinter) VisitBinaryExpr(expr *parse.BinaryExpr) any {
	return fmt.Sprintf("(%s %s %s)", expr.Operator.Lexeme, expr.Left.Accept(ap), expr.Right.Accept(ap))
}

func (ap *AstPrinter) VisitLiteralExpr(expr *parse.LiteralExpr) any {
	if expr.Obj == nil {
		return "nil"
	}
	if val, ok := expr.Obj.(float64); ok {
		return strconv.FormatFloat(val, 'f', -1, 64)
	}
	return fmt.Sprint(expr.Obj)
}

func (ap *AstPrinter) VisitGroupExpr(expr *parse.GroupingExpr) any {
	return fmt.Sprintf("(group %s)", expr.E.Accept(ap))
}

func (ap *AstPrinter) VisitUnaryExpr(expr *parse.UnaryExpr) any {
	return fmt.Sprintf("(%s %s)", expr.Operator.Lexeme, expr.Right.Accept(ap))
}
