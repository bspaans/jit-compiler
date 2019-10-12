package asm

import (
	"github.com/bspaans/jit/asm/encoding"
	"github.com/bspaans/jit/lib"
)

func JNE(dest encoding.Operand) lib.Instruction {
	return encoding.OpcodesToInstruction("jne", encoding.JNE, 1, dest)
}
func JMP(dest encoding.Operand) lib.Instruction {
	return encoding.OpcodesToInstruction("jmp", encoding.JMP, 1, dest)
}
