package asm

import (
	"errors"

	"github.com/bspaans/jit/asm/encoding"
	"github.com/bspaans/jit/lib"
)

type INC struct {
	Register *encoding.Register
}

func (i *INC) Encode() (lib.MachineCode, error) {
	if i.Register == nil {
		return nil, errors.New("Missing register")
	}
	if i.Register.Size == lib.QUADWORD {
		return encoding.INC_rm64.Encode([]encoding.Operand{i.Register})
	}
	return nil, errors.New("Unsupported register size")
}

func (i *INC) String() string {
	return "inc " + i.Register.String()
}

type DEC struct {
	Register *encoding.Register
}

func (i *DEC) Encode() (lib.MachineCode, error) {
	if i.Register == nil {
		return nil, errors.New("Missing register")
	}
	if i.Register.Size == lib.QUADWORD {
		return encoding.DEC_rm64.Encode([]encoding.Operand{i.Register})
	}
	return nil, errors.New("Unsupported register size")
}

func (i *DEC) String() string {
	return "dec " + i.Register.String()
}
