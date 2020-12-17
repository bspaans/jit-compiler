package statements

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm"
	"github.com/bspaans/jit-compiler/asm/encoding"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

type IR_ArrayAssignment struct {
	*BaseIR
	Variable string
	Index    IRExpression
	Expr     IRExpression
}

func NewIR_ArrayAssignment(variable string, index IRExpression, expr IRExpression) *IR_ArrayAssignment {
	return &IR_ArrayAssignment{
		BaseIR:   NewBaseIR(ArrayAssignment),
		Variable: variable,
		Index:    index,
		Expr:     expr,
	}
}

func (i *IR_ArrayAssignment) Encode(ctx *IR_Context) ([]lib.Instruction, error) {
	ctx.AddInstruction("array_assignment " + encoding.Comment(i.String()))

	returnType := i.Expr.ReturnType(ctx)
	itemWidth := returnType.Width()

	indexReg := ctx.AllocateRegister(TUint64)
	defer ctx.DeallocateRegister(indexReg)
	exprReg := ctx.AllocateRegister(returnType)
	defer ctx.DeallocateRegister(exprReg)

	reg, found := ctx.VariableMap[i.Variable]
	if !found {
		return nil, fmt.Errorf("Unknown array '%s'", i.Variable)
	}

	result, err := i.Index.Encode(ctx, indexReg)
	if err != nil {
		return nil, fmt.Errorf("Failed to encode array index in %s: %s", i.String(), err.Error())
	}

	exprInstr, err := i.Expr.Encode(ctx, exprReg)
	if err != nil {
		return nil, fmt.Errorf("Failed to encode expr in %s: %s", i.String(), err.Error())
	}
	result = lib.Instructions(result).Add(exprInstr)

	target := &encoding.SIBRegister{reg.(*encoding.Register), indexReg, encoding.ScaleForItemWidth(itemWidth)}
	mov := asm.MOV(exprReg.ForOperandWidth(itemWidth), target)
	ctx.AddInstruction(mov)
	result = append(result, mov)
	return result, nil
}

func (i *IR_ArrayAssignment) String() string {
	return fmt.Sprintf("%s[%s] = %s", i.Variable, i.Index.String(), i.Expr.String())
}

func (i *IR_ArrayAssignment) AddToDataSection(ctx *IR_Context) error {
	if err := i.Index.AddToDataSection(ctx); err != nil {
		return err
	}
	return i.Expr.AddToDataSection(ctx)
}

func (i *IR_ArrayAssignment) SSA_Transform(ctx *SSA_Context) IR {
	return i
}
