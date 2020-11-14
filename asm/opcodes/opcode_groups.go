package opcodes

import (
	. "github.com/bspaans/jit-compiler/asm/encoding"
)

var ADD = []*Opcode{
	ADD_rm8_r8,
	ADD_r8_rm8,
	ADD_rm32_r32,
	ADD_r32_rm32,
	ADD_rm64_r64,
	ADD_r64_rm64,
	ADD_rm64_imm32,
	ADDSD_xmm1_xmm2m64,
}
var CALL = []*Opcode{CALL_rm64}
var CMP = []*Opcode{
	CMP_rm8_imm8,
	CMP_rm8_imm8_no_rex,
	CMP_rm64_imm32,
	CMP_rm64_r64,
}
var CVTSI2SD = []*Opcode{CVTSI2SD_xmm1_rm64}
var CVTSD2SI = []*Opcode{CVTSD2SI_r64_xmm1m64}
var CVTTSD2SI = []*Opcode{CVTTSD2SI_r64_xmm1m64}
var DEC = []*Opcode{DEC_rm64}
var IDIV = []*Opcode{DIVSD_xmm1_xmm2m64}
var DIV = []*Opcode{DIV_rm8, DIV_rm16, DIV_rm32, DIV_rm64}
var IMUL = []*Opcode{MULSD_xmm1_xmm2m64, IMUL_r64_rm64}
var MUL = []*Opcode{MUL_rm8, MUL_rm16, MUL_rm32, MUL_rm64}
var INC = []*Opcode{INC_rm64}
var JMP = []*Opcode{JMP_rel8, JMP_rel32, JMP_rm64}
var JE = []*Opcode{JE_rel8}
var JNE = []*Opcode{JNE_rel8}
var LEA = []*Opcode{LEA_r64_m}
var MOV = []*Opcode{
	MOV_r8_imm8_no_rex,
	MOV_rm8_r8, MOV_r8_rm8, MOV_r8_imm8,
	MOV_rm16_r16, MOV_r16_rm16,
	MOV_r32_imm32,
	MOV_rm32_r32, MOV_r32_rm32,
	MOV_rm64_r64, MOV_r64_rm64,
	MOV_r64_imm64, MOV_rm64_imm32,
	MOVQ_xmm_rm64, MOVSD_xmm1m64_xmm2,
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
var SUB = []*Opcode{
	SUB_rm8_imm8, SUB_rm64_imm8,
	SUB_r8_rm8,
	SUB_rm32_r32, SUB_r32_rm32,
	SUB_rm64_r64, SUB_r64_rm64, SUB_rm64_imm32,
	SUBSD_xmm1_xmm2m64}
var XOR = []*Opcode{XOR_rm64_imm32, XOR_r64_rm64}
