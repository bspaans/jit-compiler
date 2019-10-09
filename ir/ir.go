package ir

import (
	"fmt"

	"github.com/bspaans/jit/asm"
	. "github.com/bspaans/jit/ir/expr"
	. "github.com/bspaans/jit/ir/shared"
	. "github.com/bspaans/jit/ir/statements"
)

func Compile(stmts []IR) (asm.MachineCode, error) {
	ctx := NewIRContext()
	result := []uint8{}
	fmt.Println(".data:")
	for _, stmt := range stmts {
		currentOffset := ctx.DataSectionOffset + len(ctx.DataSection)
		if err := stmt.AddToDataSection(ctx); err != nil {
			return nil, err
		}
		if len(ctx.DataSection) != currentOffset-2 {
			fmt.Printf("0x%x-0x%x: %s\n", currentOffset, len(ctx.DataSection), stmt.String())
			fmt.Println(asm.MachineCode(ctx.DataSection[currentOffset-ctx.DataSectionOffset : len(ctx.DataSection)-1]))
		}
	}
	fmt.Println(".start:")
	if len(ctx.DataSection) > 0 {
		jmp := &asm.JMP{asm.Uint8(len(ctx.DataSection))}
		fmt.Printf("0x%x: %s\n", 0, jmp.String())
		result_, err := jmp.Encode()
		if err != nil {
			return nil, err
		}
		result = result_
		fmt.Println(asm.MachineCode(result_))
		for _, d := range ctx.DataSection {
			result = append(result, d)
		}
	} else {
		ctx.DataSectionOffset = 0
		ctx.InstructionPointer = 0
	}
	address := uint(ctx.DataSectionOffset + len(ctx.DataSection))
	for _, stmt := range stmts {
		fmt.Printf("RIP: %d 0x%x\n", ctx.InstructionPointer, ctx.InstructionPointer)
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
	for _, stmt := range stmts {
		_, err := stmt.Encode(ctx)
		if err != nil {
			return nil, err
		}
	}
	return ctx.GetInstructions(), nil
}

func init() {
	i := []IR{
		NewIR_Assignment("a", NewIR_Function(&TFunction{TUint64, []Type{TUint64}, []string{"a"}},
			NewIR_AndThen(
				NewIR_Assignment("b", NewIR_Float64(3.0)),
				NewIR_AndThen(
					NewIR_Assignment("c", NewIR_Cast(NewIR_Variable("a"), TFloat64)),
					NewIR_AndThen(
						NewIR_Assignment("d", NewIR_Mul(NewIR_Variable("b"), NewIR_Variable("c"))),
						NewIR_Return(NewIR_Cast(NewIR_Variable("d"), TUint64))),
				)),
		)),
		NewIR_Assignment("g", NewIR_Uint64(2)),
		NewIR_Return(NewIR_Variable("g")),
		/*
			NewIR_Assignment("q", NewIR_Float64(2.1415)),
			NewIR_Assignment("q", NewIR_Add(NewIR_Variable("q"), NewIR_Float64(1.5))),
			NewIR_Assignment("i", NewIR_Uint64(0)),
			NewIR_While(NewIR_Not(NewIR_Equals(NewIR_Variable("i"), NewIR_Uint64(5))), NewIR_AndThen(
				NewIR_Assignment("g", NewIR_LinuxWrite(NewIR_Uint64(uint64(1)), []uint8("howdy\n"), 6)),
				NewIR_Assignment("i", NewIR_Add(NewIR_Variable("i"), NewIR_Uint64(1))),
			),
			),
			NewIR_Assignment("j", NewIR_LinuxOpen("/tmp/test.txt", os.O_CREATE|os.O_WRONLY, 0644)),
			NewIR_Assignment("g", NewIR_LinuxWrite(NewIR_Variable("j"), []uint8("howdy, how is it going\n"), 23)),
			NewIR_Return(NewIR_Variable("g")),
		*/
	}
	b, err := Compile(i)
	if err != nil {
		panic(err)
	}
	fmt.Println(b)
	b.Execute()
}
