package expr

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm"
	"github.com/bspaans/jit-compiler/asm/encoding"
	"github.com/bspaans/jit-compiler/ir/shared"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

type IR_LTE struct {
	*BaseIRExpression
	Op1 IRExpression
	Op2 IRExpression
}

func NewIR_LTE(op1, op2 IRExpression) *IR_LTE {
	return &IR_LTE{
		BaseIRExpression: NewBaseIRExpression(LTE),
		Op1:              op1,
		Op2:              op2,
	}
}

func (i *IR_LTE) ReturnType(ctx *IR_Context) Type {
	return TBool
}

func (i *IR_LTE) EncodeWithoutSETE(ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {
	return i.encode(ctx, target, false)
}

func (i *IR_LTE) Encode(ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {
	return i.encode(ctx, target, true)
}

func (i *IR_LTE) encode(ctx *IR_Context, target encoding.Operand, includeSETE bool) ([]lib.Instruction, error) {
	result, err := Compare(i.Op1, i.Op2, ctx)
	if err != nil {
		return nil, fmt.Errorf("%s in %s", err.Error(), i.String())
	}
	if includeSETE {
		tmpReg := ctx.AllocateRegister(TUint64)
		defer ctx.DeallocateRegister(tmpReg)
		// TODO xor tmpreg
		// TODO use right SET depending on sign
		sete := asm.SETBE(tmpReg.Get8BitRegister())
		if shared.IsSignedInteger(i.Op1.ReturnType(ctx)) {
			sete = asm.SETLE(tmpReg.Get8BitRegister())
		}
		mov := asm.MOV(tmpReg.ForOperandWidth(target.Width()), target)
		result = append(result, sete)
		result = append(result, mov)
		ctx.AddInstruction(sete)
		ctx.AddInstruction(mov)
	}
	return result, nil
}

func (i *IR_LTE) String() string {
	return fmt.Sprintf("%s <= %s", i.Op1.String(), i.Op2.String())
}

func (b *IR_LTE) AddToDataSection(ctx *IR_Context) error {
	if err := b.Op1.AddToDataSection(ctx); err != nil {
		return err
	}
	return b.Op2.AddToDataSection(ctx)
}
