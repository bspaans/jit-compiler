package ir

import (
	"errors"
	"fmt"

	"github.com/bspaans/jit/asm"
)

type IRExpressionType int
type IRExpression interface {
	Type() IRExpressionType
	AddToDataSection(ctx *IR_Context)
	Encode(ctx *IR_Context, target *asm.Register) ([]asm.Instruction, error)
	String() string
}

const (
	Uint64    IRExpressionType = iota
	ByteArray IRExpressionType = iota
	Bool      IRExpressionType = iota
	Not       IRExpressionType = iota
	Add       IRExpressionType = iota
	Variable  IRExpressionType = iota
	Equals    IRExpressionType = iota
	Syscall   IRExpressionType = iota
)

type BaseIRExpression struct {
	typ IRExpressionType
}

func NewBaseIRExpression(typ IRExpressionType) *BaseIRExpression {
	return &BaseIRExpression{
		typ: typ,
	}
}

func (b *BaseIRExpression) Type() IRExpressionType {
	return b.typ
}

func IREXpression_length(expr IRExpression, ctx *IR_Context, target *asm.Register) (int, error) {
	commit := ctx.Commit
	ctx.Commit = false
	instr, err := expr.Encode(ctx, target)
	if err != nil {
		return 0, err
	}
	code, err := asm.Instructions(instr).Encode()
	if err != nil {
		return 0, err
	}
	ctx.Commit = commit
	return len(code), nil
}

func (b *BaseIRExpression) AddToDataSection(ctx *IR_Context) {}

type IR_Uint64 struct {
	*BaseIRExpression
	Value uint64
}

func NewIR_Uint64(v uint64) *IR_Uint64 {
	return &IR_Uint64{
		BaseIRExpression: NewBaseIRExpression(Uint64),
		Value:            v,
	}
}

func (i *IR_Uint64) String() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *IR_Uint64) Encode(ctx *IR_Context, target *asm.Register) ([]asm.Instruction, error) {
	result := []asm.Instruction{&asm.MOV{asm.Uint64(i.Value), target}}
	ctx.AddInstructions(result)
	return result, nil
}

type IR_Bool struct {
	*BaseIRExpression
	Value bool
}

func NewIR_Bool(v bool) *IR_Bool {
	return &IR_Bool{
		BaseIRExpression: NewBaseIRExpression(Bool),
		Value:            v,
	}
}

func (i *IR_Bool) String() string {
	return fmt.Sprintf("%v", i.Value)
}

func (i *IR_Bool) Encode(ctx *IR_Context, target *asm.Register) ([]asm.Instruction, error) {
	value := 0
	if i.Value {
		value = 1
	}
	result := []asm.Instruction{&asm.MOV{asm.Uint64(value), target}}
	ctx.AddInstructions(result)
	return result, nil
}

type IR_Variable struct {
	*BaseIRExpression
	Value string
}

func NewIR_Variable(v string) *IR_Variable {
	return &IR_Variable{
		BaseIRExpression: NewBaseIRExpression(Variable),
		Value:            v,
	}
}

func (i *IR_Variable) String() string {
	return i.Value
}

func (i *IR_Variable) Encode(ctx *IR_Context, target *asm.Register) ([]asm.Instruction, error) {
	reg := asm.Get64BitRegisterByIndex(ctx.VariableMap[i.Value])
	result := []asm.Instruction{&asm.MOV{reg, target}}
	ctx.AddInstructions(result)
	return result, nil
}

type IR_Add struct {
	*BaseIRExpression
	Op1 IRExpression
	Op2 IRExpression
}

func NewIR_Add(op1, op2 IRExpression) *IR_Add {
	return &IR_Add{
		BaseIRExpression: NewBaseIRExpression(Add),
		Op1:              op1,
		Op2:              op2,
	}
}

func (i *IR_Add) Encode(ctx *IR_Context, target *asm.Register) ([]asm.Instruction, error) {
	result := []asm.Instruction{}
	if i.Op1.Type() == Variable {
		variable := i.Op1.(*IR_Variable).Value
		reg := asm.Get64BitRegisterByIndex(ctx.VariableMap[variable])
		result = append(result, &asm.MOV{reg, target})
	} else if i.Op1.Type() == Uint64 {
		value := i.Op1.(*IR_Uint64).Value
		result = append(result, &asm.MOV{asm.Uint64(value), target})
	} else {
		return nil, errors.New("Unsupported add IR operation")
	}

	if i.Op2.Type() == Variable {
		variable := i.Op2.(*IR_Variable).Value
		reg := asm.Get64BitRegisterByIndex(ctx.VariableMap[variable])
		result = append(result, &asm.ADD{reg, target})
	} else if i.Op2.Type() == Uint64 {
		value := i.Op2.(*IR_Uint64).Value
		reg := asm.Get64BitRegisterByIndex(ctx.AllocateRegister())
		result = append(result, &asm.MOV{asm.Uint64(value), reg})
		result = append(result, &asm.ADD{reg, target})
		ctx.DeallocateRegister(reg.Register)
	} else {
		return nil, errors.New("Unsupported add IR operation")
	}
	ctx.AddInstructions(result)
	return result, nil
}

func (i *IR_Add) String() string {
	return fmt.Sprintf("%s + %s", i.Op1.String(), i.Op2.String())
}

type IR_Equals struct {
	*BaseIRExpression
	Op1 IRExpression
	Op2 IRExpression
}

func NewIR_Equals(op1, op2 IRExpression) *IR_Equals {
	return &IR_Equals{
		BaseIRExpression: NewBaseIRExpression(Equals),
		Op1:              op1,
		Op2:              op2,
	}
}

func (i *IR_Equals) Encode(ctx *IR_Context, target *asm.Register) ([]asm.Instruction, error) {
	result := []asm.Instruction{}

	var reg1, reg2 *asm.Register

	if i.Op1.Type() == Variable {
		variable := i.Op1.(*IR_Variable).Value
		reg1 = asm.Get64BitRegisterByIndex(ctx.VariableMap[variable])
	} else if i.Op1.Type() == Uint64 {
		value := i.Op1.(*IR_Uint64).Value
		r := ctx.AllocateRegister()
		defer ctx.DeallocateRegister(r)
		reg1 = asm.Get64BitRegisterByIndex(r)
		result = append(result, &asm.MOV{asm.Uint64(value), reg1})
	} else {
		return nil, errors.New("Unsupported cmp IR operation")
	}

	if i.Op2.Type() == Variable {
		variable := i.Op2.(*IR_Variable).Value
		reg2 = asm.Get64BitRegisterByIndex(ctx.VariableMap[variable])
	} else if i.Op2.Type() == Uint64 {
		value := i.Op2.(*IR_Uint64).Value
		r := ctx.AllocateRegister()
		defer ctx.DeallocateRegister(r)
		reg2 = asm.Get64BitRegisterByIndex(r)
		result = append(result, &asm.MOV{asm.Uint64(value), reg2})
	} else {
		return nil, errors.New("Unsupported add IR operation")
	}
	result = append(result, &asm.CMP{reg1, reg2})
	result = append(result, &asm.MOV{asm.Uint64(0), target})
	result = append(result, &asm.SETE{target.Lower8BitRegister()})
	ctx.AddInstructions(result)
	return result, nil
}

func (i *IR_Equals) String() string {
	return fmt.Sprintf("%s == %s", i.Op1.String(), i.Op2.String())
}

type IR_Not struct {
	*BaseIRExpression
	Op1 IRExpression
}

func NewIR_Not(op1 IRExpression) *IR_Not {
	return &IR_Not{
		BaseIRExpression: NewBaseIRExpression(Not),
		Op1:              op1,
	}
}

func (i *IR_Not) Encode(ctx *IR_Context, target *asm.Register) ([]asm.Instruction, error) {

	var reg1 *asm.Register

	result := []asm.Instruction{}
	if i.Op1.Type() == Variable {
		variable := i.Op1.(*IR_Variable).Value
		reg1 = asm.Get64BitRegisterByIndex(ctx.VariableMap[variable])
	} else if i.Op1.Type() == Uint64 {
		value := i.Op1.(*IR_Uint64).Value
		r := ctx.AllocateRegister()
		defer ctx.DeallocateRegister(r)
		reg1 = asm.Get64BitRegisterByIndex(r)
		result = append(result, &asm.MOV{asm.Uint64(value), reg1})
	} else if i.Op1.Type() == Equals {
		result_, err := i.Op1.Encode(ctx, target)
		if err != nil {
			return nil, err
		}
		for _, r := range result_ {
			result = append(result, r)
		}
		reg1 = target
	} else {
		return nil, errors.New("Unsupported not IR operation: " + i.String())
	}

	instr := []asm.Instruction{}
	instr = append(instr, &asm.CMP{asm.Uint32(0), reg1})
	instr = append(instr, &asm.MOV{asm.Uint64(0), target})
	instr = append(instr, &asm.SETE{target.Lower8BitRegister()})
	for _, inst := range instr {
		result = append(result, inst)
		ctx.AddInstruction(inst)
	}
	return result, nil
}

func (i *IR_Not) String() string {
	return fmt.Sprintf("!(%s)", i.Op1.String())
}

type IR_ByteArray struct {
	*BaseIRExpression
	Value   []uint8
	address int
}

func NewIR_ByteArray(value []uint8) *IR_ByteArray {
	return &IR_ByteArray{
		BaseIRExpression: NewBaseIRExpression(ByteArray),
		Value:            value,
	}
}

func (i *IR_ByteArray) String() string {
	return fmt.Sprintf("%v", i.Value)
}

func (i *IR_ByteArray) Encode(ctx *IR_Context, target *asm.Register) ([]asm.Instruction, error) {
	// Calculate the displacement between RIP (the instruction pointer,
	// pointing to the *next* instruction) and the address of our byte array,
	// and load the resulting address into target using a LEA instruction.
	ownLength := uint(7)
	diff := uint(ctx.InstructionPointer+ownLength) - uint(i.address)
	result := []asm.Instruction{&asm.LEA{&asm.RIPRelative{asm.Int32(int32(-diff))}, target}}
	ctx.AddInstructions(result)
	return result, nil
}

func (b *IR_ByteArray) AddToDataSection(ctx *IR_Context) {
	b.address = ctx.AddToDataSection(b.Value)
}

type IR_Syscall struct {
	*BaseIRExpression
	Syscall uint
	Args    []IRExpression
}

func NewIR_Syscall(syscall_nr uint, args []IRExpression) *IR_Syscall {
	return &IR_Syscall{
		Syscall: syscall_nr,
		Args:    args,
	}
}

func (i *IR_Syscall) String() string {
	return fmt.Sprintf("syscall(%v, %v)", i.Syscall, i.Args)
}

func (i *IR_Syscall) Encode(ctx *IR_Context, target *asm.Register) ([]asm.Instruction, error) {

	result := []asm.Instruction{}
	if ctx.Registers[0] {
		push := &asm.PUSH{asm.Rax}
		result = append(result, push)
		ctx.AddInstruction(push)
	}
	targets := []*asm.Register{asm.Rdi, asm.Rsi, asm.Rdx, asm.R10, asm.R8, asm.R9}
	targetRegisterIndices := []uint8{7, 6, 2, 10, 8, 9}
	clobbered := 0
	for j, argTarget := range targets {
		if j < len(i.Args) {
			// Push registers on the stack if they are in use
			registerIndex := targetRegisterIndices[j]
			if ctx.Registers[registerIndex] {
				reg := asm.Get64BitRegisterByIndex(registerIndex)
				result = append(result, &asm.PUSH{reg})
				ctx.AddInstruction(&asm.PUSH{reg})
				clobbered += 1
			}
			instr, err := i.Args[j].Encode(ctx, argTarget)
			if err != nil {
				return nil, err
			}
			for _, code := range instr {
				result = append(result, code)
			}
		}
	}

	instr := []asm.Instruction{
		&asm.MOV{asm.Uint64(uint64(i.Syscall)), asm.Rax},
		&asm.SYSCALL{},
		&asm.MOV{asm.Rax, target},
	}
	// Restore registers from the stack
	for j := clobbered; j > 0; j-- {
		registerIndex := targetRegisterIndices[j-1]
		reg := asm.Get64BitRegisterByIndex(registerIndex)
		instr = append(instr, &asm.POP{reg})
	}
	// restore rax
	if ctx.Registers[0] {
		instr = append(instr, &asm.POP{asm.Rax})
	}
	for _, inst := range instr {
		result = append(result, inst)
		ctx.AddInstruction(inst)
	}
	return result, nil
}

func (b *IR_Syscall) AddToDataSection(ctx *IR_Context) {
	for _, arg := range b.Args {
		arg.AddToDataSection(ctx)
	}
}

type IR_Syscall_Linux uint

const (
	IR_Syscall_Linux_Read  IR_Syscall_Linux = 0
	IR_Syscall_Linux_Write IR_Syscall_Linux = 1
	IR_Syscall_Linux_Open  IR_Syscall_Linux = 2
	IR_Syscall_Linux_Close IR_Syscall_Linux = 3
)

func NewIR_LinuxWrite(fid uint64, b []uint8, size int) IRExpression {
	return NewIR_Syscall(uint(IR_Syscall_Linux_Write), []IRExpression{NewIR_Uint64(fid), NewIR_ByteArray(b), NewIR_Uint64(uint64(size))})
}
func NewIR_LinuxClose(fid uint64) IRExpression {
	return NewIR_Syscall(uint(IR_Syscall_Linux_Close), []IRExpression{NewIR_Uint64(fid)})
}
