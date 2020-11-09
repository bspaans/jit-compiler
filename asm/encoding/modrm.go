package encoding

type Mode uint8

const (
	IndirectRegisterMode                Mode = 0
	IndirectRegisterByteDisplacedMode   Mode = 1
	IndirectRegisterDoubleDisplacedMode Mode = 2
	DirectRegisterMode                  Mode = 3

	// If RM is set to SIBFollowsRM a SIB byte is expected to follow
	// the ModRM byte (see instruction format)
	SIBFollowsRM uint8 = 4
)

type ModRM struct {
	Mode Mode
	// The r/m field can specify a register or operand or it can be combined with the
	// mod field to encode an addressing mode.
	RM uint8
	// The reg/opcode field specifies either a register number or
	// three more bits of opcode information.
	Reg uint8
}

func NewModRM(mode Mode, rm, reg uint8) *ModRM {
	return &ModRM{mode, rm, reg}
}

func (m *ModRM) Encode() uint8 {
	result := m.RM & 7
	result += (m.Reg & 7) << 3
	return result + (uint8(m.Mode) << 6)
}
