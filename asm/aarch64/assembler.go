package aarch64

import (
	"github.com/bspaans/jit-compiler/asm/aarch64/opcodes"
	"github.com/bspaans/jit-compiler/lib"
)

func ADD(dest, src, val lib.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("add", opcodes.ADD, 3, dest, src, val)
}

func ADDS(dest, src, val lib.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("adds", opcodes.ADDS, 3, dest, src, val)
}

func MOVK(dest, val lib.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("movk", opcodes.MOVK, 2, dest, val)
}

func SUB(dest, src, val lib.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("sub", opcodes.SUB, 3, dest, src, val)
}

func SUBS(dest, src, val lib.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("subs", opcodes.SUBS, 3, dest, src, val)
}
