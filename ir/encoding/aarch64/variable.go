package aarch64

import (
	"github.com/bspaans/jit-compiler/ir/expr"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

func encode_IR_Variable(v *expr.IR_Variable, ctx *IR_Context, target lib.Operand) ([]lib.Instruction, error) {
	return nil, nil
}
