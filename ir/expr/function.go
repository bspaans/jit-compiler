package expr

import (
	"fmt"
	"strings"

	. "github.com/bspaans/jit-compiler/ir/shared"
)

type IR_Function struct {
	*BaseIRExpression
	Signature *TFunction
	Body      IR
	Address   int
}

func NewIR_Function(signature *TFunction, body IR) *IR_Function {
	return &IR_Function{
		BaseIRExpression: NewBaseIRExpression(Function),
		Signature:        signature,
		Body:             body,
	}
}

func (i *IR_Function) ReturnType(ctx *IR_Context) Type {
	return i.Signature
}

func (i *IR_Function) String() string {
	args := []string{}
	for j, arg := range i.Signature.ArgNames {
		args = append(args, arg+" "+i.Signature.Args[j].String())
	}
	return fmt.Sprintf("func(%s) %s { %s }", strings.Join(args, ", "), i.Signature.ReturnType.String(), i.Body.String())
}

func (b *IR_Function) SSA_Transform(ctx *SSA_Context) (SSA_Rewrites, IRExpression) {
	newBody := b.Body.SSA_Transform(ctx)
	return nil, NewIR_Function(b.Signature, newBody)
}
