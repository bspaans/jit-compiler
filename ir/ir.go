package ir

import (
	"fmt"

	"github.com/bspaans/jit/asm"
)

type IRExpressionType int

const (
	Uint64 IRExpressionType = iota
)

type IRType int

const (
	Assignment IRType = iota
)

type IR_Context struct {
	RegisterMap        map[string]uint8
	RegistersAllocated uint8
}

func NewIRContext() *IR_Context {
	return &IR_Context{
		RegisterMap:        map[string]uint8{},
		RegistersAllocated: 0,
	}
}

func (i *IR_Context) AllocateRegister(variable string) uint8 {
	if i.RegistersAllocated >= 8 {
		panic("Register allocation limit. Needs stack handling")
	}
	if v, found := i.RegisterMap[variable]; found {
		return v
	}
	i.RegisterMap[variable] = i.RegistersAllocated
	i.RegistersAllocated += 1
	return i.RegisterMap[variable]
}

type IR interface {
	Type() IRType
	Encode(*IR_Context) []asm.Instruction
}

type IRExpression interface {
	Type() IRExpressionType
}

type BaseIR struct {
	typ IRType
}

func NewBaseIR(typ IRType) *BaseIR {
	return &BaseIR{
		typ: typ,
	}
}
func (b *BaseIR) Type() IRType {
	return b.typ
}

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

func (i *IR_Assignment) Encode(ctx *IR_Context) []asm.Instruction {
	reg := ctx.AllocateRegister(i.Variable)
	if i.Expr.Type() == Uint64 {
		arg := asm.Uint64(i.Expr.(*IR_Uint64).Value)
		return []asm.Instruction{&asm.MOV{arg, asm.Get64BitRegisterByIndex(reg)}}
	}
	return []asm.Instruction{}
}

func init() {
	ctx := NewIRContext()
	instr := NewIR_Assignment("a", NewIR_Uint64(123)).Encode(ctx)
	for _, returnStmt := range []asm.Instruction{
		&asm.MOV{asm.Get64BitRegisterByIndex(0), &asm.DisplacedRegister{asm.Get64BitRegisterByIndex(4), 8}},
		&asm.RET{},
	} {
		instr = append(instr, returnStmt)
	}
	b, err := asm.CompileInstruction(instr)
	if err != nil {
		panic(err)
	}
	fmt.Println(b)
	b.Execute()
}
