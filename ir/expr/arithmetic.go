package expr

import (
	"github.com/bspaans/jit/asm"
)

func NewIR_Add(op1, op2 IRExpression) IRExpression {
	op := func(op1, op2 asm.Operand) asm.Instruction {
		return &asm.ADD{op1, op2}
	}
	return NewIR_Operator(op, "+", op1, op2)
}

func NewIR_Sub(op1, op2 IRExpression) IRExpression {
	op := func(op1, op2 asm.Operand) asm.Instruction {
		return &asm.SUB{op1, op2}
	}
	return NewIR_Operator(op, "-", op1, op2)
}

func NewIR_Mul(op1, op2 IRExpression) IRExpression {
	op := func(op1, op2 asm.Operand) asm.Instruction {
		return &asm.MUL{op1, op2}
	}
	return NewIR_Operator(op, "*", op1, op2)
}

func NewIR_Div(op1, op2 IRExpression) IRExpression {
	op := func(op1, op2 asm.Operand) asm.Instruction {
		return &asm.DIV{op1, op2}
	}
	return NewIR_Operator(op, "/", op1, op2)
}
