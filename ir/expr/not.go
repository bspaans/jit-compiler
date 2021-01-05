package expr

import (
	"fmt"

	. "github.com/bspaans/jit-compiler/ir/shared"
)

type IR_Not struct {
	*BaseIRExpression
	Op1 IRExpression
}

func NewIR_Not(op1 IRExpression) *IR_Not {
	return &IR_Not{
		BaseIRExpression: NewBaseIRExpression(Not),
		Op1:              op1,
	}
}

func (i *IR_Not) ReturnType(ctx *IR_Context) Type {
	return TBool
}

func (i *IR_Not) String() string {
	return fmt.Sprintf("!(%s)", i.Op1.String())
}

func (b *IR_Not) SSA_Transform(ctx *SSA_Context) (SSA_Rewrites, IRExpression) {
	if IsLiteralOrVariable(b.Op1) {
		return nil, b
	}
	rewrites, expr := b.Op1.SSA_Transform(ctx)
	v := ctx.GenerateVariable()
	rewrites = append(rewrites, NewSSA_Rewrite(v, expr))
	return rewrites, NewIR_Not(NewIR_Variable(v))
}
