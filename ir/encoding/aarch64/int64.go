package aarch64

import (
	"github.com/bspaans/jit-compiler/asm/aarch64"
	"github.com/bspaans/jit-compiler/asm/aarch64/encoding"
	"github.com/bspaans/jit-compiler/ir/expr"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

func encode_IR_Int64(i *expr.IR_Int64, ctx *IR_Context, target lib.Operand) ([]lib.Instruction, error) {

	result := []lib.Instruction{aarch64.MOVK(encoding.Uint64(uint64(i.Value)), target)}
	ctx.AddInstruction(result...)
	return result, nil
}
