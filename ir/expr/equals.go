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
	result, err := Compare(i.Op1, i.Op2, ctx)
	if err != nil {
		return nil, fmt.Errorf("%s in %s", err.Error(), i.String())
	}
	if includeSETE {
		tmpReg := ctx.AllocateRegister(TUint64)
		defer ctx.DeallocateRegister(tmpReg)
		// TODO xor tmpreg
		sete := asm.SETE(tmpReg.Get8BitRegister())
		mov := asm.MOV(tmpReg.ForOperandWidth(target.Width()), target)
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
func (b *IR_Equals) SSA_Transform(ctx *SSA_Context) (SSA_Rewrites, IRExpression) {
	return nil, b
}
