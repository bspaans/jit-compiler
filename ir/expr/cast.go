package expr

import (
	"fmt"

	. "github.com/bspaans/jit-compiler/ir/shared"
)

type IR_Cast struct {
	*BaseIRExpression
	Value      IRExpression
	CastToType Type
}

func NewIR_Cast(value IRExpression, typ Type) *IR_Cast {
	return &IR_Cast{
		BaseIRExpression: NewBaseIRExpression(Cast),
		Value:            value,
		CastToType:       typ,
	}
}

func (i *IR_Cast) ReturnType(ctx *IR_Context) Type {
	return i.CastToType
}

func (i *IR_Cast) String() string {
	return fmt.Sprintf("%s(%s)", i.CastToType.String(), i.Value.String())
}

func (b *IR_Cast) SSA_Transform(ctx *SSA_Context) (SSA_Rewrites, IRExpression) {
	if IsLiteralOrVariable(b.Value) {
		return nil, b
	}
	rewrites, expr := b.Value.SSA_Transform(ctx)
	v := ctx.GenerateVariable()
	rewrites = append(rewrites, NewSSA_Rewrite(v, expr))
	return rewrites, NewIR_Cast(NewIR_Variable(v), b.CastToType)
}
