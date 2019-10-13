package expr

import (
	"fmt"

	"github.com/bspaans/jit/asm"
	"github.com/bspaans/jit/asm/encoding"
	. "github.com/bspaans/jit/ir/shared"
	"github.com/bspaans/jit/lib"
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
	if ty.Type() != T_Array {
		panic("Not an array")
	}
	return ty.(*TArray).ItemType
}

func (i *IR_ArrayIndex) String() string {
	return fmt.Sprintf("%s[%s]", i.Array.String(), i.Index.String())
}

func (i *IR_ArrayIndex) Encode(ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {
	// tmpReg will contain the address of the array
	tmpReg := ctx.AllocateRegister(TUint64)
	defer ctx.DeallocateRegister(tmpReg)
	result, err := i.Array.Encode(ctx, tmpReg)
	if err != nil {
		return nil, err
	}
	// TODO: if i.Index == number => specialise
	ix, err := i.Index.Encode(ctx, target)
	if err != nil {
		return nil, err
	}
	result = lib.Instructions(result).Add(ix)
	itemWidth := i.ReturnType(ctx).Width()
	if itemWidth != 1 {
		mul := asm.MUL(encoding.Uint32(itemWidth), target)
		ctx.AddInstruction(mul)
		result = append(result, mul)
	}
	instr := lib.Instructions{
		asm.ADD(target, tmpReg),
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
