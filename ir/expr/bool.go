package expr

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm/x86_64"
	"github.com/bspaans/jit-compiler/asm/x86_64/encoding"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

type IR_Bool struct {
	*BaseIRExpression
	Value bool
}

func NewIR_Bool(v bool) *IR_Bool {
	return &IR_Bool{
		BaseIRExpression: NewBaseIRExpression(Bool),
		Value:            v,
	}
}

func (i *IR_Bool) ReturnType(ctx *IR_Context) Type {
	return TBool
}

func (i *IR_Bool) String() string {
	return fmt.Sprintf("%v", i.Value)
}

func (i *IR_Bool) Encode(ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {
	value := 0
	if i.Value {
		value = 1
	}
	result := []lib.Instruction{asm.MOV_immediate(uint64(value), target)}
	ctx.AddInstructions(result)
	return result, nil
}
func (b *IR_Bool) SSA_Transform(ctx *SSA_Context) (SSA_Rewrites, IRExpression) {
	return nil, b
}
