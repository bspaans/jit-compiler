package x86_64

import (
	"fmt"

	. "github.com/bspaans/jit-compiler/ir/shared"

	"github.com/bspaans/jit-compiler/asm/x86_64"
	"github.com/bspaans/jit-compiler/asm/x86_64/encoding"
	"github.com/bspaans/jit-compiler/ir/expr"
	"github.com/bspaans/jit-compiler/lib"
)

func encode_IR_Not(i *expr.IR_Not, ctx *IR_Context, target encoding.Operand, includeSETE bool) ([]lib.Instruction, error) {

	var reg1 encoding.Operand

	result := []lib.Instruction{}
	switch i.Op1.(type) {
	case *expr.IR_Variable:
		variable := i.Op1.(*expr.IR_Variable).Value
		reg1 = ctx.VariableMap[variable]
	case *expr.IR_Uint64:
		value := i.Op1.(*expr.IR_Uint64).Value
		reg1 = ctx.AllocateRegister(TUint64)
		defer ctx.DeallocateRegister(reg1.(*encoding.Register))
		result = append(result, x86_64.MOV(encoding.Uint64(value), reg1))
	case *expr.IR_Equals:
		var err error
		var eq []lib.Instruction
		if includeSETE {
			eq, err = encodeExpression(i.Op1, ctx, target)
		} else {
			eq, err = encode_IR_Equals(i.Op1.(*expr.IR_Equals), ctx, target, false)
		}
		if err != nil {
			return nil, err
		}
		for _, r := range eq {
			result = append(result, r)
		}
		reg1 = target
		// TODO: introduce logical operator interface
	case *expr.IR_LT:
		var err error
		var eq []lib.Instruction
		if includeSETE {
			eq, err = encodeExpression(i.Op1, ctx, target)
		} else {
			eq, err = encode_IR_LT(i.Op1.(*expr.IR_LT), ctx, target, false)
		}
		if err != nil {
			return nil, err
		}
		for _, r := range eq {
			result = append(result, r)
		}
		reg1 = target
		// TODO: introduce logical operator interface
	case *expr.IR_LTE:
		var err error
		var eq []lib.Instruction
		if includeSETE {
			eq, err = encodeExpression(i.Op1, ctx, target)
		} else {
			eq, err = encode_IR_LTE(i.Op1.(*expr.IR_LTE), ctx, target, false)
		}
		if err != nil {
			return nil, err
		}
		for _, r := range eq {
			result = append(result, r)
		}
		reg1 = target
	case *expr.IR_GT:
		var err error
		var eq []lib.Instruction
		if includeSETE {
			eq, err = encodeExpression(i.Op1, ctx, target)
		} else {
			eq, err = encode_IR_GT(i.Op1.(*expr.IR_GT), ctx, target, false)
		}
		if err != nil {
			return nil, err
		}
		for _, r := range eq {
			result = append(result, r)
		}
		reg1 = target
		// TODO: introduce logical operator interface
	case *expr.IR_GTE:
		var err error
		var eq []lib.Instruction
		if includeSETE {
			eq, err = encodeExpression(i.Op1, ctx, target)
		} else {
			eq, err = encode_IR_GTE(i.Op1.(*expr.IR_GTE), ctx, target, false)
		}
		if err != nil {
			return nil, err
		}
		for _, r := range eq {
			result = append(result, r)
		}
		reg1 = target
	case *expr.IR_And:
		var err error
		var eq []lib.Instruction
		if includeSETE {
			eq, err = encodeExpression(i.Op1, ctx, target)
		} else {
			eq, err = encode_IR_And(i.Op1.(*expr.IR_And), ctx, target)
		}
		if err != nil {
			return nil, err
		}
		for _, r := range eq {
			result = append(result, r)
		}
		reg1 = target
	case *expr.IR_Or:
		var err error
		var eq []lib.Instruction
		if includeSETE {
			eq, err = encodeExpression(i.Op1, ctx, target)
		} else {
			eq, err = encode_IR_Or(i.Op1.(*expr.IR_Or), ctx, target)
		}
		if err != nil {
			return nil, err
		}
		for _, r := range eq {
			result = append(result, r)
		}
		reg1 = target
	default:
		return nil, fmt.Errorf("Unsupported ! operation: %s", i.Op1.Type())
	}

	tmpReg := ctx.AllocateRegister(TUint64)
	defer ctx.DeallocateRegister(tmpReg)
	instr := []lib.Instruction{}
	if includeSETE {
		// TODO: use test?
		instr = append(instr, x86_64.CMP_immediate(0, reg1))
		instr = append(instr, x86_64.XOR(tmpReg, tmpReg))
		instr = append(instr, x86_64.SETE(tmpReg.Get8BitRegister()))
		instr = append(instr, x86_64.MOV(tmpReg.ForOperandWidth(target.Width()), target))
	}
	for _, inst := range instr {
		result = append(result, inst)
		ctx.AddInstruction(inst)
	}
	return result, nil
}
