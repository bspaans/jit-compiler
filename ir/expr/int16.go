package expr

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm"
	"github.com/bspaans/jit-compiler/asm/encoding"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

type IR_Int16 struct {
	*BaseIRExpression
	Value int16
}

func NewIR_Int16(v int16) *IR_Int16 {
	return &IR_Int16{
		BaseIRExpression: NewBaseIRExpression(Int16),
		Value:            v,
	}
}

func (i *IR_Int16) ReturnType(ctx *IR_Context) Type {
	return TInt16
}

func (i *IR_Int16) String() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *IR_Int16) Encode(ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {

	result := []lib.Instruction{asm.MOV(encoding.Uint16(uint16(i.Value)), target)}
	ctx.AddInstructions(result)
	return result, nil
}
func (b *IR_Int16) SSA_Transform(ctx *SSA_Context) (SSA_Rewrites, IRExpression) {
	return nil, b
}
