package ir

import "github.com/bspaans/jit/asm"

type IR_Context struct {
	Registers          []bool
	RegistersAllocated uint8
	VariableMap        map[string]uint8
	VariableTypes      map[string]Type
	DataSection        []uint8
	InstructionPointer uint
	DataSectionOffset  int
	Commit             bool // if false turns AddInstruction into a noop

	instructions []asm.Instruction
}

func NewIRContext() *IR_Context {
	ctx := &IR_Context{
		Registers:          make([]bool, 16),
		RegistersAllocated: 0,
		VariableMap:        map[string]uint8{},
		VariableTypes:      map[string]Type{},
		DataSection:        []uint8{},
		DataSectionOffset:  2,
		InstructionPointer: 2,
		Commit:             true,
		instructions:       []asm.Instruction{},
	}
	// Always allocate rsp
	// Should track usage?
	ctx.Registers[4] = true
	ctx.RegistersAllocated = 1
	return ctx
}

func (i *IR_Context) AddInstruction(instr asm.Instruction) {
	if i.Commit {
		i.instructions = append(i.instructions, instr)
		length, _ := asm.Instruction_Length(instr)
		i.InstructionPointer += uint(length)
	}
}

func (i *IR_Context) AddInstructions(instr []asm.Instruction) {
	for _, inst := range instr {
		i.AddInstruction(inst)
	}
}

func (i *IR_Context) GetInstructions() []asm.Instruction {
	return i.instructions
}

func (i *IR_Context) AllocateRegister() uint8 {
	if i.RegistersAllocated >= 16 {
		panic("Register allocation limit. Needs stack handling")
	}
	for j := 0; j < len(i.Registers); j++ {
		if !i.Registers[j] {
			i.Registers[j] = true
			i.RegistersAllocated += 1
			return uint8(j)
		}
	}
	panic("Register allocation limit reached with incorrect allocation counter. Needs stack handling")
}

func (i *IR_Context) DeallocateRegister(reg uint8) {
	i.Registers[reg] = false
	i.RegistersAllocated -= 1
}

func (i *IR_Context) AddToDataSection(value []uint8) int {
	address := len(i.DataSection)
	for _, v := range value {
		i.DataSection = append(i.DataSection, v)
		i.InstructionPointer += 1
	}
	return address + i.DataSectionOffset
}
