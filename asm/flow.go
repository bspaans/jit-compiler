package asm

import "errors"

type RET struct {
}

func (i *RET) Encode() (MachineCode, error) {
	return []uint8{0xc3}, nil
}

func (i *RET) String() string {
	return "ret"
}

type SYSCALL struct {
}

func (i *SYSCALL) Encode() (MachineCode, error) {
	return []uint8{0x0f, 0x05}, nil
}

func (i *SYSCALL) String() string {
	return "syscall"
}

type CALL struct {
	Dest Operand
}

func (i *CALL) Encode() (MachineCode, error) {
	if i.Dest.Type() == T_Register {
		dest := i.Dest.(*Register)
		result := []uint8{0xff}
		modrm := NewModRM(DirectRegisterMode, dest.Encode(), 2)
		result = append(result, modrm.Encode())
		return result, nil
	}
	return nil, errors.New("Unsupported call operation: " + i.String())
}

func (i *CALL) String() string {
	return "call " + i.Dest.String()
}
