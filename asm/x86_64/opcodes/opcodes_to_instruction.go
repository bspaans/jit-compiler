package opcodes

import (
	"fmt"
	"strings"

	. "github.com/bspaans/jit-compiler/asm/x86_64/encoding"
	"github.com/bspaans/jit-compiler/lib"
)

func OpcodeToInstruction(name string, opcode *Opcode, argCount int, operands ...lib.Operand) lib.Instruction {
	return OpcodesToInstruction(name, []*Opcode{opcode}, argCount, operands...)
}

func OpcodesToInstruction(name string, opcodes []*Opcode, argCount int, operands ...lib.Operand) lib.Instruction {
	maps := OpcodesToOpcodeMaps(opcodes, argCount)
	return NewOpcodeMapsInstruction(name, maps, operands, opcodes)
}

type opcodeMapsInstruction struct {
	Name       string
	opcodeMaps OpcodeMaps
	Operands   []lib.Operand
	Opcodes    []*Opcode
}

func NewOpcodeMapsInstruction(name string, maps OpcodeMaps, operands []lib.Operand, opcodes []*Opcode) lib.Instruction {
	return &opcodeMapsInstruction{name, maps, operands, opcodes}
}

func (o *opcodeMapsInstruction) Encode() (lib.MachineCode, error) {
	if len(o.Operands) == 0 {
		return o.Opcodes[0].Encode([]lib.Operand{})
	}
	opcode := o.opcodeMaps.ResolveOpcode(o.Operands)
	if opcode == nil {
		args := []string{}
		for _, arg := range o.Operands {
			if arg == nil {
				return nil, fmt.Errorf("nil arg in instruction: %s %s", o.Name, strings.Join(args, ", "))
			}
			args = append(args, arg.String())
		}
		return nil, fmt.Errorf("unsupported %s instruction: %s %s", o.Name, o.Name, strings.Join(args, ", "))
	}
	return opcode.Encode(o.Operands)
}

func (o *opcodeMapsInstruction) String() string {
	if len(o.Operands) == 0 {
		return o.Name
	}
	opcode := o.opcodeMaps.ResolveOpcode(o.Operands)
	args := []string{}
	for i := len(o.Operands) - 1; i >= 0; i-- {
		args = append(args, o.Operands[i].String())
	}
	if opcode == nil {
		return fmt.Sprintf("<unmatched %s instruction: %s %s>", o.Name, o.Name, strings.Join(args, ", "))
	}
	return opcode.Name + " " + strings.Join(args, ", ")
}
