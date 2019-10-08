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
				if dest.Register == Rsp {
					result = append(result, 0x24)
				}
				result = append(result, dest.Displacement)
				return result, nil
			}
		}
	} else if i.Source.Type() == T_RIPRelative {
		src := i.Source.(*RIPRelative)
		if i.Dest.Type() == T_Register {
			dest := i.Dest.(*Register)
			rex := REXEncode(nil, dest)
			modrm := NewModRM(IndirectRegisterMode, 5, 0).Encode()
			result := []uint8{rex, 0x8b, modrm}
			for _, c := range src.Displacement.Encode() {
				result = append(result, c)
			}
			return result, nil
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
			if dest.Register == Rsp {
				result = append(result, 0x24)
			}
			result = append(result, dest.Displacement)
			// Can only move a double
			for _, b := range src.Encode()[:4] {
				result = append(result, b)
			}
			return result, nil
		}
	} else if i.Source.Type() == T_Float64 {
		src := i.Source.(Float64)
		if i.Dest.Type() == T_Register {
			dest := i.Dest.(*Register)
			result := EncodeOpcodeWithREX(0xB8+(dest.Encode()&7), dest, nil)
			for _, b := range src.Encode() {
				result = append(result, b)
			}
			return result, nil
		}
	}
	return nil, errors.New("Unsupported mov operation: " + i.String())
}
func (i *MOV) String() string {
	return "mov " + i.Source.String() + ", " + i.Dest.String()
}

type MOVQ struct {
	Source Operand
	Dest   Operand
}

func (i *MOVQ) Encode() (MachineCode, error) {
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
			if dest.Size == QUADDOUBLE {
				if src.Size == QUADWORD {
					result := []uint8{0x66}
					rex := REXEncode(src, dest)
					result = append(result, rex)
					result = append(result, 0x0f)
					result = append(result, 0x6e)
					modrm := NewModRM(DirectRegisterMode, src.Encode(), dest.Encode())
					result = append(result, modrm.Encode())
					return result, nil

				}
			}
		}
	}
	return nil, errors.New("Unsupported movq operation: " + i.String())
}
func (i *MOVQ) String() string {
	return "movq " + i.Source.String() + ", " + i.Dest.String()
}

type MOVSD struct {
	Source Operand
	Dest   Operand
}

func (i *MOVSD) Encode() (MachineCode, error) {
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
			if dest.Size == QUADDOUBLE {
				if src.Size == QUADDOUBLE {
					result := []uint8{0xf2, 0x0f, 0x10}
					modrm := NewModRM(DirectRegisterMode, src.Encode(), dest.Encode())
					result = append(result, modrm.Encode())
					return result, nil

				}
			}
		}
	}
	return nil, errors.New("Unsupported movq operation: " + i.String())
}
func (i *MOVSD) String() string {
	return "movsd " + i.Source.String() + ", " + i.Dest.String()
}
