package opcodes

import (
	. "github.com/bspaans/jit-compiler/asm/aarch64/encoding"
)

var (
	ADD_Wd_Wn_imm12 = &Opcode{"add", []OpcodeChunk{OP_Exact(10, 0b000_100010_0), OP_Imm12, OP_Wn, OP_Wd}}
	ADD_Xd_Xn_imm12 = &Opcode{"add", []OpcodeChunk{OP_Exact(10, 0b100_100010_0), OP_Imm12, OP_Xn, OP_Xd}}
	ADD_Wd_Wn_Wm    = &Opcode{"add", []OpcodeChunk{OP_Exact(11, 0b000_01011_000), OP_Wm, OP_Exact(6, 0), OP_Wn, OP_Wd}}
	ADD_Xd_Xn_Xm    = &Opcode{"add", []OpcodeChunk{OP_Exact(11, 0b100_01011_000), OP_Xm, OP_Exact(6, 0), OP_Xn, OP_Xd}}

	ADDS_Wd_Wn_imm12 = &Opcode{"adds", []OpcodeChunk{OP_Exact(10, 0b001_100010_0), OP_Imm12, OP_Wn, OP_Wd}}
	ADDS_Xd_Xn_imm12 = &Opcode{"adds", []OpcodeChunk{OP_Exact(10, 0b101_100010_0), OP_Imm12, OP_Xn, OP_Xd}}

	MOVK_Wd_imm16 = &Opcode{"movk", []OpcodeChunk{OP_Exact(11, 0b011_100101_00), OP_Imm16, OP_Wd}}
	MOVK_Xd_imm16 = &Opcode{"movk", []OpcodeChunk{OP_Exact(11, 0b111_100101_00), OP_Imm16, OP_Xd}}

	SUB_Wd_Wn_imm12 = &Opcode{"sub", []OpcodeChunk{OP_Exact(10, 0b010_100010_0), OP_Imm12, OP_Wn, OP_Wd}}
	SUB_Xd_Xn_imm12 = &Opcode{"sub", []OpcodeChunk{OP_Exact(10, 0b110_100010_0), OP_Imm12, OP_Xn, OP_Xd}}
	SUB_Wd_Wn_Wm    = &Opcode{"sub", []OpcodeChunk{OP_Exact(10, 0b010_01011_00_0), OP_Wm, OP_Exact(6, 0), OP_Wn, OP_Wd}}
	SUB_Xd_Xn_Xm    = &Opcode{"sub", []OpcodeChunk{OP_Exact(10, 0b110_01011_00_0), OP_Xm, OP_Exact(6, 0), OP_Xn, OP_Xd}}

	SUBS_Wd_Wn_imm12 = &Opcode{"subs", []OpcodeChunk{OP_Exact(10, 0b011_100010_0), OP_Imm12, OP_Wn, OP_Wd}}
	SUBS_Xd_Xn_imm12 = &Opcode{"subs", []OpcodeChunk{OP_Exact(10, 0b111_100010_0), OP_Imm12, OP_Xn, OP_Xd}}
)
