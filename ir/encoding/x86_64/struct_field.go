package x86_64

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm/x86_64"
	"github.com/bspaans/jit-compiler/asm/x86_64/encoding"
	"github.com/bspaans/jit-compiler/ir/expr"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

func encode_IR_StructField(i *expr.IR_StructField, ctx *IR_Context, target lib.Operand) ([]lib.Instruction, error) {
	// Pointer to target struct is stored into tmpReg
	tmpReg := ctx.AllocateRegister(TUint64)
	defer ctx.DeallocateRegister(tmpReg)
	result, err := encodeExpression(i.Struct, ctx, tmpReg)

	// Calculate the offset for our field
	structType := i.Struct.ReturnType(ctx)
	str, ok := structType.(*TStruct)
	if !ok {
		return nil, fmt.Errorf("Expecting struct, got %s", structType)
	}
	offset := 0
	for j, f := range str.Fields {
		if f == i.Field {
			break
		}
		offset += int(str.FieldTypes[j].Width())
	}
	// Add offset and load value at address into target
	add := x86_64.ADD(encoding.Uint32(uint32(offset)), tmpReg)
	mov := x86_64.MOV(&encoding.IndirectRegister{tmpReg}, target)
	ctx.AddInstruction(add)
	ctx.AddInstruction(mov)
	result = append(result, add)
	result = append(result, mov)
	return result, err
}
