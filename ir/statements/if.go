package statements

import (
	"fmt"

	. "github.com/bspaans/jit-compiler/ir/shared"
)

type IR_If struct {
	*BaseIR
	Condition IRExpression
	Stmt1     IR
	Stmt2     IR
}

func NewIR_If(condition IRExpression, stmt1, stmt2 IR) *IR_If {
	return &IR_If{
		BaseIR:    NewBaseIR(If),
		Condition: condition,
		Stmt1:     stmt1,
		Stmt2:     stmt2,
	}
}

func (i *IR_If) String() string {
	return fmt.Sprintf("if %s { %s } else { %s }", i.Condition.String(), i.Stmt1.String(), i.Stmt2.String())
}

func (i *IR_If) SSA_Transform(ctx *SSA_Context) IR {
	rewrites, expr := i.Condition.SSA_Transform(ctx)
	ir := SSA_Rewrites_to_IR(rewrites)
	if ir == nil {
		return NewIR_If(i.Condition, i.Stmt1.SSA_Transform(ctx), i.Stmt2.SSA_Transform(ctx))
	} else {
		return NewIR_AndThen(ir, NewIR_If(expr, i.Stmt1.SSA_Transform(ctx), i.Stmt2.SSA_Transform(ctx)))
	}
}
