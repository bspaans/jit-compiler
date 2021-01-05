package x86_64

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm/x86_64/encoding"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/ir/statements"
	"github.com/bspaans/jit-compiler/lib"
)

// Allocates a new register and assigns it the value of the expression.
func encode_IR_Assignment(i *statements.IR_Assignment, ctx *IR_Context) ([]lib.Instruction, error) {
	ctx.AddInstruction("assignment " + encoding.Comment(i.String()))
	returnType := i.Expr.ReturnType(ctx)
	reg, found := ctx.VariableMap[i.Variable]
	if !found {
		reg = ctx.AllocateRegister(returnType)
		ctx.VariableMap[i.Variable] = reg
		ctx.VariableTypes[i.Variable] = returnType
	}
	expr, err := encodeExpression(i.Expr, ctx, reg)
	if err != nil {
		return nil, fmt.Errorf("Error in assignment: %s", err.Error())
	}
	return expr, nil
}
