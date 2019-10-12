package asm

import (
	"errors"

	"github.com/bspaans/jit/asm/encoding"
	"github.com/bspaans/jit/lib"
)

type MUL struct {
	Source encoding.Operand
	Dest   encoding.Operand
}

func (i *MUL) Encode() (lib.MachineCode, error) {
	if i.Dest == nil {
		return nil, errors.New("Missing dest")
	}
	if i.Source == nil {
		return nil, errors.New("Missing source")
	}
	if i.Source.Type() == encoding.T_Register {
		src := i.Source.(*encoding.Register)
		if i.Dest.Type() == encoding.T_Register {
			dest := i.Dest.(*encoding.Register)
			// mulsd
			if src.Size == lib.QUADDOUBLE && dest.Size == lib.QUADDOUBLE {
				return encoding.MULSD_xmm1_xmm2m64.Encode([]encoding.Operand{dest, src})
			}
		}
	}
	return nil, errors.New("Unsupported mul operation: " + i.String())
}

func (i *MUL) String() string {
	return "mul " + i.Source.String() + ", " + i.Dest.String()
}
