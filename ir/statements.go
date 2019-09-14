package ir

import (
	"errors"
	"fmt"

	"github.com/bspaans/jit/asm"
)

type IRType int

const (
	Assignment IRType = iota
)

type IR_Assignment struct {
	*BaseIR
	Variable string
	Expr     IRExpression
}

func NewIR_Assignment(variable string, expr IRExpression) *IR_Assignment {
	return &IR_Assignment{
		BaseIR:   NewBaseIR(Assignment),
		Variable: variable,
		Expr:     expr,
	}
}

func (i *IR_Assignment) Encode(ctx *IR_Context) ([]asm.Instruction, error) {
	reg := ctx.AllocateRegister()
	ctx.VariableMap[i.Variable] = reg
	// TODO expressions should always take a target registry?
	if i.Expr.Type() == Uint64 {
		arg := asm.Uint64(i.Expr.(*IR_Uint64).Value)
		return []asm.Instruction{&asm.MOV{arg, asm.Get64BitRegisterByIndex(reg)}}, nil
	} else if i.Expr.Type() == Add {
		expr := i.Expr.(*IR_Add)
		expr.SetTargetRegister(asm.Get64BitRegisterByIndex(reg))
		return expr.Encode(ctx)
	}
	return nil, errors.New("Unsupported assigment operation")
}

func (i *IR_Assignment) String() string {
	return fmt.Sprintf("%s = %s", i.Variable, i.Expr.String())
}
