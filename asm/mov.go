package asm

import (
	"errors"
	"fmt"

	"github.com/bspaans/jit/asm/encoding"
	"github.com/bspaans/jit/lib"
)

type MOV struct {
	Source encoding.Operand
	Dest   encoding.Operand
}

func (i *MOV) Encode() (lib.MachineCode, error) {
	if i.Dest == nil {
		return nil, errors.New("Missing dest")
	}
	if i.Source == nil {
		return nil, errors.New("Missing source")
	}
	if i.Source.Type() == encoding.T_Register {
		src := i.Source.(*encoding.Register)
		if i.Dest.Type() == encoding.T_Register || i.Dest.Type() == encoding.T_DisplacedRegister {
			return encoding.MOV_rm64_r64.Encode([]encoding.Operand{i.Dest, src})
		}
	} else if i.Source.Type() == encoding.T_RIPRelative {
		src := i.Source.(*encoding.RIPRelative)
		if i.Dest.Type() == encoding.T_Register {
			dest := i.Dest.(*encoding.Register)
			rex := REXEncode(nil, dest)
			modrm := encoding.NewModRM(encoding.IndirectRegisterMode, 5, 0).Encode()
			result := []uint8{rex, 0x8b, modrm}
			for _, c := range src.Displacement.Encode() {
				result = append(result, c)
			}
			return result, nil
		}
	} else if i.Source.Type() == encoding.T_Uint64 {
		src := i.Source.(encoding.Uint64)
		if i.Dest.Type() == encoding.T_Register {
			fmt.Println("hello", i.Dest, encoding.MOV_r64_imm64.Opcode)
			return encoding.MOV_r64_imm64.Encode([]encoding.Operand{i.Dest, src})
		} else if i.Dest.Type() == encoding.T_DisplacedRegister {
			dest := i.Dest.(*encoding.DisplacedRegister)
			result := EncodeOpcodeWithREX(0xC7, dest.Register, nil)
			modrm := encoding.NewModRM(encoding.IndirectRegisterByteDisplacedMode, dest.Encode(), 0).Encode()
			result = append(result, modrm)
			// Not sure why this is needed, but it is
			if dest.Register == encoding.Rsp {
				result = append(result, 0x24)
			}
			result = append(result, dest.Displacement)
			// Can only move a double
			for _, b := range src.Encode()[:4] {
				result = append(result, b)
			}
			return result, nil
		}
	} else if i.Source.Type() == encoding.T_Float64 {
		src := i.Source.(encoding.Float64)
		if i.Dest.Type() == encoding.T_Register {
			dest := i.Dest.(*encoding.Register)
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
	Source encoding.Operand
	Dest   encoding.Operand
}

func (i *MOVQ) Encode() (lib.MachineCode, error) {
	if i.Dest == nil {
		return nil, errors.New("Missing dest")
	}
	if i.Source == nil {
		return nil, errors.New("Missing source")
	}
	if i.Source.Type() == encoding.T_Register {
		src := i.Source.(*encoding.Register)
		if i.Dest.Type() == encoding.T_Register {
			dest := i.Dest.(*encoding.Register)
			if dest.Size == lib.QUADDOUBLE {
				if src.Size == lib.QUADWORD {
					result := []uint8{0x66}
					rex := REXEncode(src, dest)
					result = append(result, rex)
					result = append(result, 0x0f)
					result = append(result, 0x6e)
					modrm := encoding.NewModRM(encoding.DirectRegisterMode, src.Encode(), dest.Encode())
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
	Source encoding.Operand
	Dest   encoding.Operand
}

func (i *MOVSD) Encode() (lib.MachineCode, error) {
	if i.Dest == nil {
		return nil, errors.New("Missing dest")
	}
	if i.Source == nil {
		return nil, errors.New("Missing source")
	}
	if i.Source.Type() == encoding.T_Register {
		src := i.Source.(*encoding.Register)
		if i.Dest.Type() == encoding.T_Register {
			dest := i.Dest.(*encoding.Register)
			if dest.Size == lib.QUADDOUBLE {
				if src.Size == lib.QUADDOUBLE {
					result := []uint8{0xf2, 0x0f, 0x10}
					modrm := encoding.NewModRM(encoding.DirectRegisterMode, src.Encode(), dest.Encode())
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
