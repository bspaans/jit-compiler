package expr

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm"
	"github.com/bspaans/jit-compiler/asm/encoding"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

type IR_StructField struct {
	*BaseIRExpression
	Struct IRExpression
	Field  string
}

func NewIR_StructField(structExpr IRExpression, field string) *IR_StructField {
	return &IR_StructField{
		BaseIRExpression: NewBaseIRExpression(StructField),
		Struct:           structExpr,
		Field:            field,
	}
}

func (i *IR_StructField) ReturnType(ctx *IR_Context) Type {
	structType := i.Struct.ReturnType(ctx)
	if str, ok := structType.(*TStruct); !ok {
		panic("Not a struct")
	} else {
		for j, f := range str.Fields {
			if f == i.Field {
				return str.FieldTypes[j]
			}
		}
	}
	panic("Field not found")
}

func (i *IR_StructField) String() string {
	return i.Struct.String() + "." + i.Field
}

func (i *IR_StructField) Encode(ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {
	// Pointer to target struct is stored into tmpReg
	tmpReg := ctx.AllocateRegister(TUint64)
	defer ctx.DeallocateRegister(tmpReg)
	result, err := i.Struct.Encode(ctx, tmpReg)

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
	add := asm.ADD(encoding.Uint32(uint32(offset)), tmpReg)
	mov := asm.MOV(&encoding.IndirectRegister{tmpReg}, target)
	ctx.AddInstruction(add)
	ctx.AddInstruction(mov)
	result = append(result, add)
	result = append(result, mov)
	return result, err
}

func (b *IR_StructField) AddToDataSection(ctx *IR_Context) error {
	return b.Struct.AddToDataSection(ctx)
}

func (b *IR_StructField) SSA_Transform(ctx *SSA_Context) (SSA_Rewrites, IRExpression) {
	return nil, b
}
