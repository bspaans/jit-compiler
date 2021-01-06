package x86_64

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm/x86_64"
	"github.com/bspaans/jit-compiler/asm/x86_64/encoding"
	"github.com/bspaans/jit-compiler/ir/expr"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

func encode_IR_ArrayIndex(i *expr.IR_ArrayIndex, ctx *IR_Context, target lib.Operand) ([]lib.Instruction, error) {
	ctx.AddInstruction("array_index " + encoding.Comment(i.String()))

	itemWidth := i.ReturnType(ctx).Width()

	var arrayReg, indexReg lib.Operand
	if i.Array.Type() == Variable {
		variable := i.Array.(*expr.IR_Variable).Value
		reg, ok := ctx.VariableMap[variable]
		if !ok {
			return nil, fmt.Errorf("Unknown variable %s", variable)
		}
		arrayReg = reg
	} else {
		arrayReg = ctx.AllocateRegister(TUint64)
		defer ctx.DeallocateRegister(arrayReg)
	}

	result, err := encodeExpression(i.Array, ctx, arrayReg)
	if err != nil {
		return nil, fmt.Errorf("Array encoding issue: %s", err.Error())
	}

	// Specialise for integers
	if i.Index.Type() == Uint64 {
		op := i.Index.(*expr.IR_Uint64)
		if op.Value == 0 {
			mov := x86_64.MOV(&encoding.IndirectRegister{
				arrayReg.(*encoding.Register).ForOperandWidth(itemWidth)},
				target.(*encoding.Register).ForOperandWidth(itemWidth))
			// Move 0 into target register if going from a wider to narrower register
			mov0 := x86_64.MOV(encoding.Uint64(0), target.(*encoding.Register).Get64BitRegister()) // TODO use xor reg, reg
			if itemWidth < lib.QUADWORD {
				ctx.AddInstruction(mov0)
				result = append(result, mov0)
			}
			ctx.AddInstruction(mov)
			result = append(result, mov)
			return result, nil
		} else {
			// TODO add index*itemwidth to arrayReg
		}
	}

	if i.Index.Type() == Variable {
		variable := i.Index.(*expr.IR_Variable).Value
		reg, ok := ctx.VariableMap[variable]
		if !ok {
			return nil, fmt.Errorf("Unknown variable %s", variable)
		}
		indexReg = reg.(*encoding.Register)
	} else {
		indexReg = ctx.AllocateRegister(TUint64)
		defer ctx.DeallocateRegister(indexReg)
	}

	index, err := encodeExpression(i.Index, ctx, indexReg)
	if err != nil {
		return nil, fmt.Errorf("Array index encoding issue: %s", err.Error())
	}
	result = lib.Instructions(result).Add(index)
	mov := x86_64.MOV(&encoding.SIBRegister{
		arrayReg.(*encoding.Register),
		indexReg.(*encoding.Register),
		encoding.ScaleForItemWidth(itemWidth)},
		target.(*encoding.Register).ForOperandWidth(itemWidth))

	// Move 0 into target register if going from a wider to narrower register
	mov0 := x86_64.MOV(encoding.Uint64(0), target.(*encoding.Register).Get64BitRegister())
	if itemWidth < lib.QUADWORD {
		ctx.AddInstruction(mov0)
		result = append(result, mov0)
	}

	ctx.AddInstruction(mov)
	result = append(result, mov)
	return result, nil
}
