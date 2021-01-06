package x86_64

import (
	"github.com/bspaans/jit-compiler/asm/x86_64"
	"github.com/bspaans/jit-compiler/asm/x86_64/encoding"
	"github.com/bspaans/jit-compiler/ir/expr"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

func encode_IR_Bool(i *expr.IR_Bool, ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {
	value := 0
	if i.Value {
		value = 1
	}
	result := []lib.Instruction{x86_64.MOV_immediate(uint64(value), target)}
	ctx.AddInstruction(result...)
	return result, nil
}
