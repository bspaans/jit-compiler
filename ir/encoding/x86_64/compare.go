package x86_64

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm/x86_64"
	"github.com/bspaans/jit-compiler/asm/x86_64/encoding"
	"github.com/bspaans/jit-compiler/ir/expr"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

func compare(op1, op2 IRExpression, ctx *IR_Context) ([]lib.Instruction, error) {

	result := []lib.Instruction{}
	returnType1, returnType2 := op1.ReturnType(ctx), op2.ReturnType(ctx)
	if returnType1 != returnType2 {
		return nil, fmt.Errorf("Unsupported types (%s, %s) in compare operation", returnType1, returnType2)
	}

	var reg1, reg2 lib.Operand

	if op1.Type() == Variable {
		variable := op1.(*expr.IR_Variable).Value
		reg1 = ctx.VariableMap[variable]
	} else {
		reg1 = ctx.AllocateRegister(returnType1)
		defer ctx.DeallocateRegister(reg1.(*encoding.Register))
		expr1, err := encodeExpression(op1, ctx, reg1)
		if err != nil {
			return nil, err
		}
		result = lib.Instructions(result).Add(expr1)
	}

	if op2.Type() == Variable {
		variable := op2.(*expr.IR_Variable).Value
		reg2 = ctx.VariableMap[variable]
	} else {
		reg2 = ctx.AllocateRegister(returnType1)
		defer ctx.DeallocateRegister(reg2.(*encoding.Register))
		expr2, err := encodeExpression(op2, ctx, reg2)
		if err != nil {
			return nil, err
		}
		result = lib.Instructions(result).Add(expr2)
	}
	cmp := x86_64.CMP(reg2, reg1)
	result = append(result, cmp)
	ctx.AddInstruction(cmp)
	return result, nil
}

type orderOpcode func(lib.Operand) lib.Instruction

func order(op1, op2 IRExpression, ctx *IR_Context, target lib.Operand, includeSETE bool, repr string, unsignedOp, signedOp orderOpcode) ([]lib.Instruction, error) {

	result, err := compare(op1, op2, ctx)
	if err != nil {
		return nil, fmt.Errorf("%s in %s", err.Error(), repr)
	}
	if includeSETE {
		tmpReg := ctx.AllocateRegister(TUint64)
		defer ctx.DeallocateRegister(tmpReg)
		reg8 := tmpReg.(*encoding.Register).Get8BitRegister()
		sete := unsignedOp(reg8)
		if IsSignedInteger(op1.ReturnType(ctx)) {
			sete = signedOp(reg8)
		}
		mov := x86_64.MOV(tmpReg.(*encoding.Register).ForOperandWidth(target.Width()), target)
		result = append(result, sete)
		result = append(result, mov)
		ctx.AddInstruction(sete)
		ctx.AddInstruction(mov)
	}
	return result, nil
}
