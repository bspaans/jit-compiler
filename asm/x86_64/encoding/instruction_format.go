package encoding

import (
	"github.com/bspaans/jit-compiler/lib"
)

type InstructionFormat struct {
	Prefixes     []uint8
	VEXPrefix    *VEXPrefix
	REXPrefix    *REXPrefix
	Opcode       []uint8
	ModRM        *ModRM
	SIB          *SIB
	Displacement []uint8
	Immediate    []uint8
}

func NewInstructionFormat(opcode []uint8) *InstructionFormat {
	o := make([]uint8, len(opcode))
	for i, op := range opcode {
		o[i] = op
	}
	return &InstructionFormat{
		Prefixes:     []uint8{},
		Opcode:       o,
		Displacement: []uint8{},
		Immediate:    []uint8{},
	}
}

func (i *InstructionFormat) SetModRM(mode Mode, rm, reg uint8) {
	i.ModRM = NewModRM(mode, rm, reg)
}

func (i *InstructionFormat) SetDisplacement(op lib.Operand, displacement []uint8) {
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
	if i.VEXPrefix != nil {
		for _, b := range i.VEXPrefix.Encode() {
			result = append(result, b)
		}
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
	if i.SIB != nil {
		result = append(result, i.SIB.Encode())
	}
	for _, b := range i.Displacement {
		result = append(result, b)
	}
	for _, b := range i.Immediate {
		result = append(result, b)
	}
	return result
}
