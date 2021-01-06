package x86_64

import (
	"github.com/bspaans/jit-compiler/asm/x86_64"
	"github.com/bspaans/jit-compiler/asm/x86_64/encoding"
	"github.com/bspaans/jit-compiler/ir/expr"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

func encode_IR_Float64(i *expr.IR_Float64, ctx *IR_Context, target lib.Operand) ([]lib.Instruction, error) {
	tmp := ctx.AllocateRegister(TUint64)
	defer ctx.DeallocateRegister(tmp)

	result := []lib.Instruction{
		x86_64.MOV(encoding.Float64(i.Value), tmp),
		x86_64.MOV(tmp, target),
	}
	ctx.AddInstruction(result...)
	return result, nil
}
