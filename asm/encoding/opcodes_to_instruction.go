package encoding

import (
	"fmt"
	"strings"

	"github.com/bspaans/jit/lib"
)

func OpcodesToInstruction(opcodes []*Opcode, operands []Operand, argCount int) lib.Instruction {
	maps := OpcodesToOpcodeMaps(opcodes, argCount)
	return NewOpcodeMapsInstruction(maps, operands)
}

type opcodeMapsInstruction struct {
	opcodeMaps OpcodeMaps
	Operands   []Operand
}

func NewOpcodeMapsInstruction(maps OpcodeMaps, operands []Operand) lib.Instruction {
	return &opcodeMapsInstruction{maps, operands}
}

func (o *opcodeMapsInstruction) Encode() (lib.MachineCode, error) {
	opcode := o.opcodeMaps.ResolveOpcode(o.Operands)
	if opcode == nil {
		return nil, fmt.Errorf("unsupported instruction")
	}
	return opcode.Encode(o.Operands)
}

func (o *opcodeMapsInstruction) String() string {
	opcode := o.opcodeMaps.ResolveOpcode(o.Operands)
	if opcode == nil {
		return "<unmatched instruction>"
	}
	args := []string{}
	for _, arg := range o.Operands {
		args = append(args, arg.String())
	}
	return opcode.Name + " " + strings.Join(args, ", ")
}
