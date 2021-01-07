package encoding

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/bspaans/jit-compiler/lib"
)

type Opcode struct {
	Name     string
	Operands []OpcodeChunk
}

func (o *Opcode) Encode(ops []lib.Operand) ([]uint8, error) {
	result := uint32(0)
	offset := 0
	operandIx := len(ops) - 1
	for i := len(o.Operands) - 1; i >= 0; i-- {
		op := o.Operands[i]
		value := op.Value
		switch op.OperandType {
		case OT_Exact:
		case OT_ImmediateValue:
			operand := ops[operandIx]
			switch n := operand.(type) {
			case Uint64:
				value = uint64(n)
			default:
				return nil, fmt.Errorf("Expecting immediate value in %s, got %s", o.String(), ops)
			}
			operandIx--
		case OT_Register64, OT_Register32:
			operand := ops[operandIx]
			if reg, ok := operand.(*Register); ok {
				value = uint64(reg.Encode())
				operandIx--
			} else {
				return nil, fmt.Errorf("Expecting register in %s, got %s", o.String(), ops)
			}
		}
		result = result + (uint32(value) << offset)
		offset += int(op.Size)
	}
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, uint32(result))
	return bytes, nil
}

func (o *Opcode) GetOperands() []OperandType {
	result := []OperandType{}
	for _, op := range o.Operands {
		if op.OperandType != OT_Exact {
			result = append(result, op.OperandType)
		}
	}
	return result
}

func (o *Opcode) MatchesOperands(operands []lib.Operand) bool {
	expected := o.GetOperands()
	if len(expected) != len(operands) {
		return false
	}
	for i, exp := range expected {
		op := operands[i]
		switch exp {
		case OT_Register64:
			if op.Type() != lib.T_Register || op.Width() != lib.QUADWORD {
				return false
			}
		case OT_ImmediateValue:
			if op.Type() != lib.T_Uint64 {
				return false
			}
		default:
			return false
		}
	}
	return true
}

func (o *Opcode) String() string {
	args := []string{}
	for _, ops := range o.Operands {
		args = append(args, ops.String())
	}
	return o.Name + " " + strings.Join(args, ", ")
}

//go:generate stringer -type=OperandType
type OperandType int

const (
	OT_Exact          OperandType = iota
	OT_Register32     OperandType = iota
	OT_Register64     OperandType = iota
	OT_ImmediateValue OperandType = iota
)

type OpcodeChunk struct {
	OperandType
	Size  uint8 // Size of encoding in bits
	Value uint64
}

var (
	OP_Xd    = OpcodeChunk{OT_Register64, 5, 0}
	OP_Xn    = OpcodeChunk{OT_Register64, 5, 0}
	OP_Xm    = OpcodeChunk{OT_Register64, 5, 0}
	OP_Wd    = OpcodeChunk{OT_Register32, 5, 0}
	OP_Wn    = OpcodeChunk{OT_Register32, 5, 0}
	OP_Wm    = OpcodeChunk{OT_Register32, 5, 0}
	OP_Imm12 = OpcodeChunk{OT_ImmediateValue, 12, 0}
	OP_Imm16 = OpcodeChunk{OT_ImmediateValue, 16, 0}
)

func OP_Exact(size uint8, value uint64, description ...string) OpcodeChunk {
	return OpcodeChunk{OT_Exact, size, value}
}
