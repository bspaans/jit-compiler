package x86_64

import (
	"errors"
	"fmt"

	"github.com/bspaans/jit-compiler/asm/x86_64"
	"github.com/bspaans/jit-compiler/asm/x86_64/encoding"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/ir/statements"
	"github.com/bspaans/jit-compiler/lib"
)

func encode_IR_If(i *statements.IR_If, ctx *IR_Context) ([]lib.Instruction, error) {
	if i.Condition.ReturnType(ctx) != TBool {
		return nil, errors.New("Unsupported if IR condition")
	}
	reg := ctx.AllocateRegister(TBool)
	defer ctx.DeallocateRegister(reg)

	// Get the lengths of the true and false branches
	stmt1Len, err := IR_Length(i.Stmt1, ctx)
	if err != nil {
		return nil, err
	}
	stmt2Len, err := IR_Length(i.Stmt2, ctx)
	if err != nil {
		return nil, err
	}

	result, err := conditionalJump(ctx, i.Condition, stmt1Len)
	if err != nil {
		return nil, fmt.Errorf("%s in %s", err.Error(), i.String())
	}

	s1, err := encodeStatement(i.Stmt1, ctx)
	if err != nil {
		return nil, err
	}
	for _, instr := range s1 {
		result = append(result, instr)
	}
	jmp := x86_64.JMP(encoding.Uint8(stmt2Len))
	ctx.AddInstruction(jmp)
	result = append(result, jmp)

	s2, err := encodeStatement(i.Stmt2, ctx)
	if err != nil {
		return nil, err
	}
	for _, instr := range s2 {
		result = append(result, instr)
	}
	return result, nil
}
