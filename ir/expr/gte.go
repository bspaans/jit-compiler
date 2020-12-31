package expr

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm/x86_64"
	"github.com/bspaans/jit-compiler/asm/x86_64/encoding"
	"github.com/bspaans/jit-compiler/ir/shared"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

type IR_GTE struct {
	*BaseIRExpression
	Op1 IRExpression
	Op2 IRExpression
}

func NewIR_GTE(op1, op2 IRExpression) *IR_GTE {
	return &IR_GTE{
		BaseIRExpression: NewBaseIRExpression(GTE),
		Op1:              op1,
		Op2:              op2,
	}
}

func (i *IR_GTE) ReturnType(ctx *IR_Context) Type {
	return TBool
}

func (i *IR_GTE) EncodeWithoutSETE(ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {
	return i.encode(ctx, target, false)
}

func (i *IR_GTE) Encode(ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {
	return i.encode(ctx, target, true)
}

func (i *IR_GTE) encode(ctx *IR_Context, target encoding.Operand, includeSETE bool) ([]lib.Instruction, error) {
	result, err := Compare(i.Op1, i.Op2, ctx)
	if err != nil {
		return nil, fmt.Errorf("%s in %s", err.Error(), i.String())
	}
	if includeSETE {
		tmpReg := ctx.AllocateRegister(TUint64)
		defer ctx.DeallocateRegister(tmpReg)
		// TODO xor tmpreg
		// TODO use right SET depending on sign
		sete := asm.SETAE(tmpReg.Get8BitRegister())
		if shared.IsSignedInteger(i.Op1.ReturnType(ctx)) {
			sete = asm.SETGE(tmpReg.Get8BitRegister())
		}
		mov := asm.MOV(tmpReg.ForOperandWidth(target.Width()), target)
		result = append(result, sete)
		result = append(result, mov)
		ctx.AddInstruction(sete)
		ctx.AddInstruction(mov)
	}
	return result, nil
}

func (i *IR_GTE) String() string {
	return fmt.Sprintf("%s >= %s", i.Op1.String(), i.Op2.String())
}

func (b *IR_GTE) AddToDataSection(ctx *IR_Context) error {
	if err := b.Op1.AddToDataSection(ctx); err != nil {
		return err
	}
	return b.Op2.AddToDataSection(ctx)
}
func (b *IR_GTE) SSA_Transform(ctx *SSA_Context) (SSA_Rewrites, IRExpression) {
	if IsLiteralOrVariable(b.Op1) {
		if IsLiteralOrVariable(b.Op2) {
			return nil, b
		} else {
			rewrites, expr := b.Op2.SSA_Transform(ctx)
			v := ctx.GenerateVariable()
			rewrites = append(rewrites, NewSSA_Rewrite(v, expr))
			return rewrites, NewIR_GTE(b.Op1, NewIR_Variable(v))
		}
	} else {
		rewrites, expr := b.Op1.SSA_Transform(ctx)
		v := ctx.GenerateVariable()
		rewrites = append(rewrites, NewSSA_Rewrite(v, expr))
		if IsLiteralOrVariable(b.Op2) {
			return rewrites, NewIR_GTE(NewIR_Variable(v), b.Op2)
		} else {
			rewrites2, expr2 := b.Op2.SSA_Transform(ctx)
			for _, rw := range rewrites2 {
				rewrites = append(rewrites, rw)
			}
			v2 := ctx.GenerateVariable()
			rewrites = append(rewrites, NewSSA_Rewrite(v2, expr2))
			return rewrites, NewIR_GTE(NewIR_Variable(v), NewIR_Variable(v2))
		}
	}
}
