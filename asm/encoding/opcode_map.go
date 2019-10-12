package encoding

import (
	"fmt"

	"github.com/bspaans/jit/lib"
)

type OpcodeMap map[Type]map[lib.Size][]*Opcode

func (o OpcodeMap) add(ty Type, si lib.Size, op *Opcode) {
	arr, found := o[ty][si]
	if !found {
		arr = []*Opcode{}
	}
	arr = append(arr, op)
	o[ty][si] = arr
}

type OpcodeMaps []OpcodeMap

func (o OpcodeMaps) ResolveOpcode(operands []Operand) *Opcode {
	picks := map[*Opcode]bool{}

	for i, opcodeMap := range o {
		fmt.Println(i, opcodeMap)
		oper := operands[i]
		if oper == nil {
			return nil
		}
		matches := opcodeMap[oper.Type()][oper.Width()]
		if len(matches) == 0 {
			return nil
		}
		newPick := map[*Opcode]bool{}
		for _, opcode := range matches {
			if i == 0 {
				newPick[opcode] = true
			} else {
				if picks[opcode] {
					newPick[opcode] = true
				}
			}
		}
		picks = newPick
	}
	for pick, _ := range picks {
		return pick
	}
	return nil
}

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

func NewOpcodeMap() OpcodeMap {
	return map[Type]map[lib.Size][]*Opcode{
		T_Register:          map[lib.Size][]*Opcode{},
		T_DisplacedRegister: map[lib.Size][]*Opcode{},
		T_RIPRelative:       map[lib.Size][]*Opcode{},
		T_Uint8:             map[lib.Size][]*Opcode{},
		T_Uint16:            map[lib.Size][]*Opcode{},
		T_Uint32:            map[lib.Size][]*Opcode{},
		T_Uint64:            map[lib.Size][]*Opcode{},
		T_Int32:             map[lib.Size][]*Opcode{},
		T_Float32:           map[lib.Size][]*Opcode{},
		T_Float64:           map[lib.Size][]*Opcode{},
	}
}

func OpcodesToOpcodeMaps(opcodes []*Opcode, argCount int) OpcodeMaps {
	maps := make([]OpcodeMap, argCount)
	for i := 0; i < argCount; i++ {
		opcodeMap := OpcodesToOpcodeMap(opcodes, i)
		maps[i] = opcodeMap
	}
	return maps
}

func OpcodesToOpcodeMap(opcodes []*Opcode, operand int) OpcodeMap {
	opcodeMap := NewOpcodeMap()
	for _, opcode := range opcodes {
		if opcode.Operands[operand].Type == OT_rel8 {
			opcodeMap.add(T_Uint8, lib.BYTE, opcode)
		} else if opcode.Operands[operand].Type == OT_rel16 {
			opcodeMap.add(T_Uint16, lib.WORD, opcode)
		} else if opcode.Operands[operand].Type == OT_rel32 {
			opcodeMap.add(T_Uint32, lib.DOUBLE, opcode)
		} else if opcode.Operands[operand].Type == OT_rm8 {
			opcodeMap.add(T_Uint8, lib.BYTE, opcode)
			opcodeMap.add(T_Register, lib.BYTE, opcode)
			opcodeMap.add(T_DisplacedRegister, lib.BYTE, opcode)
			opcodeMap.add(T_RIPRelative, lib.BYTE, opcode)
		} else if opcode.Operands[operand].Type == OT_rm16 {
			opcodeMap.add(T_Uint16, lib.WORD, opcode)
			opcodeMap.add(T_Register, lib.WORD, opcode)
			opcodeMap.add(T_DisplacedRegister, lib.WORD, opcode)
			opcodeMap.add(T_RIPRelative, lib.WORD, opcode)
		} else if opcode.Operands[operand].Type == OT_rm32 {
			opcodeMap.add(T_Uint32, lib.DOUBLE, opcode)
			opcodeMap.add(T_Register, lib.DOUBLE, opcode)
			opcodeMap.add(T_DisplacedRegister, lib.DOUBLE, opcode)
			opcodeMap.add(T_RIPRelative, lib.DOUBLE, opcode)
		} else if opcode.Operands[operand].Type == OT_rm64 {
			opcodeMap.add(T_Uint64, lib.QUADWORD, opcode)
			opcodeMap.add(T_Register, lib.QUADWORD, opcode)
			opcodeMap.add(T_DisplacedRegister, lib.QUADWORD, opcode)
			opcodeMap.add(T_RIPRelative, lib.QUADWORD, opcode)
		} else if opcode.Operands[operand].Type == OT_m {
			opcodeMap.add(T_Uint64, lib.QUADWORD, opcode)
		} else if opcode.Operands[operand].Type == OT_m16 {
			opcodeMap.add(T_Uint16, lib.WORD, opcode)
		} else if opcode.Operands[operand].Type == OT_m32 {
			opcodeMap.add(T_Uint32, lib.DOUBLE, opcode)
		} else if opcode.Operands[operand].Type == OT_m64 {
			opcodeMap.add(T_Uint64, lib.QUADWORD, opcode)
		} else if opcode.Operands[operand].Type == OT_imm8 {
			opcodeMap.add(T_Uint8, lib.BYTE, opcode)
		} else if opcode.Operands[operand].Type == OT_imm16 {
			opcodeMap.add(T_Uint16, lib.WORD, opcode)
		} else if opcode.Operands[operand].Type == OT_imm32 {
			opcodeMap.add(T_Uint32, lib.DOUBLE, opcode)
		} else if opcode.Operands[operand].Type == OT_imm64 {
			opcodeMap.add(T_Uint64, lib.QUADWORD, opcode)
		} else if opcode.Operands[operand].Type == OT_r8 {
			opcodeMap.add(T_Register, lib.BYTE, opcode)
			opcodeMap.add(T_DisplacedRegister, lib.BYTE, opcode)
			opcodeMap.add(T_RIPRelative, lib.BYTE, opcode)
		} else if opcode.Operands[operand].Type == OT_r16 {
			opcodeMap.add(T_Register, lib.WORD, opcode)
			opcodeMap.add(T_DisplacedRegister, lib.WORD, opcode)
			opcodeMap.add(T_RIPRelative, lib.WORD, opcode)
		} else if opcode.Operands[operand].Type == OT_r32 {
			opcodeMap.add(T_Register, lib.DOUBLE, opcode)
			opcodeMap.add(T_DisplacedRegister, lib.DOUBLE, opcode)
			opcodeMap.add(T_RIPRelative, lib.DOUBLE, opcode)
		} else if opcode.Operands[operand].Type == OT_r64 {
			opcodeMap.add(T_Register, lib.QUADWORD, opcode)
			opcodeMap.add(T_DisplacedRegister, lib.QUADWORD, opcode)
			opcodeMap.add(T_RIPRelative, lib.QUADWORD, opcode)
		} else if opcode.Operands[operand].Type == OT_xmm1 {
			opcodeMap.add(T_Register, lib.QUADDOUBLE, opcode)
			opcodeMap.add(T_DisplacedRegister, lib.QUADDOUBLE, opcode)
			opcodeMap.add(T_RIPRelative, lib.QUADDOUBLE, opcode)
		} else if opcode.Operands[operand].Type == OT_xmm2 {
			opcodeMap.add(T_Register, lib.QUADDOUBLE, opcode)
			opcodeMap.add(T_DisplacedRegister, lib.QUADDOUBLE, opcode)
			opcodeMap.add(T_RIPRelative, lib.QUADDOUBLE, opcode)
		} else if opcode.Operands[operand].Type == OT_xmm1m64 {
			opcodeMap.add(T_Register, lib.QUADDOUBLE, opcode)
			opcodeMap.add(T_DisplacedRegister, lib.QUADDOUBLE, opcode)
			opcodeMap.add(T_RIPRelative, lib.QUADDOUBLE, opcode)
		} else if opcode.Operands[operand].Type == OT_xmm2m64 {
			opcodeMap.add(T_Register, lib.QUADDOUBLE, opcode)
			opcodeMap.add(T_DisplacedRegister, lib.QUADDOUBLE, opcode)
			opcodeMap.add(T_RIPRelative, lib.QUADDOUBLE, opcode)
		}
	}
	return opcodeMap
}

func init() {
	jmps := []*Opcode{JMP_rel8, JMP_rel32, JMP_rm64}
	mmap := OpcodesToOpcodeMaps(jmps, 1)
	mmaps := OpcodesToOpcodeMaps(jmps, 1)
	fmt.Println(JMP_OpcodeMap)
	fmt.Println(mmap)
	fmt.Println(mmaps.ResolveOpcode([]Operand{Uint8(uint8(3))}))
	fmt.Println(mmaps.ResolveOpcode([]Operand{Uint64(uint64(3))}))
}
