package asm

import (
	"errors"

	"github.com/bspaans/jit/asm/encoding"
	"github.com/bspaans/jit/lib"
)

type ADD struct {
	Source encoding.Operand
	Dest   encoding.Operand
}

func (i *ADD) Encode() (lib.MachineCode, error) {
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
				return encoding.ADD_r64_rm64.Encode([]encoding.Operand{dest, src})
			} else if src.Size == lib.QUADDOUBLE && dest.Size == lib.QUADDOUBLE {
				return encoding.ADDSD_xmm1_xmm2m64.Encode([]encoding.Operand{dest, src})
			}
		}
	} else if i.Source.Type() == encoding.T_Uint32 {
		src := i.Source.(encoding.Uint32)
		if i.Dest.Type() == encoding.T_Register {
			dest := i.Dest.(*encoding.Register)
			return encoding.ADD_rm64_imm32.Encode([]encoding.Operand{dest, src})
		}
	}
	return nil, errors.New("Unsupported add operation: " + i.String())
}

func (i *ADD) String() string {
	opcode := "add"
	if i.Source.Type() == encoding.T_Register {
		src := i.Source.(*encoding.Register)
		if i.Dest.Type() == encoding.T_Register {
			dest := i.Dest.(*encoding.Register)
			if src.Size == lib.QUADDOUBLE && dest.Size == lib.QUADDOUBLE {
				opcode = "addsd"
			}
		}
	}
	return opcode + " " + i.Source.String() + ", " + i.Dest.String()
}
