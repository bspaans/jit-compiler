package shared

import (
	"github.com/bspaans/jit-compiler/asm/x86_64/encoding"
	"github.com/bspaans/jit-compiler/lib"
)

type Architecture interface {
	EncodeExpression(expr IRExpression, ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error)
	EncodeStatement(stmt IR, ctx *IR_Context) ([]lib.Instruction, error)
	EncodeDataSection(stmts []IR, ctx *IR_Context) (*Segments, error)
	GetAllocator() Allocator
}

type Allocator interface {
	AllocateRegister(typ Type) encoding.Operand
	DeallocateRegister(encoding.Operand)
	Copy() Allocator
}

type IR_Context struct {
	Architecture Architecture
	ABI          ABI
	Allocator    Allocator

	VariableMap        map[string]encoding.Operand
	VariableTypes      map[string]Type
	ReturnOperandStack []encoding.Operand
	Segments           *Segments
	InstructionPointer uint
	StackPointer       int
	Commit             bool // if false turns AddInstruction into a noop

	instructions []lib.Instruction
}

func NewIRContext(arch Architecture, abi ABI) *IR_Context {
	ctx := &IR_Context{
		Architecture:       arch,
		ABI:                abi,
		VariableMap:        map[string]encoding.Operand{},
		VariableTypes:      map[string]Type{},
		ReturnOperandStack: []encoding.Operand{&encoding.DisplacedRegister{encoding.Rsp, 8}},
		InstructionPointer: 2,
		StackPointer:       8,
		Commit:             true,
		instructions:       []lib.Instruction{},
	}
	ctx.Allocator = arch.GetAllocator()
	return ctx
}

func (i *IR_Context) PushReturnOperand(op encoding.Operand) {
	i.ReturnOperandStack = append(i.ReturnOperandStack, op)
}
func (i *IR_Context) PeekReturn() encoding.Operand {
	return i.ReturnOperandStack[len(i.ReturnOperandStack)-1]
}

func (i *IR_Context) PopReturn() encoding.Operand {
	op := i.ReturnOperandStack[len(i.ReturnOperandStack)-1]
	i.ReturnOperandStack = i.ReturnOperandStack[:len(i.ReturnOperandStack)-1]
	return op
}

func (i *IR_Context) Copy() *IR_Context {
	variableMap := map[string]encoding.Operand{}
	for arg, reg := range i.VariableMap {
		variableMap[arg] = reg
	}
	variableTypes := map[string]Type{}
	for arg, ty := range i.VariableTypes {
		variableTypes[arg] = ty
	}
	instructions := []lib.Instruction{}
	for _, d := range i.instructions {
		instructions = append(instructions, d)
	}
	returns := []encoding.Operand{}
	for _, d := range i.ReturnOperandStack {
		returns = append(returns, d)
	}
	return &IR_Context{
		Architecture:       i.Architecture,
		ABI:                i.ABI,
		Allocator:          i.Allocator.Copy(),
		VariableMap:        variableMap,
		VariableTypes:      variableTypes,
		ReturnOperandStack: returns,
		Segments:           i.Segments,
		InstructionPointer: i.InstructionPointer,
		StackPointer:       i.StackPointer,
		Commit:             i.Commit,
		instructions:       instructions,
	}
}

func (i *IR_Context) AddInstruction(instr ...lib.Instruction) {
	if i.Commit {
		for _, in := range instr {
			i.instructions = append(i.instructions, in)
			length, _ := lib.Instruction_Length(in)
			i.InstructionPointer += uint(length)
		}
	}
}

func (i *IR_Context) GetInstructions() []lib.Instruction {
	return i.instructions
}

func (i *IR_Context) AllocateRegister(typ Type) *encoding.Register {
	return i.Allocator.AllocateRegister(typ).(*encoding.Register) // TODO: remove cast
}

func (i *IR_Context) DeallocateRegister(reg *encoding.Register) {
	i.Allocator.DeallocateRegister(reg) // TODO: remove cast
}
