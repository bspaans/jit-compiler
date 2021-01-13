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
	if w == BYTE {
		return r.Get8BitRegister()
	} else if w == WORD {
		return r.Get16BitRegister()
	} else if w == DOUBLE {
		return r.Get32BitRegister()
	}
	return r
}

func (r *Register) Get64BitRegister() *Register {
	return Get64BitRegisterByIndex(r.Register)
}

func (r *Register) Get8BitRegister() *Register {
	return Registers8[r.Register]
}

func (r *Register) Get16BitRegister() *Register {
	return Registers16[r.Register]
}

func (r *Register) Get32BitRegister() *Register {
	return Registers32[r.Register]
}

func Get64BitRegisterByIndex(ix uint8) *Register {
	return Registers64[ix]
}

func GetFloatingPointRegisterByIndex(ix uint8) *Register {
	return Registers128[ix]
}

var Registers64 []*Register = []*Register{
	Rax, Rcx, Rdx, Rbx, Rsp, Rbp, Rsi, Rdi,
	R8, R9, R10, R11, R12, R13, R14, R15,
}
var Registers32 []*Register = []*Register{
	Eax, Ecx, Edx, Ebx, Esp, Ebp, Esi, Edi,
	R8d, R9d, R10d, R11d, R12d, R13d, R14d, R15d,
}
var Registers16 []*Register = []*Register{
	Ax, Cx, Dx, Bx, Sp, Bp, Si, Di,
	R8w, R9w, R10w, R11w, R12w, R13w, R14w, R15w,
}

// Note: does not include Ah, Bh, Ch, Dh
var Registers8 []*Register = []*Register{
	Al, Cl, Dl, Bl, Spl, Bpl, Sil, Dil,
	R8b, R9b, R10b, R11b, R12b, R13b, R14b, R15b,
}

var Registers128 []*Register = []*Register{
	Xmm0, Xmm1, Xmm2, Xmm3, Xmm4, Xmm5, Xmm6, Xmm7,
	Xmm8, Xmm9, Xmm10, Xmm11, Xmm11, Xmm12, Xmm13, Xmm14, Xmm15,
}

var Registers256 []*Register = []*Register{
	Ymm0, Ymm1, Ymm2, Ymm3, Ymm4, Ymm5, Ymm6, Ymm7,
	Ymm8, Ymm9, Ymm10, Ymm11, Ymm11, Ymm12, Ymm13, Ymm14, Ymm15,
}

var Registers512 []*Register = []*Register{
	Zmm0, Zmm1, Zmm2, Zmm3, Zmm4, Zmm5, Zmm6, Zmm7,
	Zmm8, Zmm9, Zmm10, Zmm11, Zmm11, Zmm12, Zmm13, Zmm14, Zmm15,
}

var (
	Rax *Register = NewRegister("rax", 0, QUADWORD)
	Rcx *Register = NewRegister("rcx", 1, QUADWORD)
	Rdx *Register = NewRegister("rdx", 2, QUADWORD)
	Rbx *Register = NewRegister("rbx", 3, QUADWORD)
	Rsp *Register = NewRegister("rsp", 4, QUADWORD) // stack pointer
	Rbp *Register = NewRegister("rbp", 5, QUADWORD) // frame pointer
	Rsi *Register = NewRegister("rsi", 6, QUADWORD)
	Rdi *Register = NewRegister("rdi", 7, QUADWORD)
	R8  *Register = NewRegister("r8", 8, QUADWORD)
	R9  *Register = NewRegister("r9", 9, QUADWORD)
	R10 *Register = NewRegister("r10", 10, QUADWORD)
	R11 *Register = NewRegister("r11", 11, QUADWORD)
	R12 *Register = NewRegister("r12", 12, QUADWORD)
	R13 *Register = NewRegister("r13", 13, QUADWORD)
	R14 *Register = NewRegister("r14", 14, QUADWORD)
	R15 *Register = NewRegister("r15", 15, QUADWORD)

	Eax  *Register = NewRegister("eax", 0, DOUBLE)
	Ecx  *Register = NewRegister("ecx", 1, DOUBLE)
	Edx  *Register = NewRegister("edx", 2, DOUBLE)
	Ebx  *Register = NewRegister("ebx", 3, DOUBLE)
	Esp  *Register = NewRegister("esp", 4, DOUBLE)
	Ebp  *Register = NewRegister("ebp", 5, DOUBLE)
	Esi  *Register = NewRegister("esi", 6, DOUBLE)
	Edi  *Register = NewRegister("edi", 7, DOUBLE)
	R8d  *Register = NewRegister("r8d", 8, DOUBLE)
	R9d  *Register = NewRegister("r9d", 9, DOUBLE)
	R10d *Register = NewRegister("r10d", 10, DOUBLE)
	R11d *Register = NewRegister("r11d", 11, DOUBLE)
	R12d *Register = NewRegister("r12d", 12, DOUBLE)
	R13d *Register = NewRegister("r13d", 13, DOUBLE)
	R14d *Register = NewRegister("r14d", 14, DOUBLE)
	R15d *Register = NewRegister("r15d", 15, DOUBLE)

	Ax   *Register = NewRegister("ax", 0, WORD)
	Cx   *Register = NewRegister("cx", 1, WORD)
	Dx   *Register = NewRegister("dx", 2, WORD)
	Bx   *Register = NewRegister("bx", 3, WORD)
	Sp   *Register = NewRegister("sp", 4, WORD)
	Bp   *Register = NewRegister("bp", 5, WORD)
	Si   *Register = NewRegister("si", 6, WORD)
	Di   *Register = NewRegister("di", 7, WORD)
	R8w  *Register = NewRegister("r8w", 8, WORD)
	R9w  *Register = NewRegister("r9w", 9, WORD)
	R10w *Register = NewRegister("r10w", 10, WORD)
	R11w *Register = NewRegister("r11w", 11, WORD)
	R12w *Register = NewRegister("r12w", 12, WORD)
	R13w *Register = NewRegister("r13w", 13, WORD)
	R14w *Register = NewRegister("r14w", 14, WORD)
	R15w *Register = NewRegister("r15w", 15, WORD)

	Al   *Register = NewRegister("al", 0, BYTE)
	Cl   *Register = NewRegister("cl", 1, BYTE)
	Dl   *Register = NewRegister("dl", 2, BYTE)
	Bl   *Register = NewRegister("bl", 3, BYTE)
	Spl  *Register = NewRegister("spl", 4, BYTE)
	Bpl  *Register = NewRegister("bpl", 5, BYTE)
	Sil  *Register = NewRegister("sil", 6, BYTE)
	Dil  *Register = NewRegister("dil", 7, BYTE)
	R8b  *Register = NewRegister("r8b", 8, BYTE)
	R9b  *Register = NewRegister("r9b", 9, BYTE)
	R10b *Register = NewRegister("r10b", 10, BYTE)
	R11b *Register = NewRegister("r11b", 11, BYTE)
	R12b *Register = NewRegister("r12b", 12, BYTE)
	R13b *Register = NewRegister("r13b", 13, BYTE)
	R14b *Register = NewRegister("r14b", 14, BYTE)
	R15b *Register = NewRegister("r15b", 15, BYTE)

	Ah *Register = NewRegister("ah", 4, BYTE)
	Ch *Register = NewRegister("ch", 5, BYTE)
	Dh *Register = NewRegister("dh", 6, BYTE)
	Bh *Register = NewRegister("bh", 7, BYTE)

	Xmm0  *Register = NewRegister("xmm0", 0, OWORD)
	Xmm1  *Register = NewRegister("xmm1", 1, OWORD)
	Xmm2  *Register = NewRegister("xmm2", 2, OWORD)
	Xmm3  *Register = NewRegister("xmm3", 3, OWORD)
	Xmm4  *Register = NewRegister("xmm4", 4, OWORD)
	Xmm5  *Register = NewRegister("xmm5", 5, OWORD)
	Xmm6  *Register = NewRegister("xmm6", 6, OWORD)
	Xmm7  *Register = NewRegister("xmm7", 7, OWORD)
	Xmm8  *Register = NewRegister("xmm8", 8, OWORD)
	Xmm9  *Register = NewRegister("xmm9", 9, OWORD)
	Xmm10 *Register = NewRegister("xmm10", 10, OWORD)
	Xmm11 *Register = NewRegister("xmm11", 11, OWORD)
	Xmm12 *Register = NewRegister("xmm12", 12, OWORD)
	Xmm13 *Register = NewRegister("xmm13", 13, OWORD)
	Xmm14 *Register = NewRegister("xmm14", 14, OWORD)
	Xmm15 *Register = NewRegister("xmm15", 15, OWORD)

	Ymm0  *Register = NewRegister("ymm0", 0, YWORD)
	Ymm1  *Register = NewRegister("ymm1", 1, YWORD)
	Ymm2  *Register = NewRegister("ymm2", 2, YWORD)
	Ymm3  *Register = NewRegister("ymm3", 3, YWORD)
	Ymm4  *Register = NewRegister("ymm4", 4, YWORD)
	Ymm5  *Register = NewRegister("ymm5", 5, YWORD)
	Ymm6  *Register = NewRegister("ymm6", 6, YWORD)
	Ymm7  *Register = NewRegister("ymm7", 7, YWORD)
	Ymm8  *Register = NewRegister("ymm8", 8, YWORD)
	Ymm9  *Register = NewRegister("ymm9", 9, YWORD)
	Ymm10 *Register = NewRegister("ymm10", 10, YWORD)
	Ymm11 *Register = NewRegister("ymm11", 11, YWORD)
	Ymm12 *Register = NewRegister("ymm12", 12, YWORD)
	Ymm13 *Register = NewRegister("ymm13", 13, YWORD)
	Ymm14 *Register = NewRegister("ymm14", 14, YWORD)
	Ymm15 *Register = NewRegister("ymm15", 15, YWORD)

	Zmm0  *Register = NewRegister("zmm0", 0, ZWORD)
	Zmm1  *Register = NewRegister("zmm1", 1, ZWORD)
	Zmm2  *Register = NewRegister("zmm2", 2, ZWORD)
	Zmm3  *Register = NewRegister("zmm3", 3, ZWORD)
	Zmm4  *Register = NewRegister("zmm4", 4, ZWORD)
	Zmm5  *Register = NewRegister("zmm5", 5, ZWORD)
	Zmm6  *Register = NewRegister("zmm6", 6, ZWORD)
	Zmm7  *Register = NewRegister("zmm7", 7, ZWORD)
	Zmm8  *Register = NewRegister("zmm8", 8, ZWORD)
	Zmm9  *Register = NewRegister("zmm9", 9, ZWORD)
	Zmm10 *Register = NewRegister("zmm10", 10, ZWORD)
	Zmm11 *Register = NewRegister("zmm11", 11, ZWORD)
	Zmm12 *Register = NewRegister("zmm12", 12, ZWORD)
	Zmm13 *Register = NewRegister("zmm13", 13, ZWORD)
	Zmm14 *Register = NewRegister("zmm14", 14, ZWORD)
	Zmm15 *Register = NewRegister("zmm15", 15, ZWORD)
)
