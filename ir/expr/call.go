package expr

import (
	"fmt"
	"strings"

	"github.com/bspaans/jit/asm"
	. "github.com/bspaans/jit/ir/shared"
)

type IR_Call struct {
	*BaseIRExpression
	Function string
	Args     []IRExpression
}

func NewIR_Call(function string, args []IRExpression) *IR_Call {
	return &IR_Call{
		BaseIRExpression: NewBaseIRExpression(Call),
		Function:         function,
		Args:             args,
	}
}

func (i *IR_Call) ReturnType(ctx *IR_Context) Type {
	signature := ctx.VariableTypes[i.Function]
	if signature == nil {
		panic("Unknown function: " + i.Function)
	}
	if _, ok := signature.(*TFunction); !ok {
		panic("Expected function, got: " + signature.String())
	}
	return signature.(*TFunction).ReturnType
}

func (i *IR_Call) String() string {
	args := []string{}
	for _, arg := range i.Args {
		args = append(args, arg.String())
	}
	return fmt.Sprintf("%s(%s)", i.Function, strings.Join(args, ", "))
}

func (i *IR_Call) Encode(ctx *IR_Context, target *asm.Register) ([]asm.Instruction, error) {

	// TODO preserve arguments
	// TODO result, clobbered, err := ABI_Call_Setup(ctx, args, i.ReturnType(ctx))
	result := []asm.Instruction{}
	targets := []*asm.Register{asm.Rdi, asm.Rsi, asm.Rdx, asm.R10, asm.R8, asm.R9}
	for j, arg := range i.Args {
		if arg.ReturnType(ctx) == TFloat64 {
			return nil, fmt.Errorf("Float arguments not supported")
		}
		instr, err := arg.Encode(ctx, targets[j])
		if err != nil {
			return nil, err
		}
		for _, inst := range instr {
			result = append(result, inst)
		}
	}
	// TODO: variable register could have been clobbered...!
	function := ctx.VariableMap[i.Function]
	if function == nil {
		return nil, fmt.Errorf("Unknown function:" + i.Function)
	}
	call := &asm.CALL{function}
	mov := &asm.MOV{asm.Rax, target}
	ctx.AddInstruction(call)
	ctx.AddInstruction(mov)
	result = append(result, call)
	result = append(result, mov)
	return result, nil
}

// Sets up for a CALL or SYSCALL. Returns instructions and clobbered registers
func ABI_Call_Setup(abi ABI, ctx *IR_Context, args []IRExpression, returnType Type) ([]asm.Instruction, []*asm.Register, error) {
	argTypes := []Type{}
	for _, arg := range args {
		argTypes = append(argTypes, arg.ReturnType(ctx))
	}
	result, clobbered := PreserveRegisters(abi, ctx, argTypes, returnType)
	targets := abi.GetRegistersForArgs(argTypes)

	for i, arg := range args {
		instr, err := arg.Encode(ctx, targets[i])
		if err != nil {
			return nil, clobbered, err
		}
		result = result.Add(instr)
	}
	return result, clobbered, nil
}
