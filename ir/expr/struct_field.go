package expr

import (
	. "github.com/bspaans/jit-compiler/ir/shared"
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

func (b *IR_StructField) AddToDataSection(ctx *IR_Context) error {
	return b.Struct.AddToDataSection(ctx)
}

func (b *IR_StructField) SSA_Transform(ctx *SSA_Context) (SSA_Rewrites, IRExpression) {
	if IsLiteralOrVariable(b.Struct) {
		return nil, b
	}
	rewrites, expr := b.Struct.SSA_Transform(ctx)
	v := ctx.GenerateVariable()
	rewrites = append(rewrites, NewSSA_Rewrite(v, expr))
	return rewrites, NewIR_StructField(NewIR_Variable(v), b.Field)
}
