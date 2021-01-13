package encoding

type VEXOpcodeExtension uint8
type VEXLegacyByte uint8

const (
	VEXOpcodeExtension_None VEXOpcodeExtension = 0x00
	VEXOpcodeExtension_66   VEXOpcodeExtension = 0x01
	VEXOpcodeExtension_f2   VEXOpcodeExtension = 0x03
	VEXOpcodeExtension_f3   VEXOpcodeExtension = 0x02

	VEXLegacyByte_None  VEXLegacyByte = 0x00
	VEXLegacyByte_0f    VEXLegacyByte = 0x01
	VEXLegacyByte_0f_38 VEXLegacyByte = 0x02
	VEXLegacyByte_0f_3a VEXLegacyByte = 0x03
)

// The VEX prefix is encoded in either the two-byte form (the first byte must
// be C5H) or in the three-byte form (the first byte must be C4H). The two-byte
// VEX is used mainly for 128-bit, scalar, and the most common 256-bit AVX
// instructions; while the three-byte VEX provides a compact replacement of REX
// and 3-byte opcode instructions (including AVX and FMA instructions).
type VEXPrefix struct {
	// Non-destructive source register encoding.
	// It is represented by the notation, VEX.vvvv. This field is encoded using
	// 1’s complement form (inverted form), i.e. XMM0/YMM0/R0 is encoded as 1111B,
	// XMM15/YMM15/R15 is encoded as 0000B.
	Source uint8

	// Vector length encoding: This 1-bit field represented by the notation
	// VEX.L. L= 0 means vector length is 128 bits wide, L=1 means 256 bit vector.
	L bool

	// REX prefix functionality: Full REX prefix functionality is provided in
	// the three-byte form of VEX prefix. However the VEX bit fields providing
	// REX functionality are encoded using 1’s complement form.
	// Two-byte form of the VEX prefix only provides the equivalent functionality of REX.R
	R bool
	W bool
	X bool
	B bool

	// Compaction of SIMD prefix: Legacy SSE instructions effectively use SIMD
	// prefixes (66H, F2H, F3H) as an opcode extension field. VEX prefix
	// encoding allows the functional capability of such legacy SSE
	// instructions (operating on XMM registers, bits 255:128 of corresponding
	// YMM unmodified) to be encoded using the VEX.pp field without the
	// presence of any SIMD prefix. The VEX-encoded 128-bit instruction will
	// zero-out bits 255:128 of the destination register. VEX-encoded
	// instruction may have 128 bit vector length or 256 bits length.
	//
	// VEX.pp
	VEXOpcodeExtension VEXOpcodeExtension

	// Compaction of two-byte and three-byte opcode: More recently introduced
	// legacy SSE instructions employ two and three-byte opcode. The one or two
	// leading bytes are: 0FH, and 0FH 3AH/0FH 38H. The one-byte escape (0FH)
	// and two-byte escape (0FH 3AH, 0FH 38H) can also be interpreted as an
	// opcode extension field. The VEX.mmmmm field provides compaction to allow
	// many legacy instruction to be encoded without the constant byte
	// sequence, 0FH, 0FH 3AH, 0FH 38H.
	//
	// 00000: Reserved for future use (will #UD)
	// 00001: implied 0F leading opcode byte
	// 00010: implied 0F 38 leading opcode bytes
	// 00011: implied 0F 3A leading opcode bytes
	// 00100-11111: Reserved for future use (will #UD)
	//
	// VEX.mmmmm
	VEXLegacyByte VEXLegacyByte
}

func NewVEXPrefix() *VEXPrefix {
	return &VEXPrefix{
		// 2's complement, so set to true by default
		R: true,
		W: true,
		X: true,
		B: true,
	}
}

func (v *VEXPrefix) Encode() []uint8 {

	// Two byte form
	// The presence of 0F3A and 0F38 in the opcode column implies that opcode
	// can only be encoded by the three-byte form of VEX. The presence of 0F in
	// the opcode column does not preclude the opcode to be encoded by the
	// two-byte of VEX if the semantics of the opcode does not require any
	// subfield of VEX not present in the two-byte form of the VEX prefix.
	if (v.VEXLegacyByte == 0 || v.VEXLegacyByte == VEXLegacyByte_0f) && v.X && v.B {
		byte0 := uint8(0xc5)
		byte1 := uint8(0)
		if v.R {
			byte1 = 1 << 7
		}
		byte1 += (v.Source << 3)
		if v.L {
			byte1 += 1 << 2
		}
		byte1 += uint8(v.VEXOpcodeExtension)
		return []uint8{byte0, byte1}
	}

	// Three byte form
	byte0 := uint8(0xc4)

	byte1 := uint8(0)
	if v.R {
		byte1 = 1 << 7
	}
	if v.X {
		byte1 += (1 << 6)
	}
	if v.B {
		byte1 += (1 << 5)
	}
	byte1 += uint8(v.VEXLegacyByte)

	byte2 := uint8(0)
	if v.W {
		byte2 = (1 << 7)
	}
	byte2 += (v.Source << 3)
	if v.L {
		byte2 += (1 << 2)
	}
	byte2 += uint8(v.VEXOpcodeExtension)
	return []uint8{byte0, byte1, byte2}
}
