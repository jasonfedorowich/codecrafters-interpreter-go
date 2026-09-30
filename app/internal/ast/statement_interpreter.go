package ast

//go:generate mockgen -source=statement_interpreter.go -destination=mocks/mock_printer.go -package=mocks


import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/app/internal/error_code"
	"github.com/codecrafters-io/interpreter-starter-go/app/internal/parse"
)

type Printer interface {
	Print(a any)
}

type printer struct{

}

func (p *printer) Print(a any){
	fmt.Println(a)
}


type StatementInterpreter struct {
	ExprInterpreter
	p Printer
}

func NewStatementInterpreter() *StatementInterpreter {
	return &StatementInterpreter{
		eh: &error_code.ErrorHandler{},
		p: &printer{},
	}
}

func (i *StatementInterpreter) Interpret(statements []parse.Statement) int {
	for _, s := range statements {
		i.evaluate(s)
	}
	return i.eh.ErrorCode
}

func (i *StatementInterpreter) evaluate(expr parse.Expr) any {
	return expr.Accept(i)
}

func (i *StatementInterpreter) VisitExpressionStatement(expr *parse.ExpressionStatement) any {
	return i.evaluate(expr.E)
}

func (i *StatementInterpreter) VisitPrintStatement(expr *parse.PrintStatement) any {
	r := i.evaluate(expr.E)
	s := toString(r)
	i.p.Print(s)
	return nil
}
