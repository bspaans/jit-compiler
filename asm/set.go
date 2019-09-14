package asm

import "errors"

type SETE struct {
	Dest Operand
}

func (i *SETE) Encode() (MachineCode, error) {
	if i.Dest == nil {
		return nil, errors.New("Missing dest")
	}
	if i.Dest.Type() == T_Register {
		dest := i.Dest.(*Register)
		if dest.Size == BYTE {
			modrm := NewModRM(DirectRegisterMode, dest.Encode(), 0).Encode()
			return []uint8{0x0f, 0x94, modrm}, nil
		}
	}
	return nil, errors.New("Unsupported sete operation")
}

func (i *SETE) String() string {
	return "sete " + i.Dest.String()
}
