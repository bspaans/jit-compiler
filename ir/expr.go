package ir

import (
	"errors"
	"fmt"

	"github.com/bspaans/jit/asm"
)

type IRExpressionType int
type IRExpression interface {
	Type() IRExpressionType
	String() string
}

const (
	Uint64   IRExpressionType = iota
	Add      IRExpressionType = iota
	Variable IRExpressionType = iota
)

type BaseIRExpression struct {
	typ IRExpressionType
}

func NewBaseIRExpression(typ IRExpressionType) *BaseIRExpression {
	return &BaseIRExpression{
		typ: typ,
	}
}

func (b *BaseIRExpression) Type() IRExpressionType {
	return b.typ
}

type IR_Uint64 struct {
	*BaseIRExpression
	Value uint64
}

func NewIR_Uint64(v uint64) *IR_Uint64 {
	return &IR_Uint64{
		BaseIRExpression: NewBaseIRExpression(Uint64),
		Value:            v,
	}
}

func (i *IR_Uint64) String() string {
	return fmt.Sprintf("%d", i.Value)
}

type IR_Variable struct {
	*BaseIRExpression
	Value string
}

func NewIR_Variable(v string) *IR_Variable {
	return &IR_Variable{
		BaseIRExpression: NewBaseIRExpression(Variable),
		Value:            v,
	}
}

func (i *IR_Variable) String() string {
	return i.Value
}

type IR_Add struct {
	*BaseIRExpression
	Op1    IRExpression
	Op2    IRExpression
	target *asm.Register
}

func NewIR_Add(op1, op2 IRExpression) *IR_Add {
	return &IR_Add{
		BaseIRExpression: NewBaseIRExpression(Add),
		Op1:              op1,
		Op2:              op2,
	}
}

func (i *IR_Add) Encode(ctx *IR_Context) ([]asm.Instruction, error) {
	if i.target == nil {
		return nil, errors.New("Target register not set")
	}
	result := []asm.Instruction{}
	if i.Op1.Type() == Variable {
		variable := i.Op1.(*IR_Variable).Value
		reg := asm.Get64BitRegisterByIndex(ctx.VariableMap[variable])
		result = append(result, &asm.MOV{reg, i.target})
	} else if i.Op1.Type() == Uint64 {
		value := i.Op1.(*IR_Uint64).Value
		result = append(result, &asm.MOV{asm.Uint64(value), i.target})
	} else {
		return nil, errors.New("Unsupported add IR operation")
	}

	if i.Op2.Type() == Variable {
		variable := i.Op2.(*IR_Variable).Value
		reg := asm.Get64BitRegisterByIndex(ctx.VariableMap[variable])
		result = append(result, &asm.ADD{reg, i.target})
	} else if i.Op2.Type() == Uint64 {
		value := i.Op2.(*IR_Uint64).Value
		reg := asm.Get64BitRegisterByIndex(ctx.AllocateRegister())
		result = append(result, &asm.MOV{asm.Uint64(value), reg})
		result = append(result, &asm.ADD{reg, i.target})
		ctx.DeallocateRegister(reg.Register)
	} else {
		return nil, errors.New("Unsupported add IR operation")
	}
	return result, nil
}

func (i *IR_Add) SetTargetRegister(r *asm.Register) {
	i.target = r
}

func (i *IR_Add) String() string {
	return fmt.Sprintf("%s + %s", i.Op1.String(), i.Op2.String())
}
