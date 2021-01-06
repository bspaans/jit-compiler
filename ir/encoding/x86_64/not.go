package x86_64

import (
	"fmt"

	. "github.com/bspaans/jit-compiler/ir/shared"

	"github.com/bspaans/jit-compiler/asm/x86_64"
	"github.com/bspaans/jit-compiler/asm/x86_64/encoding"
	"github.com/bspaans/jit-compiler/ir/expr"
	"github.com/bspaans/jit-compiler/lib"
)

func encode_IR_Not(i *expr.IR_Not, ctx *IR_Context, target lib.Operand, includeSETE bool) ([]lib.Instruction, error) {

	var reg1 lib.Operand

	result := []lib.Instruction{}

	switch c := i.Op1.(type) {
	case *expr.IR_Variable:
		variable := c.Value
		reg1 = ctx.VariableMap[variable]
		if !includeSETE {
			cmp := x86_64.CMP_immediate(1, reg1)
			result = append(result, cmp)
			ctx.AddInstruction(cmp)
		}
	case *expr.IR_Bool:
		value := uint64(0)
		if c.Value {
			value = 1
		}
		reg1 = ctx.AllocateRegister(TUint64)
		defer ctx.DeallocateRegister(reg1.(*encoding.Register))
		mov := x86_64.MOV_immediate(value, reg1)
		cmp := x86_64.CMP_immediate(1, reg1)
		result = append(result, mov, cmp)
		ctx.AddInstruction(mov, cmp)
	case *expr.IR_Equals:
		var err error
		var eq []lib.Instruction
		eq, err = encode_IR_Equals(c, ctx, target, includeSETE)
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
		eq, err = encode_IR_LT(c, ctx, target, includeSETE)
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
		eq, err = encode_IR_LTE(c, ctx, target, includeSETE)
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
		eq, err = encode_IR_GT(c, ctx, target, includeSETE)
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
		eq, err = encode_IR_GTE(c, ctx, target, includeSETE)
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
		eq, err = encode_IR_And(c, ctx, target)
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
		eq, err = encode_IR_Or(c, ctx, target)
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

	if includeSETE {
		tmpReg := ctx.AllocateRegister(TUint64)
		defer ctx.DeallocateRegister(tmpReg)
		instr := []lib.Instruction{}
		// TODO: use test?
		instr = append(instr, x86_64.XOR(tmpReg, tmpReg))
		instr = append(instr, x86_64.CMP_immediate(0, reg1))
		instr = append(instr, x86_64.SETE(tmpReg.(*encoding.Register).Get8BitRegister()))
		instr = append(instr, x86_64.MOV(tmpReg.(*encoding.Register).ForOperandWidth(target.Width()), target))
		for _, inst := range instr {
			result = append(result, inst)
			ctx.AddInstruction(inst)
		}
	}
	return result, nil
}
