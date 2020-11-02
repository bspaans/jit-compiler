package statements

import (
	"fmt"

	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
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

func (i *IR_AndThen) Encode(ctx *IR_Context) ([]lib.Instruction, error) {
	result, err := i.Stmt1.Encode(ctx)
	if err != nil {
		return nil, err
	}
	s2, err := i.Stmt2.Encode(ctx)
	if err != nil {
		return nil, err
	}
	for _, instr := range s2 {
		result = append(result, instr)
	}
	return result, nil
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
