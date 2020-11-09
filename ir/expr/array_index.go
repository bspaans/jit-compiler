package expr

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm"
	"github.com/bspaans/jit-compiler/asm/encoding"
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

	arrayReg := ctx.AllocateRegister(TUint64)
	defer ctx.DeallocateRegister(arrayReg)
	indexReg := ctx.AllocateRegister(TUint64)
	defer ctx.DeallocateRegister(indexReg)

	itemWidth := i.ReturnType(ctx).Width()

	result, err := i.Array.Encode(ctx, arrayReg)
	if err != nil {
		return nil, fmt.Errorf("Array encoding issue: %s", err.Error())
	}

	index, err := i.Index.Encode(ctx, indexReg)
	if err != nil {
		return nil, fmt.Errorf("Array index encoding issue: %s", err.Error())
	}
	result = lib.Instructions(result).Add(index)
	mov0 := asm.MOV(encoding.Uint64(0), target)
	mov := asm.MOV(&encoding.SIBRegister{arrayReg, indexReg, encoding.ScaleForItemWidth(itemWidth)},
		target.(*encoding.Register).ForOperandWidth(itemWidth))
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
