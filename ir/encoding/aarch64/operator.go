package aarch64

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm/aarch64/encoding"
	"github.com/bspaans/jit-compiler/ir/expr"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

type op func(o1, o2, op3 lib.Operand) lib.Instruction

func encode_Operator(op1, op2 IRExpression, operator op, repr string, ctx *IR_Context, target lib.Operand) ([]lib.Instruction, error) {
	returnType1, returnType2 := op1.ReturnType(ctx), op2.ReturnType(ctx)
	if returnType1 == returnType2 && IsNumber(returnType1) {
		result, err := encodeExpression(op1, ctx, target)
		if err != nil {
			return nil, err
		}

		var reg lib.Operand
		if op2.Type() == Variable {
			variable := op2.(*expr.IR_Variable).Value
			reg = ctx.VariableMap[variable]
		} else {
			reg = ctx.AllocateRegister(returnType2)
			defer ctx.DeallocateRegister(reg.(*encoding.Register))

			expr, err := encodeExpression(op2, ctx, reg)
			if err != nil {
				return nil, err
			}
			result = lib.Instructions(result).Add(expr)
		}
		instr := operator(reg, reg, target)
		ctx.AddInstruction(instr)
		result = append(result, instr)
		return result, nil
	}
	return nil, fmt.Errorf("Unsupported types (%s, %s) in IR operation: %s", returnType1, returnType2, repr)
}
