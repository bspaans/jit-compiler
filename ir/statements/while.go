package statements

import (
	"errors"
	"fmt"

	"github.com/bspaans/jit/asm"
	. "github.com/bspaans/jit/ir/shared"
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

func (i *IR_While) Encode(ctx *IR_Context) ([]asm.Instruction, error) {
	if i.Condition.ReturnType(ctx) == TBool {
		reg := ctx.AllocateRegister(TBool)
		defer ctx.DeallocateRegister(reg)
		beginning := ctx.InstructionPointer
		result, err := i.Condition.Encode(ctx, reg)
		if err != nil {
			return nil, err
		}
		stmtLen, err := IR_Length(i.Stmt, ctx)
		if err != nil {
			return nil, err
		}
		instr := []asm.Instruction{
			&asm.CMP{asm.Uint32(1), reg},
			&asm.JNE{asm.Uint8(stmtLen + 2)},
		}
		for _, inst := range instr {
			ctx.AddInstruction(inst)
			result = append(result, inst)
		}
		s1, err := i.Stmt.Encode(ctx)
		if err != nil {
			return nil, err
		}
		for _, instr := range s1 {
			result = append(result, instr)
		}
		jmp := &asm.JMP{asm.Uint8(uint8(0xff - (int(ctx.InstructionPointer+1) - int(beginning))))}
		result = append(result, jmp)
		ctx.AddInstruction(jmp)
		fmt.Println("InstructionPointer", ctx.InstructionPointer, beginning, ctx.InstructionPointer-beginning)
		return result, nil
	}
	return nil, errors.New("Unsupported if IR expression")
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
