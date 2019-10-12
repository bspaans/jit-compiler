package expr

import (
	"fmt"
	"strings"

	"github.com/bspaans/jit/asm"
	"github.com/bspaans/jit/asm/encoding"
	. "github.com/bspaans/jit/ir/shared"
	"github.com/bspaans/jit/lib"
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

func (i *IR_Call) Encode(ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {
	result, mapping, clobbered, err := ABI_Call_Setup(ctx, i.Args, ctx.VariableTypes[i.Function].(*TFunction).ReturnType)
	if err != nil {
		return nil, err
	}

	function := ctx.VariableMap[i.Function]
	if function == nil {
		return nil, fmt.Errorf("Unknown function:" + i.Function)
	}

	// Use a different address for function if its location has been clobbered
	if movedTarget, found := mapping[function]; found {
		function = movedTarget
	}
	fmt.Println("Calling", i.Function, "at", function)
	for op1, op2 := range mapping {
		fmt.Println(op1, "=>", op2)
	}

	call := &asm.CALL{function}
	mov := &asm.MOV{encoding.Rax, target}
	ctx.AddInstruction(call)
	ctx.AddInstruction(mov)
	result = append(result, call)
	result = append(result, mov)

	restore := RestoreRegisters(ctx, clobbered)
	result = result.Add(restore)
	return result, nil
}
