package expr

import (
	"fmt"
	"strings"

	"github.com/bspaans/jit-compiler/asm/x86_64"
	"github.com/bspaans/jit-compiler/asm/x86_64/encoding"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

type IR_StaticArray struct {
	*BaseIRExpression
	ElemType Type
	Value    []IRExpression // only literals are supported
	address  int
}

func NewIR_StaticArray(elemType Type, value []IRExpression) *IR_StaticArray {
	return &IR_StaticArray{
		BaseIRExpression: NewBaseIRExpression(StaticArray),
		ElemType:         elemType,
		Value:            value,
	}
}

func (i *IR_StaticArray) ReturnType(ctx *IR_Context) Type {
	return &TArray{i.ElemType, len(i.Value)}
}

func (i *IR_StaticArray) String() string {
	elems := []string{}
	for _, v := range i.Value {
		elems = append(elems, v.String())
	}
	return fmt.Sprintf("[]%s{%s}", i.ElemType.String(), strings.Join(elems, ", "))
}

func (i *IR_StaticArray) Encode(ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {
	// Calculate the displacement between RIP (the instruction pointer,
	// pointing to the *next* instruction) and the address of our byte array,
	// and load the resulting address into target using a LEA instruction.
	ownLength := uint(7)
	diff := uint(ctx.InstructionPointer+ownLength) - uint(i.address)
	result := []lib.Instruction{asm.LEA(&encoding.RIPRelative{encoding.Int32(int32(-diff))}, target)}
	ctx.AddInstructions(result)
	return result, nil
}

func (b *IR_StaticArray) AddToDataSection(ctx *IR_Context) error {
	b.address = -1
	for _, v := range b.Value {
		bytes := []uint8{}
		if b.ElemType == TUint8 {
			switch v.(type) {
			case *IR_Uint8:
				bytes = []uint8{uint8(v.(*IR_Uint8).Value)}
			case *IR_Uint64:
				bytes = []uint8{uint8(v.(*IR_Uint64).Value)}
			default:
				return fmt.Errorf("Unsupport uint8 array type %s in %s", b.ElemType, b.String())
			}
		} else if b.ElemType == TUint16 {
			switch v.(type) {
			case *IR_Uint16:
				bytes = encoding.Uint16(v.(*IR_Uint16).Value).Encode()
			case *IR_Uint64:
				bytes = encoding.Uint16(uint16(v.(*IR_Uint64).Value)).Encode()
			default:
				return fmt.Errorf("Unsupport uint16 array type %s in %s", b.ElemType, b.String())
			}
		} else if b.ElemType == TUint32 {
			switch v.(type) {
			case *IR_Uint32:
				bytes = encoding.Uint32(v.(*IR_Uint32).Value).Encode()
			case *IR_Uint64:
				bytes = encoding.Uint32(uint32(v.(*IR_Uint64).Value)).Encode()
			default:
				return fmt.Errorf("Unsupport uint32 array type %s in %s", b.ElemType, b.String())
			}
		} else if b.ElemType == TUint64 {
			ir := v.(*IR_Uint64)
			bytes = encoding.Uint64(ir.Value).Encode()
		} else if b.ElemType == TFloat64 {
			ir := v.(*IR_Float64)
			bytes = encoding.Float64(ir.Value).Encode()
		} else if b.ElemType == TInt64 {
			ir := v.(*IR_Int64)
			bytes = encoding.Uint64(ir.Value).Encode()
		} else {
			return fmt.Errorf("Unsupported array type %s", v.Type().String())
		}
		addr := ctx.AddToDataSection(bytes)
		if b.address == -1 {
			b.address = addr
		}
	}
	return nil
}

func (b *IR_StaticArray) SSA_Transform(ctx *SSA_Context) (SSA_Rewrites, IRExpression) {
	return nil, b
}
