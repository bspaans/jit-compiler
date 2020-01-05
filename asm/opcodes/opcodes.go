package opcodes

import (
	. "github.com/bspaans/jit/asm/encoding"
)

var (
	ADD_rm64_r64 = &Opcode{"add", []uint8{}, []uint8{0x01}, []OpcodeExtensions{RexW, SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm64, ModRM_rm_rw},
			OpcodeOperand{OT_r64, ModRM_reg_r},
		},
	}
	ADD_r64_rm64 = &Opcode{"add", []uint8{}, []uint8{0x03}, []OpcodeExtensions{RexW, SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r64, ModRM_reg_rw},
			OpcodeOperand{OT_rm64, ModRM_rm_r},
		},
	}
	ADD_rm64_imm32 = &Opcode{"add", []uint8{}, []uint8{0x81}, []OpcodeExtensions{RexW, Slash0, ImmediateDouble},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm64, ModRM_rm_rw},
			OpcodeOperand{OT_imm32, ImmediateValue},
		},
	}
	// Add packed double-precision floating point values from xmm2/mem to xmm1 and store result in xmm1
	ADDPD_xmm1_xmm2m128 = &Opcode{"addpd", []uint8{}, []uint8{0x66, 0x0f, 0x58}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_xmm1, ModRM_reg_rw},
			OpcodeOperand{OT_xmm2m128, ModRM_rm_r},
		},
	}
	// Add the low double-precision floating-point value from xmm2/mem to xmm1 and store the result in xmm1
	ADDSD_xmm1_xmm2m64 = &Opcode{"addsd", []uint8{}, []uint8{0xf2, 0x0f, 0x58}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_xmm1, ModRM_reg_rw},
			OpcodeOperand{OT_xmm2m64, ModRM_rm_r},
		},
	}
	CALL_rm64 = &Opcode{"call", []uint8{}, []uint8{0xff}, []OpcodeExtensions{Slash2},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm64, ModRM_rm_rw},
		},
	}
	CMP_rm64_imm32 = &Opcode{"cmp", []uint8{}, []uint8{0x81}, []OpcodeExtensions{RexW, Slash7},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm64, ModRM_rm_r},
			OpcodeOperand{OT_imm32, ImmediateValue},
		},
	}
	CMP_rm64_r64 = &Opcode{"cmp", []uint8{}, []uint8{0x39}, []OpcodeExtensions{RexW, SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm64, ModRM_rm_r},
			OpcodeOperand{OT_r64, ModRM_reg_r},
		},
	}
	// Convert Doubleword integer to Scalar Double-precision floating-point value
	CVTSI2SD_xmm1_rm64 = &Opcode{"cvtsi2sd", []uint8{0xf2}, []uint8{0x0f, 0x2a}, []OpcodeExtensions{RexW, SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_xmm1, ModRM_reg_rw},
			OpcodeOperand{OT_rm64, ModRM_rm_r},
		},
	}
	// Convert Scalar Double precision floating point to Doubleword integer
	CVTSD2SI_r64_xmm1m64 = &Opcode{"cvtsd2si", []uint8{0xf2}, []uint8{0x0f, 0x2d}, []OpcodeExtensions{RexW, SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r64, ModRM_reg_rw},
			OpcodeOperand{OT_xmm2m64, ModRM_rm_r},
		},
	}
	// Convert with truncation Scalar Double precision floating point to Signed Integer
	CVTTSD2SI_r64_xmm1m64 = &Opcode{"cvttsd2si", []uint8{0xf2}, []uint8{0x0f, 0x2c}, []OpcodeExtensions{RexW, SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r64, ModRM_reg_rw},
			OpcodeOperand{OT_xmm2m64, ModRM_rm_r},
		},
	}
	DEC_rm64 = &Opcode{"dec", []uint8{}, []uint8{0xff}, []OpcodeExtensions{RexW, Slash1},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm64, ModRM_rm_rw},
		},
	}
	// Divides unsigned the value in the AX, DX:AX, EDX:EAX, or RDX:RAX
	// registers (dividend) by the source operand (divisor) and stores the
	// result in the AX (AH:AL), DX:AX, EDX:EAX, or RDX:RAX registers.
	DIV_rm64 = &Opcode{"div", []uint8{}, []uint8{0xf7}, []OpcodeExtensions{RexW, Slash6},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm64, ModRM_rm_r},
		},
	}
	DIVSD_xmm1_xmm2m64 = &Opcode{"divsd", []uint8{}, []uint8{0xf2, 0x0f, 0x5e}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_xmm1, ModRM_reg_rw},
			OpcodeOperand{OT_xmm2m64, ModRM_rm_r},
		},
	}
	INC_rm64 = &Opcode{"inc", []uint8{}, []uint8{0xff}, []OpcodeExtensions{RexW, Slash0},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm64, ModRM_rm_rw},
		},
	}
	IMUL_r64_rm64 = &Opcode{"imul", []uint8{}, []uint8{0x0f, 0xaf}, []OpcodeExtensions{RexW, Slash6},
		[]OpcodeOperand{
			OpcodeOperand{OT_r64, ModRM_reg_rw},
			OpcodeOperand{OT_rm64, ModRM_rm_rw},
		},
	}
	// Jump short, RIP = RIP + 8 bit displacement sign extended to 64 bits
	JMP_rel8 = &Opcode{"jmp", []uint8{}, []uint8{0xeb}, []OpcodeExtensions{ImmediateByte},
		[]OpcodeOperand{
			OpcodeOperand{OT_rel8, ImmediateValue},
		},
	}
	// Jump near, relative, RIP = RIP + 32 bit displacement sign extended to 64 bits
	JMP_rel32 = &Opcode{"jmp", []uint8{}, []uint8{0xe9}, []OpcodeExtensions{ImmediateDouble},
		[]OpcodeOperand{
			OpcodeOperand{OT_rel32, ImmediateValue},
		},
	}
	// Jump near, absolute indirect, RIP = 64-Bit offset from register or memory
	JMP_rm64 = &Opcode{"jmp", []uint8{}, []uint8{0xff}, []OpcodeExtensions{Slash4},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm64, ModRM_rm_r},
		},
	}
	// Jump short if not equal (ZF=0)
	JNE_rel8 = &Opcode{"jne", []uint8{}, []uint8{0x75}, []OpcodeExtensions{ImmediateByte},
		[]OpcodeOperand{
			OpcodeOperand{OT_rel8, ImmediateValue},
		},
	}
	LEA_r64_m = &Opcode{"lea", []uint8{}, []uint8{0x8d}, []OpcodeExtensions{RexW, SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r64, ModRM_reg_rw},
			OpcodeOperand{OT_m, ModRM_rm_r},
		},
	}
	MOV_rm8_r8 = &Opcode{"mov", []uint8{}, []uint8{0x88}, []OpcodeExtensions{Rex, SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm8, ModRM_rm_rw},
			OpcodeOperand{OT_r8, ModRM_reg_r},
		},
	}
	MOV_r8_rm8 = &Opcode{"mov", []uint8{}, []uint8{0x8a}, []OpcodeExtensions{Rex, SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r8, ModRM_reg_rw},
			OpcodeOperand{OT_rm8, ModRM_rm_r},
		},
	}
	MOV_rm16_r16 = &Opcode{"mov", []uint8{}, []uint8{0x89}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm16, ModRM_rm_rw},
			OpcodeOperand{OT_r16, ModRM_reg_r},
		},
	}
	MOV_r16_rm16 = &Opcode{"mov", []uint8{}, []uint8{0x8b}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r16, ModRM_reg_rw},
			OpcodeOperand{OT_rm16, ModRM_rm_r},
		},
	}
	MOV_rm32_r32 = &Opcode{"mov", []uint8{}, []uint8{0x89}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm32, ModRM_rm_rw},
			OpcodeOperand{OT_r32, ModRM_reg_r},
		},
	}
	MOV_r32_rm32 = &Opcode{"mov", []uint8{}, []uint8{0x8b}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r32, ModRM_reg_rw},
			OpcodeOperand{OT_rm32, ModRM_rm_r},
		},
	}
	MOV_rm64_r64 = &Opcode{"mov", []uint8{}, []uint8{0x89}, []OpcodeExtensions{RexW, SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm64, ModRM_rm_rw},
			OpcodeOperand{OT_r64, ModRM_reg_r},
		},
	}
	MOV_r64_rm64 = &Opcode{"mov", []uint8{}, []uint8{0x8b}, []OpcodeExtensions{RexW, SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r64, ModRM_reg_rw},
			OpcodeOperand{OT_rm64, ModRM_rm_r},
		},
	}
	MOV_r64_imm64 = &Opcode{"mov", []uint8{}, []uint8{0xb8}, []OpcodeExtensions{RexW},
		[]OpcodeOperand{
			OpcodeOperand{OT_r64, Opcode_plus_rd_r},
			OpcodeOperand{OT_imm64, ImmediateValue},
		},
	}
	MOV_rm64_imm32 = &Opcode{"mov", []uint8{}, []uint8{0xc7}, []OpcodeExtensions{RexW, Slash0, ImmediateDouble},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm64, ModRM_rm_rw},
			OpcodeOperand{OT_imm32, ImmediateValue},
		},
	}
	MOVQ_xmm_rm64 = &Opcode{"movq", []uint8{0x66}, []uint8{0x0f, 0x6e}, []OpcodeExtensions{RexW, SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_xmm1, ModRM_reg_rw},
			OpcodeOperand{OT_rm64, ModRM_rm_r},
		},
	}
	MOVSD_xmm1m64_xmm2 = &Opcode{"movsd", []uint8{}, []uint8{0xf2, 0x0f, 0x11}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_xmm1m64, ModRM_rm_rw},
			OpcodeOperand{OT_xmm2, ModRM_reg_r},
		},
	}
	// Performs an unsigned multiplication of the first operand (destination
	// operand) and the second operand (source operand) and stores the result
	// in the destination operand. The destination operand is an implied
	// operand located in register AL, AX or EAX (depending on the size of the
	// operand);
	MUL_rm64 = &Opcode{"mul", []uint8{}, []uint8{0xf7}, []OpcodeExtensions{RexW, Slash4},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm64, ModRM_rm_r},
		},
	}
	MULSD_xmm1_xmm2m64 = &Opcode{"mulsd", []uint8{}, []uint8{0xf2, 0x0f, 0x59}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_xmm1, ModRM_reg_rw},
			OpcodeOperand{OT_xmm2m64, ModRM_rm_r},
		},
	}
	PUSH_imm32 = &Opcode{"push", []uint8{}, []uint8{0x68}, []OpcodeExtensions{},
		[]OpcodeOperand{
			OpcodeOperand{OT_imm32, ImmediateValue},
		},
	}
	PUSH_r64 = &Opcode{"push", []uint8{}, []uint8{0x50}, []OpcodeExtensions{},
		[]OpcodeOperand{
			OpcodeOperand{OT_r64, Opcode_plus_rd_r},
		},
	}
	// Push EFLAGS register onto the stack
	PUSHFQ = &Opcode{"pushfq", []uint8{}, []uint8{0x9c}, []OpcodeExtensions{},
		[]OpcodeOperand{},
	}
	POP_r64 = &Opcode{"pop", []uint8{}, []uint8{0x58}, []OpcodeExtensions{},
		[]OpcodeOperand{
			OpcodeOperand{OT_r64, Opcode_plus_rd_r},
		},
	}
	RETURN = &Opcode{"return", []uint8{}, []uint8{0xc3}, []OpcodeExtensions{},
		[]OpcodeOperand{},
	}
	// Set byte if above (CF=0, ZF=0)
	SETA_rm8 = &Opcode{"seta", []uint8{}, []uint8{0x0f, 0x97}, []OpcodeExtensions{},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm8, ModRM_rm_r},
		},
	}
	// Set byte if above or equal (CF=0, ZF=0)
	SETAE_rm8 = &Opcode{"setae", []uint8{}, []uint8{0x0f, 0x93}, []OpcodeExtensions{},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm8, ModRM_rm_r},
		},
	}
	// Set byte if below (CF=1)
	SETB_rm8 = &Opcode{"setb", []uint8{}, []uint8{0x0f, 0x92}, []OpcodeExtensions{},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm8, ModRM_rm_r},
		},
	}
	// Set byte if below or equal (CF=1 or ZF=1)
	SETBE_rm8 = &Opcode{"setbe", []uint8{}, []uint8{0x0f, 0x96}, []OpcodeExtensions{},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm8, ModRM_rm_r},
		},
	}
	// Set byte if carry (ZC=1)
	SETC_rm8 = &Opcode{"setc", []uint8{}, []uint8{0x0f, 0x92}, []OpcodeExtensions{},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm8, ModRM_rm_r},
		},
	}
	// Set byte if equal (ZF=1)
	SETE_rm8 = &Opcode{"sete", []uint8{}, []uint8{0x0f, 0x94}, []OpcodeExtensions{},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm8, ModRM_rm_r},
		},
	}
	// Set byte if not equal (ZF=0)
	SETNE_rm8 = &Opcode{"setne", []uint8{}, []uint8{0x0f, 0x95}, []OpcodeExtensions{},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm8, ModRM_rm_r},
		},
	}
	SUB_rm64_r64 = &Opcode{"sub", []uint8{}, []uint8{0x29}, []OpcodeExtensions{RexW, SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm64, ModRM_rm_rw},
			OpcodeOperand{OT_r64, ModRM_reg_r},
		},
	}
	SUB_r64_rm64 = &Opcode{"sub", []uint8{}, []uint8{0x2b}, []OpcodeExtensions{RexW, SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r64, ModRM_reg_rw},
			OpcodeOperand{OT_rm64, ModRM_rm_r},
		},
	}
	SUB_rm64_imm32 = &Opcode{"sub", []uint8{}, []uint8{0x81}, []OpcodeExtensions{RexW, Slash5, ImmediateDouble},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm64, ModRM_rm_rw},
			OpcodeOperand{OT_imm32, ImmediateValue},
		},
	}
	SUBSD_xmm1_xmm2m64 = &Opcode{"subsd", []uint8{}, []uint8{0xf2, 0x0f, 0x5c}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_xmm1, ModRM_reg_rw},
			OpcodeOperand{OT_xmm2m64, ModRM_rm_r},
		},
	}
	SYSCALL = &Opcode{"syscall", []uint8{}, []uint8{0x0f, 0x05}, []OpcodeExtensions{},
		[]OpcodeOperand{},
	}
	XOR_rm64_imm32 = &Opcode{"xor", []uint8{}, []uint8{0x81}, []OpcodeExtensions{RexW, Slash6, ImmediateDouble},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm64, ModRM_rm_rw},
			OpcodeOperand{OT_imm32, ImmediateValue},
		},
	}
	XOR_rm64_r64 = &Opcode{"xor", []uint8{}, []uint8{0x31}, []OpcodeExtensions{RexW, SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm64, ModRM_rm_rw},
			OpcodeOperand{OT_r64, ModRM_reg_r},
		},
	}
	XOR_r64_rm64 = &Opcode{"xor", []uint8{}, []uint8{0x33}, []OpcodeExtensions{RexW, SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r64, ModRM_reg_rw},
			OpcodeOperand{OT_rm64, ModRM_rm_r},
		},
	}
)
