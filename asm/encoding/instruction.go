package encoding

import (
	"github.com/bspaans/jit/lib"
)

type InstructionFormat struct {
	Prefixes     []uint8
	REXPrefix    *REXPrefix
	Opcode       []uint8
	ModRM        *ModRM
	Displacement []uint8
	Immediate    []uint8
}

func NewInstructionFormat(opcode []uint8) *InstructionFormat {
	return &InstructionFormat{
		Prefixes:     []uint8{},
		Opcode:       opcode,
		Displacement: []uint8{},
		Immediate:    []uint8{},
	}
}

func (i *InstructionFormat) SetModRM(mode Mode, rm, reg uint8) {
	i.ModRM = NewModRM(mode, rm, reg)
}

func (i *InstructionFormat) SetDisplacement(op Operand, displacement []uint8) {
	// Not sure why this is needed, but it is
	if _, ok := op.(*Register); ok && op.(*Register) == Rsp {
		i.Displacement = append(i.Displacement, 0x24)
	}
	for _, d := range displacement {
		i.Displacement = append(i.Displacement, d)
	}
}

func (i *InstructionFormat) Encode() lib.MachineCode {
	result := []uint8{}
	for _, b := range i.Prefixes {
		result = append(result, b)
	}
	if i.REXPrefix != nil {
		result = append(result, i.REXPrefix.Encode())
	}
	for _, b := range i.Opcode {
		result = append(result, b)
	}
	if i.ModRM != nil {
		result = append(result, i.ModRM.Encode())
	}
	for _, b := range i.Displacement {
		result = append(result, b)
	}
	for _, b := range i.Immediate {
		result = append(result, b)
	}
	return result
}
