package x86_64

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm/x86_64"
	"github.com/bspaans/jit-compiler/asm/x86_64/encoding"
	"github.com/bspaans/jit-compiler/ir/expr"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

func encode_IR_Call(i *expr.IR_Call, ctx *IR_Context, target lib.Operand) ([]lib.Instruction, error) {
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

	call := x86_64.CALL(function)
	tmpReg := ctx.AllocateRegister(TUint64)
	defer ctx.DeallocateRegister(tmpReg)
	mov := x86_64.MOV(encoding.Rax, tmpReg)
	ctx.AddInstruction(call)
	ctx.AddInstruction(mov)
	result = append(result, call)
	result = append(result, mov)

	restore := RestoreRegisters(ctx, clobbered)
	result = result.Add(restore)

	mov = x86_64.MOV(tmpReg, target)
	ctx.AddInstruction(mov)
	result = append(result, mov)
	return result, nil
}
