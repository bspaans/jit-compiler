package statements

import (
	"fmt"

	. "github.com/bspaans/jit-compiler/ir/shared"
)

type IR_While struct {
	*BaseIR
	Condition IRExpression
	Stmt      IR
}

func NewIR_While(condition IRExpression, stmt IR) *IR_While {
	return &IR_While{
		BaseIR:    NewBaseIR(While),
		Condition: condition,
		Stmt:      stmt,
	}
}

func (i *IR_While) String() string {
	return fmt.Sprintf("while %s { %s }", i.Condition.String(), i.Stmt.String())
}

func (i *IR_While) AddToDataSection(ctx *IR_Context) error {
	if err := i.Condition.AddToDataSection(ctx); err != nil {
		return err
	}
	if err := i.Stmt.AddToDataSection(ctx); err != nil {
		return err
	}
	return nil
}

func (i *IR_While) SSA_Transform(ctx *SSA_Context) IR {
	// TODO: transform i.Condition => changes the encoding though
	i.Stmt = i.Stmt.SSA_Transform(ctx)
	return i
}
