package expr

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm"
	"github.com/bspaans/jit-compiler/asm/encoding"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
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

func (i *IR_Int32) Encode(ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {

	result := []lib.Instruction{asm.MOV(encoding.Uint32(uint32(i.Value)), target)}
	ctx.AddInstructions(result)
	return result, nil
}
