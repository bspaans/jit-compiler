package expr

import (
	"fmt"

	. "github.com/bspaans/jit-compiler/ir/shared"
)

type IR_Bool struct {
	*BaseIRExpression
	Value bool
}

func NewIR_Bool(v bool) *IR_Bool {
	return &IR_Bool{
		BaseIRExpression: NewBaseIRExpression(Bool),
		Value:            v,
	}
}

func (i *IR_Bool) ReturnType(ctx *IR_Context) Type {
	return TBool
}

func (i *IR_Bool) String() string {
	return fmt.Sprintf("%v", i.Value)
}

func (b *IR_Bool) SSA_Transform(ctx *SSA_Context) (SSA_Rewrites, IRExpression) {
	return nil, b
}
