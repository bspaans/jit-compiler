package asm

type REXPrefix struct {
	// When true a 64-bit operand size is used. Otherwise the default size is used.
	W bool
	// This is an extension to the ModRM.Reg field.
	R bool
	// This is an extension to the SIB.index field.
	X bool
	// This is an extension to the ModRM.RM field or the SIB.base field.
	B bool
}

func NewREXPrefix(w, r, x, b bool) *REXPrefix {
	return &REXPrefix{w, r, x, b}
}

func (r *REXPrefix) Encode() uint8 {
	result := uint8(0)
	if r.B {
		result = 1
	}
	if r.X {
		result += 1 << 1
	}
	if r.R {
		result += 1 << 2
	}
	if r.W {
		result += 1 << 3
	}
	return result + (1 << 6)
}
