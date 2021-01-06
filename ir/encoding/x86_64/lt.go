package x86_64

import (
	"github.com/bspaans/jit-compiler/asm/x86_64"
	"github.com/bspaans/jit-compiler/ir/expr"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

func encode_IR_LT(i *expr.IR_LT, ctx *IR_Context, target lib.Operand, includeSETE bool) ([]lib.Instruction, error) {
	return order(i.Op1, i.Op2, ctx, target, includeSETE, i.String(), x86_64.SETB, x86_64.SETL)
}
