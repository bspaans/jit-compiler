package aarch64

import (
	"github.com/bspaans/jit-compiler/asm/aarch64/opcodes"
	"github.com/bspaans/jit-compiler/lib"
)

func ADD(src, dest, val lib.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("add", opcodes.ADD, dest, src, val)
}

func ADDS(src, dest, val lib.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("adds", opcodes.ADDS, dest, src, val)
}

func MOVK(val, dest lib.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("movk", opcodes.MOVK, val, dest)
}

func SUB(src, dest, val lib.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("sub", opcodes.SUB, dest, src, val)
}

func SUBS(src, dest, val lib.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("subs", opcodes.SUBS, dest, src, val)
}
