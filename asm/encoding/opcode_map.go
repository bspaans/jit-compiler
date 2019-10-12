package encoding

import "fmt"

// TODO make into a []*Opcode because we need to match on size
// Or: map[Type]map[Size]*Opcode?
type OpcodeMap map[Type]*Opcode
type OpcodeMaps []map[Type]*Opcode

var JMP_OpcodeMap = map[Type]*Opcode{
	T_Register:          JMP_rm64,
	T_DisplacedRegister: JMP_rm64,
	T_RIPRelative:       JMP_rm64,
	T_Uint8:             JMP_rel8,
	T_Uint16:            nil,
	T_Uint32:            JMP_rel32,
	T_Uint64:            JMP_rm64,
	T_Int32:             nil,
	T_Float32:           nil,
	T_Float64:           nil,
}

func OpcodesToOpcodeMap(opcodes []*Opcode, operand int) OpcodeMap {
	opcodeMap := OpcodeMap{}
	for _, opcode := range opcodes {
		if opcode.Operands[operand].Type == OT_rel8 {
			opcodeMap[T_Uint8] = opcode
		} else if opcode.Operands[operand].Type == OT_rel16 {
			opcodeMap[T_Uint16] = opcode
		} else if opcode.Operands[operand].Type == OT_rel32 {
			opcodeMap[T_Uint32] = opcode
		} else if opcode.Operands[operand].Type == OT_rm8 {
			opcodeMap[T_Uint8] = opcode
			opcodeMap[T_Register] = opcode
			opcodeMap[T_DisplacedRegister] = opcode
			opcodeMap[T_RIPRelative] = opcode
		} else if opcode.Operands[operand].Type == OT_rm16 {
			opcodeMap[T_Uint16] = opcode
			opcodeMap[T_Register] = opcode
			opcodeMap[T_DisplacedRegister] = opcode
			opcodeMap[T_RIPRelative] = opcode
		} else if opcode.Operands[operand].Type == OT_rm32 {
			opcodeMap[T_Uint32] = opcode
			opcodeMap[T_Register] = opcode
			opcodeMap[T_DisplacedRegister] = opcode
			opcodeMap[T_RIPRelative] = opcode
		} else if opcode.Operands[operand].Type == OT_rm64 {
			opcodeMap[T_Uint64] = opcode
			opcodeMap[T_Register] = opcode
			opcodeMap[T_DisplacedRegister] = opcode
			opcodeMap[T_RIPRelative] = opcode
		} else if opcode.Operands[operand].Type == OT_m {
			opcodeMap[T_Uint64] = opcode
		} else if opcode.Operands[operand].Type == OT_m16 {
			opcodeMap[T_Uint16] = opcode
		} else if opcode.Operands[operand].Type == OT_m32 {
			opcodeMap[T_Uint32] = opcode
		} else if opcode.Operands[operand].Type == OT_m64 {
			opcodeMap[T_Uint64] = opcode
		} else if opcode.Operands[operand].Type == OT_imm8 {
			opcodeMap[T_Uint8] = opcode
		} else if opcode.Operands[operand].Type == OT_imm16 {
			opcodeMap[T_Uint16] = opcode
		} else if opcode.Operands[operand].Type == OT_imm32 {
			opcodeMap[T_Uint32] = opcode
		} else if opcode.Operands[operand].Type == OT_imm64 {
			opcodeMap[T_Uint64] = opcode
		} else if opcode.Operands[operand].Type == OT_r8 {
			opcodeMap[T_Register] = opcode
			opcodeMap[T_DisplacedRegister] = opcode
			opcodeMap[T_RIPRelative] = opcode
		} else if opcode.Operands[operand].Type == OT_r16 {
			opcodeMap[T_Register] = opcode
			opcodeMap[T_DisplacedRegister] = opcode
			opcodeMap[T_RIPRelative] = opcode
		} else if opcode.Operands[operand].Type == OT_r32 {
			opcodeMap[T_Register] = opcode
			opcodeMap[T_DisplacedRegister] = opcode
			opcodeMap[T_RIPRelative] = opcode
		} else if opcode.Operands[operand].Type == OT_r64 {
			opcodeMap[T_Register] = opcode
			opcodeMap[T_DisplacedRegister] = opcode
			opcodeMap[T_RIPRelative] = opcode
		} else if opcode.Operands[operand].Type == OT_xmm1 {
			opcodeMap[T_Register] = opcode
			opcodeMap[T_DisplacedRegister] = opcode
			opcodeMap[T_RIPRelative] = opcode
		} else if opcode.Operands[operand].Type == OT_xmm2 {
			opcodeMap[T_Register] = opcode
			opcodeMap[T_DisplacedRegister] = opcode
			opcodeMap[T_RIPRelative] = opcode
		} else if opcode.Operands[operand].Type == OT_xmm1m64 {
			opcodeMap[T_Register] = opcode
			opcodeMap[T_DisplacedRegister] = opcode
			opcodeMap[T_RIPRelative] = opcode
		} else if opcode.Operands[operand].Type == OT_xmm2m64 {
			opcodeMap[T_Register] = opcode
			opcodeMap[T_DisplacedRegister] = opcode
			opcodeMap[T_RIPRelative] = opcode
		}
	}
	return opcodeMap
}

func init() {
	jmps := []*Opcode{JMP_rel8, JMP_rel32, JMP_rm64}
	mmap := OpcodesToOpcodeMap(jmps, 0)
	fmt.Println(JMP_OpcodeMap)
	fmt.Println(mmap)
}
