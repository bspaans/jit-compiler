package asm

import "errors"

type MUL struct {
	Source Operand
	Dest   Operand
}

func (i *MUL) Encode() (MachineCode, error) {
	if i.Dest == nil {
		return nil, errors.New("Missing dest")
	}
	if i.Source == nil {
		return nil, errors.New("Missing source")
	}
	if i.Source.Type() == T_Register {
		src := i.Source.(*Register)
		if i.Dest.Type() == T_Register {
			dest := i.Dest.(*Register)
			// mulsd
			if src.Size == QUADDOUBLE && dest.Size == QUADDOUBLE {
				result := []uint8{0xf2, 0x0f, 0x59}
				modrm := NewModRM(DirectRegisterMode, src.Encode(), dest.Encode())
				result = append(result, modrm.Encode())
				return result, nil
			}
		}
	}
	return nil, errors.New("Unsupported mul operation: " + i.String())
}

func (i *MUL) String() string {
	return "mul " + i.Source.String() + ", " + i.Dest.String()
}
