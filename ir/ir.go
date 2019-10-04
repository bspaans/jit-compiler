package ir

import (
	"fmt"

	"github.com/bspaans/jit/asm"
)

type IR interface {
	Type() IRType
	String() string
	AddToDataSection(ctx *IR_Context)
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
func (b *BaseIR) AddToDataSection(ctx *IR_Context) {}

func Compile(stmts []IR) (asm.MachineCode, error) {
	ctx := NewIRContext()
	address := uint(0)
	result := []uint8{}
	for _, stmt := range stmts {
		stmt.AddToDataSection(ctx)
	}
	if len(ctx.DataSection) > 0 {
		jmp := &asm.JMP{asm.Uint8(len(ctx.DataSection))}
		fmt.Printf("0x%x: %s\n", address, jmp.String())
		result_, err := jmp.Encode()
		if err != nil {
			return nil, err
		}
		result = result_
		fmt.Println(asm.MachineCode(result_))
		for _, d := range ctx.DataSection {
			result = append(result, d)
		}
		address += uint(len(ctx.DataSection))
	}
	for _, stmt := range stmts {
		ctx.InstructionPointer = address
		code, err := stmt.Encode(ctx)
		if err != nil {
			return nil, err
		}
		fmt.Println("\n:: " + stmt.String() + "\n")
		for _, i := range code {
			b, err := i.Encode()
			if err != nil {
				return nil, err
			}
			fmt.Printf("0x%x: %s\n", address, i.String())
			address += uint(len(b))
			fmt.Println(asm.MachineCode(b))
			for _, code := range b {
				result = append(result, code)
			}
		}
	}
	fmt.Println()
	return result, nil
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
		NewIR_Assignment("b", NewIR_ByteArray([]uint8("test"))),
		NewIR_If(NewIR_Equals(NewIR_Uint64(53), NewIR_Uint64(53)),
			NewIR_Assignment("f", NewIR_Uint64(42)),
			NewIR_Assignment("f", NewIR_Uint64(53)),
		),
		NewIR_Return(NewIR_Variable("f")),
	}
	b, err := Compile(i)
	if err != nil {
		panic(err)
	}
	fmt.Println(b)
	b.Execute()
}
