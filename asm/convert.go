package asm

import (
	"errors"

	"github.com/bspaans/jit/asm/encoding"
	"github.com/bspaans/jit/lib"
)

// Convert signed integer to scalar double-precision floating point (float64)
type CVTSI2SD struct {
	Source encoding.Operand
	Dest   encoding.Operand
}

func (i *CVTSI2SD) Encode() (lib.MachineCode, error) {
	if i.Source.Type() == encoding.T_Register {
		src := i.Source.(*encoding.Register)
		if i.Dest.Type() == encoding.T_Register {
			dest := i.Dest.(*encoding.Register)
			if dest.Size == lib.QUADDOUBLE && src.Size == lib.QUADWORD {
				return encoding.CVTSI2SD_xmm1_rm64.Encode([]encoding.Operand{dest, src})
			}
		}
	}
	return nil, errors.New("Unsupported cvtsi2sd operation: " + i.String())
}
func (j *CVTSI2SD) String() string {
	return "cvtsi2sd " + j.Source.String() + " " + j.Dest.String()
}

// Convert double precision float to signed integer
type CVTTSD2SI struct {
	Source encoding.Operand
	Dest   encoding.Operand
}

func (i *CVTTSD2SI) Encode() (lib.MachineCode, error) {
	if i.Source.Type() == encoding.T_Register {
		src := i.Source.(*encoding.Register)
		if i.Dest.Type() == encoding.T_Register {
			dest := i.Dest.(*encoding.Register)
			if src.Size == lib.QUADDOUBLE && dest.Size == lib.QUADWORD {
				return encoding.CVTTSD2SI_r64_xmm1m64.Encode([]encoding.Operand{dest, src})
			}
		}
	}
	return nil, errors.New("Unsupported cvttsd2si operation: " + i.String())
}
func (j *CVTTSD2SI) String() string {
	return "cvttsd2si " + j.Source.String() + " " + j.Dest.String()
}
