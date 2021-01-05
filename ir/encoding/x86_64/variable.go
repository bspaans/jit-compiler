package x86_64

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm/x86_64"
	"github.com/bspaans/jit-compiler/asm/x86_64/encoding"
	"github.com/bspaans/jit-compiler/ir/expr"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

func encode_IR_Variable(i *expr.IR_Variable, ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {
	reg, ok := ctx.VariableMap[i.Value]
	if !ok || reg == nil {
		return nil, fmt.Errorf("Unknown variable '%s'", i.Value)
	}
	result := []lib.Instruction{x86_64.MOV(reg, target)}
	ctx.AddInstructions(result)
	return result, nil
}
