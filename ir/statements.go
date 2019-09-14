package ir

import (
	"fmt"

	"github.com/bspaans/jit/asm"
)

type IRType int

const (
	Assignment IRType = iota
)

type IR_Assignment struct {
	*BaseIR
	Variable string
	Expr     IRExpression
}

func NewIR_Assignment(variable string, expr IRExpression) *IR_Assignment {
	return &IR_Assignment{
		BaseIR:   NewBaseIR(Assignment),
		Variable: variable,
		Expr:     expr,
	}
}

// Allocates a new register and assigns it the value of the expression.
func (i *IR_Assignment) Encode(ctx *IR_Context) ([]asm.Instruction, error) {
	r := ctx.AllocateRegister()
	ctx.VariableMap[i.Variable] = r
	reg := asm.Get64BitRegisterByIndex(r)
	return i.Expr.Encode(ctx, reg)
}

func (i *IR_Assignment) String() string {
	return fmt.Sprintf("%s = %s", i.Variable, i.Expr.String())
}
