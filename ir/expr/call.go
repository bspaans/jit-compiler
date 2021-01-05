package expr

import (
	"fmt"
	"strings"

	. "github.com/bspaans/jit-compiler/ir/shared"
)

type IR_Call struct {
	*BaseIRExpression
	Function string
	Args     []IRExpression
}

func NewIR_Call(function string, args []IRExpression) *IR_Call {
	return &IR_Call{
		BaseIRExpression: NewBaseIRExpression(Call),
		Function:         function,
		Args:             args,
	}
}

func (i *IR_Call) ReturnType(ctx *IR_Context) Type {
	signature := ctx.VariableTypes[i.Function]
	if signature == nil {
		panic("Unknown function: " + i.Function)
	}
	if _, ok := signature.(*TFunction); !ok {
		panic("Expected function, got: " + signature.String())
	}
	return signature.(*TFunction).ReturnType
}

func (i *IR_Call) String() string {
	args := []string{}
	for _, arg := range i.Args {
		args = append(args, arg.String())
	}
	return fmt.Sprintf("%s(%s)", i.Function, strings.Join(args, ", "))
}

func (b *IR_Call) SSA_Transform(ctx *SSA_Context) (SSA_Rewrites, IRExpression) {
	newArgs := make([]IRExpression, len(b.Args))
	rewrites := SSA_Rewrites{}
	for i, arg := range b.Args {
		if IsLiteralOrVariable(arg) {
			newArgs[i] = arg
		} else {
			rw, expr := arg.SSA_Transform(ctx)
			for _, rewrite := range rw {
				rewrites = append(rewrites, rewrite)
			}
			v := ctx.GenerateVariable()
			rewrites = append(rewrites, NewSSA_Rewrite(v, expr))
			newArgs[i] = NewIR_Variable(v)
		}
	}
	return rewrites, NewIR_Call(b.Function, newArgs)
}
