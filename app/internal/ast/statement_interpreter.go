package ast

//go:generate mockgen -source=statement_interpreter.go -destination=mocks/mock_printer.go -package=mocks

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/codecrafters-io/interpreter-starter-go/app/internal/error_code"
	"github.com/codecrafters-io/interpreter-starter-go/app/internal/parse"
)

type Printer interface {
	Print(a any)
}

type printer struct {
}

func (p *printer) Print(a any) {
	fmt.Println(a)
}

type StatementInterpreter struct {
	ExprInterpreter
	p Printer
}

func NewStatementInterpreter() *StatementInterpreter {
	return &StatementInterpreter{
		eh: &error_code.ErrorHandler{},
		p:  &printer{},
	}
}

func (i *StatementInterpreter) Interpret(statements []parse.Statement) int {
	for _, s := range statements {
		if i.eh.InError() {
			return i.eh.ErrorCode
		}
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
	if i.eh.InError() {
		return nil
	}
	s := stringify(r)
	i.p.Print(s)
	return nil
}

func stringify(a any) string {
	if a == nil {
		return "nil"
	}
	if val, ok := a.(float64); ok {
		s := strconv.FormatFloat(val, 'f', -1, 64)
		if strings.HasSuffix(s, ".0") {
			return strings.TrimSuffix(s, ".0")
		}
	}
	return fmt.Sprint(a)
}
