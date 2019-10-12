package asm

import (
	"errors"

	"github.com/bspaans/jit/asm/encoding"
	"github.com/bspaans/jit/lib"
)

type SETE struct {
	Dest encoding.Operand
}

func (i *SETE) Encode() (lib.MachineCode, error) {
	if i.Dest == nil {
		return nil, errors.New("Missing dest")
	}
	if i.Dest.Type() == encoding.T_Register {
		dest := i.Dest.(*encoding.Register)
		if dest.Size == lib.BYTE {
			return encoding.SETE_rm8.Encode([]encoding.Operand{dest})
		}
	}
	return nil, errors.New("Unsupported sete operation")
}

func (i *SETE) String() string {
	return "sete " + i.Dest.String()
}
