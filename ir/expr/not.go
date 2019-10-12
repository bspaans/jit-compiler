package expr

import (
	"errors"
	"fmt"

	"github.com/bspaans/jit/asm"
	"github.com/bspaans/jit/asm/encoding"
	. "github.com/bspaans/jit/ir/shared"
	"github.com/bspaans/jit/lib"
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

func (i *IR_Not) Encode(ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {

	var reg1 encoding.Operand

	result := []lib.Instruction{}
	if i.Op1.Type() == Variable {
		variable := i.Op1.(*IR_Variable).Value
		reg1 = ctx.VariableMap[variable]
	} else if i.Op1.Type() == Uint64 {
		value := i.Op1.(*IR_Uint64).Value
		reg1 = ctx.AllocateRegister(TUint64)
		defer ctx.DeallocateRegister(reg1.(*encoding.Register))
		result = append(result, asm.MOV(encoding.Uint64(value), reg1))
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

	tmpReg := ctx.AllocateRegister(TUint64)
	defer ctx.DeallocateRegister(tmpReg)
	instr := []lib.Instruction{}
	instr = append(instr, asm.CMP(encoding.Uint32(0), reg1))
	instr = append(instr, asm.MOV(encoding.Uint64(0), tmpReg))
	instr = append(instr, asm.SETE(tmpReg.Lower8BitRegister()))
	instr = append(instr, asm.MOV(tmpReg, target))
	for _, inst := range instr {
		result = append(result, inst)
		ctx.AddInstruction(inst)
	}
	return result, nil
}

func (i *IR_Not) String() string {
	return fmt.Sprintf("!(%s)", i.Op1.String())
}
