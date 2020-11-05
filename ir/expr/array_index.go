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
	// tmpReg will contain the address of the array
	tmpReg := ctx.AllocateRegister(TUint64)
	defer ctx.DeallocateRegister(tmpReg)
	result, err := i.Array.Encode(ctx, tmpReg)
	if err != nil {
		return nil, err
	}

	// TODO: if i.Index == number => specialise
	// SIB encoding?
	tmpIndexReg := ctx.AllocateRegister(TUint64)
	defer ctx.DeallocateRegister(tmpIndexReg)
	ix, err := i.Index.Encode(ctx, tmpIndexReg)
	if err != nil {
		return nil, err
	}
	result = lib.Instructions(result).Add(ix)
	returnType := i.ReturnType(ctx)
	if returnType == nil {
		return nil, fmt.Errorf("Return type nil. wut")
	}
	itemWidth := i.ReturnType(ctx).Width()
	if itemWidth != 1 {
		tmpReg3 := ctx.AllocateRegister(TUint64)
		defer ctx.DeallocateRegister(tmpReg3)
		mov := asm.MOV(encoding.Uint32(itemWidth), tmpReg3)
		mul := asm.MUL(tmpReg3, tmpIndexReg)
		ctx.AddInstruction(mov)
		ctx.AddInstruction(mul)
		result = append(result, mov)
		result = append(result, mul)
	}
	instr := lib.Instructions{
		asm.ADD(tmpIndexReg, tmpReg),
	}
	if itemWidth == 1 {
		var tmpReg2 *encoding.Register
		if target.Type() == encoding.T_Register {
			tmpReg2 = target.(*encoding.Register)
		} else {
			tmpReg2 = ctx.AllocateRegister(TUint64)
			defer ctx.DeallocateRegister(tmpReg2)
		}

		instr = append(instr, asm.MOV(encoding.Uint64(0), tmpReg2))
		// TODO: replace Lower8BitRegister with GetRegisterForWidth(itemWidth)
		instr = append(instr, asm.MOV(&encoding.IndirectRegister{tmpReg.Lower8BitRegister()}, tmpReg2.Lower8BitRegister()))
		instr = append(instr, asm.MOV(tmpReg2, target))
	} else {
		instr = append(instr, asm.MOV(&encoding.IndirectRegister{tmpReg}, target))
	}
	ctx.AddInstructions(instr)
	return lib.Instructions(result).Add(instr), nil
}

func (b *IR_ArrayIndex) AddToDataSection(ctx *IR_Context) error {
	if err := b.Array.AddToDataSection(ctx); err != nil {
		return err
	}
	return b.Index.AddToDataSection(ctx)
}
