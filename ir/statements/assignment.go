package statements

import (
	"fmt"

	. "github.com/bspaans/jit-compiler/ir/shared"
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

func (i *IR_Assignment) String() string {
	return fmt.Sprintf("%s = %s", i.Variable, i.Expr.String())
}

func (i *IR_Assignment) AddToDataSection(ctx *IR_Context) error {
	return i.Expr.AddToDataSection(ctx)
}

func (i *IR_Assignment) SSA_Transform(ctx *SSA_Context) IR {
	rewrites, expr := i.Expr.SSA_Transform(ctx)
	ir := SSA_Rewrites_to_IR(rewrites)
	if ir == nil {
		return i
	}
	return NewIR_AndThen(ir, NewIR_Assignment(i.Variable, expr))
}
