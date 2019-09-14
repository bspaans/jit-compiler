package ir

import (
	"errors"
	"fmt"

	"github.com/bspaans/jit/asm"
)

type IRType int

const (
	Assignment IRType = iota
	If         IRType = iota
	Return     IRType = iota
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

// Allocates a new register and assigns it the value of the expression.
func (i *IR_Assignment) Encode(ctx *IR_Context) ([]asm.Instruction, error) {
	r, found := ctx.VariableMap[i.Variable]
	if !found {
		r = ctx.AllocateRegister()
		ctx.VariableMap[i.Variable] = r
	}
	reg := asm.Get64BitRegisterByIndex(r)
	return i.Expr.Encode(ctx, reg)
}

func (i *IR_Assignment) String() string {
	return fmt.Sprintf("%s = %s", i.Variable, i.Expr.String())
}

type IR_If struct {
	*BaseIR
	Condition IRExpression
	Stmt1     IR
	Stmt2     IR
}

func NewIR_If(condition IRExpression, stmt1, stmt2 IR) *IR_If {
	return &IR_If{
		BaseIR:    NewBaseIR(If),
		Condition: condition,
		Stmt1:     stmt1,
		Stmt2:     stmt2,
	}
}

func (i *IR_If) Encode(ctx *IR_Context) ([]asm.Instruction, error) {
	// TODO: introduce return type function
	if i.Condition.Type() == Bool || i.Condition.Type() == Equals {
		r := ctx.AllocateRegister()
		defer ctx.DeallocateRegister(r)
		reg := asm.Get64BitRegisterByIndex(r)
		result, err := i.Condition.Encode(ctx, reg)
		if err != nil {
			return nil, err
		}
		s1, err := i.Stmt1.Encode(ctx)
		if err != nil {
			return nil, err
		}
		asmS1, err := asm.Instructions(s1).Encode()
		if err != nil {
			return nil, err
		}
		s2, err := i.Stmt2.Encode(ctx)
		if err != nil {
			return nil, err
		}
		asmS2, err := asm.Instructions(s2).Encode()
		if err != nil {
			return nil, err
		}
		result = append(result, &asm.CMP{asm.Uint32(1), reg})
		result = append(result, &asm.JNE{asm.Uint8(len(asmS1))})
		for _, instr := range s1 {
			result = append(result, instr)
		}
		result = append(result, &asm.JMP{asm.Uint8(len(asmS2))})
		for _, instr := range s2 {
			result = append(result, instr)
		}
		return result, nil
	}
	return nil, errors.New("Unsupported if IR expression")
}

func (i *IR_If) String() string {
	return fmt.Sprintf("if %s; then %s; else %s;", i.Condition.String(), i.Stmt1.String(), i.Stmt2.String())
}

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

func (i *IR_Return) Encode(ctx *IR_Context) ([]asm.Instruction, error) {
	result, err := i.Expr.Encode(ctx, asm.Rax)
	if err != nil {
		return nil, err
	}
	result = append(result, &asm.MOV{asm.Rax, &asm.DisplacedRegister{asm.Rsp, 8}})
	result = append(result, &asm.RET{})
	return result, nil
}

func (i *IR_Return) String() string {
	return fmt.Sprintf("return %s", i.Expr.String())
}
