package opcodes

import (
	. "github.com/bspaans/jit-compiler/asm/encoding"
)

var ADD = []*Opcode{
	ADD_rm8_r8,
	ADD_r8_rm8,
	ADD_rm16_r16,
	ADD_r16_rm16,
	ADD_rm32_r32,
	ADD_r32_rm32,
	ADD_rm64_r64,
	ADD_r64_rm64,
	ADD_rm64_imm32,
	ADDSD_xmm1_xmm2m64,
}
var AND = []*Opcode{
	AND_r8_rm8,
	AND_r8_rm8_no_rex,
	AND_rm8_r8,
	AND_rm8_r8_no_rex,
	AND_r16_rm16,
	AND_rm16_r16,
	AND_r32_rm32,
	AND_rm32_r32,
	AND_r64_rm64,
	AND_rm64_r64,
}
var CALL = []*Opcode{CALL_rm64}
var CMP = []*Opcode{
	CMP_rm8_imm8,
	CMP_rm8_imm8_no_rex,
	CMP_r8_rm8,
	CMP_r8_rm8_no_rex,
	CMP_rm8_r8,
	CMP_rm8_r8_no_rex,
	CMP_r16_rm16,
	CMP_rm16_r16,
	CMP_r32_rm32,
	CMP_rm32_r32,
	CMP_r64_rm64,
	CMP_rm64_r64,
	CMP_rm64_imm32,
}
var CVTSI2SD = []*Opcode{CVTSI2SD_xmm1_rm64}
var CVTSD2SI = []*Opcode{CVTSD2SI_r64_xmm1m64}
var CVTTSD2SI = []*Opcode{CVTTSD2SI_r64_xmm1m64}
var DEC = []*Opcode{DEC_rm64}
var IDIV1 = []*Opcode{
	IDIV_rm8,
	IDIV_rm8_no_rex,
	IDIV_rm16,
	IDIV_rm32,
	IDIV_rm64}
var IDIV2 = []*Opcode{DIVSD_xmm1_xmm2m64}
var DIV = []*Opcode{
	DIV_rm8,
	DIV_rm16,
	DIV_rm32,
	DIV_rm64,
}
var IMUL1 = []*Opcode{
	IMUL_rm8,
	IMUL_rm8_no_rex,
	IMUL_rm16,
	IMUL_rm32,
	IMUL_rm64,
}
var IMUL2 = []*Opcode{
	MULSD_xmm1_xmm2m64,
	IMUL_r64_rm64,
}
var MUL = []*Opcode{
	MUL_rm8,
	MUL_rm16,
	MUL_rm32,
	MUL_rm64,
}
var INC = []*Opcode{INC_rm64}
var JMP = []*Opcode{JMP_rel8, JMP_rel32, JMP_rm64}
var JA = []*Opcode{JA_rel8}
var JAE = []*Opcode{JAE_rel8}
var JB = []*Opcode{JB_rel8}
var JBE = []*Opcode{JBE_rel8}
var JE = []*Opcode{JE_rel8}
var JG = []*Opcode{JG_rel8}
var JGE = []*Opcode{JGE_rel8}
var JL = []*Opcode{JL_rel8}
var JLE = []*Opcode{JLE_rel8}
var JNA = []*Opcode{JNA_rel8}
var JNAE = []*Opcode{JNAE_rel8}
var JNB = []*Opcode{JNB_rel8}
var JNBE = []*Opcode{JNBE_rel8}
var JNE = []*Opcode{JNE_rel8}
var JNG = []*Opcode{JNG_rel8}
var JNGE = []*Opcode{JNGE_rel8}
var JNL = []*Opcode{JNL_rel8}
var JNLE = []*Opcode{JNLE_rel8}
var LEA = []*Opcode{LEA_r64_m}
var MOV = []*Opcode{
	MOV_r8_imm8_no_rex,
	MOV_rm8_r8, MOV_r8_rm8, MOV_r8_imm8,
	MOV_rm16_r16, MOV_r16_rm16,
	MOV_r16_imm16,
	MOV_r32_imm32,
	MOV_rm32_r32, MOV_r32_rm32,
	MOV_rm64_r64, MOV_r64_rm64,
	MOV_r64_imm64, MOV_rm64_imm32,
	MOVQ_xmm_rm64, MOVSD_xmm1m64_xmm2,
}
var MOVSX = []*Opcode{
	MOVSX_r16_rm8,
	MOVSX_r32_rm8,
	MOVSX_r32_rm16,
	MOVSX_r64_rm8,
	MOVSX_r64_rm16,
	MOVSX_r64_rm32,
}
var MOVZX = []*Opcode{
	MOVZX_r16_rm8,
	MOVZX_r32_rm8,
	MOVZX_r64_rm8,
	MOVZX_r32_rm16,
	MOVZX_r64_rm16,
}

var OR = []*Opcode{
	OR_r8_rm8,
	OR_r8_rm8_no_rex,
	OR_rm8_r8,
	OR_rm8_r8_no_rex,
	OR_r16_rm16,
	OR_rm16_r16,
	OR_r32_rm32,
	OR_rm32_r32,
	OR_r64_rm64,
	OR_rm64_r64,
}
var POP = []*Opcode{POP_r64}
var PUSH = []*Opcode{PUSH_imm32, PUSH_r64}
var SETA = []*Opcode{SETA_rm8}
var SETAE = []*Opcode{SETAE_rm8}
var SETB = []*Opcode{SETB_rm8}
var SETBE = []*Opcode{SETBE_rm8}
var SETC = []*Opcode{SETC_rm8}
var SETE = []*Opcode{SETE_rm8}
var SETNE = []*Opcode{SETNE_rm8}

var SHL = []*Opcode{
	SHL_rm8_imm8,
	SHL_rm8_imm8_no_rex,
	SHL_rm16_imm8,
	SHL_rm32_imm8,
	SHL_rm64_imm8,
}
var SHR = []*Opcode{
	SHR_rm8_imm8,
	SHR_rm8_imm8_no_rex,
	SHR_rm16_imm8,
	SHR_rm32_imm8,
	SHR_rm64_imm8,
}
var SUB = []*Opcode{
	SUB_rm8_imm8, SUB_rm64_imm8,
	SUB_r8_rm8,
	SUB_rm16_r16, SUB_r16_rm16,
	SUB_rm32_r32, SUB_r32_rm32,
	SUB_rm64_r64, SUB_r64_rm64, SUB_rm64_imm32,
	SUBSD_xmm1_xmm2m64,
}

var XOR = []*Opcode{
	XOR_r8_rm8,
	XOR_r8_rm8_no_rex,
	XOR_rm8_r8,
	XOR_rm8_r8_no_rex,
	XOR_r16_rm16,
	XOR_rm16_r16,
	XOR_r32_rm32,
	XOR_rm32_r32,
	XOR_rm64_imm32,
	XOR_r64_rm64, XOR_rm64_r64}
