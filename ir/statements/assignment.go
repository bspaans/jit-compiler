package statements

import (
	"fmt"

	"github.com/bspaans/jit/asm"
	. "github.com/bspaans/jit/ir/shared"
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
	returnType := i.Expr.ReturnType(ctx)
	reg, found := ctx.VariableMap[i.Variable]
	if !found {
		reg = ctx.AllocateRegister(returnType)
		ctx.VariableMap[i.Variable] = reg
		ctx.VariableTypes[i.Variable] = returnType
	}
	return i.Expr.Encode(ctx, reg)
}

func (i *IR_Assignment) String() string {
	return fmt.Sprintf("%s = %s", i.Variable, i.Expr.String())
}

func (i *IR_Assignment) AddToDataSection(ctx *IR_Context) error {
	return i.Expr.AddToDataSection(ctx)
}
