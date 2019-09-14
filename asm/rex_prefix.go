package asm

type REXPrefix struct {
	// When true a 64-bit operand size is used. Otherwise the default size is used.
	W bool
	// This is an extension to the ModRM.Reg field.
	R bool
	// This is an extension to the SIB.index field.
	X bool
	// This is an extension to the ModRM.RM field or the SIB.base field.
	B bool
}

func NewREXPrefix(w, r, x, b bool) *REXPrefix {
	return &REXPrefix{w, r, x, b}
}

func (r *REXPrefix) Encode() uint8 {
	result := uint8(0)
	if r.B {
		result = 1
	}
	if r.X {
		result += 1 << 1
	}
	if r.R {
		result += 1 << 2
	}
	if r.W {
		result += 1 << 3
	}
	return result + (1 << 6)
}

func REXEncode(opR, opB *Register) uint8 {
	rexR := opR != nil && opR.Register > 7
	rexB := opB != nil && opB.Register > 7
	return NewREXPrefix(true, rexR, false, rexB).Encode()
}

func EncodeOpcodeWithREX(opcode uint8, op1, op2 *Register) []uint8 {
	rex := REXEncode(op2, op1)
	return []uint8{rex, opcode}
}

func EncodeOpcodeWithREXAndModRM(opcode uint8, op1, op2 *Register, reg uint8) []uint8 {
	rex := REXEncode(op2, op1)
	if op2 != nil {
		reg = op2.Encode()
	}
	modrm := NewModRM(DirectRegisterMode, op1.Encode(), reg).Encode()
	return []uint8{rex, opcode, modrm}
}

func EncodeOpcodeWithREXAndModRMAndImm(opcode uint8, op1, op2 *Register, reg uint8, value Value) []uint8 {
	result := EncodeOpcodeWithREXAndModRM(0x81, op1, op2, reg)
	for _, b := range value.Encode() {
		result = append(result, b)
	}
	return result
}
