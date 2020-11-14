package asm

/*
	Instructions

	The below instructions implement a simple Instruction interface that allows
	them to be combined, optimised and encoded into machine code.

	x86_64 is a bit interesting in that instructions with different sets of
	operands might require subtly different machine code opcodes, even though
	they do the same thing functionally speaking, so the below instructions
	make use of a thing called OpcodeMaps that will match the given arguments
	to an appropriate opcode. For more on OpcodeMaps see asm/opcodes/
*/

import (
	"github.com/bspaans/jit-compiler/asm/encoding"
	"github.com/bspaans/jit-compiler/asm/opcodes"
	"github.com/bspaans/jit-compiler/lib"
)

func ADD(src, dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("add", opcodes.ADD, 2, dest, src)
}
func CALL(dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("call", opcodes.CALL, 1, dest)
}
func CMP(src, dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("cmp", opcodes.CMP, 2, dest, src)
}
func CMP_immediate(v uint64, dest encoding.Operand) lib.Instruction {
	if reg, ok := dest.(*encoding.Register); ok && reg.Width() == lib.BYTE {
		return CMP(encoding.Uint8(v), dest)
	}
	return opcodes.OpcodesToInstruction("cmp", opcodes.CMP, 2, dest, encoding.Uint32(v))
}

// Convert signed integer to scalar double-precision floating point (float64)
func CVTSI2SD(src, dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("cvtsi2sd", opcodes.CVTSI2SD, 2, dest, src)
}

// Convert double precision float to signed integer
func CVTTSD2SI(src, dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("cvttsd2si", opcodes.CVTTSD2SI, 2, dest, src)
}
func DEC(dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("dec", opcodes.DEC, 1, dest)
}
func DIV(src encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("div", opcodes.DIV, 1, src)
}
func IDIV(src, dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("div", opcodes.IDIV, 2, dest, src)
}
func INC(dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("inc", opcodes.INC, 1, dest)
}
func JE(dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("je", opcodes.JE, 1, dest)
}
func JNE(dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("jne", opcodes.JNE, 1, dest)
}
func JMP(dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("jmp", opcodes.JMP, 1, dest)
}
func LEA(src, dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("lea", opcodes.LEA, 2, dest, src)
}
func MOV(src, dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("mov", opcodes.MOV, 2, dest, src)
}
func MOV_immediate(v uint64, dest encoding.Operand) lib.Instruction {
	if reg, ok := dest.(*encoding.Register); ok && reg.Width() == lib.BYTE {
		return MOV(encoding.Uint8(v), dest)
	}
	if reg, ok := dest.(*encoding.Register); ok && reg.Width() == lib.WORD {
		return MOV(encoding.Uint16(v), dest)
	}
	if v < (1 << 32) {
		return MOV(encoding.Uint32(v), dest)
	}
	return MOV(encoding.Uint64(v), dest)
}
func MOVZX(src, dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("movzx", opcodes.MOV, 2, dest, src)
}
func IMUL(src, dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("imul", opcodes.IMUL, 2, dest, src)
}
func MUL(src encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("mul", opcodes.MUL, 1, src)
}
func POP(dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("pop", opcodes.POP, 1, dest)
}
func PUSH(dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("push", opcodes.PUSH, 1, dest)
}
func PUSHFQ() lib.Instruction {
	return opcodes.OpcodeToInstruction("pushfq", opcodes.PUSHFQ, 0)
}
func RETURN() lib.Instruction {
	return opcodes.OpcodeToInstruction("return", opcodes.RETURN, 0)
}
func SETA(dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("seta", opcodes.SETA, 1, dest)
}
func SETAE(dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("setae", opcodes.SETAE, 1, dest)
}
func SETB(dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("seta", opcodes.SETB, 1, dest)
}
func SETBE(dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("setae", opcodes.SETBE, 1, dest)
}
func SETC(dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("seta", opcodes.SETC, 1, dest)
}
func SETE(dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("sete", opcodes.SETE, 1, dest)
}
func SETNE(dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("setne", opcodes.SETNE, 1, dest)
}
func SUB(src, dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("sub", opcodes.SUB, 2, dest, src)
}
func SYSCALL() lib.Instruction {
	return opcodes.OpcodeToInstruction("syscall", opcodes.SYSCALL, 0)
}
func XOR(src, dest encoding.Operand) lib.Instruction {
	return opcodes.OpcodesToInstruction("xor", opcodes.XOR, 2, dest, src)
}
