package ast

//go:generate mockgen -source=interpreter.go -destination=mocks/mock_printer.go -package=mocks

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/codecrafters-io/interpreter-starter-go/app/internal/env"
	"github.com/codecrafters-io/interpreter-starter-go/app/internal/error_code"
	"github.com/codecrafters-io/interpreter-starter-go/app/internal/parse"
	"github.com/codecrafters-io/interpreter-starter-go/app/internal/scanner"
	"github.com/codecrafters-io/interpreter-starter-go/app/internal/types"
)

type Printer interface {
	Print(a any)
}

type printer struct {
}

func (p *printer) Print(a any) {
	fmt.Println(a)
}

type Interpreter struct {
	eh      *error_code.ErrorHandler
	p       Printer
	e       *env.Env
	locals  map[parse.Expr]int
	globals *env.Env
}

func NewInterpreter() *Interpreter {
	g := env.NewEnv(loadNativeFunctions())
	return &Interpreter{
		eh:      &error_code.ErrorHandler{},
		p:       &printer{},
		e:       g,
		locals:  make(map[parse.Expr]int),
		globals: g,
	}
}

func (i *Interpreter) InterpretStatements(statements []parse.Statement) int {
	for _, s := range statements {
		if i.eh.InError() {
			return i.eh.ErrorCode
		}
		i.evaluate(s)
	}
	return i.eh.ErrorCode
}

func (i *Interpreter) InterpretExpr(expr parse.Expr) (any, int) {
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

func (i *Interpreter) InError() bool {
	return i.eh.InError()
}

func (i *Interpreter) VisitThisExpr(te *parse.ThisExpr) any {
	return i.lookupVariable(te.Keyword, te)
}

func (i *Interpreter) VisitSetExpr(se *parse.SetExpr) any {
	obj := i.evaluate(se.Object)
	v, ok := obj.(*instance)
	if ok {
		value := i.evaluate(se.Value)
		v.set(se.Name, value)
		return value
	}
	i.eh.Error(70, "Only instances have properties.")
	return nil
}

func (i *Interpreter) VisitGetExpr(ge *parse.GetExpr) any {
	obj := i.evaluate(ge.Object)
	v, ok := obj.(*instance)
	if ok {
		r, err := v.get(ge.Name)
		if err != nil {
			i.eh.FromError(70, err)
			return nil
		}
		return r
	}
	i.eh.Error(70, "Only instances have properties.")
	return nil
}

func (i *Interpreter) VisitClassStatement(c *parse.ClassStatement) any {
	i.e.Put(c.Name.Lexeme, nil)

	var methods = make(map[string]*function)
	for _, method := range c.Statements {
		methods[method.Name.Lexeme] = newFunction(method, i.e)
	}
	cl := newClass(c.Name.Lexeme, methods)
	i.e.Assign(c.Name.Lexeme, cl)
	return nil
}

func (i *Interpreter) VisitFuncStatement(f *parse.FuncStatement) any {
	i.e.Put(f.Name.Lexeme, newFunction(f, i.e))
	return nil
}
func (i *Interpreter) VisitReturnStmt(r *parse.ReturnStmt) any {
	var e any
	if r.E != nil {
		e = i.evaluate(r.E)
	}
	panic(&returnValue{e})
	return nil
}
func (i *Interpreter) VisitCallExpr(expr *parse.CallExpr) any {
	f := i.evaluate(expr.Callee)
	var args []any
	for _, a := range expr.Arguments {
		args = append(args, i.evaluate(a))
	}
	function, ok := f.(Callable)
	if !ok {
		i.eh.Error(70, "Can only call functions and classes.")
		return nil
	}
	if function.arity() != len(args) {
		i.eh.Error(70, "Wrong number of arguments.")
		return nil
	}
	return function.call(i, args)
}

func (i *Interpreter) VisitLogicalExpr(expr *parse.LogicalExpr) any {
	l := i.evaluate(expr.Left)
	switch expr.Operator.Type {
	case scanner.OR:
		if i.isBool(l) {
			return l
		}
	case scanner.AND:
		if !i.isBool(l) {
			return l
		}
	}
	return i.evaluate(expr.Right)
}

func (i *Interpreter) VisitBlockStatement(b *parse.BlockStatement) any {
	return i.blockStatement(b.Statements, env.NewWrappedEnv(i.e))
}

func (i *Interpreter) blockStatement(statements []parse.Statement, wrappedEnv *env.Env) any {
	prev := i.e
	i.e = wrappedEnv
	defer func() {
		i.e = prev
	}()
	for _, s := range statements {
		if i.eh.InError() {
			return nil
		}
		i.evaluate(s)
	}
	return nil
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
	case scanner.EQUAL_EQUAL, scanner.BANG_EQUAL:
		return i.equal(expr.Operator.Type, left, right)

	}
	return nil
}

func (i *Interpreter) VisitLiteralExpr(expr *parse.LiteralExpr) any {
	return expr.Obj
}

func (i *Interpreter) VisitGroupExpr(expr *parse.GroupingExpr) any {
	return i.evaluate(expr.E)
}

func (i *Interpreter) VisitUnaryExpr(expr *parse.UnaryExpr) any {
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

func (i *Interpreter) VisitWhileStatement(while *parse.WhileStatement) any {
	for i.isBool(i.evaluate(while.Cond)) {
		i.evaluate(while.Body)
	}
	return nil
}

func (i *Interpreter) VisitIfStatement(ifStmt *parse.IfStatement) any {
	if i.isBool(i.evaluate(ifStmt.Cond)) {
		return i.evaluate(ifStmt.Then)
	} else if ifStmt.Else != nil {
		return i.evaluate(ifStmt.Else)
	}
	return nil
}

func (i *Interpreter) VisitAssignmentExpr(expr *parse.AssignmentExpr) any {
	value := i.evaluate(expr.Value)
	d, ok := i.locals[expr]
	if ok {
		i.e.AssignAt(d, expr.Name, value)
	} else {
		err := i.globals.Assign(expr.Name.Lexeme, value)
		if err != nil {
			i.eh.Error(70, "Undefined variable: %s", expr.Name.Lexeme)
			fmt.Printf("[line %d]", expr.Name.Line)
			return nil
		}
	}
	return value
}

func (i *Interpreter) evaluate(expr parse.Expr) any {
	return expr.Accept(i)
}

func (i *Interpreter) parseNumberTypes(numbers ...any) []float64 {
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

func (i *Interpreter) add(left any, right any) any {
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

func (i *Interpreter) sub(left any, right any) any {
	arr := i.parseNumberTypes(left, right)
	if arr != nil {
		return arr[0] - arr[1]
	}
	return nil
}

func (i *Interpreter) mul(left any, right any) any {
	arr := i.parseNumberTypes(left, right)
	if arr != nil {
		return arr[0] * arr[1]
	}
	return nil
}

func (i *Interpreter) div(left any, right any) any {
	arr := i.parseNumberTypes(left, right)
	if arr != nil {
		return arr[0] / arr[1]
	}
	return nil
}

func (i *Interpreter) logical(tokenType scanner.TokenType, left any, right any) any {
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

func (i *Interpreter) equal(tokenType scanner.TokenType, left any, right any) any {
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

func (i *Interpreter) isBool(e any) bool {
	if e == nil {
		return false
	}
	if v, ok := types.IsBoolean(e); ok {
		return v
	}
	return true
}

func (i *Interpreter) VisitIdentifierExpr(ie *parse.IdentifierExpr) any {
	return i.lookupVariable(ie.Name, ie)
}

func (i *Interpreter) VisitVarStatement(v *parse.VarStatement) any {
	if v.Value == nil {
		i.e.Put(v.Name.Lexeme, nil)
		return nil
	}
	rhs := i.evaluate(v.Value)
	i.e.Put(v.Name.Lexeme, rhs)
	return nil
}

func (i *Interpreter) VisitExpressionStatement(expr *parse.ExpressionStatement) any {
	return i.evaluate(expr.E)
}

func (i *Interpreter) VisitPrintStatement(expr *parse.PrintStatement) any {
	r := i.evaluate(expr.E)
	if i.eh.InError() {
		return nil
	}
	s := stringify(r)
	i.p.Print(s)
	return nil
}

func (i *Interpreter) resolve(expr parse.Expr, depth int) {
	i.locals[expr] = depth
}

func (i *Interpreter) lookupVariable(name *scanner.Token, expr parse.Expr) any {
	d, ok := i.locals[expr]
	if !ok {
		v, ok := i.globals.Get(name.Lexeme)
		if !ok {
			i.eh.Error(70, "Undefined variable '%s'", name.Lexeme)
			return nil
		}
		return v
	}
	return i.e.GetAt(d, name.Lexeme)
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
