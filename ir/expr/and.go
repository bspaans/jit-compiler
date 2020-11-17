package expr

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm"
	"github.com/bspaans/jit-compiler/asm/encoding"
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
