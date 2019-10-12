package asm

import (
	"errors"

	"github.com/bspaans/jit/asm/encoding"
	"github.com/bspaans/jit/lib"
)

type JNE struct {
	Dest encoding.Value
}

func (j *JNE) Encode() (lib.MachineCode, error) {
	if j.Dest == nil {
		return nil, errors.New("Missing destination")
	}
	var result []uint8
	if j.Dest.Type() == encoding.T_Uint8 {
		return encoding.JNE_rel8.Encode([]encoding.Operand{j.Dest})
	} else {
		return nil, errors.New("Unsupported destination")
	}
	return result, nil

}
func (j *JNE) String() string {
	return "jne " + j.Dest.String()
}

type JMP struct {
	Dest encoding.Value
}

func (j *JMP) Encode() (lib.MachineCode, error) {
	if j.Dest == nil {
		return nil, errors.New("Missing destination")
	}
	var result []uint8
	if j.Dest.Type() == encoding.T_Uint8 {
		return encoding.JMP_rel8.Encode([]encoding.Operand{j.Dest})
	} else if j.Dest.Type() == encoding.T_Uint32 {
		return encoding.JMP_rel32.Encode([]encoding.Operand{j.Dest})
	} else {
		return nil, errors.New("Unsupported destination")
	}
	return result, nil

}
func (j *JMP) String() string {
	return "jmp " + j.Dest.String()
}
