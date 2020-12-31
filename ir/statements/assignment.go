package statements

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm/x86_64/encoding"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

type IR_Assignment struct {
	*BaseIR
	Variable string
	Expr     IRExpression
}

func NewIR_Assignment(variable string, expr IRExpression) *IR_Assignment {
	return &IR_Assignment{
		BaseIR:   NewBaseIR(Assignment),
		Variable: variable,
		Expr:     expr,
	}
}

// Allocates a new register and assigns it the value of the expression.
func (i *IR_Assignment) Encode(ctx *IR_Context) ([]lib.Instruction, error) {
	ctx.AddInstruction("assignment " + encoding.Comment(i.String()))
	returnType := i.Expr.ReturnType(ctx)
	reg, found := ctx.VariableMap[i.Variable]
	if !found {
		reg = ctx.AllocateRegister(returnType)
		ctx.VariableMap[i.Variable] = reg
		ctx.VariableTypes[i.Variable] = returnType
	}
	expr, err := i.Expr.Encode(ctx, reg)
	if err != nil {
		return nil, fmt.Errorf("Error in assignment: %s", err.Error())
	}
	return expr, nil
}

func (i *IR_Assignment) String() string {
	return fmt.Sprintf("%s = %s", i.Variable, i.Expr.String())
}

func (i *IR_Assignment) AddToDataSection(ctx *IR_Context) error {
	return i.Expr.AddToDataSection(ctx)
}

func (i *IR_Assignment) SSA_Transform(ctx *SSA_Context) IR {
	rewrites, expr := i.Expr.SSA_Transform(ctx)
	ir := SSA_Rewrites_to_IR(rewrites)
	if ir == nil {
		return i
	}
	return NewIR_AndThen(ir, NewIR_Assignment(i.Variable, expr))
}
