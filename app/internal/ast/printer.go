package ast

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/app/internal/parse"
)

type AstPrinter struct {
}

func (ap *AstPrinter) VisitFuncStatement(f *parse.FuncStatement) any {
	//TODO implement me
	panic("implement me")
}

func (ap *AstPrinter) VisitCallExpr(ce *parse.CallExpr) any {
	//TODO implement me
	panic("implement me")
}

func (ap *AstPrinter) VisitWhileStatement(i *parse.WhileStatement) any {
	//TODO implement me
	panic("implement me")
}

func (ap *AstPrinter) VisitLogicalExpr(le *parse.LogicalExpr) any {
	//TODO implement me
	panic("implement me")
}

func (ap *AstPrinter) VisitIfStatement(i *parse.IfStatement) any {
	//TODO implement me
	panic("implement me")
}

func (ap *AstPrinter) VisitBlockStatement(b *parse.BlockStatement) any {
	//TODO implement me
	panic("implement me")
}

func (ap *AstPrinter) VisitAssignmentExpr(ae *parse.AssignmentExpr) any {
	//TODO implement me
	panic("implement me")
}

func (ap *AstPrinter) VisitIdentifierExpr(ie *parse.IdentifierExpr) any {
	//TODO implement me
	panic("implement	 me")
}

func (ap *AstPrinter) VisitVarStatement(v *parse.VarStatement) any {
	//TODO implement me
	panic("implement me")
}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (ap *AstPrinter) VisitBinaryExpr(expr *parse.BinaryExpr) any {
	return fmt.Sprintf("(%s %s %s)", expr.Operator.Lexeme, expr.Left.Accept(ap), expr.Right.Accept(ap))
}

func (ap *AstPrinter) VisitLiteralExpr(expr *parse.LiteralExpr) any {
	return toString(expr.Obj)
}

func (ap *AstPrinter) VisitGroupExpr(expr *parse.GroupingExpr) any {
	return fmt.Sprintf("(group %s)", expr.E.Accept(ap))
}

func (ap *AstPrinter) VisitUnaryExpr(expr *parse.UnaryExpr) any {
	return fmt.Sprintf("(%s %s)", expr.Operator.Lexeme, expr.Right.Accept(ap))
}

func (ap *AstPrinter) VisitPrintStatement(expr *parse.PrintStatement) any {
	r := expr.E.Accept(ap)
	s := toString(r)
	fmt.Println(s)
	return nil
}

func (ap *AstPrinter) VisitExpressionStatement(expr *parse.ExpressionStatement) any {
	return expr.E.Accept(ap)
}

func toString(expr any) string {
	if expr == nil {
		return "nil"
	}
	if val, ok := expr.(float64); ok && val == float64(int64(val)) {
		return fmt.Sprintf("%.1f", val)
	}
	return fmt.Sprint(expr)
}
