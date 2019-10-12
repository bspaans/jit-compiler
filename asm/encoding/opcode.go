package encoding

import (
	"fmt"
	"strings"
)

type OperandType int

const (
	OT_rel8    OperandType = iota
	OT_rel16   OperandType = iota
	OT_rel32   OperandType = iota
	OT_rm8     OperandType = iota
	OT_rm16    OperandType = iota
	OT_rm32    OperandType = iota
	OT_rm64    OperandType = iota
	OT_m       OperandType = iota
	OT_m16     OperandType = iota
	OT_m32     OperandType = iota
	OT_m64     OperandType = iota
	OT_m128    OperandType = iota
	OT_r8      OperandType = iota
	OT_r16     OperandType = iota
	OT_r32     OperandType = iota
	OT_r64     OperandType = iota
	OT_imm8    OperandType = iota
	OT_imm16   OperandType = iota
	OT_imm32   OperandType = iota
	OT_imm64   OperandType = iota
	OT_xmm1    OperandType = iota
	OT_xmm1m64 OperandType = iota
	OT_xmm2    OperandType = iota
	OT_xmm2m64 OperandType = iota
)

type OperandEncoding int

const (
	// Register will be read by the processor
	ModRM_rm_r OperandEncoding = iota
	// Register will be updated by the processor
	ModRM_rm_rw OperandEncoding = iota
	// Register will be read by the processor
	ModRM_reg_r OperandEncoding = iota
	// Register will be updated by the processor
	ModRM_reg_rw OperandEncoding = iota
	// Immediate value
	ImmediateValue OperandEncoding = iota
	// Add register to opcode
	Opcode_plus_rd_r = iota
)

type OpcodeExtensions int

const (
	NoExtensions    OpcodeExtensions = iota
	ImmediateByte   OpcodeExtensions = iota
	ImmediateWord   OpcodeExtensions = iota
	ImmediateDouble OpcodeExtensions = iota
	Slash0          OpcodeExtensions = iota
	Slash1          OpcodeExtensions = iota
	Slash2          OpcodeExtensions = iota
	Slash3          OpcodeExtensions = iota
	Slash4          OpcodeExtensions = iota
	Slash5          OpcodeExtensions = iota
	Slash6          OpcodeExtensions = iota
	Slash7          OpcodeExtensions = iota
	SlashR          OpcodeExtensions = iota
	RexW            OpcodeExtensions = iota
)

type OpcodeOperand struct {
	Type     OperandType
	Encoding OperandEncoding
}

func (o *OpcodeOperand) String() string {
	switch o.Type {
	case OT_rm64:
		return "r/m64"
	case OT_r64:
		return "r64"
	case OT_imm32:
		return "imm32"
	}
	return "?"
}

func (o OpcodeOperand) TypeCheck(op Operand) bool {
	return true
}

type Opcode struct {
	Name             string
	Prefixes         []uint8
	Opcode           []uint8
	OpcodeExtensions []OpcodeExtensions
	Operands         []OpcodeOperand
}

func (o *Opcode) Encode(ops []Operand) ([]uint8, error) {
	fmt.Println(o.Opcode)
	instr := NewInstructionFormat(o.Opcode)
	exts := map[OpcodeExtensions]bool{}
	for _, p := range o.Prefixes {
		instr.Prefixes = append(instr.Prefixes, p)
	}
	for _, ext := range o.OpcodeExtensions {
		if ext == Slash0 {
			if instr.ModRM == nil {
				instr.ModRM = &ModRM{}
			}
			instr.ModRM.Reg = 0
		} else if ext == Slash1 {
			if instr.ModRM == nil {
				instr.ModRM = &ModRM{}
			}
			instr.ModRM.Reg = 1
		} else if ext == Slash2 {
			if instr.ModRM == nil {
				instr.ModRM = &ModRM{}
			}
			instr.ModRM.Reg = 2
		} else if ext == Slash3 {
			if instr.ModRM == nil {
				instr.ModRM = &ModRM{}
			}
			instr.ModRM.Reg = 3
		} else if ext == Slash4 {
			if instr.ModRM == nil {
				instr.ModRM = &ModRM{}
			}
			instr.ModRM.Reg = 4
		} else if ext == Slash5 {
			if instr.ModRM == nil {
				instr.ModRM = &ModRM{}
			}
			instr.ModRM.Reg = 5
		} else if ext == Slash6 {
			if instr.ModRM == nil {
				instr.ModRM = &ModRM{}
			}
			instr.ModRM.Reg = 6
		} else if ext == Slash7 {
			if instr.ModRM == nil {
				instr.ModRM = &ModRM{}
			}
			instr.ModRM.Reg = 7
		} else if ext == RexW {
			if instr.REXPrefix == nil {
				instr.REXPrefix = &REXPrefix{}
			}
			instr.REXPrefix.W = true
		} else if ext == SlashR {
			if instr.ModRM == nil {
				instr.ModRM = &ModRM{}
			}
		}
		exts[ext] = true
	}
	for i, opcodeOperand := range o.Operands {
		op := ops[i]
		if opcodeOperand.TypeCheck(op) {
			if op.Type() == T_Register {
				oper := op.(*Register)
				if opcodeOperand.Encoding == ModRM_rm_r || opcodeOperand.Encoding == ModRM_rm_rw {
					if instr.ModRM == nil {
						instr.ModRM = &ModRM{}
					}
					instr.ModRM.Mode = DirectRegisterMode
					instr.ModRM.RM = oper.Encode()
					if exts[RexW] {
						instr.REXPrefix.B = oper.Register > 7
					}
				} else if opcodeOperand.Encoding == ModRM_reg_r || opcodeOperand.Encoding == ModRM_reg_rw {
					if instr.ModRM == nil {
						instr.ModRM = &ModRM{}
						instr.ModRM.Mode = DirectRegisterMode
					}
					instr.ModRM.Reg = oper.Encode()
					if exts[RexW] {
						instr.REXPrefix.R = oper.Register > 7
					}
				} else if opcodeOperand.Encoding == Opcode_plus_rd_r {
					fmt.Println("encoding", instr.Opcode)
					instr.Opcode[0] += op.(*Register).Register & 7
					if exts[RexW] {
						instr.REXPrefix.B = op.(*Register).Register > 7
					}
				} else {
					return nil, fmt.Errorf("Unsupported encoding [%d] in %s", opcodeOperand.Encoding, o.String())
				}
			} else if op.Type() == T_DisplacedRegister {
				oper := op.(*DisplacedRegister)
				if opcodeOperand.Encoding == ModRM_rm_r || opcodeOperand.Encoding == ModRM_rm_rw {
					if instr.ModRM == nil {
						instr.ModRM = &ModRM{}
					}
					instr.ModRM.Mode = IndirectRegisterByteDisplacedMode
					instr.ModRM.RM = oper.Encode()
					instr.SetDisplacement(oper.Register, []uint8{oper.Displacement})

					if exts[RexW] {
						instr.REXPrefix.B = oper.Register.Register > 7
					}
				} else if opcodeOperand.Encoding == ModRM_reg_r || opcodeOperand.Encoding == ModRM_reg_rw {
					if instr.ModRM == nil {
						instr.ModRM = &ModRM{}
						instr.ModRM.Mode = IndirectRegisterByteDisplacedMode
					}
					instr.ModRM.Reg = oper.Encode()
					instr.SetDisplacement(oper.Register, []uint8{oper.Displacement})
					if exts[RexW] {
						instr.REXPrefix.R = oper.Register.Register > 7
					}
				} else {
					return nil, fmt.Errorf("Unsupported encoding [%d] in %s", opcodeOperand.Encoding, o.String())
				}
			} else if op.Type() == T_RIPRelative {
				oper := op.(*RIPRelative)
				if opcodeOperand.Encoding == ModRM_rm_r || opcodeOperand.Encoding == ModRM_rm_rw {
					if instr.ModRM == nil {
						instr.ModRM = &ModRM{}
					}
					instr.ModRM.Mode = IndirectRegisterMode
					instr.ModRM.RM = 5
					instr.SetDisplacement(op, oper.Displacement.Encode())
				} else if opcodeOperand.Encoding == ModRM_reg_r || opcodeOperand.Encoding == ModRM_reg_rw {
					if instr.ModRM == nil {
						instr.ModRM = &ModRM{}
					}
					instr.ModRM.Mode = IndirectRegisterMode
					instr.ModRM.Reg = 5
					instr.SetDisplacement(op, oper.Displacement.Encode())
				}

			} else if op.Type() == T_Uint64 {
				for _, b := range op.(Uint64).Encode() {
					instr.Immediate = append(instr.Immediate, b)
				}
			} else if op.Type() == T_Float64 {
				for _, b := range op.(Float64).Encode() {
					instr.Immediate = append(instr.Immediate, b)
				}
			} else if op.Type() == T_Uint32 {
				for _, b := range op.(Uint32).Encode() {
					instr.Immediate = append(instr.Immediate, b)
				}
			} else if op.Type() == T_Uint16 {
				for _, b := range op.(Uint16).Encode() {
					instr.Immediate = append(instr.Immediate, b)
				}
			} else if op.Type() == T_Uint8 {
				for _, b := range op.(Uint8).Encode() {
					instr.Immediate = append(instr.Immediate, b)
				}
			} else {
				return nil, fmt.Errorf("Unsupported type [%d] in %s", op.Type(), o.String())
			}
		} else {
			return nil, fmt.Errorf("Unsupported type")
		}
	}
	return instr.Encode(), nil
}

func (o *Opcode) String() string {
	args := []string{}
	for _, ops := range o.Operands {
		args = append(args, ops.String())
	}
	return o.Name + " " + strings.Join(args, ", ")
}
