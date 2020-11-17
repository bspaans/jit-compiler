package statements

import (
	"errors"
	"fmt"

	"github.com/bspaans/jit-compiler/asm"
	"github.com/bspaans/jit-compiler/asm/encoding"
	"github.com/bspaans/jit-compiler/ir/expr"
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

	reg := ctx.AllocateRegister(TBool)
	defer ctx.DeallocateRegister(reg)

	// Get the length of the loop statement
	stmtLen, err := IR_Length(i.Stmt, ctx)
	if err != nil {
		return nil, err
	}

	beginning := ctx.InstructionPointer

	var result []lib.Instruction
	switch i.Condition.(type) {
	case *expr.IR_Equals:
		result, err = i.Condition.(*expr.IR_Equals).EncodeWithoutSETE(ctx, reg)
		if err != nil {
			return nil, err
		}
		instr := []lib.Instruction{
			asm.JNE(encoding.Uint8(stmtLen + int(jmpSize))),
		}
		for _, inst := range instr {
			ctx.AddInstruction(inst)
			result = append(result, inst)
		}
	case *expr.IR_Not:
		result, err = i.Condition.(*expr.IR_Not).EncodeWithoutSETE(ctx, reg)
		if err != nil {
			return nil, err
		}
		instr := []lib.Instruction{
			asm.JE(encoding.Uint8(stmtLen + int(jmpSize))),
		}
		for _, inst := range instr {
			ctx.AddInstruction(inst)
			result = append(result, inst)
		}
	case *expr.IR_Bool:
		result, err = i.Condition.Encode(ctx, reg)
		if err != nil {
			return nil, err
		}
		instr := []lib.Instruction{
			asm.CMP(encoding.Uint32(1), reg),
			asm.JNE(encoding.Uint8(stmtLen + int(jmpSize))),
		}
		for _, inst := range instr {
			ctx.AddInstruction(inst)
			result = append(result, inst)
		}
	default:
		return nil, fmt.Errorf("Unsupported while condition %s", i.Condition.String())
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
