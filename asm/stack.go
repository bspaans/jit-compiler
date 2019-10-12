package asm

import (
	"errors"

	"github.com/bspaans/jit/asm/encoding"
	"github.com/bspaans/jit/lib"
)

type PUSH struct {
	Source encoding.Operand
}

func (i *PUSH) Encode() (lib.MachineCode, error) {
	if i.Source == nil {
		return nil, errors.New("Missing source")
	}
	if i.Source.Type() == encoding.T_Register {
		src := i.Source.(*encoding.Register)
		if src.Size == lib.QUADWORD {
			return encoding.PUSH_r64.Encode([]encoding.Operand{src})
		}
	} else if i.Source.Type() == encoding.T_Uint32 {
		return encoding.PUSH_imm32.Encode([]encoding.Operand{i.Source})
	}
	return nil, errors.New("Unsupported push operation")
}

func (i *PUSH) String() string {
	return "push " + i.Source.String()
}

type POP struct {
	Source encoding.Operand
}

func (i *POP) Encode() (lib.MachineCode, error) {
	if i.Source == nil {
		return nil, errors.New("Missing source")
	}
	if i.Source.Type() == encoding.T_Register {
		src := i.Source.(*encoding.Register)
		if src.Size == lib.QUADWORD {
			return encoding.POP_r64.Encode([]encoding.Operand{src})
		}
	}
	return nil, errors.New("Unsupported pop operation")
}

func (i *POP) String() string {
	return "pop " + i.Source.String()
}

type PUSHFQ struct {
}

func (i *PUSHFQ) Encode() (lib.MachineCode, error) {
	return []uint8{0x9C}, nil
}

func (i *PUSHFQ) String() string {
	return "pushfq"
}
