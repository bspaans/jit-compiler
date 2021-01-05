package expr

import (
	"fmt"

	. "github.com/bspaans/jit-compiler/ir/shared"
)

type IR_Int8 struct {
	*BaseIRExpression
	Value int8
}

func NewIR_Int8(v int8) *IR_Int8 {
	return &IR_Int8{
		BaseIRExpression: NewBaseIRExpression(Int8),
		Value:            v,
	}
}

func (i *IR_Int8) ReturnType(ctx *IR_Context) Type {
	return TInt8
}

func (i *IR_Int8) String() string {
	return fmt.Sprintf("%d", i.Value)
}

func (b *IR_Int8) SSA_Transform(ctx *SSA_Context) (SSA_Rewrites, IRExpression) {
	return nil, b
}
