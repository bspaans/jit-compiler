package expr

import (
	"fmt"

	. "github.com/bspaans/jit-compiler/ir/shared"
)

type IR_Int64 struct {
	*BaseIRExpression
	Value int64
}

func NewIR_Int64(v int64) *IR_Int64 {
	return &IR_Int64{
		BaseIRExpression: NewBaseIRExpression(Int64),
		Value:            v,
	}
}

func (i *IR_Int64) ReturnType(ctx *IR_Context) Type {
	return TInt64
}

func (i *IR_Int64) String() string {
	return fmt.Sprintf("%d", i.Value)
}

func (b *IR_Int64) SSA_Transform(ctx *SSA_Context) (SSA_Rewrites, IRExpression) {
	return nil, b
}
