package lib

import (
	"fmt"
	"strings"
)

type Instruction interface {
	Encode() (MachineCode, error)
	String() string
}

type Instructions []Instruction

func (i Instructions) Encode() (MachineCode, error) {
	result := []uint8{}
	for _, instr := range i {
		b, err := instr.Encode()
		if err != nil {
			return nil, err
		}
		for _, code := range b {
			result = append(result, code)
		}
	}
	return result, nil
}

func (i Instructions) Add(i2 []Instruction) Instructions {
	for _, instr := range i2 {
		i = append(i, instr)
	}
	return i
}

func (i Instructions) String() string {
	result := make([]string, len(i))
	for j, instr := range i {
		result[j] = instr.String()
	}
	return strings.Join(result, "\n")
}

func CompileInstruction(instr []Instruction, debug bool) (MachineCode, error) {
	result := []uint8{}
	address := 0
	for _, i := range instr {
		b, err := i.Encode()
		if err != nil {
			return nil, err
		}
		if debug {
			fmt.Printf("0x%x: %s\n", address, i.String())
		}
		address += len(b)
		if debug {
			fmt.Println(MachineCode(b))
		}
		for _, code := range b {
			result = append(result, code)
		}
	}
	return result, nil
}

func Instruction_Length(instr Instruction) (int, error) {
	b, err := instr.Encode()
	if err != nil {
		return 0, err
	}
	return len(b), nil

}
