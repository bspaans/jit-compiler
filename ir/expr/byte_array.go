package expr

import (
	"fmt"

	"github.com/bspaans/jit/asm"
	"github.com/bspaans/jit/asm/encoding"
	. "github.com/bspaans/jit/ir/shared"
	"github.com/bspaans/jit/lib"
)

type IR_ByteArray struct {
	*BaseIRExpression
	Value   []uint8
	address int
}

func NewIR_ByteArray(value []uint8) *IR_ByteArray {
	return &IR_ByteArray{
		BaseIRExpression: NewBaseIRExpression(ByteArray),
		Value:            value,
	}
}

func (i *IR_ByteArray) ReturnType(ctx *IR_Context) Type {
	return &TArray{TUint8, len(i.Value)}
}

func (i *IR_ByteArray) String() string {
	return fmt.Sprintf("%v", i.Value)
}

func (i *IR_ByteArray) Encode(ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {
	// Calculate the displacement between RIP (the instruction pointer,
	// pointing to the *next* instruction) and the address of our byte array,
	// and load the resulting address into target using a LEA instruction.
	fmt.Println("Loading from ", i.address, ctx.InstructionPointer)
	ownLength := uint(7)
	diff := uint(ctx.InstructionPointer+ownLength) - uint(i.address)
	result := []lib.Instruction{asm.LEA(&encoding.RIPRelative{encoding.Int32(int32(-diff))}, target)}
	ctx.AddInstructions(result)
	return result, nil
}

func (b *IR_ByteArray) AddToDataSection(ctx *IR_Context) error {
	b.address = ctx.AddToDataSection(b.Value)
	return nil
}
