package asm

import "errors"

type ADD struct {
	Source Operand
	Dest   Operand
}

func (i *ADD) Encode() (MachineCode, error) {
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
			if src.Size == QUADWORD && dest.Size == QUADWORD {
				return EncodeOpcodeWithREXAndModRM(0x03, src, dest, 0), nil
			} else if src.Size == QUADDOUBLE && dest.Size == QUADDOUBLE {
				// addsd
				result := []uint8{0xf2, 0x0f, 0x58}
				modrm := NewModRM(DirectRegisterMode, src.Encode(), dest.Encode())
				result = append(result, modrm.Encode())
				return result, nil
			}
		}
	} else if i.Source.Type() == T_Uint32 {
		src := i.Source.(Uint32)
		if i.Dest.Type() == T_Register {
			dest := i.Dest.(*Register)
			return EncodeOpcodeWithREXAndModRMAndImm(0x81, dest, nil, 0, src), nil
		}
	}
	return nil, errors.New("Unsupported add operation: " + i.String())
}

func (i *ADD) String() string {
	opcode := "add"
	if i.Source.Type() == T_Register {
		src := i.Source.(*Register)
		if i.Dest.Type() == T_Register {
			dest := i.Dest.(*Register)
			if src.Size == QUADDOUBLE && dest.Size == QUADDOUBLE {
				opcode = "addsd"
			}
		}
	}
	return opcode + " " + i.Source.String() + ", " + i.Dest.String()
}

func (i *ADD) TwoOperands() (Operand, Operand) {
	return i.Source, i.Dest
}
