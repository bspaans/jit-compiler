package expr

import (
	. "github.com/bspaans/jit-compiler/ir/shared"
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
	if _, ok := ctx.VariableTypes[i.Value]; !ok {
		panic("Unknown variable: " + i.Value)
	}
	return ctx.VariableTypes[i.Value]
}

func (i *IR_Variable) String() string {
	return i.Value
}

func (b *IR_Variable) SSA_Transform(ctx *SSA_Context) (SSA_Rewrites, IRExpression) {
	return nil, b
}
