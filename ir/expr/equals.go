package expr

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm"
	"github.com/bspaans/jit-compiler/asm/encoding"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
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

func (i *IR_Equals) EncodeWithoutSETE(ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {
	return i.encode(ctx, target, false)
}

func (i *IR_Equals) Encode(ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {
	return i.encode(ctx, target, true)
}

func (i *IR_Equals) encode(ctx *IR_Context, target encoding.Operand, includeSETE bool) ([]lib.Instruction, error) {
	result := []lib.Instruction{}

	returnType1, returnType2 := i.Op1.ReturnType(ctx), i.Op2.ReturnType(ctx)
	if returnType1 != returnType2 {
		return nil, fmt.Errorf("Unsupported types (%s, %s) in == IR operation: %s", returnType1, returnType2, i.String())
	}

	var reg1, reg2 encoding.Operand

	if i.Op1.Type() == Variable {
		variable := i.Op1.(*IR_Variable).Value
		reg1 = ctx.VariableMap[variable]
	} else {
		reg1 = ctx.AllocateRegister(returnType1)
		defer ctx.DeallocateRegister(reg1.(*encoding.Register))
		expr1, err := i.Op1.Encode(ctx, reg1)
		if err != nil {
			return nil, err
		}
		result = lib.Instructions(result).Add(expr1)
	}

	if i.Op2.Type() == Variable {
		variable := i.Op2.(*IR_Variable).Value
		reg2 = ctx.VariableMap[variable]
	} else {
		reg2 = ctx.AllocateRegister(returnType1)
		defer ctx.DeallocateRegister(reg2.(*encoding.Register))
		expr2, err := i.Op2.Encode(ctx, reg2)
		if err != nil {
			return nil, err
		}
		result = lib.Instructions(result).Add(expr2)
	}
	cmp := asm.CMP(reg1, reg2)
	result = append(result, cmp)
	ctx.AddInstruction(cmp)
	if includeSETE {
		tmpReg := ctx.AllocateRegister(TUint64)
		defer ctx.DeallocateRegister(tmpReg)
		// TODO xor tmpreg
		sete := asm.SETE(tmpReg.Get8BitRegister())
		mov := asm.MOV(tmpReg, target)
		result = append(result, sete)
		result = append(result, mov)
		ctx.AddInstruction(sete)
		ctx.AddInstruction(mov)
	}
	return result, nil
}

func (i *IR_Equals) String() string {
	return fmt.Sprintf("%s == %s", i.Op1.String(), i.Op2.String())
}

func (b *IR_Equals) AddToDataSection(ctx *IR_Context) error {
	if err := b.Op1.AddToDataSection(ctx); err != nil {
		return err
	}
	return b.Op2.AddToDataSection(ctx)
}
