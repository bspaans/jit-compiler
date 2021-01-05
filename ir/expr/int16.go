package expr

import (
	"fmt"

	. "github.com/bspaans/jit-compiler/ir/shared"
)

type IR_Int16 struct {
	*BaseIRExpression
	Value int16
}

func NewIR_Int16(v int16) *IR_Int16 {
	return &IR_Int16{
		BaseIRExpression: NewBaseIRExpression(Int16),
		Value:            v,
	}
}

func (i *IR_Int16) ReturnType(ctx *IR_Context) Type {
	return TInt16
}

func (i *IR_Int16) String() string {
	return fmt.Sprintf("%d", i.Value)
}

func (b *IR_Int16) SSA_Transform(ctx *SSA_Context) (SSA_Rewrites, IRExpression) {
	return nil, b
}
