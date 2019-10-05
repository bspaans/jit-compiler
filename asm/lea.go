package asm

import "errors"

type LEA struct {
	Source Operand
	Dest   Operand
}

func (i *LEA) Encode() (MachineCode, error) {
	if i.Dest == nil {
		return nil, errors.New("Missing dest")
	}
	if i.Source == nil {
		return nil, errors.New("Missing source")
	}
	if i.Source.Type() == T_RIPRelative {
		src := i.Source.(*RIPRelative)
		if i.Dest.Type() == T_Register {
			dest := i.Dest.(*Register)
			rex := REXEncode(nil, dest)
			modrm := NewModRM(IndirectRegisterMode, 5, dest.Encode()).Encode()
			result := []uint8{rex, 0x8d, modrm}
			for _, c := range src.Displacement.Encode() {
				result = append(result, c)
			}
			return result, nil
		}
	}
	return nil, errors.New("Unsupported lea operation: " + i.String())
}

func (i *LEA) String() string {
	return "lea " + i.Source.String() + ", " + i.Dest.String()
}
