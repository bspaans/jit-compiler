package x86_64

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm/x86_64"
	"github.com/bspaans/jit-compiler/asm/x86_64/encoding"
	"github.com/bspaans/jit-compiler/ir/expr"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

func encode_IR_StaticArray(i *expr.IR_StaticArray, ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {
	// Calculate the displacement between RIP (the instruction pointer,
	// pointing to the *next* instruction) and the address of our byte array,
	// and load the resulting address into target using a LEA instruction.
	ownLength := uint(7)
	// TODO: instead of address + 2; use datasectionOffset + len(ctx.ReadonlySegment) + address; in JIT mode
	diff := uint(ctx.InstructionPointer+ownLength) - uint(ctx.Segments.GetAddress(i.Address))
	result := []lib.Instruction{x86_64.LEA(&encoding.RIPRelative{encoding.Int32(int32(-diff))}, target)}
	ctx.AddInstructions(result)
	return result, nil
}

func encode_IR_StaticArray_for_DataSection(b *expr.IR_StaticArray, segments *Segments) error {
	b.Address = nil
	for _, v := range b.Value {
		bytes := []uint8{}
		if b.ElemType == TUint8 {
			switch c := v.(type) {
			case *expr.IR_Uint8:
				bytes = []uint8{uint8(c.Value)}
			case *expr.IR_Uint64:
				bytes = []uint8{uint8(c.Value)}
			default:
				return fmt.Errorf("Unsupport uint8 array type %s in %s", b.ElemType, b.String())
			}
		} else if b.ElemType == TUint16 {
			switch c := v.(type) {
			case *expr.IR_Uint16:
				bytes = encoding.Uint16(c.Value).Encode()
			case *expr.IR_Uint64:
				bytes = encoding.Uint16(uint16(c.Value)).Encode()
			default:
				return fmt.Errorf("Unsupport uint16 array type %s in %s", b.ElemType, b.String())
			}
		} else if b.ElemType == TUint32 {
			switch c := v.(type) {
			case *expr.IR_Uint32:
				bytes = encoding.Uint32(c.Value).Encode()
			case *expr.IR_Uint64:
				bytes = encoding.Uint32(uint32(c.Value)).Encode()
			default:
				return fmt.Errorf("Unsupport uint32 array type %s in %s", b.ElemType, b.String())
			}
		} else if b.ElemType == TUint64 {
			ir := v.(*expr.IR_Uint64)
			bytes = encoding.Uint64(ir.Value).Encode()
		} else if b.ElemType == TFloat64 {
			ir := v.(*expr.IR_Float64)
			bytes = encoding.Float64(ir.Value).Encode()
		} else if b.ElemType == TInt64 {
			ir := v.(*expr.IR_Int64)
			bytes = encoding.Uint64(ir.Value).Encode()
		} else {
			return fmt.Errorf("Unsupported array type %s", v.Type().String())
		}
		addr := segments.Add(ReadWrite, bytes...)
		if b.Address == nil {
			b.Address = addr
		}
	}
	return nil
}
