package expr

import (
	"fmt"

	"github.com/bspaans/jit/asm"
	. "github.com/bspaans/jit/ir/shared"
)

type IR_Uint64 struct {
	*BaseIRExpression
	Value uint64
}

func NewIR_Uint64(v uint64) *IR_Uint64 {
	return &IR_Uint64{
		BaseIRExpression: NewBaseIRExpression(Uint64),
		Value:            v,
	}
}

func (i *IR_Uint64) ReturnType(ctx *IR_Context) Type {
	return TUint64
}

func (i *IR_Uint64) String() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *IR_Uint64) Encode(ctx *IR_Context, target *asm.Register) ([]asm.Instruction, error) {
	result := []asm.Instruction{&asm.MOV{asm.Uint64(i.Value), target}}
	ctx.AddInstructions(result)
	return result, nil
}
