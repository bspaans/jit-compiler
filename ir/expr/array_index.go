package expr

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm"
	"github.com/bspaans/jit-compiler/asm/encoding"
	"github.com/bspaans/jit-compiler/ir/shared"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

type IR_ArrayIndex struct {
	*BaseIRExpression
	Array IRExpression
	Index IRExpression
}

func NewIR_ArrayIndex(array, index IRExpression) *IR_ArrayIndex {
	return &IR_ArrayIndex{
		BaseIRExpression: NewBaseIRExpression(ArrayIndex),
		Array:            array,
		Index:            index,
	}
}

func (i *IR_ArrayIndex) ReturnType(ctx *IR_Context) Type {
	ty := i.Array.ReturnType(ctx)
	if ty == nil {
		panic("Type is nil")
	}
	if ty.Type() != T_Array {
		panic("Not an array")
	}
	return ty.(*TArray).ItemType
}

func (i *IR_ArrayIndex) String() string {
	return fmt.Sprintf("%s[%s]", i.Array.String(), i.Index.String())
}

func (i *IR_ArrayIndex) Encode(ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {
	ctx.AddInstruction("array_index " + encoding.Comment(i.String()))

	itemWidth := i.ReturnType(ctx).Width()

	var arrayReg *encoding.Register
	var indexReg *encoding.Register
	if i.Array.Type() == shared.Variable {
		variable := i.Array.(*IR_Variable).Value
		reg, ok := ctx.VariableMap[variable]
		if !ok {
			return nil, fmt.Errorf("Unknown variable %s", variable)
		}
		arrayReg = reg.(*encoding.Register)
	} else {
		arrayReg = ctx.AllocateRegister(TUint64)
		defer ctx.DeallocateRegister(arrayReg)
	}

	result, err := i.Array.Encode(ctx, arrayReg)
	if err != nil {
		return nil, fmt.Errorf("Array encoding issue: %s", err.Error())
	}

	// Specialise for integers
	if i.Index.Type() == shared.Uint64 {
		op := i.Index.(*IR_Uint64)
		if op.Value == 0 {
			mov := asm.MOV(&encoding.IndirectRegister{arrayReg.ForOperandWidth(itemWidth)}, target.(*encoding.Register).ForOperandWidth(itemWidth))
			// Move 0 into target register if going from a wider to narrower register
			mov0 := asm.MOV(encoding.Uint64(0), target.(*encoding.Register).Get64BitRegister()) // TODO use xor reg, reg
			if itemWidth < lib.QUADWORD {
				ctx.AddInstruction(mov0)
				result = append(result, mov0)
			}
			ctx.AddInstruction(mov)
			result = append(result, mov)
			return result, nil
		} else {
			// TODO add index*itemwidth to arrayReg
		}
	}

	if i.Index.Type() == shared.Variable {
		variable := i.Index.(*IR_Variable).Value
		reg, ok := ctx.VariableMap[variable]
		if !ok {
			return nil, fmt.Errorf("Unknown variable %s", variable)
		}
		indexReg = reg.(*encoding.Register)
	} else {
		indexReg = ctx.AllocateRegister(TUint64)
		defer ctx.DeallocateRegister(indexReg)
	}

	index, err := i.Index.Encode(ctx, indexReg)
	if err != nil {
		return nil, fmt.Errorf("Array index encoding issue: %s", err.Error())
	}
	result = lib.Instructions(result).Add(index)
	mov := asm.MOV(&encoding.SIBRegister{arrayReg, indexReg, encoding.ScaleForItemWidth(itemWidth)},
		target.(*encoding.Register).ForOperandWidth(itemWidth))

	// Move 0 into target register if going from a wider to narrower register
	mov0 := asm.MOV(encoding.Uint64(0), target.(*encoding.Register).Get64BitRegister())
	if itemWidth < lib.QUADWORD {
		ctx.AddInstruction(mov0)
		result = append(result, mov0)
	}

	ctx.AddInstruction(mov)
	result = append(result, mov)
	return result, nil
}

func (b *IR_ArrayIndex) AddToDataSection(ctx *IR_Context) error {
	if err := b.Array.AddToDataSection(ctx); err != nil {
		return err
	}
	return b.Index.AddToDataSection(ctx)
}
