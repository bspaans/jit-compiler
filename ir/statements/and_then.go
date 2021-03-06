package statements

import (
	"fmt"

	. "github.com/bspaans/jit-compiler/ir/shared"
)

type IR_AndThen struct {
	*BaseIR
	Stmt1 IR
	Stmt2 IR
}

func NewIR_AndThen(stmt1, stmt2 IR) *IR_AndThen {
	return &IR_AndThen{
		BaseIR: NewBaseIR(AndThen),
		Stmt1:  stmt1,
		Stmt2:  stmt2,
	}
}

func (i *IR_AndThen) String() string {
	return fmt.Sprintf("%s ; %s", i.Stmt1.String(), i.Stmt2.String())
}

func (i *IR_AndThen) AddToDataSection(ctx *IR_Context) error {
	if err := i.Stmt1.AddToDataSection(ctx); err != nil {
		return err
	}
	if err := i.Stmt2.AddToDataSection(ctx); err != nil {
		return err
	}
	return nil
}

func (i *IR_AndThen) SSA_Transform(ctx *SSA_Context) IR {
	i.Stmt1 = i.Stmt1.SSA_Transform(ctx)
	i.Stmt2 = i.Stmt2.SSA_Transform(ctx)
	return i
}
