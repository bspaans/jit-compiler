package expr

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm"
	"github.com/bspaans/jit-compiler/asm/encoding"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
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

func (i *IR_Int64) Encode(ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {

	result := []lib.Instruction{asm.MOV_immediate(uint64(i.Value), target)}
	ctx.AddInstructions(result)
	return result, nil
}
