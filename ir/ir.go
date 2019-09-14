package ir

import (
	"fmt"

	"github.com/bspaans/jit/asm"
)

type IR interface {
	Type() IRType
	Encode(*IR_Context) ([]asm.Instruction, error)
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

func CompileIR(stmts []IR) ([]asm.Instruction, error) {
	ctx := NewIRContext()
	result := []asm.Instruction{}
	for _, stmt := range stmts {
		code, err := stmt.Encode(ctx)
		if err != nil {
			return nil, err
		}
		fmt.Println(stmt)
		fmt.Println(code)
		for _, i := range code {
			result = append(result, i)
		}
	}
	return result, nil
}

func init() {
	i := []IR{
		NewIR_Assignment("a", NewIR_Uint64(123)),
		NewIR_Assignment("b", NewIR_Uint64(123)),
		NewIR_Assignment("c", NewIR_Add(NewIR_Variable("a"), NewIR_Variable("b"))),
		NewIR_Assignment("d", NewIR_Add(NewIR_Uint64(15), NewIR_Variable("b"))),
		NewIR_Assignment("e", NewIR_Add(NewIR_Variable("b"), NewIR_Uint64(53))),
		NewIR_Assignment("f", NewIR_Equals(NewIR_Variable("b"), NewIR_Uint64(53))),
	}
	instr, err := CompileIR(i)
	if err != nil {
		panic(err)
	}
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
