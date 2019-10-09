package asm

import "errors"

type DIV struct {
	Source Operand
	Dest   Operand
}

func (i *DIV) Encode() (MachineCode, error) {
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
			// divsd
			if src.Size == QUADDOUBLE && dest.Size == QUADDOUBLE {
				result := []uint8{0xf2, 0x0f, 0x5e}
				modrm := NewModRM(DirectRegisterMode, src.Encode(), dest.Encode())
				result = append(result, modrm.Encode())
				return result, nil
			}
		}
	}
	return nil, errors.New("Unsupported divsd operation")
}

func (i *DIV) String() string {
	return "div " + i.Source.String() + ", " + i.Dest.String()
}
