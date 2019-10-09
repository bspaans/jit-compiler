package asm

import "errors"

type SUB struct {
	Source Operand
	Dest   Operand
}

func (i *SUB) Encode() (MachineCode, error) {
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
			// subsd
			if src.Size == QUADDOUBLE && dest.Size == QUADDOUBLE {
				result := []uint8{0xf2, 0x0f, 0x5c}
				modrm := NewModRM(DirectRegisterMode, src.Encode(), dest.Encode())
				result = append(result, modrm.Encode())
				return result, nil
			}
		}
	}
	if i.Source.Type() == T_Uint32 {
		src := i.Source.(Uint32)
		if i.Dest.Type() == T_Register {
			dest := i.Dest.(*Register)
			return EncodeOpcodeWithREXAndModRMAndImm(0x81, dest, nil, 5, src), nil
		}
	}
	return nil, errors.New("Unsupported sub operation")
}

func (i *SUB) String() string {
	cmd := "sub"
	if i.Source.Type() == T_Register {
		src := i.Source.(*Register)
		if i.Dest.Type() == T_Register {
			dest := i.Dest.(*Register)
			if src.Size == QUADDOUBLE && dest.Size == QUADDOUBLE {
				cmd = "subsd"
			}
		}
	}
	return cmd + " " + i.Source.String() + ", " + i.Dest.String()
}
