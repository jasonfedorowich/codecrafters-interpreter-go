package ast

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/codecrafters-io/interpreter-starter-go/app/internal/error_code"
	"github.com/codecrafters-io/interpreter-starter-go/app/internal/parse"
	"github.com/codecrafters-io/interpreter-starter-go/app/internal/scanner"
	"github.com/codecrafters-io/interpreter-starter-go/app/internal/types"
)

type ExprInterpreter struct {
	eh *error_code.ErrorHandler
}

func (i *ExprInterpreter) VisitIdentifierExpr(ie *parse.IdentifierExpr) any {
	//TODO implement me
	panic("implement me")
}

func (i *ExprInterpreter) VisitVarStatement(v *parse.VarStatement) any {
	//TODO implement me
	panic("implement me")
}

func NewExprInterpreter() *ExprInterpreter {
	return &ExprInterpreter{
		eh: &error_code.ErrorHandler{},
	}
}

func (i *ExprInterpreter) Interpret(expr parse.Expr) (any, int) {
	res := i.evaluate(expr)
	if res == nil {
		return "nil", i.eh.ErrorCode
	}
	if val, ok := res.(float64); ok {
		s := strconv.FormatFloat(val, 'f', -1, 64)
		if strings.HasSuffix(s, ".0") {
			return strings.TrimSuffix(s, ".0"), i.eh.ErrorCode
		}
	}
	return fmt.Sprint(res), i.eh.ErrorCode
}

func (i *ExprInterpreter) InError() bool {
	return i.eh.InError()
}

func (i *ExprInterpreter) VisitBinaryExpr(expr *parse.BinaryExpr) any {
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
	case scanner.EQUAL_EQUAL, scanner.BANG_EQUAL:
		return i.equal(expr.Operator.Type, left, right)

	}
	return nil
}

func (i *ExprInterpreter) VisitLiteralExpr(expr *parse.LiteralExpr) any {
	return expr.Obj
}

func (i *ExprInterpreter) VisitGroupExpr(expr *parse.GroupingExpr) any {
	return i.evaluate(expr.E)
}

func (i *ExprInterpreter) VisitUnaryExpr(expr *parse.UnaryExpr) any {
	switch expr.Operator.Type {
	case scanner.MINUS:
		if ok, v := types.IsNumber(i.evaluate(expr.Right)); ok {
			return -v
		}
		i.eh.Error(70, "Operand must be a number")
		return nil
	case scanner.BANG:
		return !i.isBool(i.evaluate(expr.Right))
	default:
		i.eh.Error(70, "Unknown unary operator")
		return nil
	}
}

func (i *ExprInterpreter) VisitExpressionStatement(expr *parse.ExpressionStatement) any {
	return i.evaluate(expr.E)
}

func (i *ExprInterpreter) VisitPrintStatement(expr *parse.PrintStatement) any {
	return nil
}

func (i *ExprInterpreter) evaluate(expr parse.Expr) any {
	return expr.Accept(i)
}

func (i *ExprInterpreter) parseNumberTypes(numbers ...any) []float64 {
	var nums []float64
	for _, n := range numbers {
		if number, ok := n.(float64); ok {
			nums = append(nums, number)
		} else {
			i.eh.Error(70, "Operands must be numbers")
			return nil
		}
	}
	return nums
}

func (i *ExprInterpreter) add(left any, right any) any {
	if i.eh.InError() {
		return nil
	}
	if strLeft, ok := left.(string); ok {
		if strRight, ok := right.(string); ok {
			return strLeft + strRight
		}
		i.eh.Error(70, "Operator + can only be used on two strings")
		return nil
	}
	arr := i.parseNumberTypes(left, right)
	if arr != nil {
		return arr[0] + arr[1]
	}
	return nil
}

func (i *ExprInterpreter) sub(left any, right any) any {
	arr := i.parseNumberTypes(left, right)
	if arr != nil {
		return arr[0] - arr[1]
	}
	return nil
}

func (i *ExprInterpreter) mul(left any, right any) any {
	arr := i.parseNumberTypes(left, right)
	if arr != nil {
		return arr[0] * arr[1]
	}
	return nil
}

func (i *ExprInterpreter) div(left any, right any) any {
	arr := i.parseNumberTypes(left, right)
	if arr != nil {
		return arr[0] / arr[1]
	}
	return nil
}

func (i *ExprInterpreter) logical(tokenType scanner.TokenType, left any, right any) any {
	arr := i.parseNumberTypes(left, right)
	if arr == nil {
		return nil
	}
	l, r := arr[0], arr[1]
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

func (i *ExprInterpreter) equal(tokenType scanner.TokenType, left any, right any) any {
	//numbers := types.ParseNumbers(left, right)
	//if numbers != nil {
	//	switch tokenType {
	//	case scanner.EQUAL_EQUAL:
	//		return numbers[0] == numbers[1]
	//	case scanner.BANG_EQUAL:
	//		return numbers[0] != numbers[1]
	//	}
	//}
	//strs := types.ParseStrings(left, right)
	//if strs != nil {
	//	switch tokenType {
	//	case scanner.EQUAL_EQUAL:
	//		return strs[0] == strs[1]
	//	case scanner.BANG_EQUAL:
	//		return strs[0] != strs[1]
	//	}
	//}
	//
	//i.eh.Error(65, "Expected number or string type")
	//return nil
	switch tokenType {
	case scanner.EQUAL_EQUAL:
		return left == right
	case scanner.BANG_EQUAL:
		return left != right
	}
	return nil
}

func (i *ExprInterpreter) isBool(e any) bool {
	if e == nil {
		return false
	}
	if v, ok := types.IsBoolean(e); ok {
		return v
	}
	return true
}
