package opcodes

import (
	. "github.com/bspaans/jit-compiler/asm/aarch64/encoding"
	"github.com/bspaans/jit-compiler/lib"
)

func OpcodeToInstruction(name string, opcode *Opcode, argCount int, operands ...lib.Operand) lib.Instruction {
	return OpcodesToInstruction(name, []*Opcode{opcode}, argCount, operands...)
}

func OpcodesToInstruction(name string, opcodes []*Opcode, argCount int, operands ...lib.Operand) lib.Instruction {
	return nil
}
