package expr

import (
	"github.com/bspaans/jit/asm"
	"github.com/bspaans/jit/asm/encoding"
	. "github.com/bspaans/jit/ir/shared"
	"github.com/bspaans/jit/lib"
)

type IR_Variable struct {
	*BaseIRExpression
	Value string
}

func NewIR_Variable(v string) *IR_Variable {
	return &IR_Variable{
		BaseIRExpression: NewBaseIRExpression(Variable),
		Value:            v,
	}
}

func (i *IR_Variable) ReturnType(ctx *IR_Context) Type {
	return ctx.VariableTypes[i.Value]
}

func (i *IR_Variable) String() string {
	return i.Value
}

func (i *IR_Variable) Encode(ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {
	reg := ctx.VariableMap[i.Value]
	var result []lib.Instruction
	if i.ReturnType(ctx) == TFloat64 {
		result = []lib.Instruction{&asm.MOVSD{reg, target}}
	} else {
		result = []lib.Instruction{&asm.MOV{reg, target}}
	}
	ctx.AddInstructions(result)
	return result, nil
}
