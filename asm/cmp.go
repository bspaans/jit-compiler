package asm

import (
	"errors"

	"github.com/bspaans/jit/asm/encoding"
	"github.com/bspaans/jit/lib"
)

type CMP struct {
	Source encoding.Operand
	Dest   encoding.Operand
}

func (i *CMP) Encode() (lib.MachineCode, error) {
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
			if src.Size == lib.QUADWORD && dest.Size == lib.QUADWORD {
				return encoding.CMP_rm64_r64.Encode([]encoding.Operand{dest, src})
			}
		}

	} else if i.Source.Type() == encoding.T_Uint32 {
		src := i.Source.(encoding.Uint32)
		if i.Dest.Type() == encoding.T_Register {
			return encoding.CMP_rm64_imm32.Encode([]encoding.Operand{i.Dest, src})
		}
	}
	return nil, errors.New("Unsupported cmp operation: " + i.String())
}

func (i *CMP) String() string {
	return "cmp " + i.Source.String() + " " + i.Dest.String()
}
