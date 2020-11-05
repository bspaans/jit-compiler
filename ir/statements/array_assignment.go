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
	tmpReg := ctx.AllocateRegister(TUint64)
	defer ctx.DeallocateRegister(tmpReg)

	// Calculate the index offset and add the address of
	// the array to it (TODO: only works for quadwords)
	// TODO: check returntype of i.Index
	result, err := i.Index.Encode(ctx, tmpReg)
	if err != nil {
		return nil, err
	}
	reg, found := ctx.VariableMap[i.Variable]
	if !found {
		return nil, fmt.Errorf("Unknown array '%s'", i.Variable)
	}
	add := asm.ADD(reg, tmpReg)
	ctx.AddInstruction(add)
	result = append(result, add)

	// Encode the expression

	returnType := i.Expr.ReturnType(ctx)
	// TODO write directly to location?
	tmpReg2 := ctx.AllocateRegister(returnType)
	defer ctx.DeallocateRegister(tmpReg2)
	expr, err := i.Expr.Encode(ctx, tmpReg2)
	if err != nil {
		return nil, err
	}
	ctx.AddInstructions(expr)
	result = lib.Instructions(result).Add(expr)

	// Move the expr result into the array
	itemWidth := returnType.Width()
	if itemWidth == 1 {
		mov := asm.MOV(tmpReg2.Lower8BitRegister(), &encoding.IndirectRegister{tmpReg.Lower8BitRegister()})
		ctx.AddInstruction(mov)
		result = append(result, mov)
	} else if itemWidth == 8 {
		mov := asm.MOV(tmpReg2, &encoding.IndirectRegister{tmpReg})
		ctx.AddInstruction(mov)
		result = append(result, mov)

	} else {
		return nil, fmt.Errorf("Assigning to non 64 bit arrays not supported at this time [TODO]")
	}
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
