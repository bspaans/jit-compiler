package x86_64

import (
	"github.com/bspaans/jit-compiler/asm/x86_64"
	"github.com/bspaans/jit-compiler/asm/x86_64/encoding"
	"github.com/bspaans/jit-compiler/ir/expr"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

func encode_IR_Add(i *expr.IR_Add, ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {
	return encode_Operator(i.Op1, i.Op2, x86_64.ADD, i.String(), ctx, target)
}
