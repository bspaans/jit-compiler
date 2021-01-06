package x86_64

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm/x86_64"
	"github.com/bspaans/jit-compiler/asm/x86_64/encoding"
	"github.com/bspaans/jit-compiler/ir/expr"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/ir/statements"
	"github.com/bspaans/jit-compiler/lib"
)

func encode_IR_Return(i *statements.IR_Return, ctx *IR_Context) ([]lib.Instruction, error) {
	result := []lib.Instruction{}
	var reg lib.Operand
	var ok bool
	if i.Expr.Type() == Variable {
		reg, ok = ctx.VariableMap[i.Expr.(*expr.IR_Variable).Value]
		if !ok {
			return nil, fmt.Errorf("Unknown variable '%s' in return expression: %s", i.Expr.(*expr.IR_Variable).Value, i.String())
		}
	} else {
		reg = ctx.AllocateRegister(i.Expr.ReturnType(ctx))
		defer ctx.DeallocateRegister(reg)
		result_, err := encodeExpression(i.Expr, ctx, reg)
		if err != nil {
			return nil, err
		}
		result = result_
	}
	if reg.Width() != lib.QUADWORD {
		cast := ctx.AllocateRegister(TUint64)
		defer ctx.DeallocateRegister(cast)

		if reg.Width() == lib.BYTE || reg.Width() == lib.WORD {
			movzx := x86_64.MOVZX(reg, cast)
			// TODO? use movsx for signed integers?
			//if shared.IsSignedInteger(i.Expr.ReturnType(ctx)) {
			//	movzx = x86_64.MOVSX(reg, cast)
			//}
			result = append(result, movzx)
			ctx.AddInstruction(movzx)
		} else {
			xor := x86_64.XOR(cast, cast)
			mov := x86_64.MOV(reg, cast.(*encoding.Register).ForOperandWidth(reg.Width()))
			result = append(result, xor)
			result = append(result, mov)
			ctx.AddInstruction(mov)
			ctx.AddInstruction(xor)
		}
		reg = cast
	}
	target := ctx.PeekReturn()
	instr := []lib.Instruction{
		x86_64.MOV(reg.(*encoding.Register).Get64BitRegister(), target),
		x86_64.RETURN(),
	}
	for _, inst := range instr {
		ctx.AddInstruction(inst)
		result = append(result, inst)
	}
	return result, nil
}
