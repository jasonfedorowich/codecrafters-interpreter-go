package ast

import (
	"github.com/codecrafters-io/interpreter-starter-go/app/internal/error_code"
	"github.com/codecrafters-io/interpreter-starter-go/app/internal/parse"
	"github.com/codecrafters-io/interpreter-starter-go/app/internal/scanner"
)

const (
	IN_FUNC = iota
	NO_FUNC
	IN_METHOD
)

type Resolver struct {
	i      *Interpreter
	scopes []map[string]bool
	eh     *error_code.ErrorHandler
	fn     int
}

func NewResolver(i *Interpreter) *Resolver {
	return &Resolver{
		i:  i,
		eh: &error_code.ErrorHandler{},
		fn: NO_FUNC,
	}
}

func (r *Resolver) ResolveStatements(statements []parse.Statement) int {
	for _, statement := range statements {
		statement.Accept(r)
	}
	return r.eh.ErrorCode
}

func (r *Resolver) resolveExpr(expr parse.Expr) {
	expr.Accept(r)
}

func (r *Resolver) resolveStatement(statement parse.Statement) {
	statement.Accept(r)
}

func (r *Resolver) VisitSetExpr(se *parse.SetExpr) any {
	r.resolveExpr(se.Value)
	r.resolveExpr(se.Object)
	return nil
}

func (r *Resolver) VisitGetExpr(ge *parse.GetExpr) any {
	r.resolveExpr(ge.Object)
	return nil
}

func (r *Resolver) VisitBinaryExpr(expr *parse.BinaryExpr) any {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitClassStatement(c *parse.ClassStatement) any {
	r.declare(c.Name)
	r.define(c.Name)
	for _, method := range c.Statements {
		r.resolveFunc(method, IN_METHOD)
	}
	return nil
}

func (r *Resolver) VisitLiteralExpr(expr *parse.LiteralExpr) any {
	return nil
}

func (r *Resolver) VisitGroupExpr(expr *parse.GroupingExpr) any {
	r.resolveExpr(expr.E)
	return nil
}

func (r *Resolver) VisitUnaryExpr(expr *parse.UnaryExpr) any {
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitExpressionStatement(expr *parse.ExpressionStatement) any {
	r.resolveExpr(expr.E)
	return nil
}

func (r *Resolver) VisitPrintStatement(expr *parse.PrintStatement) any {
	r.resolveExpr(expr.E)
	return nil
}

func (r *Resolver) VisitIdentifierExpr(ie *parse.IdentifierExpr) any {
	if len(r.scopes) != 0 {
		v, ok := r.scopes[len(r.scopes)-1][ie.Name.Lexeme]
		if ok && !v {
			r.eh.Error(65, "Can't read local variable in its own initializer.")
		}
	}
	r.resolveLocal(ie, ie.Name)
	return nil
}

func (r *Resolver) VisitVarStatement(v *parse.VarStatement) any {
	r.declare(v.Name)
	if v.Value != nil {
		r.resolveExpr(v.Value)
	}
	r.define(v.Name)
	return nil
}

func (r *Resolver) VisitAssignmentExpr(ae *parse.AssignmentExpr) any {
	r.resolveExpr(ae.Value)
	r.resolveLocal(ae, ae.Name)
	return nil
}

func (r *Resolver) VisitBlockStatement(b *parse.BlockStatement) any {
	r.beginScope()
	r.ResolveStatements(b.Statements)
	r.endScope()
	return nil
}

func (r *Resolver) VisitIfStatement(i *parse.IfStatement) any {
	r.resolveExpr(i.Cond)
	r.resolveStatement(i.Then)
	if i.Else != nil {
		r.resolveStatement(i.Else)
	}
	return nil
}

func (r *Resolver) VisitLogicalExpr(le *parse.LogicalExpr) any {
	r.resolveExpr(le.Left)
	r.resolveExpr(le.Right)
	return nil
}

func (r *Resolver) VisitWhileStatement(i *parse.WhileStatement) any {
	r.resolveExpr(i.Cond)
	r.resolveStatement(i.Body)
	return nil
}

func (r *Resolver) VisitCallExpr(ce *parse.CallExpr) any {
	r.resolveExpr(ce.Callee)
	for _, arg := range ce.Arguments {
		r.resolveExpr(arg)
	}
	return nil
}

func (r *Resolver) VisitFuncStatement(f *parse.FuncStatement) any {
	r.declare(f.Name)
	r.define(f.Name)
	r.resolveFunc(f, IN_FUNC)
	return nil
}

func (r *Resolver) VisitReturnStmt(re *parse.ReturnStmt) any {
	if r.fn == NO_FUNC {
		r.eh.Error(65, "Cannot return from top-level code.")
		return nil
	}
	if re.E != nil {
		r.resolveExpr(re.E)
	}
	return nil
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, make(map[string]bool))
}

func (r *Resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) declare(name *scanner.Token) {
	if len(r.scopes) == 0 {
		return
	}
	if _, ok := r.scopes[len(r.scopes)-1][name.Lexeme]; ok {
		r.eh.Error(65, "Already a variable with this name in this scope.")
		return
	}
	r.scopes[len(r.scopes)-1][name.Lexeme] = false
}

func (r *Resolver) define(name *scanner.Token) {
	if len(r.scopes) == 0 {
		return
	}
	r.scopes[len(r.scopes)-1][name.Lexeme] = true
}

func (r *Resolver) resolveLocal(expr parse.Expr, name *scanner.Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name.Lexeme]; ok {
			r.i.resolve(expr, len(r.scopes)-1-i)
			return
		}
	}
}

func (r *Resolver) resolveFunc(f *parse.FuncStatement, fn int) {
	prev := r.fn
	r.fn = fn
	r.beginScope()
	for _, param := range f.Params {
		r.declare(param)
		r.define(param)
	}
	r.ResolveStatements(f.Body)
	r.endScope()
	r.fn = prev
}
