package expr

import (
	"errors"
	"fmt"

	"github.com/bspaans/jit/asm"
	. "github.com/bspaans/jit/ir/shared"
)

type IR_Add struct {
	*BaseIRExpression
	Op1 IRExpression
	Op2 IRExpression
}

func NewIR_Add(op1, op2 IRExpression) *IR_Add {
	return &IR_Add{
		BaseIRExpression: NewBaseIRExpression(Add),
		Op1:              op1,
		Op2:              op2,
	}
}

func (i *IR_Add) ReturnType(ctx *IR_Context) Type {
	return TUint64
}

func (i *IR_Add) Encode(ctx *IR_Context, target *asm.Register) ([]asm.Instruction, error) {
	result := []asm.Instruction{}
	if i.Op1.Type() == Variable {
		variable := i.Op1.(*IR_Variable).Value
		reg := asm.Get64BitRegisterByIndex(ctx.VariableMap[variable])
		result = append(result, &asm.MOV{reg, target})
	} else if i.Op1.Type() == Uint64 {
		value := i.Op1.(*IR_Uint64).Value
		result = append(result, &asm.MOV{asm.Uint64(value), target})
	} else {
		return nil, errors.New("Unsupported add IR operation")
	}

	if i.Op2.Type() == Variable {
		variable := i.Op2.(*IR_Variable).Value
		reg := asm.Get64BitRegisterByIndex(ctx.VariableMap[variable])
		result = append(result, &asm.ADD{reg, target})
	} else if i.Op2.Type() == Uint64 {
		value := i.Op2.(*IR_Uint64).Value
		reg := asm.Get64BitRegisterByIndex(ctx.AllocateRegister())
		result = append(result, &asm.MOV{asm.Uint64(value), reg})
		result = append(result, &asm.ADD{reg, target})
		ctx.DeallocateRegister(reg.Register)
	} else {
		return nil, errors.New("Unsupported add IR operation")
	}
	ctx.AddInstructions(result)
	return result, nil
}

func (i *IR_Add) String() string {
	return fmt.Sprintf("%s + %s", i.Op1.String(), i.Op2.String())
}
