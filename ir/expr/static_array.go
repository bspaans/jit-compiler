package expr

import (
	"fmt"
	"strings"

	"github.com/bspaans/jit-compiler/asm"
	"github.com/bspaans/jit-compiler/asm/encoding"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

type IR_StaticArray struct {
	*BaseIRExpression
	ElemType Type
	Value    []IRExpression
	address  int
}

func NewIR_StaticArray(elemType Type, value []IRExpression) *IR_StaticArray {
	return &IR_StaticArray{
		BaseIRExpression: NewBaseIRExpression(StaticArray),
		ElemType:         elemType,
		Value:            value,
	}
}

func (i *IR_StaticArray) ReturnType(ctx *IR_Context) Type {
	return &TArray{i.ElemType, len(i.Value)}
}

func (i *IR_StaticArray) String() string {
	elems := []string{}
	for _, v := range i.Value {
		elems = append(elems, v.String())
	}
	return fmt.Sprintf("[]%s{%s}", i.ElemType.String(), strings.Join(elems, ", "))
}

func (i *IR_StaticArray) Encode(ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {
	// Calculate the displacement between RIP (the instruction pointer,
	// pointing to the *next* instruction) and the address of our byte array,
	// and load the resulting address into target using a LEA instruction.
	ownLength := uint(7)
	diff := uint(ctx.InstructionPointer+ownLength) - uint(i.address)
	result := []lib.Instruction{asm.LEA(&encoding.RIPRelative{encoding.Int32(int32(-diff))}, target)}
	ctx.AddInstructions(result)
	return result, nil
}

func (b *IR_StaticArray) AddToDataSection(ctx *IR_Context) error {
	b.address = -1
	for _, v := range b.Value {
		bytes := []uint8{}
		if ir, ok := v.(*IR_Uint64); ok {
			bytes = encoding.Uint64(ir.Value).Encode()
		} else if ir, ok := v.(*IR_Float64); ok {
			bytes = encoding.Float64(ir.Value).Encode()
		} else {
			return fmt.Errorf("Unsupported array type %s", v.Type().String())
		}
		addr := ctx.AddToDataSection(bytes)
		if b.address == -1 {
			b.address = addr
		}
	}
	return nil
}
