package expr

import (
	"fmt"
	"strings"

	"github.com/bspaans/jit/asm"
	"github.com/bspaans/jit/asm/encoding"
	. "github.com/bspaans/jit/ir/shared"
	"github.com/bspaans/jit/lib"
)

type IR_StaticArray struct {
	*BaseIRExpression
	ElemType Type
	Value    []encoding.Value
	address  int
}

func NewIR_StaticArray(elemType Type, value []encoding.Value) *IR_StaticArray {
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
	return fmt.Sprintf("[%s]", strings.Join(elems, ", "))
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
		bytes := v.Encode()
		addr := ctx.AddToDataSection(bytes)
		if b.address == -1 {
			b.address = addr
		}
	}
	return nil
}
