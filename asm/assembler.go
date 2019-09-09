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
	String() string
}

type INC struct {
	Register *Register
}

func (i *INC) Encode() (MachineCode, error) {
	if i.Register == nil {
		return nil, errors.New("Missing register")
	}
	if i.Register.Size == QUADWORD {
		rexB := i.Register.Register > 7
		rex := NewREXPrefix(true, false, false, rexB).Encode()
		modrm := NewModRM(DirectRegisterMode, i.Register.Encode(), 0).Encode()
		return []uint8{rex, 0xff, modrm}, nil
	}
	return nil, errors.New("Unsupported register size")
}

func (i *INC) String() string {
	return "inc " + i.Register.String()
}

type DEC struct {
	Register *Register
}

func (i *DEC) Encode() (MachineCode, error) {
	if i.Register == nil {
		return nil, errors.New("Missing register")
	}
	if i.Register.Size == QUADWORD {
		rexB := i.Register.Register > 7
		rex := NewREXPrefix(true, false, false, rexB).Encode()
		modrm := NewModRM(DirectRegisterMode, i.Register.Encode(), 1).Encode()
		return []uint8{rex, 0xff, modrm}, nil
	}
	return nil, errors.New("Unsupported register size")
}

func (i *DEC) String() string {
	return "dec " + i.Register.String()
}

type MOV struct {
	Dest   *Register
	Source *Register
}

func (i *MOV) Encode() (MachineCode, error) {
	if i.Dest == nil {
		return nil, errors.New("Missing register")
	}
	if i.Source != nil {
		if i.Source.Size == QUADWORD && i.Dest.Size == QUADWORD {
			rexB := i.Source.Register > 7
			rexR := i.Dest.Register > 7
			rex := NewREXPrefix(true, rexR, false, rexB).Encode()
			modrm := NewModRM(DirectRegisterMode, i.Source.Encode(), i.Dest.Encode()).Encode()
			return []uint8{rex, 0x89, modrm}, nil
		}
		return nil, errors.New("Unsupported register size")
	}
	return nil, errors.New("Unsupported mov operation")
}
func (i *MOV) String() string {
	return "mov " + i.Source.String() + " " + i.Dest.String()
}

func CompileInstruction(instr []Instruction) (MachineCode, error) {
	result := []uint8{}
	for _, i := range instr {
		b, err := i.Encode()
		if err != nil {
			return nil, err
		}
		fmt.Println(i)
		fmt.Println(MachineCode(b))
		for _, code := range b {
			result = append(result, code)
		}
	}
	return result, nil
}

func init() {
	b, err := CompileInstruction([]Instruction{
		&INC{rax},
		&INC{rcx},
		&INC{r14},
		&DEC{rax},
		&DEC{r14},
		&MOV{rax, rax},
		&MOV{rax, rcx},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(b)
}
