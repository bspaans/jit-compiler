package expr

import (
	"github.com/bspaans/jit/asm"
	"github.com/bspaans/jit/asm/encoding"
	. "github.com/bspaans/jit/ir/shared"
	"github.com/bspaans/jit/lib"
)

func NewIR_Add(op1, op2 IRExpression) IRExpression {
	op := func(op1, op2 encoding.Operand) lib.Instruction {
		return asm.ADD(op1, op2)
	}
	return NewIR_Operator(op, "+", op1, op2)
}

func NewIR_Sub(op1, op2 IRExpression) IRExpression {
	op := func(op1, op2 encoding.Operand) lib.Instruction {
		return asm.SUB(op1, op2)
	}
	return NewIR_Operator(op, "-", op1, op2)
}

func NewIR_Mul(op1, op2 IRExpression) IRExpression {
	op := func(op1, op2 encoding.Operand) lib.Instruction {
		return asm.MUL(op1, op2)
	}
	return NewIR_Operator(op, "*", op1, op2)
}

func NewIR_Div(op1, op2 IRExpression) IRExpression {
	op := func(op1, op2 encoding.Operand) lib.Instruction {
		return asm.DIV(op1, op2)
	}
	return NewIR_Operator(op, "/", op1, op2)
}
