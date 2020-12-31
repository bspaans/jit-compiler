package expr

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm/x86_64"
	"github.com/bspaans/jit-compiler/asm/x86_64/encoding"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

type IR_And struct {
	*BaseIRExpression
	Op1 IRExpression
	Op2 IRExpression
}

func NewIR_And(op1, op2 IRExpression) *IR_And {
	return &IR_And{
		BaseIRExpression: NewBaseIRExpression(And),
		Op1:              op1,
		Op2:              op2,
	}
}

func (i *IR_And) ReturnType(ctx *IR_Context) Type {
	return i.Op1.ReturnType(ctx)
}

func (i *IR_And) EncodeWithoutSETE(ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {
	return i.Encode(ctx, target)
}

func (i *IR_And) Encode(ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {
	ctx.AddInstruction("operator " + encoding.Comment(i.String()))
	returnType1, returnType2 := i.Op1.ReturnType(ctx), i.Op2.ReturnType(ctx)
	if returnType1 != returnType2 {
		return nil, fmt.Errorf("Unsupported types (%s, %s) in && IR operation: %s", returnType1, returnType2, i.String())
	}
	if returnType1 != TBool {
		return nil, fmt.Errorf("Unsupported types (%s, %s) in && IR operation: %s", returnType1, returnType2, i.String())
	}

	reg := ctx.AllocateRegister(returnType1)
	defer ctx.DeallocateRegister(reg)

	result, err := i.Op1.Encode(ctx, reg)
	if err != nil {
		return nil, err
	}
	expr2, err := i.Op2.Encode(ctx, target)
	if err != nil {
		return nil, err
	}
	result = lib.Instructions(result).Add(expr2)
	// TODO: should be using test?
	and := asm.AND(reg, target)
	result = append(result, and)
	ctx.AddInstruction(and)
	return result, nil
}

func (i *IR_And) String() string {
	return fmt.Sprintf("%s && %s", i.Op1.String(), i.Op2.String())
}

func (b *IR_And) AddToDataSection(ctx *IR_Context) error {
	if err := b.Op1.AddToDataSection(ctx); err != nil {
		return err
	}
	return b.Op2.AddToDataSection(ctx)
}

func (b *IR_And) SSA_Transform(ctx *SSA_Context) (SSA_Rewrites, IRExpression) {
	if IsLiteralOrVariable(b.Op1) {
		if IsLiteralOrVariable(b.Op2) {
			return nil, b
		} else {
			rewrites, expr := b.Op2.SSA_Transform(ctx)
			v := ctx.GenerateVariable()
			rewrites = append(rewrites, NewSSA_Rewrite(v, expr))
			return rewrites, NewIR_And(b.Op1, NewIR_Variable(v))
		}
	} else {
		rewrites, expr := b.Op1.SSA_Transform(ctx)
		v := ctx.GenerateVariable()
		rewrites = append(rewrites, NewSSA_Rewrite(v, expr))
		if IsLiteralOrVariable(b.Op2) {
			return rewrites, NewIR_And(NewIR_Variable(v), b.Op2)
		} else {
			rewrites2, expr2 := b.Op2.SSA_Transform(ctx)
			for _, rw := range rewrites2 {
				rewrites = append(rewrites, rw)
			}
			v2 := ctx.GenerateVariable()
			rewrites = append(rewrites, NewSSA_Rewrite(v2, expr2))
			return rewrites, NewIR_And(NewIR_Variable(v), NewIR_Variable(v2))
		}

	}
	return nil, b
}
