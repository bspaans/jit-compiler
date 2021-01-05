package expr

import (
	"fmt"

	. "github.com/bspaans/jit-compiler/ir/shared"
)

type IR_Int32 struct {
	*BaseIRExpression
	Value int32
}

func NewIR_Int32(v int32) *IR_Int32 {
	return &IR_Int32{
		BaseIRExpression: NewBaseIRExpression(Int32),
		Value:            v,
	}
}

func (i *IR_Int32) ReturnType(ctx *IR_Context) Type {
	return TInt32
}

func (i *IR_Int32) String() string {
	return fmt.Sprintf("%d", i.Value)
}

func (b *IR_Int32) SSA_Transform(ctx *SSA_Context) (SSA_Rewrites, IRExpression) {
	return nil, b
}
