package statements

import (
	"errors"
	"fmt"

	"github.com/bspaans/jit-compiler/asm"
	"github.com/bspaans/jit-compiler/asm/encoding"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

type IR_While struct {
	*BaseIR
	Condition IRExpression
	Stmt      IR
}

func NewIR_While(condition IRExpression, stmt IR) *IR_While {
	return &IR_While{
		BaseIR:    NewBaseIR(While),
		Condition: condition,
		Stmt:      stmt,
	}
}

func (i *IR_While) Encode(ctx *IR_Context) ([]lib.Instruction, error) {
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

	result, err := ConditionalJump(ctx, i.Condition, stmtLen)
	if err != nil {
		return nil, fmt.Errorf("%s in %s", err.Error(), i.String())
	}
	s1, err := i.Stmt.Encode(ctx)
	if err != nil {
		return nil, err
	}
	result = lib.Instructions(result).Add(s1)
	jump := uint8((ctx.InstructionPointer + jmpSize) - beginning)
	// two's complement
	jump = (^jump) + 1
	jmp := asm.JMP(encoding.Uint8(jump))
	result = append(result, jmp)
	ctx.AddInstruction(jmp)
	return result, nil
}

func (i *IR_While) String() string {
	return fmt.Sprintf("while %s { %s }", i.Condition.String(), i.Stmt.String())
}

func (i *IR_While) AddToDataSection(ctx *IR_Context) error {
	if err := i.Condition.AddToDataSection(ctx); err != nil {
		return err
	}
	if err := i.Stmt.AddToDataSection(ctx); err != nil {
		return err
	}
	return nil
}

func (i *IR_While) SSA_Transform(ctx *SSA_Context) IR {
	return i
}
