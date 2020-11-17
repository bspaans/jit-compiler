package opcodes

import (
	. "github.com/bspaans/jit-compiler/asm/encoding"
)

var (
	ADD_rm8_r8 = &Opcode{"add", []uint8{}, []uint8{0x00}, []OpcodeExtensions{Rex, SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm8, ModRM_rm_rw},
			OpcodeOperand{OT_r8, ModRM_reg_r},
		},
	}
	ADD_r8_rm8 = &Opcode{"add", []uint8{}, []uint8{0x02}, []OpcodeExtensions{Rex, SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r8, ModRM_reg_rw},
			OpcodeOperand{OT_rm8, ModRM_rm_r},
		},
	}
	ADD_rm16_r16 = &Opcode{"add", []uint8{0x66}, []uint8{0x01}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm16, ModRM_rm_rw},
			OpcodeOperand{OT_r16, ModRM_reg_r},
		},
	}
	ADD_r16_rm16 = &Opcode{"add", []uint8{0x66}, []uint8{0x03}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r16, ModRM_reg_rw},
			OpcodeOperand{OT_rm16, ModRM_rm_r},
		},
	}
	ADD_rm32_r32 = &Opcode{"add", []uint8{}, []uint8{0x01}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm32, ModRM_rm_rw},
			OpcodeOperand{OT_r32, ModRM_reg_r},
		},
	}
	ADD_r32_rm32 = &Opcode{"add", []uint8{}, []uint8{0x03}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r32, ModRM_reg_rw},
			OpcodeOperand{OT_rm32, ModRM_rm_r},
		},
	}
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
	// Logical AND
	AND_r8_rm8 = &Opcode{"and", []uint8{}, []uint8{0x22}, []OpcodeExtensions{Rex, SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r8, ModRM_reg_rw},
			OpcodeOperand{OT_rm8, ModRM_rm_r},
		},
	}
	AND_r8_rm8_no_rex = &Opcode{"and", []uint8{}, []uint8{0x22}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r8, ModRM_reg_rw},
			OpcodeOperand{OT_rm8, ModRM_rm_r},
		},
	}
	AND_rm8_r8 = &Opcode{"and", []uint8{}, []uint8{0x20}, []OpcodeExtensions{Rex, SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm8, ModRM_rm_rw},
			OpcodeOperand{OT_r8, ModRM_reg_r},
		},
	}
	AND_rm8_r8_no_rex = &Opcode{"and", []uint8{}, []uint8{0x20}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm8, ModRM_rm_rw},
			OpcodeOperand{OT_r8, ModRM_reg_r},
		},
	}
	AND_rm16_r16 = &Opcode{"and", []uint8{0xff}, []uint8{0x21}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm16, ModRM_rm_rw},
			OpcodeOperand{OT_r16, ModRM_reg_r},
		},
	}
	AND_r16_rm16 = &Opcode{"and", []uint8{0x66}, []uint8{0x23}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r16, ModRM_reg_rw},
			OpcodeOperand{OT_rm16, ModRM_rm_r},
		},
	}
	AND_rm32_r32 = &Opcode{"and", []uint8{}, []uint8{0x21}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm32, ModRM_rm_rw},
			OpcodeOperand{OT_r32, ModRM_reg_r},
		},
	}
	AND_r32_rm32 = &Opcode{"and", []uint8{}, []uint8{0x23}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r32, ModRM_reg_rw},
			OpcodeOperand{OT_rm32, ModRM_rm_r},
		},
	}
	AND_rm64_r64 = &Opcode{"and", []uint8{}, []uint8{0x21}, []OpcodeExtensions{RexW, SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm64, ModRM_rm_rw},
			OpcodeOperand{OT_r64, ModRM_reg_r},
		},
	}
	AND_r64_rm64 = &Opcode{"and", []uint8{}, []uint8{0x23}, []OpcodeExtensions{RexW, SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r64, ModRM_reg_rw},
			OpcodeOperand{OT_rm64, ModRM_rm_r},
		},
	}
	CALL_rm64 = &Opcode{"call", []uint8{}, []uint8{0xff}, []OpcodeExtensions{Slash2},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm64, ModRM_rm_rw},
		},
	}
	CMP_rm8_imm8 = &Opcode{"cmp", []uint8{}, []uint8{0x80}, []OpcodeExtensions{Rex, Slash7, ImmediateByte},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm8, ModRM_rm_r},
			OpcodeOperand{OT_imm8, ImmediateValue},
		},
	}
	CMP_rm8_imm8_no_rex = &Opcode{"cmp", []uint8{}, []uint8{0x80}, []OpcodeExtensions{Slash7, ImmediateByte},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm8, ModRM_rm_r},
			OpcodeOperand{OT_imm8, ImmediateValue},
		},
	}
	CMP_rm64_imm32 = &Opcode{"cmp", []uint8{}, []uint8{0x81}, []OpcodeExtensions{RexW, Slash7, ImmediateDouble},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm64, ModRM_rm_r},
			OpcodeOperand{OT_imm32, ImmediateValue},
		},
	}
	CMP_r8_rm8 = &Opcode{"cmp", []uint8{}, []uint8{0x3a}, []OpcodeExtensions{Rex, SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r8, ModRM_reg_r},
			OpcodeOperand{OT_rm8, ModRM_rm_r},
		},
	}
	CMP_r8_rm8_no_rex = &Opcode{"cmp", []uint8{}, []uint8{0x3a}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r8, ModRM_reg_r},
			OpcodeOperand{OT_rm8, ModRM_rm_r},
		},
	}
	CMP_rm8_r8 = &Opcode{"cmp", []uint8{}, []uint8{0x38}, []OpcodeExtensions{Rex, SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm8, ModRM_rm_r},
			OpcodeOperand{OT_r8, ModRM_reg_r},
		},
	}
	CMP_rm8_r8_no_rex = &Opcode{"cmp", []uint8{}, []uint8{0x38}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm8, ModRM_rm_r},
			OpcodeOperand{OT_r8, ModRM_reg_r},
		},
	}
	CMP_r16_rm16 = &Opcode{"cmp", []uint8{0x66}, []uint8{0x3b}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r16, ModRM_reg_r},
			OpcodeOperand{OT_rm16, ModRM_rm_r},
		},
	}
	CMP_rm16_r16 = &Opcode{"cmp", []uint8{0x66}, []uint8{0x39}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm16, ModRM_rm_r},
			OpcodeOperand{OT_r16, ModRM_reg_r},
		},
	}
	CMP_r32_rm32 = &Opcode{"cmp", []uint8{}, []uint8{0x3b}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r32, ModRM_reg_r},
			OpcodeOperand{OT_rm32, ModRM_rm_r},
		},
	}
	CMP_rm32_r32 = &Opcode{"cmp", []uint8{}, []uint8{0x39}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm32, ModRM_rm_r},
			OpcodeOperand{OT_r32, ModRM_reg_r},
		},
	}
	CMP_r64_rm64 = &Opcode{"cmp", []uint8{}, []uint8{0x3b}, []OpcodeExtensions{RexW, SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r64, ModRM_reg_r},
			OpcodeOperand{OT_rm64, ModRM_rm_r},
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
	// Convert Byte to Word; ax = sign extend(al)
	CBW = &Opcode{"cbw", []uint8{0x66}, []uint8{0x98}, []OpcodeExtensions{}, []OpcodeOperand{}}
	// Convert Word to Doubleword; dx:ax = sign extend(ax)
	CWD = &Opcode{"cwd", []uint8{0x66}, []uint8{0x99}, []OpcodeExtensions{}, []OpcodeOperand{}}
	// Convert Double word to Quadword; edx:eax = sign extend(eax)
	CDQ = &Opcode{"cdq", []uint8{}, []uint8{0x99}, []OpcodeExtensions{}, []OpcodeOperand{}}
	// Convert Quad word to double quad word; rdx:rax = sign extend(rax)
	CQO = &Opcode{"cqo", []uint8{}, []uint8{0x99}, []OpcodeExtensions{RexW}, []OpcodeOperand{}}

	DEC_rm64 = &Opcode{"dec", []uint8{}, []uint8{0xff}, []OpcodeExtensions{RexW, Slash1},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm64, ModRM_rm_rw},
		},
	}
	DIV_rm8 = &Opcode{"div", []uint8{}, []uint8{0xf6}, []OpcodeExtensions{Rex, Slash6},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm8, ModRM_rm_r},
		},
	}
	DIV_rm16 = &Opcode{"div", []uint8{0x66}, []uint8{0xf7}, []OpcodeExtensions{Slash6},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm16, ModRM_rm_r},
		},
	}
	DIV_rm32 = &Opcode{"div", []uint8{}, []uint8{0xf7}, []OpcodeExtensions{Slash6},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm32, ModRM_rm_r},
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
	IDIV_rm8 = &Opcode{"idiv", []uint8{}, []uint8{0xf6}, []OpcodeExtensions{RexW, Slash7},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm8, ModRM_rm_r},
		},
	}
	IDIV_rm8_no_rex = &Opcode{"idiv", []uint8{}, []uint8{0xf6}, []OpcodeExtensions{Slash7},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm8, ModRM_rm_r},
		},
	}
	IDIV_rm16 = &Opcode{"idiv", []uint8{0x66}, []uint8{0xf7}, []OpcodeExtensions{Slash7},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm16, ModRM_rm_r},
		},
	}
	IDIV_rm32 = &Opcode{"idiv", []uint8{}, []uint8{0xf7}, []OpcodeExtensions{Slash7},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm32, ModRM_rm_r},
		},
	}
	IDIV_rm64 = &Opcode{"idiv", []uint8{}, []uint8{0xf7}, []OpcodeExtensions{RexW, Slash7},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm64, ModRM_rm_r},
		},
	}
	INC_rm64 = &Opcode{"inc", []uint8{}, []uint8{0xff}, []OpcodeExtensions{RexW, Slash0},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm64, ModRM_rm_rw},
		},
	}
	IMUL_rm8 = &Opcode{"imul", []uint8{}, []uint8{0xf6}, []OpcodeExtensions{RexW, Slash5},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm8, ModRM_rm_rw},
		},
	}
	IMUL_rm8_no_rex = &Opcode{"imul", []uint8{}, []uint8{0xf6}, []OpcodeExtensions{Slash5},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm8, ModRM_rm_rw},
		},
	}
	IMUL_rm16 = &Opcode{"imul", []uint8{0x66}, []uint8{0xf7}, []OpcodeExtensions{Slash5},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm16, ModRM_rm_rw},
		},
	}
	IMUL_rm32 = &Opcode{"imul", []uint8{}, []uint8{0xf7}, []OpcodeExtensions{Slash5},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm32, ModRM_rm_rw},
		},
	}
	IMUL_rm64 = &Opcode{"imul", []uint8{}, []uint8{0xf7}, []OpcodeExtensions{RexW, Slash5},
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
	// Jump short if above (CF=0 or ZF=0) (for unsigned)
	JA_rel8 = &Opcode{"ja", []uint8{}, []uint8{0x77}, []OpcodeExtensions{ImmediateByte},
		[]OpcodeOperand{
			OpcodeOperand{OT_rel8, ImmediateValue},
		},
	}
	// Jump short if above or equal (CF=0) (for unsigned)
	JAE_rel8 = &Opcode{"jae", []uint8{}, []uint8{0x73}, []OpcodeExtensions{ImmediateByte},
		[]OpcodeOperand{
			OpcodeOperand{OT_rel8, ImmediateValue},
		},
	}
	// Jump short if below (CF=1)
	JB_rel8 = &Opcode{"jb", []uint8{}, []uint8{0x72}, []OpcodeExtensions{ImmediateByte},
		[]OpcodeOperand{
			OpcodeOperand{OT_rel8, ImmediateValue},
		},
	}
	// Jump short if below (CF=1 or ZF=0)
	JBE_rel8 = &Opcode{"jbe", []uint8{}, []uint8{0x76}, []OpcodeExtensions{ImmediateByte},
		[]OpcodeOperand{
			OpcodeOperand{OT_rel8, ImmediateValue},
		},
	}
	// Jump short if equal (ZF=1)
	JE_rel8 = &Opcode{"je", []uint8{}, []uint8{0x74}, []OpcodeExtensions{ImmediateByte},
		[]OpcodeOperand{
			OpcodeOperand{OT_rel8, ImmediateValue},
		},
	}
	// Jump short if greater (ZF=0 and SF=OF) (for signed)
	JG_rel8 = &Opcode{"jg", []uint8{}, []uint8{0x7f}, []OpcodeExtensions{ImmediateByte},
		[]OpcodeOperand{
			OpcodeOperand{OT_rel8, ImmediateValue},
		},
	}
	// Jump short if greater or equal (SF=OF) (for signed)
	JGE_rel8 = &Opcode{"jge", []uint8{}, []uint8{0x7d}, []OpcodeExtensions{ImmediateByte},
		[]OpcodeOperand{
			OpcodeOperand{OT_rel8, ImmediateValue},
		},
	}
	// Jump short if less (SF!=OF) (for signed)
	JL_rel8 = &Opcode{"jl", []uint8{}, []uint8{0x7c}, []OpcodeExtensions{ImmediateByte},
		[]OpcodeOperand{
			OpcodeOperand{OT_rel8, ImmediateValue},
		},
	}
	// Jump short if less or equal (SF!=OF) (for signed)
	JLE_rel8 = &Opcode{"jle", []uint8{}, []uint8{0x7e}, []OpcodeExtensions{ImmediateByte},
		[]OpcodeOperand{
			OpcodeOperand{OT_rel8, ImmediateValue},
		},
	}
	// Jump short if not above (ZF=0)
	JNA_rel8 = &Opcode{"jna", []uint8{}, []uint8{0x76}, []OpcodeExtensions{ImmediateByte},
		[]OpcodeOperand{
			OpcodeOperand{OT_rel8, ImmediateValue},
		},
	}
	// Jump short if not above or equal (CF=1)
	JNAE_rel8 = &Opcode{"jnae", []uint8{}, []uint8{0x72}, []OpcodeExtensions{ImmediateByte},
		[]OpcodeOperand{
			OpcodeOperand{OT_rel8, ImmediateValue},
		},
	}
	// Jump short if not below (CF=0)
	JNB_rel8 = &Opcode{"jnb", []uint8{}, []uint8{0x73}, []OpcodeExtensions{ImmediateByte},
		[]OpcodeOperand{
			OpcodeOperand{OT_rel8, ImmediateValue},
		},
	}
	// Jump short if not below or equal (CF=0 or ZF=0)
	JNBE_rel8 = &Opcode{"jnbe", []uint8{}, []uint8{0x77}, []OpcodeExtensions{ImmediateByte},
		[]OpcodeOperand{
			OpcodeOperand{OT_rel8, ImmediateValue},
		},
	}
	// Jump short if not equal (ZF=0)
	JNE_rel8 = &Opcode{"jne", []uint8{}, []uint8{0x75}, []OpcodeExtensions{ImmediateByte},
		[]OpcodeOperand{
			OpcodeOperand{OT_rel8, ImmediateValue},
		},
	}
	// Jump short if not greater (ZF=1 or SF!=0)
	JNG_rel8 = &Opcode{"jng", []uint8{}, []uint8{0x7e}, []OpcodeExtensions{ImmediateByte},
		[]OpcodeOperand{
			OpcodeOperand{OT_rel8, ImmediateValue},
		},
	}
	// Jump short if not greater or equal (SF!=0)
	JNGE_rel8 = &Opcode{"jnge", []uint8{}, []uint8{0x7c}, []OpcodeExtensions{ImmediateByte},
		[]OpcodeOperand{
			OpcodeOperand{OT_rel8, ImmediateValue},
		},
	}
	// Jump short if not less (SF=OF)
	JNL_rel8 = &Opcode{"jnl", []uint8{}, []uint8{0x7d}, []OpcodeExtensions{ImmediateByte},
		[]OpcodeOperand{
			OpcodeOperand{OT_rel8, ImmediateValue},
		},
	}
	// Jump short if not less or equal (ZF=0 and SF=OF)
	JNLE_rel8 = &Opcode{"jnle", []uint8{}, []uint8{0x7f}, []OpcodeExtensions{ImmediateByte},
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
	MOV_r8_imm8_no_rex = &Opcode{"mov", []uint8{}, []uint8{0xb0}, []OpcodeExtensions{ImmediateByte},
		[]OpcodeOperand{
			OpcodeOperand{OT_r8, Opcode_plus_rd_r},
			OpcodeOperand{OT_imm8, ImmediateValue},
		},
	}
	MOV_r8_imm8 = &Opcode{"mov", []uint8{}, []uint8{0xb0}, []OpcodeExtensions{Rex, ImmediateByte},
		[]OpcodeOperand{
			OpcodeOperand{OT_r8, Opcode_plus_rd_r},
			OpcodeOperand{OT_imm8, ImmediateValue},
		},
	}
	MOV_rm16_r16 = &Opcode{"mov", []uint8{0x66}, []uint8{0x89}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm16, ModRM_rm_rw},
			OpcodeOperand{OT_r16, ModRM_reg_r},
		},
	}
	MOV_r16_rm16 = &Opcode{"mov", []uint8{0x66}, []uint8{0x8b}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r16, ModRM_reg_rw},
			OpcodeOperand{OT_rm16, ModRM_rm_r},
		},
	}
	MOV_r16_imm16 = &Opcode{"mov", []uint8{0x66}, []uint8{0xb8}, []OpcodeExtensions{ImmediateWord},
		[]OpcodeOperand{
			OpcodeOperand{OT_r16, Opcode_plus_rd_r},
			OpcodeOperand{OT_imm16, ImmediateValue},
		},
	}
	MOV_r32_imm32 = &Opcode{"mov", []uint8{}, []uint8{0xb8}, []OpcodeExtensions{ImmediateDouble},
		[]OpcodeOperand{
			OpcodeOperand{OT_r32, Opcode_plus_rd_r},
			OpcodeOperand{OT_imm32, ImmediateValue},
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
	// Move or Merge Scalar Double-Precision Floating-Point Value
	MOVSD_xmm1m64_xmm2 = &Opcode{"movsd", []uint8{}, []uint8{0xf2, 0x0f, 0x11}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_xmm1m64, ModRM_rm_rw},
			OpcodeOperand{OT_xmm2, ModRM_reg_r},
		},
	}
	// Move with sign-extend
	MOVSX_r16_rm8 = &Opcode{"movsx", []uint8{0x66}, []uint8{0x0f, 0xbe}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r32, ModRM_reg_rw},
			OpcodeOperand{OT_rm8, ModRM_rm_r},
		},
	}
	MOVSX_r32_rm8 = &Opcode{"movsx", []uint8{}, []uint8{0x0f, 0xbe}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r32, ModRM_reg_rw},
			OpcodeOperand{OT_rm8, ModRM_rm_r},
		},
	}
	MOVSX_r32_rm16 = &Opcode{"movsx", []uint8{}, []uint8{0x0f, 0xbf}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r32, ModRM_reg_rw},
			OpcodeOperand{OT_rm32, ModRM_rm_r},
		},
	}
	MOVSX_r64_rm8 = &Opcode{"movsx", []uint8{}, []uint8{0x0f, 0xbe}, []OpcodeExtensions{RexW, SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r64, ModRM_reg_rw},
			OpcodeOperand{OT_rm8, ModRM_rm_r},
		},
	}
	MOVSX_r64_rm16 = &Opcode{"movsx", []uint8{}, []uint8{0x0f, 0xbf}, []OpcodeExtensions{RexW, SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r64, ModRM_reg_rw},
			OpcodeOperand{OT_rm16, ModRM_rm_r},
		},
	}
	MOVSX_r64_rm32 = &Opcode{"movsx", []uint8{}, []uint8{0x63}, []OpcodeExtensions{RexW, SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r64, ModRM_reg_rw},
			OpcodeOperand{OT_rm32, ModRM_rm_r},
		},
	}
	// Move with zero-extend
	MOVZX_r16_rm8 = &Opcode{"movzx", []uint8{0x66}, []uint8{0x0f, 0xb6}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r16, ModRM_reg_rw},
			OpcodeOperand{OT_rm8, ModRM_rm_r},
		},
	}
	MOVZX_r32_rm8 = &Opcode{"movzx", []uint8{}, []uint8{0x0f, 0xb6}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r32, ModRM_reg_rw},
			OpcodeOperand{OT_rm8, ModRM_rm_r},
		},
	}
	MOVZX_r64_rm8 = &Opcode{"movzx", []uint8{}, []uint8{0x0f, 0xb6}, []OpcodeExtensions{RexW, SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r64, ModRM_reg_rw},
			OpcodeOperand{OT_rm8, ModRM_rm_r},
		},
	}
	MOVZX_r32_rm16 = &Opcode{"movzx", []uint8{}, []uint8{0x0f, 0xb7}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r32, ModRM_reg_rw},
			OpcodeOperand{OT_rm16, ModRM_rm_r},
		},
	}
	MOVZX_r64_rm16 = &Opcode{"movzx", []uint8{}, []uint8{0x0f, 0xb7}, []OpcodeExtensions{RexW, SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r64, ModRM_reg_rw},
			OpcodeOperand{OT_rm16, ModRM_rm_r},
		},
	}
	MUL_rm8 = &Opcode{"mul", []uint8{}, []uint8{0xf6}, []OpcodeExtensions{Rex, Slash4},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm8, ModRM_rm_r},
		},
	}
	MUL_rm16 = &Opcode{"mul", []uint8{0x66}, []uint8{0xf7}, []OpcodeExtensions{Rex, Slash4},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm16, ModRM_rm_r},
		},
	}
	MUL_rm32 = &Opcode{"mul", []uint8{}, []uint8{0xf7}, []OpcodeExtensions{Rex, Slash4},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm32, ModRM_rm_r},
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
	// Logical OR
	OR_r8_rm8 = &Opcode{"or", []uint8{}, []uint8{0x0a}, []OpcodeExtensions{Rex, SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r8, ModRM_reg_rw},
			OpcodeOperand{OT_rm8, ModRM_rm_r},
		},
	}
	OR_r8_rm8_no_rex = &Opcode{"or", []uint8{}, []uint8{0x0a}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r8, ModRM_reg_rw},
			OpcodeOperand{OT_rm8, ModRM_rm_r},
		},
	}
	OR_rm8_r8 = &Opcode{"or", []uint8{}, []uint8{0x08}, []OpcodeExtensions{Rex, SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm8, ModRM_rm_rw},
			OpcodeOperand{OT_r8, ModRM_reg_r},
		},
	}
	OR_rm8_r8_no_rex = &Opcode{"or", []uint8{}, []uint8{0x08}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm8, ModRM_rm_rw},
			OpcodeOperand{OT_r8, ModRM_reg_r},
		},
	}
	OR_rm16_r16 = &Opcode{"or", []uint8{0x66}, []uint8{0x09}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm16, ModRM_rm_rw},
			OpcodeOperand{OT_r16, ModRM_reg_r},
		},
	}
	OR_r16_rm16 = &Opcode{"or", []uint8{0x66}, []uint8{0x0b}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r16, ModRM_reg_rw},
			OpcodeOperand{OT_rm16, ModRM_rm_r},
		},
	}
	OR_rm32_r32 = &Opcode{"or", []uint8{}, []uint8{0x09}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm32, ModRM_rm_rw},
			OpcodeOperand{OT_r32, ModRM_reg_r},
		},
	}
	OR_r32_rm32 = &Opcode{"or", []uint8{}, []uint8{0x0b}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r32, ModRM_reg_rw},
			OpcodeOperand{OT_rm32, ModRM_rm_r},
		},
	}
	OR_rm64_r64 = &Opcode{"or", []uint8{}, []uint8{0x09}, []OpcodeExtensions{RexW, SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm64, ModRM_rm_rw},
			OpcodeOperand{OT_r64, ModRM_reg_r},
		},
	}
	OR_r64_rm64 = &Opcode{"or", []uint8{}, []uint8{0x0b}, []OpcodeExtensions{RexW, SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r64, ModRM_reg_rw},
			OpcodeOperand{OT_rm64, ModRM_rm_r},
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
	SHL_rm8_imm8 = &Opcode{"shl", []uint8{}, []uint8{0xc0}, []OpcodeExtensions{RexW, Slash4, ImmediateByte},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm8, ModRM_rm_rw},
			OpcodeOperand{OT_imm8, ImmediateValue},
		},
	}
	SHL_rm8_imm8_no_rex = &Opcode{"shl", []uint8{}, []uint8{0xc0}, []OpcodeExtensions{Slash4, ImmediateByte},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm8, ModRM_rm_rw},
			OpcodeOperand{OT_imm8, ImmediateValue},
		},
	}
	SHL_rm16_imm8 = &Opcode{"shl", []uint8{0x66}, []uint8{0xc1}, []OpcodeExtensions{Slash4, ImmediateByte},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm16, ModRM_rm_rw},
			OpcodeOperand{OT_imm8, ImmediateValue},
		},
	}
	SHL_rm32_imm8 = &Opcode{"shl", []uint8{}, []uint8{0xc1}, []OpcodeExtensions{Slash4, ImmediateByte},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm32, ModRM_rm_rw},
			OpcodeOperand{OT_imm8, ImmediateValue},
		},
	}
	SHL_rm64_imm8 = &Opcode{"shl", []uint8{}, []uint8{0xc1}, []OpcodeExtensions{RexW, Slash4, ImmediateByte},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm64, ModRM_rm_rw},
			OpcodeOperand{OT_imm8, ImmediateValue},
		},
	}
	SHR_rm8_imm8 = &Opcode{"shr", []uint8{}, []uint8{0xc0}, []OpcodeExtensions{RexW, Slash5, ImmediateByte},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm8, ModRM_rm_rw},
			OpcodeOperand{OT_imm8, ImmediateValue},
		},
	}
	SHR_rm8_imm8_no_rex = &Opcode{"shr", []uint8{}, []uint8{0xc0}, []OpcodeExtensions{Slash5, ImmediateByte},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm8, ModRM_rm_rw},
			OpcodeOperand{OT_imm8, ImmediateValue},
		},
	}
	SHR_rm16_imm8 = &Opcode{"shr", []uint8{0x66}, []uint8{0xc1}, []OpcodeExtensions{Slash5, ImmediateByte},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm16, ModRM_rm_rw},
			OpcodeOperand{OT_imm8, ImmediateValue},
		},
	}
	SHR_rm32_imm8 = &Opcode{"shr", []uint8{}, []uint8{0xc1}, []OpcodeExtensions{Slash5, ImmediateByte},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm32, ModRM_rm_rw},
			OpcodeOperand{OT_imm8, ImmediateValue},
		},
	}
	SHR_rm64_imm8 = &Opcode{"shr", []uint8{}, []uint8{0xc1}, []OpcodeExtensions{RexW, Slash5, ImmediateByte},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm64, ModRM_rm_rw},
			OpcodeOperand{OT_imm8, ImmediateValue},
		},
	}
	SUB_rm8_imm8 = &Opcode{"sub", []uint8{}, []uint8{0x80}, []OpcodeExtensions{Rex, Slash5, ImmediateByte},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm8, ModRM_rm_rw},
			OpcodeOperand{OT_imm8, ImmediateValue},
		},
	}
	SUB_r8_rm8 = &Opcode{"sub", []uint8{}, []uint8{0x2a}, []OpcodeExtensions{Rex, SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r8, ModRM_reg_rw},
			OpcodeOperand{OT_rm8, ModRM_rm_r},
		},
	}
	SUB_rm8_r8 = &Opcode{"sub", []uint8{}, []uint8{0x28}, []OpcodeExtensions{Rex, SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm8, ModRM_rm_rw},
			OpcodeOperand{OT_r8, ModRM_reg_r},
		},
	}
	SUB_rm16_r16 = &Opcode{"sub", []uint8{0x66}, []uint8{0x29}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm16, ModRM_rm_rw},
			OpcodeOperand{OT_r16, ModRM_reg_r},
		},
	}
	SUB_r16_rm16 = &Opcode{"sub", []uint8{0x66}, []uint8{0x2b}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r16, ModRM_reg_rw},
			OpcodeOperand{OT_rm16, ModRM_rm_r},
		},
	}
	SUB_rm32_r32 = &Opcode{"sub", []uint8{}, []uint8{0x29}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm32, ModRM_rm_rw},
			OpcodeOperand{OT_r32, ModRM_reg_r},
		},
	}
	SUB_r32_rm32 = &Opcode{"sub", []uint8{}, []uint8{0x2b}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r32, ModRM_reg_rw},
			OpcodeOperand{OT_rm32, ModRM_rm_r},
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
	SUB_rm64_imm8 = &Opcode{"sub", []uint8{}, []uint8{0x83}, []OpcodeExtensions{Rex, Slash5, ImmediateByte},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm64, ModRM_rm_rw},
			OpcodeOperand{OT_imm8, ImmediateValue},
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
	XOR_r8_rm8 = &Opcode{"xor", []uint8{}, []uint8{0x32}, []OpcodeExtensions{Rex, SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r8, ModRM_reg_rw},
			OpcodeOperand{OT_rm8, ModRM_rm_r},
		},
	}
	XOR_r8_rm8_no_rex = &Opcode{"xor", []uint8{}, []uint8{0x32}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r8, ModRM_reg_rw},
			OpcodeOperand{OT_rm8, ModRM_rm_r},
		},
	}
	XOR_rm8_r8 = &Opcode{"xor", []uint8{}, []uint8{0x30}, []OpcodeExtensions{RexW, SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm8, ModRM_rm_rw},
			OpcodeOperand{OT_r8, ModRM_reg_r},
		},
	}
	XOR_rm8_r8_no_rex = &Opcode{"xor", []uint8{}, []uint8{0x30}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm8, ModRM_rm_rw},
			OpcodeOperand{OT_r8, ModRM_reg_r},
		},
	}
	XOR_r16_rm16 = &Opcode{"xor", []uint8{0x66}, []uint8{0x33}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r16, ModRM_reg_rw},
			OpcodeOperand{OT_rm16, ModRM_rm_r},
		},
	}
	XOR_rm16_r16 = &Opcode{"xor", []uint8{0x66}, []uint8{0x31}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm16, ModRM_rm_rw},
			OpcodeOperand{OT_r16, ModRM_reg_r},
		},
	}
	XOR_r32_rm32 = &Opcode{"xor", []uint8{}, []uint8{0x33}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_r32, ModRM_reg_rw},
			OpcodeOperand{OT_rm32, ModRM_rm_r},
		},
	}
	XOR_rm32_r32 = &Opcode{"xor", []uint8{}, []uint8{0x31}, []OpcodeExtensions{SlashR},
		[]OpcodeOperand{
			OpcodeOperand{OT_rm32, ModRM_rm_rw},
			OpcodeOperand{OT_r32, ModRM_reg_r},
		},
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
