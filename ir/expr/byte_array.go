package expr

import (
	"fmt"

	. "github.com/bspaans/jit-compiler/ir/shared"
)

type IR_ByteArray struct {
	*BaseIRExpression
	Value []uint8

	// Set during EncodeDataSection
	Address *SegmentPointer
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

func (b *IR_ByteArray) SSA_Transform(ctx *SSA_Context) (SSA_Rewrites, IRExpression) {
	return nil, b
}
