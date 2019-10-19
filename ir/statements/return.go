package statements

import (
	"fmt"

	"github.com/bspaans/jit/asm"
	"github.com/bspaans/jit/asm/encoding"
	"github.com/bspaans/jit/ir/expr"
	. "github.com/bspaans/jit/ir/shared"
	"github.com/bspaans/jit/lib"
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
	if i.Expr.Type() == Variable {
		reg = ctx.VariableMap[i.Expr.(*expr.IR_Variable).Value]
	} else {
		reg = ctx.AllocateRegister(i.Expr.ReturnType(ctx))
		defer ctx.DeallocateRegister(reg.(*encoding.Register))
		result_, err := i.Expr.Encode(ctx, reg)
		if err != nil {
			return nil, err
		}
		result = result_
	}
	target := ctx.PeekReturn()
	instr := []lib.Instruction{
		asm.MOV(reg, target),
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
