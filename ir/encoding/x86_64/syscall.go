package x86_64

import (
	"github.com/bspaans/jit-compiler/asm/x86_64"
	"github.com/bspaans/jit-compiler/asm/x86_64/encoding"
	"github.com/bspaans/jit-compiler/ir/expr"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

func encode_IR_Syscall(i *expr.IR_Syscall, ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {

	result, _, clobbered, err := ABI_Call_Setup(ctx, i.Args, TUint64)
	if err != nil {
		return nil, err
	}
	instr, err := encodeExpression(i.Syscall, ctx, encoding.Rax)
	if err != nil {
		return nil, err
	}
	for _, inst := range instr {
		result = append(result, inst)
	}
	tmpTarget := ctx.AllocateRegister(TUint64)
	defer ctx.DeallocateRegister(tmpTarget)

	instr = []lib.Instruction{
		x86_64.SYSCALL(),
		x86_64.MOV(encoding.Rax, tmpTarget),
	}
	for _, inst := range instr {
		result = append(result, inst)
		ctx.AddInstruction(inst)
	}
	restore := RestoreRegisters(ctx, clobbered)
	result = result.Add(restore)
	mov := x86_64.MOV(tmpTarget, target)
	ctx.AddInstruction(mov)
	result = append(result, mov)
	return result, nil
}
