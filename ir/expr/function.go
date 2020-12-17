package expr

import (
	"fmt"
	"strings"

	"github.com/bspaans/jit-compiler/asm"
	"github.com/bspaans/jit-compiler/asm/encoding"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

type IR_Function struct {
	*BaseIRExpression
	Signature *TFunction
	Body      IR
	Address   int
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

func (i *IR_Function) Encode(ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {
	ownLength := uint(7)
	diff := uint(ctx.InstructionPointer+ownLength) - uint(i.Address)
	result := []lib.Instruction{asm.LEA(&encoding.RIPRelative{encoding.Int32(int32(-diff))}, target)}
	ctx.AddInstructions(result)
	return result, nil
}

func (b *IR_Function) encodeFunction(ctx *IR_Context) ([]uint8, error) {

	// TODO: restore rbx, rbp, r12-r15
	targets := []*encoding.Register{encoding.Rdi, encoding.Rsi, encoding.Rdx, encoding.R10, encoding.R8, encoding.R9}
	returnTarget := encoding.Rax
	registers := make([]bool, 16)
	registers[returnTarget.Register] = true
	variableMap := map[string]encoding.Operand{}
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
	return lib.Instructions(instr).Encode()
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
	b.Address = ctx.AddToDataSection(code)
	return nil
}
func (b *IR_Function) SSA_Transform(ctx *SSA_Context) (SSA_Rewrites, IRExpression) {
	newBody := b.Body.SSA_Transform(ctx)
	return nil, NewIR_Function(b.Signature, newBody)
}
