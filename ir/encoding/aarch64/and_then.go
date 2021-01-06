package aarch64

import (
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/ir/statements"
	"github.com/bspaans/jit-compiler/lib"
)

func encode_IR_AndThen(i *statements.IR_AndThen, ctx *IR_Context) ([]lib.Instruction, error) {
	result, err := encodeStatement(i.Stmt1, ctx)
	if err != nil {
		return nil, err
	}
	s2, err := encodeStatement(i.Stmt2, ctx)
	if err != nil {
		return nil, err
	}
	for _, instr := range s2 {
		result = append(result, instr)
	}
	return result, nil
}
