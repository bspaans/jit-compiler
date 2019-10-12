package asm

import (
	"errors"

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
		if i.Dest.Type() == encoding.T_Register || i.Dest.Type() == encoding.T_DisplacedRegister {
			return encoding.MOV_rm64_r64.Encode([]encoding.Operand{i.Dest, i.Source})
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
	} else if i.Source.Type() == encoding.T_Uint64 || i.Source.Type() == encoding.T_Float64 {
		if i.Dest.Type() == encoding.T_Register {
			return encoding.MOV_r64_imm64.Encode([]encoding.Operand{i.Dest, i.Source})
		}
	} else if i.Source.Type() == encoding.T_Uint32 {
		if i.Dest.Type() == encoding.T_DisplacedRegister || i.Dest.Type() == encoding.T_Register {
			return encoding.MOV_rm64_imm32.Encode([]encoding.Operand{i.Dest, i.Source})
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
		if i.Dest.Type() == encoding.T_Register {
			return encoding.MOVQ_xmm_rm64.Encode([]encoding.Operand{i.Dest, i.Source})
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
					return encoding.MOVSD_xmm1m64_xmm2.Encode([]encoding.Operand{i.Dest, i.Source})
				}
			}
		}
	}
	return nil, errors.New("Unsupported movsd operation: " + i.String())
}
func (i *MOVSD) String() string {
	return "movsd " + i.Source.String() + ", " + i.Dest.String()
}
