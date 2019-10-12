package asm

import (
	"errors"

	"github.com/bspaans/jit/asm/encoding"
	"github.com/bspaans/jit/lib"
)

type LEA struct {
	Source encoding.Operand
	Dest   encoding.Operand
}

func (i *LEA) Encode() (lib.MachineCode, error) {
	if i.Dest == nil {
		return nil, errors.New("Missing dest")
	}
	if i.Source == nil {
		return nil, errors.New("Missing source")
	}
	if i.Source.Type() == encoding.T_RIPRelative {
		src := i.Source.(*encoding.RIPRelative)
		if i.Dest.Type() == encoding.T_Register {
			dest := i.Dest.(*encoding.Register)
			rex := REXEncode(nil, dest)
			modrm := encoding.NewModRM(encoding.IndirectRegisterMode, 5, dest.Encode()).Encode()
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
