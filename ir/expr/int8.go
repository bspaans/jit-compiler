package expr

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm/x86_64"
	"github.com/bspaans/jit-compiler/asm/x86_64/encoding"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

type IR_Int8 struct {
	*BaseIRExpression
	Value int8
}

func NewIR_Int8(v int8) *IR_Int8 {
	return &IR_Int8{
		BaseIRExpression: NewBaseIRExpression(Int8),
		Value:            v,
	}
}

func (i *IR_Int8) ReturnType(ctx *IR_Context) Type {
	return TInt8
}

func (i *IR_Int8) String() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *IR_Int8) Encode(ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {
	result := []lib.Instruction{asm.MOV(encoding.Uint8(int8(i.Value)), target)}
	ctx.AddInstructions(result)
	return result, nil
}
func (b *IR_Int8) SSA_Transform(ctx *SSA_Context) (SSA_Rewrites, IRExpression) {
	return nil, b
}
