package expr

import (
	"fmt"

	. "github.com/bspaans/jit-compiler/ir/shared"
)

type IR_Uint16 struct {
	*BaseIRExpression
	Value uint16
}

func NewIR_Uint16(v uint16) *IR_Uint16 {
	return &IR_Uint16{
		BaseIRExpression: NewBaseIRExpression(Uint16),
		Value:            v,
	}
}

func (i *IR_Uint16) ReturnType(ctx *IR_Context) Type {
	return TUint16
}

func (i *IR_Uint16) String() string {
	return fmt.Sprintf("%d", i.Value)
}

func (b *IR_Uint16) SSA_Transform(ctx *SSA_Context) (SSA_Rewrites, IRExpression) {
	return nil, b
}
