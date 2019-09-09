package asm

import (
	"encoding/binary"
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

type Type uint8

const (
	T_Register Type = 0
	T_Uint64   Type = 1
)

type Value interface {
	Type() Type
	String() string
}

type MOV struct {
	Source Value
	Dest   *Register
}

type Uint64 uint64

func (i Uint64) Type() Type {
	return T_Uint64
}
func (i Uint64) String() string {
	return fmt.Sprintf("%d", i)
}
func (i Uint64) Encode() []uint8 {
	result := make([]byte, 8)
	binary.LittleEndian.PutUint64(result, uint64(i))
	return result
}

func (i *MOV) Encode() (MachineCode, error) {
	if i.Dest == nil {
		return nil, errors.New("Missing register")
	}
	if i.Source == nil {
		return nil, errors.New("Missing source")
	}
	if i.Source.Type() == T_Register {
		src := i.Source.(*Register)
		if src.Size == QUADWORD && i.Dest.Size == QUADWORD {
			rexB := i.Dest.Register > 7
			rexR := src.Register > 7
			rex := NewREXPrefix(true, rexR, false, rexB).Encode()
			modrm := NewModRM(DirectRegisterMode, i.Dest.Encode(), src.Encode()).Encode()
			return []uint8{rex, 0x89, modrm}, nil
		}
		return nil, errors.New("Unsupported register size")
	} else if i.Source.Type() == T_Uint64 {
		src := i.Source.(Uint64)
		rexB := i.Dest.Register > 7
		rex := NewREXPrefix(true, false, false, rexB).Encode()
		result := []uint8{rex, 0xB8 + (i.Dest.Encode() & 7)}
		for _, b := range src.Encode() {
			result = append(result, b)
		}
		return result, nil

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
		&MOV{Uint64(0), rcx},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(b)
}
