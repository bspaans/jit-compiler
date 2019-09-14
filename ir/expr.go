package ir

import (
	"errors"
	"fmt"

	"github.com/bspaans/jit/asm"
)

type IRExpressionType int
type IRExpression interface {
	Type() IRExpressionType
	Encode(ctx *IR_Context, target *asm.Register) ([]asm.Instruction, error)
	String() string
}

const (
	Uint64   IRExpressionType = iota
	Bool     IRExpressionType = iota
	Add      IRExpressionType = iota
	Variable IRExpressionType = iota
	Equals   IRExpressionType = iota
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

func (i *IR_Uint64) Encode(ctx *IR_Context, target *asm.Register) ([]asm.Instruction, error) {
	return []asm.Instruction{&asm.MOV{asm.Uint64(i.Value), target}}, nil
}

type IR_Bool struct {
	*BaseIRExpression
	Value bool
}

func NewIR_Bool(v bool) *IR_Bool {
	return &IR_Bool{
		BaseIRExpression: NewBaseIRExpression(Bool),
		Value:            v,
	}
}

func (i *IR_Bool) String() string {
	return fmt.Sprintf("%v", i.Value)
}

func (i *IR_Bool) Encode(ctx *IR_Context, target *asm.Register) ([]asm.Instruction, error) {
	value := 0
	if i.Value {
		value = 1
	}
	return []asm.Instruction{&asm.MOV{asm.Uint64(value), target}}, nil
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

func (i *IR_Variable) Encode(ctx *IR_Context, target *asm.Register) ([]asm.Instruction, error) {
	reg := asm.Get64BitRegisterByIndex(ctx.VariableMap[i.Value])
	return []asm.Instruction{&asm.MOV{reg, target}}, nil
}

type IR_Add struct {
	*BaseIRExpression
	Op1 IRExpression
	Op2 IRExpression
}

func NewIR_Add(op1, op2 IRExpression) *IR_Add {
	return &IR_Add{
		BaseIRExpression: NewBaseIRExpression(Add),
		Op1:              op1,
		Op2:              op2,
	}
}

func (i *IR_Add) Encode(ctx *IR_Context, target *asm.Register) ([]asm.Instruction, error) {
	result := []asm.Instruction{}
	if i.Op1.Type() == Variable {
		variable := i.Op1.(*IR_Variable).Value
		reg := asm.Get64BitRegisterByIndex(ctx.VariableMap[variable])
		result = append(result, &asm.MOV{reg, target})
	} else if i.Op1.Type() == Uint64 {
		value := i.Op1.(*IR_Uint64).Value
		result = append(result, &asm.MOV{asm.Uint64(value), target})
	} else {
		return nil, errors.New("Unsupported add IR operation")
	}

	if i.Op2.Type() == Variable {
		variable := i.Op2.(*IR_Variable).Value
		reg := asm.Get64BitRegisterByIndex(ctx.VariableMap[variable])
		result = append(result, &asm.ADD{reg, target})
	} else if i.Op2.Type() == Uint64 {
		value := i.Op2.(*IR_Uint64).Value
		reg := asm.Get64BitRegisterByIndex(ctx.AllocateRegister())
		result = append(result, &asm.MOV{asm.Uint64(value), reg})
		result = append(result, &asm.ADD{reg, target})
		ctx.DeallocateRegister(reg.Register)
	} else {
		return nil, errors.New("Unsupported add IR operation")
	}
	return result, nil
}

func (i *IR_Add) String() string {
	return fmt.Sprintf("%s + %s", i.Op1.String(), i.Op2.String())
}

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

func (i *IR_Equals) Encode(ctx *IR_Context, target *asm.Register) ([]asm.Instruction, error) {
	result := []asm.Instruction{}

	var reg1, reg2 *asm.Register

	if i.Op1.Type() == Variable {
		variable := i.Op1.(*IR_Variable).Value
		reg1 = asm.Get64BitRegisterByIndex(ctx.VariableMap[variable])
	} else if i.Op1.Type() == Uint64 {
		value := i.Op1.(*IR_Uint64).Value
		r := ctx.AllocateRegister()
		defer ctx.DeallocateRegister(r)
		reg1 = asm.Get64BitRegisterByIndex(r)
		result = append(result, &asm.MOV{asm.Uint64(value), reg1})
	} else {
		return nil, errors.New("Unsupported cmp IR operation")
	}

	if i.Op2.Type() == Variable {
		variable := i.Op2.(*IR_Variable).Value
		reg2 = asm.Get64BitRegisterByIndex(ctx.VariableMap[variable])
	} else if i.Op2.Type() == Uint64 {
		value := i.Op2.(*IR_Uint64).Value
		r := ctx.AllocateRegister()
		defer ctx.DeallocateRegister(r)
		reg2 = asm.Get64BitRegisterByIndex(r)
		result = append(result, &asm.MOV{asm.Uint64(value), reg2})
	} else {
		return nil, errors.New("Unsupported add IR operation")
	}
	result = append(result, &asm.CMP{reg1, reg2})
	result = append(result, &asm.MOV{asm.Uint64(0), target})
	result = append(result, &asm.SETE{target.Lower8BitRegister()})
	return result, nil
}

func (i *IR_Equals) String() string {
	return fmt.Sprintf("%s == %s", i.Op1.String(), i.Op2.String())
}
