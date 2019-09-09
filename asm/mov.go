package asm

import "errors"

type MOV struct {
	Source Operand
	Dest   Operand
}

func (i *MOV) Encode() (MachineCode, error) {
	if i.Dest == nil {
		return nil, errors.New("Missing dest")
	}
	if i.Source == nil {
		return nil, errors.New("Missing source")
	}
	if i.Source.Type() == T_Register {
		src := i.Source.(*Register)
		if i.Dest.Type() == T_Register {
			dest := i.Dest.(*Register)
			if src.Size == QUADWORD && dest.Size == QUADWORD {
				return EncodeOpcodeWithREXAndModRM(0x89, dest, src, 0), nil
			}
		} else if i.Dest.Type() == T_DisplacedRegister {
			dest := i.Dest.(*DisplacedRegister)
			if src.Size == QUADWORD && dest.Size == QUADWORD {
				rex := REXEncode(src, dest.Register)
				modrm := NewModRM(IndirectRegisterByteDisplacedMode, dest.Encode(), src.Encode()).Encode()
				result := []uint8{rex, 0x89, modrm}
				// Not sure why this is needed, but it is.
				if dest.Register == rsp {
					result = append(result, 0x24)
				}
				result = append(result, dest.Displacement)
				return result, nil
			}
		}
	} else if i.Source.Type() == T_Uint64 {
		src := i.Source.(Uint64)
		if i.Dest.Type() == T_Register {
			dest := i.Dest.(*Register)
			result := EncodeOpcodeWithREX(0xB8+(dest.Encode()&7), dest, nil)
			for _, b := range src.Encode() {
				result = append(result, b)
			}
			return result, nil
		} else if i.Dest.Type() == T_DisplacedRegister {
			dest := i.Dest.(*DisplacedRegister)
			result := EncodeOpcodeWithREX(0xC7, dest.Register, nil)
			modrm := NewModRM(IndirectRegisterByteDisplacedMode, dest.Encode(), 0).Encode()
			result = append(result, modrm)
			// Not sure why this is needed, but it is
			if dest.Register == rsp {
				result = append(result, 0x24)
			}
			result = append(result, dest.Displacement)
			// Can only move a double
			for _, b := range src.Encode()[:4] {
				result = append(result, b)
			}
			return result, nil
		}
	}
	return nil, errors.New("Unsupported mov operation")
}
func (i *MOV) String() string {
	return "mov " + i.Source.String() + " " + i.Dest.String()
}
