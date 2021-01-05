package expr

import (
	"fmt"

	. "github.com/bspaans/jit-compiler/ir/shared"
)

type IR_Uint8 struct {
	*BaseIRExpression
	Value uint8
}

func NewIR_Uint8(v uint8) *IR_Uint8 {
	return &IR_Uint8{
		BaseIRExpression: NewBaseIRExpression(Uint8),
		Value:            v,
	}
}

func (i *IR_Uint8) ReturnType(ctx *IR_Context) Type {
	return TUint8
}

func (i *IR_Uint8) String() string {
	return fmt.Sprintf("%d", i.Value)
}

func (b *IR_Uint8) SSA_Transform(ctx *SSA_Context) (SSA_Rewrites, IRExpression) {
	return nil, b
}
