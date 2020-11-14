package statements

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm"
	"github.com/bspaans/jit-compiler/asm/encoding"
	"github.com/bspaans/jit-compiler/ir/expr"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

type IR_Return struct {
	*BaseIR
	Expr IRExpression
}

func NewIR_Return(expr IRExpression) *IR_Return {
	return &IR_Return{
		BaseIR: NewBaseIR(Return),
		Expr:   expr,
	}
}

func (i *IR_Return) Encode(ctx *IR_Context) ([]lib.Instruction, error) {
	result := []lib.Instruction{}
	var reg encoding.Operand
	var ok bool
	if i.Expr.Type() == Variable {
		reg, ok = ctx.VariableMap[i.Expr.(*expr.IR_Variable).Value]
		if !ok {
			return nil, fmt.Errorf("Unknown variable '%s' in return expression: %s", i.Expr.(*expr.IR_Variable).Value, i.String())
		}
	} else {
		reg = ctx.AllocateRegister(i.Expr.ReturnType(ctx))
		defer ctx.DeallocateRegister(reg.(*encoding.Register))
		result_, err := i.Expr.Encode(ctx, reg)
		if err != nil {
			return nil, err
		}
		result = result_
	}
	if reg.Width() != lib.QUADWORD {
		cast := ctx.AllocateRegister(TUint64)
		defer ctx.DeallocateRegister(cast)

		if reg.Width() == lib.BYTE || reg.Width() == lib.WORD {
			movzx := asm.MOVZX(reg, cast)
			result = append(result, movzx)
			ctx.AddInstruction(movzx)
		} else {
			xor := asm.XOR(cast, cast)
			mov := asm.MOV(reg, cast.ForOperandWidth(reg.Width()))
			result = append(result, xor)
			result = append(result, mov)
			ctx.AddInstruction(mov)
			ctx.AddInstruction(xor)
		}
		reg = cast
	}
	target := ctx.PeekReturn()
	instr := []lib.Instruction{
		asm.MOV(reg.(*encoding.Register).Get64BitRegister(), target),
		asm.RETURN(),
	}
	for _, inst := range instr {
		ctx.AddInstruction(inst)
		result = append(result, inst)
	}
	return result, nil
}

func (i *IR_Return) String() string {
	return fmt.Sprintf("return %s", i.Expr.String())
}

func (i *IR_Return) AddToDataSection(ctx *IR_Context) error {
	return i.Expr.AddToDataSection(ctx)
}
