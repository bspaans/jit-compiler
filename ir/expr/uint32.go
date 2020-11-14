package expr

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm"
	"github.com/bspaans/jit-compiler/asm/encoding"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
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

func (i *IR_Uint32) Encode(ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {

	result := []lib.Instruction{asm.MOV(encoding.Uint32(i.Value), target)}
	ctx.AddInstructions(result)
	return result, nil
}
