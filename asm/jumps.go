package asm

import (
	"github.com/bspaans/jit/asm/encoding"
	"github.com/bspaans/jit/lib"
)

func ADD(src, dest encoding.Operand) lib.Instruction {
	return encoding.OpcodesToInstruction("add", encoding.ADD, 2, dest, src)
}
func CALL(dest encoding.Operand) lib.Instruction {
	return encoding.OpcodesToInstruction("call", encoding.CALL, 1, dest)
}
func CMP(src, dest encoding.Operand) lib.Instruction {
	return encoding.OpcodesToInstruction("cmp", encoding.CMP, 2, dest, src)
}

// Convert signed integer to scalar double-precision floating point (float64)
func CVTSI2SD(src, dest encoding.Operand) lib.Instruction {
	return encoding.OpcodesToInstruction("cvtsi2sd", encoding.CVTSI2SD, 2, dest, src)
}

// Convert double precision float to signed integer
func CVTTSD2SI(src, dest encoding.Operand) lib.Instruction {
	return encoding.OpcodesToInstruction("cvttsd2si", encoding.CVTTSD2SI, 2, dest, src)
}
func DEC(dest encoding.Operand) lib.Instruction {
	return encoding.OpcodesToInstruction("dec", encoding.DEC, 1, dest)
}
func DIV(src, dest encoding.Operand) lib.Instruction {
	return encoding.OpcodesToInstruction("div", encoding.DIV, 2, dest, src)
}
func INC(dest encoding.Operand) lib.Instruction {
	return encoding.OpcodesToInstruction("inc", encoding.INC, 1, dest)
}
func JNE(dest encoding.Operand) lib.Instruction {
	return encoding.OpcodesToInstruction("jne", encoding.JNE, 1, dest)
}
func JMP(dest encoding.Operand) lib.Instruction {
	return encoding.OpcodesToInstruction("jmp", encoding.JMP, 1, dest)
}
func LEA(src, dest encoding.Operand) lib.Instruction {
	return encoding.OpcodesToInstruction("lea", encoding.LEA, 2, dest, src)
}
func MOV(src, dest encoding.Operand) lib.Instruction {
	return encoding.OpcodesToInstruction("mov", encoding.MOV, 2, dest, src)
}
func MUL(src, dest encoding.Operand) lib.Instruction {
	return encoding.OpcodesToInstruction("mul", encoding.MUL, 2, dest, src)
}
func POP(dest encoding.Operand) lib.Instruction {
	return encoding.OpcodesToInstruction("pop", encoding.POP, 1, dest)
}
func PUSH(dest encoding.Operand) lib.Instruction {
	return encoding.OpcodesToInstruction("push", encoding.PUSH, 1, dest)
}
func PUSHFQ() lib.Instruction {
	return encoding.OpcodeToInstruction("pushfq", encoding.PUSHFQ, 0)
}
func RETURN() lib.Instruction {
	return encoding.OpcodeToInstruction("return", encoding.RETURN, 0)
}
func SETA(dest encoding.Operand) lib.Instruction {
	return encoding.OpcodesToInstruction("seta", encoding.SETA, 1, dest)
}
func SETAE(dest encoding.Operand) lib.Instruction {
	return encoding.OpcodesToInstruction("setae", encoding.SETAE, 1, dest)
}
func SETB(dest encoding.Operand) lib.Instruction {
	return encoding.OpcodesToInstruction("seta", encoding.SETB, 1, dest)
}
func SETBE(dest encoding.Operand) lib.Instruction {
	return encoding.OpcodesToInstruction("setae", encoding.SETBE, 1, dest)
}
func SETC(dest encoding.Operand) lib.Instruction {
	return encoding.OpcodesToInstruction("seta", encoding.SETC, 1, dest)
}
func SETE(dest encoding.Operand) lib.Instruction {
	return encoding.OpcodesToInstruction("sete", encoding.SETE, 1, dest)
}
func SETNE(dest encoding.Operand) lib.Instruction {
	return encoding.OpcodesToInstruction("setne", encoding.SETNE, 1, dest)
}
func SUB(src, dest encoding.Operand) lib.Instruction {
	return encoding.OpcodesToInstruction("sub", encoding.SUB, 2, dest, src)
}
func SYSCALL() lib.Instruction {
	return encoding.OpcodeToInstruction("syscall", encoding.SYSCALL, 0)
}
