package aarch64

import (
	"github.com/bspaans/jit-compiler/asm/aarch64"
	"github.com/bspaans/jit-compiler/ir/expr"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

func encode_IR_Add(i *expr.IR_Add, ctx *IR_Context, target lib.Operand) ([]lib.Instruction, error) {

	return encode_Operator(i.Op1, i.Op2, aarch64.ADD, i.String(), ctx, target)
}
