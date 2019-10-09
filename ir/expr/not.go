package expr

import (
	"errors"
	"fmt"

	"github.com/bspaans/jit/asm"
	. "github.com/bspaans/jit/ir/shared"
)

type IR_Not struct {
	*BaseIRExpression
	Op1 IRExpression
}

func NewIR_Not(op1 IRExpression) *IR_Not {
	return &IR_Not{
		BaseIRExpression: NewBaseIRExpression(Not),
		Op1:              op1,
	}
}

func (i *IR_Not) ReturnType(ctx *IR_Context) Type {
	return TBool
}

func (i *IR_Not) Encode(ctx *IR_Context, target *asm.Register) ([]asm.Instruction, error) {

	var reg1 *asm.Register

	result := []asm.Instruction{}
	if i.Op1.Type() == Variable {
		variable := i.Op1.(*IR_Variable).Value
		reg1 = ctx.VariableMap[variable]
	} else if i.Op1.Type() == Uint64 {
		value := i.Op1.(*IR_Uint64).Value
		reg1 = ctx.AllocateRegister(TUint64)
		defer ctx.DeallocateRegister(reg1)
		result = append(result, &asm.MOV{asm.Uint64(value), reg1})
	} else if i.Op1.Type() == Equals {
		result_, err := i.Op1.Encode(ctx, target)
		if err != nil {
			return nil, err
		}
		for _, r := range result_ {
			result = append(result, r)
		}
		reg1 = target
	} else {
		return nil, errors.New("Unsupported not IR operation: " + i.String())
	}

	instr := []asm.Instruction{}
	instr = append(instr, &asm.CMP{asm.Uint32(0), reg1})
	instr = append(instr, &asm.MOV{asm.Uint64(0), target})
	instr = append(instr, &asm.SETE{target.Lower8BitRegister()})
	for _, inst := range instr {
		result = append(result, inst)
		ctx.AddInstruction(inst)
	}
	return result, nil
}

func (i *IR_Not) String() string {
	return fmt.Sprintf("!(%s)", i.Op1.String())
}
