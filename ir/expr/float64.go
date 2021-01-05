package expr

import (
	"fmt"

	. "github.com/bspaans/jit-compiler/ir/shared"
)

type IR_Float64 struct {
	*BaseIRExpression
	Value float64
}

func NewIR_Float64(v float64) *IR_Float64 {
	return &IR_Float64{
		BaseIRExpression: NewBaseIRExpression(Float64),
		Value:            v,
	}
}

func (i *IR_Float64) ReturnType(ctx *IR_Context) Type {
	return TFloat64
}

func (i *IR_Float64) String() string {
	return fmt.Sprintf("%f", i.Value)
}

func (b *IR_Float64) SSA_Transform(ctx *SSA_Context) (SSA_Rewrites, IRExpression) {
	return nil, b
}
