package asm

import "errors"

type PUSH struct {
	Source Operand
}

func (i *PUSH) Encode() (MachineCode, error) {
	if i.Source == nil {
		return nil, errors.New("Missing source")
	}
	if i.Source.Type() == T_Register {
		src := i.Source.(*Register)
		if src.Size == QUADWORD {
			return []uint8{0x50 + src.Register}, nil
		}
	}
	return nil, errors.New("Unsupported push operation")
}

func (i *PUSH) String() string {
	return "push " + i.Source.String()
}

type POP struct {
	Source Operand
}

func (i *POP) Encode() (MachineCode, error) {
	if i.Source == nil {
		return nil, errors.New("Missing source")
	}
	if i.Source.Type() == T_Register {
		src := i.Source.(*Register)
		if src.Size == QUADWORD {
			return []uint8{0x58 + src.Register}, nil
		}
	}
	return nil, errors.New("Unsupported pop operation")
}

func (i *POP) String() string {
	return "pop " + i.Source.String()
}

type PUSHFQ struct {
}

func (i *PUSHFQ) Encode() (MachineCode, error) {
	return []uint8{0x9C}, nil
}

func (i *PUSHFQ) String() string {
	return "pushfq"
}
