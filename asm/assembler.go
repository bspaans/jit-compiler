package asm

import (
	"encoding/hex"
	"errors"
	"fmt"
)

type Size uint8

const (
	BYTE     Size = 1
	WORD     Size = 2
	DOUBLE   Size = 4
	QUADWORD Size = 8
)

type MachineCode []uint8

func (m MachineCode) String() string {
	return hex.EncodeToString(m)
}

func EncodeModRM(mod, reg, rm uint8) uint8 {
	return rm
}

type Instruction interface {
	Encode() (MachineCode, error)
}

type INC struct {
	Register *Register
}

func (i *INC) Encode() (MachineCode, error) {
	if i.Register == nil {
		return nil, errors.New("Missing register")
	}
	if i.Register.Size == QUADWORD {
		if i.Register.Register > 7 {
			return nil, fmt.Errorf("Unsupported register %s", i.Register.Name)
		}
		rex := NewREXPrefix(true, false, false, false).Encode()
		modrm := NewModRM(DirectRegisterMode, i.Register.Encode(), 0).Encode()
		return []uint8{rex, 0xff, modrm}, nil
	}
	return nil, errors.New("Unsupported register size")
}

func CompileInstruction(instr []Instruction) (MachineCode, error) {
	result := []uint8{}
	for _, i := range instr {
		b, err := i.Encode()
		if err != nil {
			return nil, err
		}
		for _, code := range b {
			result = append(result, code)
		}
	}
	fmt.Println(result)
	return result, nil
}

func init() {
	b, err := CompileInstruction([]Instruction{
		&INC{rax},
		&INC{rcx},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(b)
}
