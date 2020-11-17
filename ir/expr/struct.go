package expr

import (
	"fmt"
	"strings"

	"github.com/bspaans/jit-compiler/asm"
	"github.com/bspaans/jit-compiler/asm/encoding"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

type IR_Struct struct {
	*BaseIRExpression
	StructType *TStruct
	Values     []IRExpression
	address    int
}

func NewIR_Struct(ty *TStruct, values []IRExpression) *IR_Struct {
	return &IR_Struct{
		BaseIRExpression: NewBaseIRExpression(Struct),
		StructType:       ty,
		Values:           values,
	}
}

func (i *IR_Struct) ReturnType(ctx *IR_Context) Type {
	return i.StructType
}

func (i *IR_Struct) String() string {
	values := []string{}
	for _, v := range i.Values {
		values = append(values, v.String())
	}
	return i.StructType.String() + "{" + strings.Join(values, ", ") + "}"
}

func (i *IR_Struct) Encode(ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {
	// Calculate the displacement between RIP (the instruction pointer,
	// pointing to the *next* instruction) and the address of our byte array,
	// and load the resulting address into target using a LEA instruction.
	ownLength := uint(7)
	diff := uint(ctx.InstructionPointer+ownLength) - uint(i.address)
	result := []lib.Instruction{asm.LEA(&encoding.RIPRelative{encoding.Int32(int32(-diff))}, target)}
	ctx.AddInstructions(result)
	return result, nil
}

func (b *IR_Struct) AddToDataSection(ctx *IR_Context) error {
	b.address = -1
	for _, v := range b.Values {
		bytes := []uint8{}
		if ir, ok := v.(*IR_Uint64); ok {
			bytes = encoding.Uint64(ir.Value).Encode()
		} else if ir, ok := v.(*IR_Int64); ok {
			bytes = encoding.Uint64(ir.Value).Encode()
		} else if ir, ok := v.(*IR_Float64); ok {
			bytes = encoding.Float64(ir.Value).Encode()
		} else {
			return fmt.Errorf("Unsupported struct type %s", v.Type())
		}
		addr := ctx.AddToDataSection(bytes)
		if b.address == -1 {
			b.address = addr
		}
	}
	return nil
}
