package expr

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm"
	"github.com/bspaans/jit-compiler/asm/encoding"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
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
	reg, ok := ctx.VariableMap[i.Value]
	if !ok || reg == nil {
		return nil, fmt.Errorf("Unknown variable '%s'", i.Value)
	}
	result := []lib.Instruction{asm.MOV(reg, target)}
	ctx.AddInstructions(result)
	return result, nil
}
func (b *IR_Variable) SSA_Transform(ctx *SSA_Context) (SSA_Rewrites, IRExpression) {
	return nil, b
}
