package ir

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm"
	"github.com/bspaans/jit-compiler/asm/encoding"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

func Compile(stmts []IR, debug bool) (lib.MachineCode, error) {
	ctx := NewIRContext()
	result := []uint8{}
	if debug {
		fmt.Println(".data:")
	}
	for _, stmt := range stmts {
		currentOffset := ctx.DataSectionOffset + len(ctx.DataSection)
		if err := stmt.AddToDataSection(ctx); err != nil {
			return nil, err
		}
		if len(ctx.DataSection) != currentOffset-2 {
			if debug {
				fmt.Printf("0x%x-0x%x (0x%x): %s\n",
					currentOffset,
					len(ctx.DataSection)+ctx.DataSectionOffset,
					len(ctx.DataSection)+ctx.DataSectionOffset-currentOffset, stmt.String())
				fmt.Println(lib.MachineCode(ctx.DataSection[currentOffset-ctx.DataSectionOffset : len(ctx.DataSection)]))
			}
		}
	}
	if debug {
		fmt.Println(".start:")
	}
	if len(ctx.DataSection) > 0 {
		jmp := asm.JMP(encoding.Uint8(len(ctx.DataSection)))
		if debug {
			fmt.Printf("0x%x: %s\n", 0, jmp.String())
		}
		result_, err := jmp.Encode()
		if err != nil {
			return nil, err
		}
		result = result_
		if debug {
			fmt.Println(lib.MachineCode(result_))
		}
		for _, d := range ctx.DataSection {
			result = append(result, d)
		}
	} else {
		ctx.DataSectionOffset = 0
		ctx.InstructionPointer = 0
	}
	address := uint(ctx.DataSectionOffset + len(ctx.DataSection))
	for _, stmt := range stmts {
		code, err := stmt.Encode(ctx)
		if err != nil {
			return nil, err
		}
		if debug {
			fmt.Println("\n:: " + stmt.String() + "\n")
		}
		for _, i := range code {
			b, err := i.Encode()
			if err != nil {
				return nil, err
			}
			if debug {
				fmt.Printf("0x%x-0x%x 0x%x: %s\n", address, address+uint(len(b)), ctx.InstructionPointer, i.String())
			}
			address += uint(len(b))
			if debug {
				fmt.Println(lib.MachineCode(b))
			}
			for _, code := range b {
				result = append(result, code)
			}
		}
	}
	if debug {
		fmt.Println()
	}
	return result, nil
}

func CompileIR(stmts []IR) (lib.Instructions, error) {
	ctx := NewIRContext()
	for _, stmt := range stmts {
		_, err := stmt.Encode(ctx)
		if err != nil {
			return nil, err
		}
	}
	return ctx.GetInstructions(), nil
}
