package x86_64

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm/x86_64"
	"github.com/bspaans/jit-compiler/asm/x86_64/encoding"
	"github.com/bspaans/jit-compiler/ir/expr"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

func encode_IR_And(i *expr.IR_And, ctx *IR_Context, target lib.Operand) ([]lib.Instruction, error) {
	ctx.AddInstruction("operator " + encoding.Comment(i.String()))
	returnType1, returnType2 := i.Op1.ReturnType(ctx), i.Op2.ReturnType(ctx)
	if returnType1 != returnType2 {
		return nil, fmt.Errorf("Unsupported types (%s, %s) in && IR operation: %s", returnType1, returnType2, i.String())
	}
	if returnType1 != TBool {
		return nil, fmt.Errorf("Unsupported types (%s, %s) in && IR operation: %s", returnType1, returnType2, i.String())
	}

	reg := ctx.AllocateRegister(returnType1)
	defer ctx.DeallocateRegister(reg)

	result, err := encodeExpression(i.Op1, ctx, reg)
	if err != nil {
		return nil, err
	}
	expr2, err := encodeExpression(i.Op2, ctx, target)
	if err != nil {
		return nil, err
	}
	result = lib.Instructions(result).Add(expr2)
	// TODO: should be using test?
	and := x86_64.AND(reg, target)
	result = append(result, and)
	ctx.AddInstruction(and)
	return result, nil
}
