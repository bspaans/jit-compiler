package expr

import (
	"fmt"

	"github.com/bspaans/jit/asm"
	. "github.com/bspaans/jit/ir/shared"
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

func (i *IR_Float64) Encode(ctx *IR_Context, target asm.Operand) ([]asm.Instruction, error) {
	tmp := ctx.AllocateRegister(TUint64)
	defer ctx.DeallocateRegister(tmp)

	result := []asm.Instruction{
		&asm.MOV{asm.Float64(i.Value), tmp},
		&asm.MOVQ{tmp, target}}
	ctx.AddInstructions(result)
	return result, nil
}
