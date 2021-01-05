package expr

import (
	"strings"

	. "github.com/bspaans/jit-compiler/ir/shared"
)

type IR_Struct struct {
	*BaseIRExpression
	StructType *TStruct
	Values     []IRExpression // only literals are supported
	// Set during EncodeDataSection
	Address int
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

func (b *IR_Struct) SSA_Transform(ctx *SSA_Context) (SSA_Rewrites, IRExpression) {
	return nil, b
}
