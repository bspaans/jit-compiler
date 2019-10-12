package asm

import (
	"errors"

	"github.com/bspaans/jit/asm/encoding"
	"github.com/bspaans/jit/lib"
)

type DIV struct {
	Source encoding.Operand
	Dest   encoding.Operand
}

func (i *DIV) Encode() (lib.MachineCode, error) {
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
			// divsd
			if src.Size == lib.QUADDOUBLE && dest.Size == lib.QUADDOUBLE {
				return encoding.DIVSD_xmm1_xmm2m64.Encode([]encoding.Operand{dest, src})
			}
		}
	}
	return nil, errors.New("Unsupported divsd operation")
}

func (i *DIV) String() string {
	return "div " + i.Source.String() + ", " + i.Dest.String()
}
