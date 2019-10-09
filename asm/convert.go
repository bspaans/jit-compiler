package asm

import "errors"

// Convert signed integer to scalar double-precision floating point (float64)
type CVTSI2SS struct {
	Source Operand
	Dest   Operand
}

func (i *CVTSI2SS) Encode() (MachineCode, error) {
	if i.Source.Type() == T_Register {
		src := i.Source.(*Register)
		if i.Dest.Type() == T_Register {
			dest := i.Dest.(*Register)
			if dest.Size == QUADDOUBLE {
				if src.Size == QUADWORD {
					result := []uint8{0xf2}
					rex := REXEncode(src, dest)
					result = append(result, rex)
					result = append(result, 0x0f)
					result = append(result, 0x2a)
					modrm := NewModRM(DirectRegisterMode, src.Encode(), dest.Encode())
					result = append(result, modrm.Encode())
					return result, nil
				} else if src.Size == DOUBLE {
				}
			}
		}
	}
	return nil, errors.New("Unsupported cvtsi2ss operation: " + i.String())
}
func (j *CVTSI2SS) String() string {
	return "cvtsi2ss " + j.Source.String() + " " + j.Dest.String()
}

// Convert double precision float to signed integer
type CVTTSD2SI struct {
	Source Operand
	Dest   Operand
}

func (i *CVTTSD2SI) Encode() (MachineCode, error) {
	if i.Source.Type() == T_Register {
		src := i.Source.(*Register)
		if i.Dest.Type() == T_Register {
			dest := i.Dest.(*Register)
			if src.Size == QUADDOUBLE {
				if dest.Size == QUADWORD {
					result := []uint8{0xf2}
					rex := REXEncode(src, dest)
					result = append(result, rex)
					result = append(result, 0x0f)
					result = append(result, 0x2c)
					modrm := NewModRM(DirectRegisterMode, src.Encode(), dest.Encode())
					result = append(result, modrm.Encode())
					return result, nil
				}
			}
		}
	}
	return nil, errors.New("Unsupported cvttsd2si operation: " + i.String())
}
func (j *CVTTSD2SI) String() string {
	return "cvttsd2si " + j.Source.String() + " " + j.Dest.String()
}
