package expr

import (
	"fmt"

	. "github.com/bspaans/jit-compiler/ir/shared"
)

type IR_Uint32 struct {
	*BaseIRExpression
	Value uint32
}

func NewIR_Uint32(v uint32) *IR_Uint32 {
	return &IR_Uint32{
		BaseIRExpression: NewBaseIRExpression(Uint32),
		Value:            v,
	}
}

func (i *IR_Uint32) ReturnType(ctx *IR_Context) Type {
	return TUint32
}

func (i *IR_Uint32) String() string {
	return fmt.Sprintf("%d", i.Value)
}

func (b *IR_Uint32) SSA_Transform(ctx *SSA_Context) (SSA_Rewrites, IRExpression) {
	return nil, b
}
