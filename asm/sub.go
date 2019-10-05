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
	if i.Source.Type() == T_Uint32 {
		src := i.Source.(Uint32)
		if i.Dest.Type() == T_Register {
			dest := i.Dest.(*Register)
			return EncodeOpcodeWithREXAndModRMAndImm(0x81, dest, nil, 5, src), nil
		}
	}
	return nil, errors.New("Unsupported add operation")
}

func (i *SUB) String() string {
	return "sub " + i.Source.String() + ", " + i.Dest.String()
}
