package expr

import (
	"github.com/bspaans/jit/asm"
	. "github.com/bspaans/jit/ir/shared"
)

func NewIR_Add(op1, op2 IRExpression) IRExpression {
	return NewIR_Operator(asm.ADD, "+", op1, op2)
}

func NewIR_Sub(op1, op2 IRExpression) IRExpression {
	return NewIR_Operator(asm.SUB, "-", op1, op2)
}

func NewIR_Mul(op1, op2 IRExpression) IRExpression {
	return NewIR_Operator(asm.MUL, "*", op1, op2)
}

func NewIR_Div(op1, op2 IRExpression) IRExpression {
	return NewIR_Operator(asm.DIV, "/", op1, op2)
}
