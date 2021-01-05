package x86_64

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm/x86_64"
	"github.com/bspaans/jit-compiler/asm/x86_64/encoding"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/ir/statements"
	"github.com/bspaans/jit-compiler/lib"
)

func encode_IR_ArrayAssignment(i *statements.IR_ArrayAssignment, ctx *IR_Context) ([]lib.Instruction, error) {
	ctx.AddInstruction("array_assignment " + encoding.Comment(i.String()))

	returnType := i.Expr.ReturnType(ctx)
	itemWidth := returnType.Width()

	indexReg := ctx.AllocateRegister(TUint64)
	defer ctx.DeallocateRegister(indexReg)
	exprReg := ctx.AllocateRegister(returnType)
	defer ctx.DeallocateRegister(exprReg)

	reg, found := ctx.VariableMap[i.Variable]
	if !found {
		return nil, fmt.Errorf("Unknown array '%s'", i.Variable)
	}

	result, err := encodeExpression(i.Index, ctx, indexReg)
	if err != nil {
		return nil, fmt.Errorf("Failed to encode array index in %s: %s", i.String(), err.Error())
	}

	exprInstr, err := encodeExpression(i.Expr, ctx, exprReg)
	if err != nil {
		return nil, fmt.Errorf("Failed to encode expr in %s: %s", i.String(), err.Error())
	}
	result = lib.Instructions(result).Add(exprInstr)

	target := &encoding.SIBRegister{reg.(*encoding.Register), indexReg, encoding.ScaleForItemWidth(itemWidth)}
	mov := x86_64.MOV(exprReg.ForOperandWidth(itemWidth), target)
	ctx.AddInstruction(mov)
	result = append(result, mov)
	return result, nil
}
