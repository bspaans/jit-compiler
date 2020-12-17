package expr

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm"
	"github.com/bspaans/jit-compiler/asm/encoding"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
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

func (i *IR_Uint16) Encode(ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {

	result := []lib.Instruction{asm.MOV(encoding.Uint16(i.Value), target)}
	ctx.AddInstructions(result)
	return result, nil
}

func (b *IR_Uint16) SSA_Transform(ctx *SSA_Context) (SSA_Rewrites, IRExpression) {
	return nil, b
}
