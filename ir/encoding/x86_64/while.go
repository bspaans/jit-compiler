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

func encode_IR_While(i *statements.IR_While, ctx *IR_Context) ([]lib.Instruction, error) {
	jmpSize := uint(2)
	if i.Condition.ReturnType(ctx) != TBool {
		return nil, errors.New("Unsupported if IR expression")
	}

	// Get the length of the loop statement
	stmtLen, err := IR_Length(i.Stmt, ctx)
	if err != nil {
		return nil, err
	}

	beginning := ctx.InstructionPointer

	result, err := conditionalJump(ctx, i.Condition, stmtLen)
	if err != nil {
		return nil, fmt.Errorf("%s in %s", err.Error(), i.String())
	}
	s1, err := encodeStatement(i.Stmt, ctx)
	if err != nil {
		return nil, err
	}
	result = lib.Instructions(result).Add(s1)
	jump := uint8((ctx.InstructionPointer + jmpSize) - beginning)
	// two's complement
	jump = (^jump) + 1
	jmp := x86_64.JMP(encoding.Uint8(jump))
	result = append(result, jmp)
	ctx.AddInstruction(jmp)
	return result, nil
}
