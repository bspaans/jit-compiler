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
			}
		}
	} else if i.Source.Type() == T_Uint32 {
		src := i.Source.(Uint32)
		if i.Dest.Type() == T_Register {
			dest := i.Dest.(*Register)
			return EncodeOpcodeWithREXAndModRMAndImm(0x81, dest, nil, 0, src), nil
		}
	}
	return nil, errors.New("Unsupported sub operation")
}

func (i *ADD) String() string {
	return "add " + i.Source.String() + ", " + i.Dest.String()
}
