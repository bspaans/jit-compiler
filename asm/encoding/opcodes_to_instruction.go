package encoding

import (
	"fmt"
	"strings"

	"github.com/bspaans/jit/lib"
)

func OpcodesToInstruction(name string, opcodes []*Opcode, argCount int, operands ...Operand) lib.Instruction {
	maps := OpcodesToOpcodeMaps(opcodes, argCount)
	return NewOpcodeMapsInstruction(name, maps, operands)
}

type opcodeMapsInstruction struct {
	Name       string
	opcodeMaps OpcodeMaps
	Operands   []Operand
}

func NewOpcodeMapsInstruction(name string, maps OpcodeMaps, operands []Operand) lib.Instruction {
	return &opcodeMapsInstruction{name, maps, operands}
}

func (o *opcodeMapsInstruction) Encode() (lib.MachineCode, error) {
	opcode := o.opcodeMaps.ResolveOpcode(o.Operands)
	if opcode == nil {
		return nil, fmt.Errorf("unsupported %s instruction", o.Name)
	}
	return opcode.Encode(o.Operands)
}

func (o *opcodeMapsInstruction) String() string {
	opcode := o.opcodeMaps.ResolveOpcode(o.Operands)
	if opcode == nil {
		return "<unmatched " + o.Name + " instruction>"
	}
	args := []string{}
	for _, arg := range o.Operands {
		args = append(args, arg.String())
	}
	return opcode.Name + " " + strings.Join(args, ", ")
}
