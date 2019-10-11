package expr

import (
	"fmt"
	"strings"

	"github.com/bspaans/jit/asm"
	. "github.com/bspaans/jit/ir/shared"
)

type IR_Function struct {
	*BaseIRExpression
	Signature *TFunction
	Body      IR
	address   int
}

func NewIR_Function(signature *TFunction, body IR) *IR_Function {
	return &IR_Function{
		BaseIRExpression: NewBaseIRExpression(Function),
		Signature:        signature,
		Body:             body,
	}
}

func (i *IR_Function) ReturnType(ctx *IR_Context) Type {
	return i.Signature
}

func (i *IR_Function) String() string {
	args := []string{}
	for j, arg := range i.Signature.ArgNames {
		args = append(args, arg+" "+i.Signature.Args[j].String())
	}
	return fmt.Sprintf("func(%s) %s { %s }", strings.Join(args, ", "), i.Signature.ReturnType.String(), i.Body.String())
}

func (i *IR_Function) Encode(ctx *IR_Context, target asm.Operand) ([]asm.Instruction, error) {
	ownLength := uint(7)
	diff := uint(ctx.InstructionPointer+ownLength) - uint(i.address)
	result := []asm.Instruction{&asm.LEA{&asm.RIPRelative{asm.Int32(int32(-diff))}, target}}
	ctx.AddInstructions(result)
	return result, nil
}

func (b *IR_Function) encodeFunction(ctx *IR_Context) ([]uint8, error) {

	// TODO: restore rbx, rbp, r12-r15
	targets := []*asm.Register{asm.Rdi, asm.Rsi, asm.Rdx, asm.R10, asm.R8, asm.R9}
	returnTarget := asm.Rax
	registers := make([]bool, 16)
	registers[returnTarget.Register] = true
	variableMap := map[string]asm.Operand{}
	variableTypes := map[string]Type{}
	for i, arg := range b.Signature.Args {
		if arg.Type() == T_Float64 {
			return nil, fmt.Errorf("Float arguments not supported")
		}
		v := b.Signature.ArgNames[i]
		registers[targets[i].Register] = true
		variableMap[v] = targets[i]
		variableTypes[v] = arg
	}

	ctx_ := ctx.Copy()
	ctx_.PushReturnOperand(returnTarget)
	ctx_.Commit = false
	ctx_.Registers = registers
	ctx_.RegistersAllocated = uint8(len(b.Signature.Args) + 1)
	ctx_.VariableMap = variableMap
	ctx_.VariableTypes = variableTypes
	instr, err := b.Body.Encode(ctx_)
	if err != nil {
		return nil, err
	}
	//result := []uint8{}
	for _, instr := range instr {
		fmt.Println(instr)
	}
	return asm.Instructions(instr).Encode()
}

func (b *IR_Function) AddToDataSection(ctx *IR_Context) error {
	if err := b.Body.AddToDataSection(ctx); err != nil {
		return err
	}

	code, err := b.encodeFunction(ctx)
	if err != nil {
		return err
	}
	_ = code
	b.address = ctx.AddToDataSection(code)
	return nil
}
