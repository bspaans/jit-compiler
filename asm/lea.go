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
		if i.Dest.Type() == encoding.T_Register {
			return encoding.LEA_r64_m.Encode([]encoding.Operand{i.Dest, i.Source})
		}
	}
	return nil, errors.New("Unsupported lea operation: " + i.String())
}

func (i *LEA) String() string {
	return "lea " + i.Source.String() + ", " + i.Dest.String()
}
