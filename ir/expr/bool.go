package expr

import (
	"fmt"

	"github.com/bspaans/jit/asm"
	"github.com/bspaans/jit/asm/encoding"
	. "github.com/bspaans/jit/ir/shared"
	"github.com/bspaans/jit/lib"
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
	result := []lib.Instruction{&asm.MOV{encoding.Uint64(value), target}}
	ctx.AddInstructions(result)
	return result, nil
}
