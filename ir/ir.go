package ir

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm/x86_64"
	"github.com/bspaans/jit-compiler/asm/x86_64/encoding"
	"github.com/bspaans/jit-compiler/elf"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

func Compile(targetArchitecture Architecture, stmts []IR, debug bool) (lib.MachineCode, error) {
	ctx := NewIRContext(targetArchitecture)
	return CompileWithContext(stmts, debug, ctx)
}

func CompileToBinary(targetArchitecture Architecture, stmts []IR, debug bool, path string) error {
	ctx := NewIRContext(targetArchitecture)
	ctx.ReturnOperandStack = []encoding.Operand{encoding.Rax}
	code, err := CompileWithContext(stmts, debug, ctx)
	if err != nil {
		return err
	}
	return elf.CreateTinyBinary(code, path)
}

func CompileWithContext(stmts []IR, debug bool, ctx *IR_Context) (lib.MachineCode, error) {
	result := []uint8{}
	readonly, rw, err := ctx.Architecture.EncodeDataSection(stmts, ctx)
	if err != nil {
		return nil, err
	}
	if debug {
		fmt.Println(".rodata:")
		fmt.Println(readonly)
		fmt.Println(".data:")
		fmt.Println(rw)
	}
	// TODO: do this properly
	ctx.DataSection = append(readonly, rw...)
	ctx.InstructionPointer += uint(len(ctx.DataSection))
	if debug {
		fmt.Println("_start:")
	}
	if len(ctx.DataSection) > 0 {
		// TODO make Architecture dependent
		jmp := x86_64.JMP(encoding.Uint8(len(ctx.DataSection)))
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
		code, err := ctx.Architecture.EncodeStatement(stmt, ctx)
		if err != nil {
			return nil, fmt.Errorf("Error encoding %s: %s", stmt, err.Error())
		}
		if debug {
			fmt.Println("\n:: " + stmt.String() + "\n")
		}
		for _, i := range code {
			b, err := i.Encode()
			if err != nil {
				return nil, fmt.Errorf("Failed to encode %s: %s\n%s", stmt, err.Error(), lib.Instructions(code).String())
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
