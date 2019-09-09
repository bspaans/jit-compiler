package asm

import "errors"

type MOV struct {
	Source Value
	Dest   *Register
}

func (i *MOV) Encode() (MachineCode, error) {
	if i.Dest == nil {
		return nil, errors.New("Missing register")
	}
	if i.Source == nil {
		return nil, errors.New("Missing source")
	}
	if i.Source.Type() == T_Register {
		src := i.Source.(*Register)
		if src.Size == QUADWORD && i.Dest.Size == QUADWORD {
			rexB := i.Dest.Register > 7
			rexR := src.Register > 7
			rex := NewREXPrefix(true, rexR, false, rexB).Encode()
			modrm := NewModRM(DirectRegisterMode, i.Dest.Encode(), src.Encode()).Encode()
			return []uint8{rex, 0x89, modrm}, nil
		}
		return nil, errors.New("Unsupported register size")
	} else if i.Source.Type() == T_Uint64 {
		src := i.Source.(Uint64)
		rexB := i.Dest.Register > 7
		rex := NewREXPrefix(true, false, false, rexB).Encode()
		result := []uint8{rex, 0xB8 + (i.Dest.Encode() & 7)}
		for _, b := range src.Encode() {
			result = append(result, b)
		}
		return result, nil

	}
	return nil, errors.New("Unsupported mov operation")
}
func (i *MOV) String() string {
	return "mov " + i.Source.String() + " " + i.Dest.String()
}
