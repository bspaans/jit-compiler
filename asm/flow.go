package asm

import (
	"errors"

	"github.com/bspaans/jit/asm/encoding"
	"github.com/bspaans/jit/lib"
)

type RET struct {
}

func (i *RET) Encode() (lib.MachineCode, error) {
	return encoding.RETURN.Encode([]encoding.Operand{})
}

func (i *RET) String() string {
	return "ret"
}

type SYSCALL struct {
}

func (i *SYSCALL) Encode() (lib.MachineCode, error) {
	return encoding.SYSCALL.Encode([]encoding.Operand{})
}

func (i *SYSCALL) String() string {
	return "syscall"
}

type CALL struct {
	Dest encoding.Operand
}

func (i *CALL) Encode() (lib.MachineCode, error) {
	if i.Dest.Type() == encoding.T_Register || i.Dest.Type() == encoding.T_DisplacedRegister {
		return encoding.CALL_rm64.Encode([]encoding.Operand{i.Dest})
	}
	return nil, errors.New("Unsupported call operation: " + i.String())
}

func (i *CALL) String() string {
	return "call " + i.Dest.String()
}
