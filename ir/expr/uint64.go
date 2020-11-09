package expr

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm"
	"github.com/bspaans/jit-compiler/asm/encoding"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
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

func (i *IR_Uint64) Encode(ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {

	result := []lib.Instruction{asm.MOV_immediate(i.Value, target)}
	ctx.AddInstructions(result)
	return result, nil
}
