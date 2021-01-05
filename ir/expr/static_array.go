package expr

import (
	"fmt"
	"strings"

	. "github.com/bspaans/jit-compiler/ir/shared"
)

type IR_StaticArray struct {
	*BaseIRExpression
	ElemType Type
	Value    []IRExpression // only literals are supported
	// Set during EncodeDataSection
	Address *SegmentPointer
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

func (b *IR_StaticArray) SSA_Transform(ctx *SSA_Context) (SSA_Rewrites, IRExpression) {
	return nil, b
}
