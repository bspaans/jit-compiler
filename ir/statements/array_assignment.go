package statements

import (
	"fmt"

	. "github.com/bspaans/jit-compiler/ir/shared"
)

type IR_ArrayAssignment struct {
	*BaseIR
	Variable string
	Index    IRExpression
	Expr     IRExpression
}

func NewIR_ArrayAssignment(variable string, index IRExpression, expr IRExpression) *IR_ArrayAssignment {
	return &IR_ArrayAssignment{
		BaseIR:   NewBaseIR(ArrayAssignment),
		Variable: variable,
		Index:    index,
		Expr:     expr,
	}
}

func (i *IR_ArrayAssignment) String() string {
	return fmt.Sprintf("%s[%s] = %s", i.Variable, i.Index.String(), i.Expr.String())
}

func (i *IR_ArrayAssignment) AddToDataSection(ctx *IR_Context) error {
	if err := i.Index.AddToDataSection(ctx); err != nil {
		return err
	}
	return i.Expr.AddToDataSection(ctx)
}

func (i *IR_ArrayAssignment) SSA_Transform(ctx *SSA_Context) IR {
	rewrites, expr := i.Index.SSA_Transform(ctx)
	rewrites2, expr2 := i.Expr.SSA_Transform(ctx)
	for _, rw := range rewrites2 {
		rewrites = append(rewrites, rw)
	}
	ir := SSA_Rewrites_to_IR(rewrites)
	if ir == nil {
		return i
	}
	return NewIR_AndThen(ir, NewIR_ArrayAssignment(i.Variable, expr, expr2))
}
