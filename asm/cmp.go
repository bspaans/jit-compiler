package asm

import "errors"

type CMP struct {
	Source Operand
	Dest   Operand
}

func (i *CMP) Encode() (MachineCode, error) {
	if i.Dest == nil {
		return nil, errors.New("Missing dest")
	}
	if i.Source == nil {
		return nil, errors.New("Missing source")
	}
	if i.Source.Type() == T_Uint64 {
		src := i.Source.(Uint64)
		if i.Dest.Type() == T_Register {
			dest := i.Dest.(*Register)
			rexB := dest.Register > 7
			rex := NewREXPrefix(true, false, false, rexB).Encode()
			modrm := NewModRM(DirectRegisterMode, dest.Encode(), 7).Encode()
			result := []uint8{rex, 0x81, modrm}
			// Can only cmp to a double
			for _, b := range src.Encode()[:4] {
				result = append(result, b)
			}
			return result, nil
		}
	}
	return nil, errors.New("Unsupported cmp operation")
}

func (i *CMP) String() string {
	return "cmp " + i.Source.String() + " " + i.Dest.String()
}
