package ast

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/interpreter-starter-go/app/internal/error_code"
	"github.com/codecrafters-io/interpreter-starter-go/app/internal/parse"
	"github.com/codecrafters-io/interpreter-starter-go/app/internal/scanner"
)

type Interpreter struct {
	eh *error_code.ErrorHandler
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		eh: &error_code.ErrorHandler{},
	}
}

func (i *Interpreter) Interpret(expr parse.Expr) any {
	res := i.evaluate(expr)
	if res == nil {
		return "nil"
	}
	if val, ok := res.(float64); ok && val == float64(int64(val)) {
		s := fmt.Sprintf("%.1f", val)
		if strings.HasPrefix(s, ".0") {
			return strings.TrimLeft(s, ".")
		}
		return s
	}
	return fmt.Sprint(res)
}

func (i *Interpreter) InError() bool {
	return i.eh.InError()
}

func (i *Interpreter) VisitBinaryExpr(expr *parse.BinaryExpr) any {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case scanner.PLUS:
		return i.add(left, right)
	case scanner.MINUS:
		return i.sub(left, right)
	case scanner.STAR:
		return i.mul(left, right)
	case scanner.SLASH:
		return i.div(left, right)
	case scanner.GREATER, scanner.GREATER_EQUAL, scanner.LESS, scanner.LESS_EQUAL:
		return i.logical(expr.Operator.Type, left, right)

	}
	return nil
}

func (i *Interpreter) add(left any, right any) any {
	if i.eh.InError() {
		return nil
	}
	if strLeft, ok := left.(string); ok {
		if strRight, ok := right.(string); ok {
			return strLeft + strRight
		}
		i.eh.Error(65, "Operator + can only be used on two strings")
		return nil
	}
	l, r := i.parseNumberTypes(left, right)
	return l + r
}

func (i *Interpreter) VisitLiteralExpr(expr *parse.LiteralExpr) any {
	return expr.Obj
}

func (i *Interpreter) VisitGroupExpr(expr *parse.GroupingExpr) any {
	//TODO implement me
	panic("implement me")
}

func (i *Interpreter) VisitUnaryExpr(expr *parse.UnaryExpr) any {
	//TODO implement me
	panic("implement me")
}

func (i *Interpreter) evaluate(expr parse.Expr) any {
	return expr.Accept(i)
}

func (i *Interpreter) parseNumberTypes(left any, right any) (float64, float64) {
	if leftNum, ok := left.(float64); ok {
		if rightNum, ok := right.(float64); ok {
			return leftNum, rightNum
		}
	}
	i.eh.Error(65, "Operands must be numbers")
	return 0, 0
}
func (i *Interpreter) sub(left any, right any) any {
	l, r := i.parseNumberTypes(left, right)
	return l - r
}

func (i *Interpreter) mul(left any, right any) any {
	l, r := i.parseNumberTypes(left, right)
	return l * r
}

func (i *Interpreter) div(left any, right any) any {
	l, r := i.parseNumberTypes(left, right)
	return l / r
}

func (i *Interpreter) logical(tokenType scanner.TokenType, left any, right any) any {
	l, r := i.parseNumberTypes(left, right)
	switch tokenType {
	case scanner.GREATER:
		return l > r
	case scanner.GREATER_EQUAL:
		return l >= r
	case scanner.LESS:
		return l < r
	case scanner.LESS_EQUAL:
		return l <= r
	}
	return nil
}
