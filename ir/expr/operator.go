package expr

import (
	"fmt"

	"github.com/bspaans/jit/asm"
	. "github.com/bspaans/jit/ir/shared"
)

type operator func(op1, op2 asm.Operand) asm.Instruction

type IR_Operator struct {
	*BaseIRExpression
	Opcode operator
	Repr   string
	Op1    IRExpression
	Op2    IRExpression
}

func NewIR_Operator(opcode operator, repr string, op1, op2 IRExpression) *IR_Operator {
	return &IR_Operator{
		BaseIRExpression: NewBaseIRExpression(Add),
		Opcode:           opcode,
		Repr:             repr,
		Op1:              op1,
		Op2:              op2,
	}
}

func (i *IR_Operator) ReturnType(ctx *IR_Context) Type {
	return i.Op1.ReturnType(ctx)
}

func (i *IR_Operator) Encode(ctx *IR_Context, target *asm.Register) ([]asm.Instruction, error) {
	returnType1, returnType2 := i.Op1.ReturnType(ctx), i.Op2.ReturnType(ctx)
	if returnType1 == returnType2 && (returnType1 == TFloat64 || returnType1 == TUint64) {
		result, err := i.Op1.Encode(ctx, target)
		if err != nil {
			return nil, err
		}

		var reg *asm.Register
		if i.Op2.Type() == Variable {
			variable := i.Op2.(*IR_Variable).Value
			reg = ctx.VariableMap[variable]
		} else {
			reg = ctx.AllocateRegister(returnType1)
			defer ctx.DeallocateRegister(reg)

			expr, err := i.Op2.Encode(ctx, reg)
			if err != nil {
				return nil, err
			}
			result = asm.Instructions(result).Add(expr)
		}
		instr := i.Opcode(reg, target)
		ctx.AddInstruction(instr)
		result = append(result, instr)
		return result, nil
	}
	return nil, fmt.Errorf("Unsupported types (%s, %s) in %s IR operation: %s", returnType1, returnType2, i.Repr, i.String())
}

func (i *IR_Operator) String() string {
	return fmt.Sprintf("%s %s %s", i.Op1.String(), i.Repr, i.Op2.String())
}
