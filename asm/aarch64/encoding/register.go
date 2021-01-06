package encoding

import (
	. "github.com/bspaans/jit-compiler/lib"
)

type Register struct {
	Name     string
	Register uint8
	Size     Size
}

func NewRegister(name string, register uint8, size Size) *Register {
	return &Register{
		Name:     name,
		Register: register,
		Size:     size,
	}
}

func (r *Register) Encode() uint8 {
	return r.Register
}
func (r *Register) String() string {
	return "%" + r.Name
}
func (r *Register) Type() Type {
	return T_Register
}
func (r *Register) Width() Size {
	return r.Size
}
func (r *Register) ForOperandWidth(w Size) *Register {
	// TODO: support more sizes
	if w == DOUBLE {
		return r.Get32BitRegister()
	}
	return r
}
func (r *Register) Get32BitRegister() *Register {
	return Registers32[r.Register]
}

func Get64BitRegisterByIndex(ix uint8) *Register {
	return Registers64[ix]
}

// TODO
func GetFloatingPointRegisterByIndex(ix uint8) *Register {
	return nil
}

var Registers64 []*Register = []*Register{
	X0, X1, X2, X3, X4, X5, X6, X7, X8,
	X9, X10, X11, X12, X13, X14, X15, X16,
	X17, X18, X19, X20, X21, X22, X23, X24,
	X25, X26, X27, X28, X29, X30,
}

var Registers32 []*Register = []*Register{
	W0, W1, W2, W3, W4, W5, W6, W7, W8,
	W9, W10, W11, W12, W13, W14, W15, W16,
	W17, W18, W19, W20, W21, W22, W23, W24,
	W25, W26, W27, W28, W29, W30,
}

var (
	X0  *Register = NewRegister("x0", 0, QUADWORD)
	X1  *Register = NewRegister("x1", 1, QUADWORD)
	X2  *Register = NewRegister("x2", 2, QUADWORD)
	X3  *Register = NewRegister("x3", 3, QUADWORD)
	X4  *Register = NewRegister("x4", 4, QUADWORD)
	X5  *Register = NewRegister("x5", 5, QUADWORD)
	X6  *Register = NewRegister("x6", 6, QUADWORD)
	X7  *Register = NewRegister("x7", 7, QUADWORD)
	X8  *Register = NewRegister("x8", 8, QUADWORD)
	X9  *Register = NewRegister("x9", 9, QUADWORD)
	X10 *Register = NewRegister("x10", 10, QUADWORD)
	X11 *Register = NewRegister("x11", 11, QUADWORD)
	X12 *Register = NewRegister("x12", 12, QUADWORD)
	X13 *Register = NewRegister("x13", 13, QUADWORD)
	X14 *Register = NewRegister("x14", 14, QUADWORD)
	X15 *Register = NewRegister("x15", 15, QUADWORD)
	X16 *Register = NewRegister("x16", 16, QUADWORD)
	X17 *Register = NewRegister("x17", 17, QUADWORD)
	X18 *Register = NewRegister("x18", 18, QUADWORD)
	X19 *Register = NewRegister("x19", 19, QUADWORD)
	X20 *Register = NewRegister("x20", 20, QUADWORD)
	X21 *Register = NewRegister("x21", 21, QUADWORD)
	X22 *Register = NewRegister("x22", 22, QUADWORD)
	X23 *Register = NewRegister("x23", 23, QUADWORD)
	X24 *Register = NewRegister("x24", 24, QUADWORD)
	X25 *Register = NewRegister("x25", 25, QUADWORD)
	X26 *Register = NewRegister("x26", 26, QUADWORD)
	X27 *Register = NewRegister("x27", 27, QUADWORD)
	X28 *Register = NewRegister("x28", 28, QUADWORD)
	X29 *Register = NewRegister("x29", 29, QUADWORD)
	X30 *Register = NewRegister("x30", 30, QUADWORD)

	// TODO program counter

	SP  *Register = NewRegister("sp", 31, QUADWORD) // stack pointer
	XZR *Register = NewRegister("zr", 31, QUADWORD) // zero register

	W0  *Register = NewRegister("w0", 0, DOUBLE)
	W1  *Register = NewRegister("w1", 1, DOUBLE)
	W2  *Register = NewRegister("w2", 2, DOUBLE)
	W3  *Register = NewRegister("w3", 3, DOUBLE)
	W4  *Register = NewRegister("w4", 4, DOUBLE)
	W5  *Register = NewRegister("w5", 5, DOUBLE)
	W6  *Register = NewRegister("w6", 6, DOUBLE)
	W7  *Register = NewRegister("w7", 7, DOUBLE)
	W8  *Register = NewRegister("w8", 8, DOUBLE)
	W9  *Register = NewRegister("w9", 9, DOUBLE)
	W10 *Register = NewRegister("w10", 10, DOUBLE)
	W11 *Register = NewRegister("w11", 11, DOUBLE)
	W12 *Register = NewRegister("w12", 12, DOUBLE)
	W13 *Register = NewRegister("w13", 13, DOUBLE)
	W14 *Register = NewRegister("w14", 14, DOUBLE)
	W15 *Register = NewRegister("w15", 15, DOUBLE)
	W16 *Register = NewRegister("w16", 16, DOUBLE)
	W17 *Register = NewRegister("w17", 17, DOUBLE)
	W18 *Register = NewRegister("w18", 18, DOUBLE)
	W19 *Register = NewRegister("w19", 19, DOUBLE)
	W20 *Register = NewRegister("w20", 20, DOUBLE)
	W21 *Register = NewRegister("w21", 21, DOUBLE)
	W22 *Register = NewRegister("w22", 22, DOUBLE)
	W23 *Register = NewRegister("w23", 23, DOUBLE)
	W24 *Register = NewRegister("w24", 24, DOUBLE)
	W25 *Register = NewRegister("w25", 25, DOUBLE)
	W26 *Register = NewRegister("w26", 26, DOUBLE)
	W27 *Register = NewRegister("w27", 27, DOUBLE)
	W28 *Register = NewRegister("w28", 28, DOUBLE)
	W29 *Register = NewRegister("w29", 29, DOUBLE)
	W30 *Register = NewRegister("w30", 30, DOUBLE)
	WSP *Register = NewRegister("wsp", 31, DOUBLE) // current stack pointer
	WZR *Register = NewRegister("wzr", 31, DOUBLE) // zero register

	// TODO simd and FP registers
)
