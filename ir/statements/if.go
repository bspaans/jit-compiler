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

type IR_If struct {
	*BaseIR
	Condition IRExpression
	Stmt1     IR
	Stmt2     IR
}

func NewIR_If(condition IRExpression, stmt1, stmt2 IR) *IR_If {
	return &IR_If{
		BaseIR:    NewBaseIR(If),
		Condition: condition,
		Stmt1:     stmt1,
		Stmt2:     stmt2,
	}
}

func (i *IR_If) Encode(ctx *IR_Context) ([]lib.Instruction, error) {
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

	// Add the condition and jump
	var result []lib.Instruction
	switch i.Condition.(type) {
	case *expr.IR_Equals:
		result, err = i.Condition.(*expr.IR_Equals).EncodeWithoutSETE(ctx, reg)
		instr := []lib.Instruction{
			asm.JNE(encoding.Uint8(stmt1Len)),
		}
		for _, inst := range instr {
			ctx.AddInstruction(inst)
			result = append(result, inst)
		}
	case *expr.IR_Not:
		result, err = i.Condition.(*expr.IR_Not).EncodeWithoutSETE(ctx, reg)
		instr := []lib.Instruction{
			asm.JNE(encoding.Uint8(stmt1Len)),
		}
		for _, inst := range instr {
			ctx.AddInstruction(inst)
			result = append(result, inst)
		}
	case *expr.IR_Bool:
		// If it's just a boolean compare it to 1
		result, err = i.Condition.Encode(ctx, reg)
		instr := []lib.Instruction{
			asm.CMP(encoding.Uint32(1), reg),
			asm.JNE(encoding.Uint8(stmt1Len)),
		}
		for _, inst := range instr {
			ctx.AddInstruction(inst)
			result = append(result, inst)
		}
	default:
		return nil, fmt.Errorf("Unsupported if condition %s", i.Condition.String())
	}
	if err != nil {
		return nil, err
	}
	s1, err := i.Stmt1.Encode(ctx)
	if err != nil {
		return nil, err
	}
	for _, instr := range s1 {
		result = append(result, instr)
	}
	jmp := asm.JMP(encoding.Uint8(stmt2Len))
	ctx.AddInstruction(jmp)
	result = append(result, jmp)

	s2, err := i.Stmt2.Encode(ctx)
	if err != nil {
		return nil, err
	}
	for _, instr := range s2 {
		result = append(result, instr)
	}
	return result, nil
}

func (i *IR_If) String() string {
	return fmt.Sprintf("if %s; then %s; else %s;", i.Condition.String(), i.Stmt1.String(), i.Stmt2.String())
}

func (i *IR_If) AddToDataSection(ctx *IR_Context) error {
	if err := i.Condition.AddToDataSection(ctx); err != nil {
		return err
	}
	if err := i.Stmt1.AddToDataSection(ctx); err != nil {
		return err
	}
	if err := i.Stmt2.AddToDataSection(ctx); err != nil {
		return err
	}
	return nil
}
