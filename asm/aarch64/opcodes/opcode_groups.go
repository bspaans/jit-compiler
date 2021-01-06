package opcodes

import (
	. "github.com/bspaans/jit-compiler/asm/aarch64/encoding"
)

var ADD = []*Opcode{
	ADD_Wd_Wn_imm12,
	ADD_Xd_Xn_imm12,
	ADD_Wd_Wn_Wm,
	ADD_Xd_Xn_Xm,
}

var ADDS = []*Opcode{
	ADDS_Wd_Wn_imm12,
	ADDS_Xd_Xn_imm12,
}

var MOVK = []*Opcode{
	MOVK_Wd_imm16,
	MOVK_Xd_imm16,
}

var SUB = []*Opcode{
	SUB_Wd_Wn_imm12,
	SUB_Xd_Xn_imm12,
}

var SUBS = []*Opcode{
	SUBS_Wd_Wn_imm12,
	SUBS_Xd_Xn_imm12,
}
