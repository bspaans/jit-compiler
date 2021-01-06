package encoding

import "github.com/bspaans/jit-compiler/lib"

type Opcode struct {
	Name     string
	Operands []OpcodeChunk
}

func (o *Opcode) Encode(ops []lib.Operand) ([]uint8, error) {
	return []uint8{}, nil

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
