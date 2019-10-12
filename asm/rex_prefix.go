package asm

import "github.com/bspaans/jit/asm/encoding"

func REXEncode(opR, opB *encoding.Register) uint8 {
	rexR := opR != nil && opR.Register > 7
	rexB := opB != nil && opB.Register > 7
	return encoding.NewREXPrefix(true, rexR, false, rexB).Encode()
}

func EncodeOpcodeWithREX(opcode uint8, op1, op2 *encoding.Register) []uint8 {
	rex := REXEncode(op2, op1)
	return []uint8{rex, opcode}
}

func EncodeOpcodeWithREXAndModRM(opcode uint8, op1, op2 *encoding.Register, reg uint8) []uint8 {
	rex := REXEncode(op2, op1)
	if op2 != nil {
		reg = op2.Encode()
	}
	modrm := encoding.NewModRM(encoding.DirectRegisterMode, op1.Encode(), reg).Encode()
	return []uint8{rex, opcode, modrm}
}

func EncodeOpcodeWithREXAndModRMAndImm(opcode uint8, op1, op2 *encoding.Register, reg uint8, value encoding.Value) []uint8 {
	result := EncodeOpcodeWithREXAndModRM(0x81, op1, op2, reg)
	for _, b := range value.Encode() {
		result = append(result, b)
	}
	return result
}
