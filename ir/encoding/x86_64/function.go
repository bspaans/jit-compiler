package x86_64

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm/x86_64"
	"github.com/bspaans/jit-compiler/asm/x86_64/encoding"
	"github.com/bspaans/jit-compiler/ir/expr"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

func encode_IR_Function(i *expr.IR_Function, ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {
	ownLength := uint(7)
	diff := uint(ctx.InstructionPointer+ownLength) - uint(i.Address)
	result := []lib.Instruction{x86_64.LEA(&encoding.RIPRelative{encoding.Int32(int32(-diff))}, target)}
	ctx.AddInstructions(result)
	return result, nil
}

func encode_IR_Function_for_DataSection(b *expr.IR_Function, ctx *IR_Context, readonly, rw []uint8) ([]uint8, []uint8, error) {

	// TODO: restore rbx, rbp, r12-r15
	targets := []*encoding.Register{encoding.Rdi, encoding.Rsi, encoding.Rdx, encoding.R10, encoding.R8, encoding.R9}
	returnTarget := encoding.Rax
	registers := make([]bool, 16)
	registers[returnTarget.Register] = true
	variableMap := map[string]encoding.Operand{}
	variableTypes := map[string]Type{}
	for i, arg := range b.Signature.Args {
		if arg.Type() == T_Float64 {
			return nil, nil, fmt.Errorf("Float arguments not supported")
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
	instr, err := encodeStatement(b.Body, ctx_)
	if err != nil {
		return nil, nil, err
	}
	for _, i := range instr {
		fmt.Println(i)
	}
	// TODO: should go to executable segment
	bytes, err := lib.Instructions(instr).Encode()
	if err != nil {
		return nil, nil, err
	}
	b.Address = len(rw) + 2
	rw = append(rw, bytes...)
	return readonly, rw, nil
}
