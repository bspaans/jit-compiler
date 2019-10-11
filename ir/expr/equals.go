package expr

import (
	"errors"
	"fmt"

	"github.com/bspaans/jit/asm"
	. "github.com/bspaans/jit/ir/shared"
)

type IR_Equals struct {
	*BaseIRExpression
	Op1 IRExpression
	Op2 IRExpression
}

func NewIR_Equals(op1, op2 IRExpression) *IR_Equals {
	return &IR_Equals{
		BaseIRExpression: NewBaseIRExpression(Equals),
		Op1:              op1,
		Op2:              op2,
	}
}

func (i *IR_Equals) ReturnType(ctx *IR_Context) Type {
	return TBool
}

func (i *IR_Equals) Encode(ctx *IR_Context, target asm.Operand) ([]asm.Instruction, error) {
	result := []asm.Instruction{}

	var reg1, reg2 asm.Operand

	if i.Op1.Type() == Variable {
		variable := i.Op1.(*IR_Variable).Value
		reg1 = ctx.VariableMap[variable]
	} else if i.Op1.Type() == Uint64 {
		value := i.Op1.(*IR_Uint64).Value
		reg1 = ctx.AllocateRegister(TUint64)
		defer ctx.DeallocateRegister(reg1.(*asm.Register))
		result = append(result, &asm.MOV{asm.Uint64(value), reg1})
	} else {
		return nil, errors.New("Unsupported cmp IR operation")
	}

	if i.Op2.Type() == Variable {
		variable := i.Op2.(*IR_Variable).Value
		reg2 = ctx.VariableMap[variable]
	} else if i.Op2.Type() == Uint64 {
		value := i.Op2.(*IR_Uint64).Value
		reg2 = ctx.AllocateRegister(TUint64)
		defer ctx.DeallocateRegister(reg2.(*asm.Register))
		result = append(result, &asm.MOV{asm.Uint64(value), reg2})
	} else {
		return nil, errors.New("Unsupported add IR operation")
	}
	tmpReg := ctx.AllocateRegister(TUint64)
	defer ctx.DeallocateRegister(tmpReg)
	result = append(result, &asm.CMP{reg1, reg2})
	result = append(result, &asm.MOV{asm.Uint64(0), tmpReg})
	result = append(result, &asm.SETE{tmpReg.Lower8BitRegister()})
	result = append(result, &asm.MOV{tmpReg, target})
	ctx.AddInstructions(result)
	return result, nil
}

func (i *IR_Equals) String() string {
	return fmt.Sprintf("%s == %s", i.Op1.String(), i.Op2.String())
}
