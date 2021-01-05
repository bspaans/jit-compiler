package x86_64

import (
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/ir/statements"
	"github.com/bspaans/jit-compiler/lib"
)

// Allocates a new register and assigns it the value of the expression.
func encode_IR_FunctionDef(i *statements.IR_FunctionDef, ctx *IR_Context) ([]lib.Instruction, error) {
	reg := ctx.AllocateRegister(TUint64)
	returnType := i.Expr.ReturnType(ctx)
	ctx.VariableTypes[i.Name] = returnType
	ctx.VariableMap[i.Name] = reg
	return encodeExpression(i.Expr, ctx, reg)
}
