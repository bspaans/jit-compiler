package asm

import "errors"

type INC struct {
	Register *Register
}

func (i *INC) Encode() (MachineCode, error) {
	if i.Register == nil {
		return nil, errors.New("Missing register")
	}
	if i.Register.Size == QUADWORD {
		return EncodeOpcodeWithREXAndModRM(0xff, i.Register, nil, 0), nil
	}
	return nil, errors.New("Unsupported register size")
}

func (i *INC) String() string {
	return "inc " + i.Register.String()
}

type DEC struct {
	Register *Register
}

func (i *DEC) Encode() (MachineCode, error) {
	if i.Register == nil {
		return nil, errors.New("Missing register")
	}
	if i.Register.Size == QUADWORD {
		return EncodeOpcodeWithREXAndModRM(0xff, i.Register, nil, 1), nil
	}
	return nil, errors.New("Unsupported register size")
}

func (i *DEC) String() string {
	return "dec " + i.Register.String()
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
