package expr

import (
	"fmt"

	. "github.com/bspaans/jit-compiler/ir/shared"
)

type IR_Mul struct {
	*BaseIRExpression
	Op1 IRExpression
	Op2 IRExpression
}

func NewIR_Mul(op1, op2 IRExpression) *IR_Mul {
	return &IR_Mul{
		BaseIRExpression: NewBaseIRExpression(Mul),
		Op1:              op1,
		Op2:              op2,
	}
}

func (i *IR_Mul) ReturnType(ctx *IR_Context) Type {
	return i.Op1.ReturnType(ctx)
}

func (i *IR_Mul) String() string {
	return fmt.Sprintf("%s * %s", i.Op1.String(), i.Op2.String())
}

func (b *IR_Mul) SSA_Transform(ctx *SSA_Context) (SSA_Rewrites, IRExpression) {
	if IsLiteralOrVariable(b.Op1) {
		if IsLiteralOrVariable(b.Op2) {
			return nil, b
		} else {
			rewrites, expr := b.Op2.SSA_Transform(ctx)
			v := ctx.GenerateVariable()
			rewrites = append(rewrites, NewSSA_Rewrite(v, expr))
			return rewrites, NewIR_Mul(b.Op1, NewIR_Variable(v))
		}
	} else {
		rewrites, expr := b.Op1.SSA_Transform(ctx)
		v := ctx.GenerateVariable()
		rewrites = append(rewrites, NewSSA_Rewrite(v, expr))
		if IsLiteralOrVariable(b.Op2) {
			return rewrites, NewIR_Mul(NewIR_Variable(v), b.Op2)
		} else {
			rewrites2, expr2 := b.Op2.SSA_Transform(ctx)
			for _, rw := range rewrites2 {
				rewrites = append(rewrites, rw)
			}
			v2 := ctx.GenerateVariable()
			rewrites = append(rewrites, NewSSA_Rewrite(v2, expr2))
			return rewrites, NewIR_Mul(NewIR_Variable(v), NewIR_Variable(v2))
		}
	}
}
