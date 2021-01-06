package opcodes

import (
	"fmt"
	"strings"

	. "github.com/bspaans/jit-compiler/asm/aarch64/encoding"
	"github.com/bspaans/jit-compiler/lib"
)

func OpcodeToInstruction(name string, opcode *Opcode, operands ...lib.Operand) lib.Instruction {
	return OpcodesToInstruction(name, []*Opcode{opcode}, operands...)
}

func OpcodesToInstruction(name string, opcodes []*Opcode, operands ...lib.Operand) lib.Instruction {
	opcode, err := resolveOpcode(name, opcodes, operands)
	if err != nil {
		panic(err)
	}
	return NewOpcodeInstruction(name, opcode, operands)
}

type opcodeInstruction struct {
	Name     string
	Opcode   *Opcode
	Operands []lib.Operand
}

func NewOpcodeInstruction(name string, o *Opcode, args []lib.Operand) *opcodeInstruction {
	return &opcodeInstruction{name, o, args}
}

func (o *opcodeInstruction) Encode() (lib.MachineCode, error) {
	return o.Opcode.Encode(o.Operands)
}

func (o *opcodeInstruction) String() string {
	if len(o.Operands) == 0 {
		return o.Name
	}
	args := []string{}
	for _, op := range o.Operands {
		args = append(args, op.String())
	}
	return o.Name + " " + strings.Join(args, ", ")
}

func resolveOpcode(name string, opcodes []*Opcode, operands []lib.Operand) (*Opcode, error) {

	matches := []*Opcode{}
	for _, opcode := range opcodes {
		if opcode.MatchesOperands(operands) {
			matches = append(matches, opcode)
		}
	}
	if len(matches) == 0 {
		args := []string{}
		for _, arg := range operands {
			if arg == nil {
				return nil, fmt.Errorf("nil arg in instruction: %s %s", name, strings.Join(args, ", "))
			}
			args = append(args, arg.String())
		}
		return nil, fmt.Errorf("unsupported %s instruction: %s %s", name, name, strings.Join(args, ", "))
	}
	return matches[0], nil
}
