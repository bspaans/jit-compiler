package statements

import (
	"fmt"
	"strings"

	"github.com/bspaans/jit/ir/expr"
	. "github.com/bspaans/jit/ir/shared"
	"github.com/bspaans/jit/lib"
)

type IR_FunctionDef struct {
	*BaseIR
	Name    string
	Expr    *expr.IR_Function
	Address int
}

func NewIR_FunctionDef(name string, expr *expr.IR_Function) *IR_FunctionDef {
	return &IR_FunctionDef{
		BaseIR: NewBaseIR(FunctionDef),
		Name:   name,
		Expr:   expr,
	}
}

// Allocates a new register and assigns it the value of the expression.
func (i *IR_FunctionDef) Encode(ctx *IR_Context) ([]lib.Instruction, error) {
	reg := ctx.AllocateRegister(TUint64)
	returnType := i.Expr.ReturnType(ctx)
	ctx.VariableTypes[i.Name] = returnType
	ctx.VariableMap[i.Name] = reg
	return i.Expr.Encode(ctx, reg)
}

func (i *IR_FunctionDef) String() string {
	args := []string{}
	for j, arg := range i.Expr.Signature.ArgNames {
		args = append(args, arg+" "+i.Expr.Signature.Args[j].String())
	}
	return fmt.Sprintf("func %s(%s) %s { %s }", i.Name, strings.Join(args, ", "), i.Expr.Signature.ReturnType.String(), i.Expr.Body.String())
}

func (i *IR_FunctionDef) AddToDataSection(ctx *IR_Context) error {
	return i.Expr.AddToDataSection(ctx)
}
