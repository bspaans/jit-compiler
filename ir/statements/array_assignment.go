package statements

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm"
	"github.com/bspaans/jit-compiler/asm/encoding"
	"github.com/bspaans/jit-compiler/ir/expr"
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

func (i *IR_ArrayAssignment) encodeIndex(ctx *IR_Context, arrayReg encoding.Operand) ([]lib.Instruction, error) {

	// Optimisation: if we're getting a[0] we don't have to do anything.
	// The address in arrayReg will already be correct
	if ix, ok := i.Index.(*expr.IR_Uint64); ok && ix.Value == 0 {
		return []lib.Instruction{}, nil
	}

	indexReg := ctx.AllocateRegister(TUint64)
	defer ctx.DeallocateRegister(indexReg)

	// Calculate the index offset
	result, err := i.Index.Encode(ctx, indexReg)
	if err != nil {
		return nil, err
	}

	// If the item width is not 1 byte wide we need to scale up the
	// index (TODO: can we use SIB encoding for this?)
	returnType := i.Expr.ReturnType(ctx)
	itemWidth := returnType.Width()
	if itemWidth != 1 {
		mulReg := ctx.AllocateRegister(TUint64)
		defer ctx.DeallocateRegister(mulReg)
		mov := asm.MOV_immediate(uint64(itemWidth), mulReg)
		mul := asm.MUL(mulReg, indexReg)
		ctx.AddInstruction(mov)
		ctx.AddInstruction(mul)
		result = append(result, mov)
		result = append(result, mul)
	}
	add := asm.ADD(indexReg, arrayReg)
	ctx.AddInstruction(add)
	result = append(result, add)
	return result, err
}

func (i *IR_ArrayAssignment) encodeExpr(ctx *IR_Context, arrayReg *encoding.Register) ([]lib.Instruction, error) {

	returnType := i.Expr.ReturnType(ctx)
	itemWidth := returnType.Width()

	exprReg := ctx.AllocateRegister(returnType)
	defer ctx.DeallocateRegister(exprReg)

	// TODO write directly to location?
	result, err := i.Expr.Encode(ctx, exprReg)
	if err != nil {
		return nil, err
	}

	// Move the expr result into the array
	if itemWidth == 1 {
		mov := asm.MOV(exprReg.Lower8BitRegister(), &encoding.IndirectRegister{arrayReg.Lower8BitRegister()})
		ctx.AddInstruction(mov)
		result = append(result, mov)
	} else if itemWidth == 8 {
		mov := asm.MOV(exprReg, &encoding.IndirectRegister{arrayReg})
		ctx.AddInstruction(mov)
		result = append(result, mov)

	} else {
		return nil, fmt.Errorf("Assigning to arrays of type %s is not supported at this time [TODO]", returnType)
	}
	return result, nil
}

func (i *IR_ArrayAssignment) Encode(ctx *IR_Context) ([]lib.Instruction, error) {
	ctx.AddInstruction("array_assignment " + encoding.Comment(i.String()))

	arrayReg := ctx.AllocateRegister(TUint64)
	defer ctx.DeallocateRegister(arrayReg)

	// Move the address of the Array to the register
	reg, found := ctx.VariableMap[i.Variable]
	if !found {
		return nil, fmt.Errorf("Unknown array '%s'", i.Variable)
	}
	mov := asm.MOV(reg, arrayReg)
	ctx.AddInstruction(mov)
	result := []lib.Instruction{mov}

	instr, err := i.encodeIndex(ctx, arrayReg)
	if err != nil {
		return nil, fmt.Errorf("Failed to encode array index in %s: %s", i.String(), err.Error())
	}
	result = lib.Instructions(result).Add(instr)

	exprInstr, err := i.encodeExpr(ctx, arrayReg)
	if err != nil {
		return nil, fmt.Errorf("Failed to encode expr in %s: %s", i.String(), err.Error())
	}
	result = lib.Instructions(result).Add(exprInstr)
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
