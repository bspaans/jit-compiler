package shared

import (
	"fmt"

	"github.com/bspaans/jit/asm"
)

type IR_Context struct {
	ABI                     ABI
	Registers               []bool
	RegistersAllocated      uint8
	FloatRegisters          []bool
	FloatRegistersAllocated uint8
	VariableMap             map[string]*asm.Register
	VariableTypes           map[string]Type
	ReturnOperandStack      []asm.Operand
	DataSection             []uint8
	InstructionPointer      uint
	DataSectionOffset       int
	Commit                  bool // if false turns AddInstruction into a noop

	instructions []asm.Instruction
}

func NewIRContext() *IR_Context {
	ctx := &IR_Context{
		ABI:                     NewABI_AMDSystemV(),
		Registers:               make([]bool, 16),
		RegistersAllocated:      0,
		FloatRegisters:          make([]bool, 16),
		FloatRegistersAllocated: 0,
		VariableMap:             map[string]*asm.Register{},
		VariableTypes:           map[string]Type{},
		ReturnOperandStack:      []asm.Operand{&asm.DisplacedRegister{asm.Rsp, 8}},
		DataSection:             []uint8{},
		DataSectionOffset:       2,
		InstructionPointer:      2,
		Commit:                  true,
		instructions:            []asm.Instruction{},
	}
	// Always allocate rsp
	// Should track usage?
	ctx.Registers[4] = true
	ctx.RegistersAllocated = 1
	return ctx
}

func (i *IR_Context) PushReturnOperand(op asm.Operand) {
	i.ReturnOperandStack = append(i.ReturnOperandStack, op)
	fmt.Println(i.ReturnOperandStack)
}
func (i *IR_Context) PeekReturn() asm.Operand {
	return i.ReturnOperandStack[len(i.ReturnOperandStack)-1]
}

func (i *IR_Context) PopReturn() asm.Operand {
	fmt.Println(i.ReturnOperandStack)
	op := i.ReturnOperandStack[len(i.ReturnOperandStack)-1]
	i.ReturnOperandStack = i.ReturnOperandStack[:len(i.ReturnOperandStack)-1]
	fmt.Println(i.ReturnOperandStack)
	return op
}

func (i *IR_Context) Copy() *IR_Context {
	regs := make([]bool, 16)
	floatRegs := make([]bool, 16)
	for j := 0; j < 16; j++ {
		regs[j] = i.Registers[j]
		floatRegs[j] = i.FloatRegisters[j]
	}
	variableMap := map[string]*asm.Register{}
	for arg, reg := range i.VariableMap {
		variableMap[arg] = reg
	}
	variableTypes := map[string]Type{}
	for arg, ty := range i.VariableTypes {
		variableTypes[arg] = ty
	}
	ds := []uint8{}
	for _, d := range i.DataSection {
		ds = append(ds, d)
	}
	instructions := []asm.Instruction{}
	for _, d := range i.instructions {
		instructions = append(instructions, d)
	}
	returns := []asm.Operand{}
	for _, d := range i.ReturnOperandStack {
		returns = append(returns, d)
	}
	return &IR_Context{
		Registers:               regs,
		RegistersAllocated:      i.RegistersAllocated,
		FloatRegisters:          floatRegs,
		FloatRegistersAllocated: i.FloatRegistersAllocated,
		VariableMap:             variableMap,
		VariableTypes:           variableTypes,
		ReturnOperandStack:      returns,
		DataSection:             ds,
		DataSectionOffset:       i.DataSectionOffset,
		InstructionPointer:      i.InstructionPointer,
		Commit:                  i.Commit,
		instructions:            instructions,
	}
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

func (i *IR_Context) AllocateRegister(typ Type) *asm.Register {
	if typ == TFloat64 {
		return asm.GetFloatingPointRegisterByIndex(i.allocateFloatRegister())
	}
	return asm.Get64BitRegisterByIndex(i.allocateRegister())
}

func (i *IR_Context) DeallocateRegister(reg *asm.Register) {
	if reg.Size == asm.QUADDOUBLE {
		i.deallocateFloatRegister(reg.Register)
		return
	}
	i.deallocateRegister(reg.Register)
}

func (i *IR_Context) allocateRegister() uint8 {
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

func (i *IR_Context) deallocateRegister(reg uint8) {
	i.Registers[reg] = false
	i.RegistersAllocated -= 1
}

func (i *IR_Context) allocateFloatRegister() uint8 {
	if i.FloatRegistersAllocated >= 16 {
		panic("FloatRegister allocation limit. Needs stack handling")
	}
	for j := 0; j < len(i.FloatRegisters); j++ {
		if !i.FloatRegisters[j] {
			i.FloatRegisters[j] = true
			i.FloatRegistersAllocated += 1
			return uint8(j)
		}
	}
	panic("FloatRegister allocation limit reached with incorrect allocation counter. Needs stack handling")
}

func (i *IR_Context) deallocateFloatRegister(reg uint8) {
	i.FloatRegisters[reg] = false
	i.FloatRegistersAllocated -= 1
}

func (i *IR_Context) AddToDataSection(value []uint8) int {
	address := len(i.DataSection)
	for _, v := range value {
		i.DataSection = append(i.DataSection, v)
		i.InstructionPointer += 1
	}
	return address + i.DataSectionOffset
}
