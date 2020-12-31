package expr

import (
	"fmt"
	"strings"

	"github.com/bspaans/jit-compiler/asm/x86_64"
	"github.com/bspaans/jit-compiler/asm/x86_64/encoding"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
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

	call := asm.CALL(function)
	tmpReg := ctx.AllocateRegister(TUint64)
	defer ctx.DeallocateRegister(tmpReg)
	mov := asm.MOV(encoding.Rax, tmpReg)
	ctx.AddInstruction(call)
	ctx.AddInstruction(mov)
	result = append(result, call)
	result = append(result, mov)

	restore := RestoreRegisters(ctx, clobbered)
	result = result.Add(restore)

	mov = asm.MOV(tmpReg, target)
	ctx.AddInstruction(mov)
	result = append(result, mov)
	return result, nil
}

func (b *IR_Call) AddToDataSection(ctx *IR_Context) error {
	for _, arg := range b.Args {
		if err := arg.AddToDataSection(ctx); err != nil {
			return err
		}
	}
	return nil
}
func (b *IR_Call) SSA_Transform(ctx *SSA_Context) (SSA_Rewrites, IRExpression) {
	newArgs := make([]IRExpression, len(b.Args))
	rewrites := SSA_Rewrites{}
	for i, arg := range b.Args {
		if IsLiteralOrVariable(arg) {
			newArgs[i] = arg
		} else {
			rw, expr := arg.SSA_Transform(ctx)
			for _, rewrite := range rw {
				rewrites = append(rewrites, rewrite)
			}
			v := ctx.GenerateVariable()
			rewrites = append(rewrites, NewSSA_Rewrite(v, expr))
			newArgs[i] = NewIR_Variable(v)
		}
	}
	return rewrites, NewIR_Call(b.Function, newArgs)
}
