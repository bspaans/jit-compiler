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

func (i *IR_ArrayIndex) encodeIndex(ctx *IR_Context, arrayReg encoding.Operand) ([]lib.Instruction, error) {

	// Optimisation: if we're getting a[0] we don't have to do anything.
	// The address in arrayReg will already be correct
	if ix, ok := i.Index.(*IR_Uint64); ok && ix.Value == 0 {
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
	// index (TODO: can we use SIB or displacement encoding for this?)
	itemWidth := i.ReturnType(ctx).Width()
	if itemWidth != 1 {
		mulReg := ctx.AllocateRegister(TUint64)
		defer ctx.DeallocateRegister(mulReg)
		mov := asm.MOV(encoding.Uint32(itemWidth), mulReg)
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

func (i *IR_ArrayIndex) encodeMove(ctx *IR_Context, arrayReg, target encoding.Operand) ([]lib.Instruction, error) {

	itemWidth := i.ReturnType(ctx).Width()

	if itemWidth == 1 {
		var tmpReg2 *encoding.Register
		if target.Type() == encoding.T_Register {
			tmpReg2 = target.(*encoding.Register)
		} else {
			tmpReg2 = ctx.AllocateRegister(TUint64)
			defer ctx.DeallocateRegister(tmpReg2)
		}

		mov0 := asm.MOV(encoding.Uint64(0), tmpReg2)
		mov := asm.MOV(&encoding.IndirectRegister{arrayReg.(*encoding.Register).Lower8BitRegister()}, tmpReg2.Lower8BitRegister())
		movTarget := asm.MOV(tmpReg2, target)
		ctx.AddInstruction(mov0)
		ctx.AddInstruction(mov)
		if tmpReg2 == target {
			return []lib.Instruction{mov0, mov}, nil
		}
		ctx.AddInstruction(movTarget)
		return []lib.Instruction{mov0, mov, movTarget}, nil
	} else {
		mov := asm.MOV(&encoding.IndirectRegister{arrayReg.(*encoding.Register)}, target)
		ctx.AddInstruction(mov)
		return []lib.Instruction{mov}, nil
	}
	return nil, nil
}

func (i *IR_ArrayIndex) Encode(ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {
	ctx.AddInstruction("array_index " + encoding.Comment(i.String()))

	arrayReg := ctx.AllocateRegister(TUint64)
	defer ctx.DeallocateRegister(arrayReg)

	result, err := i.Array.Encode(ctx, arrayReg)
	if err != nil {
		return nil, err
	}

	index, err := i.encodeIndex(ctx, arrayReg)
	if err != nil {
		return nil, err
	}
	result = lib.Instructions(result).Add(index)
	mov, err := i.encodeMove(ctx, arrayReg, target)
	if err != nil {
		return nil, err
	}
	result = lib.Instructions(result).Add(mov)
	return result, nil
}

func (b *IR_ArrayIndex) AddToDataSection(ctx *IR_Context) error {
	if err := b.Array.AddToDataSection(ctx); err != nil {
		return err
	}
	return b.Index.AddToDataSection(ctx)
}
