package x86_64

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm/x86_64"
	"github.com/bspaans/jit-compiler/asm/x86_64/encoding"
	"github.com/bspaans/jit-compiler/ir/expr"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

func encode_IR_GT(i *expr.IR_GT, ctx *IR_Context, target encoding.Operand, includeSETE bool) ([]lib.Instruction, error) {
	resugt, err := Compare(i.Op1, i.Op2, ctx)
	if err != nil {
		return nil, fmt.Errorf("%s in %s", err.Error(), i.String())
	}
	if includeSETE {
		tmpReg := ctx.AllocateRegister(TUint64)
		defer ctx.DeallocateRegister(tmpReg)
		// TODO xor tmpreg
		// TODO use right SET depending on sign
		sete := x86_64.SETA(tmpReg.Get8BitRegister())
		if IsSignedInteger(i.Op1.ReturnType(ctx)) {
			sete = x86_64.SETG(tmpReg.Get8BitRegister())
		}
		mov := x86_64.MOV(tmpReg.ForOperandWidth(target.Width()), target)
		resugt = append(resugt, sete)
		resugt = append(resugt, mov)
		ctx.AddInstruction(sete)
		ctx.AddInstruction(mov)
	}
	return resugt, nil
}
