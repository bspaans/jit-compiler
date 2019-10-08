package asm

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
func (r *Register) Lower8BitRegister() *Register {
	for _, reg := range Registers8 {
		if reg.Register == r.Register {
			return reg
		}
	}
	return nil
}

func Get64BitRegisterByIndex(ix uint8) *Register {
	for _, reg := range Registers64 {
		if reg.Register == ix {
			return reg
		}
	}
	return nil
}

func GetFloatingPointRegisterByIndex(ix uint8) *Register {
	for _, reg := range RegistersSSE {
		if reg.Register == ix {
			return reg
		}
	}
	return nil
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
var Registers8 []*Register = []*Register{
	Al, Cl, Dl, Bl, Spl, Bpl, Sil, Dil,
	R8b, R9b, R10b, R11b, R12b, R13b, R14b, R15b,
}

var RegistersSSE []*Register = []*Register{
	Xmm0, Xmm1, Xmm2, Xmm3, Xmm4, Xmm5, Xmm6, Xmm7,
	Xmm8, Xmm9, Xmm10, Xmm11, Xmm11, Xmm12, Xmm13, Xmm14, Xmm15,
}

var (
	Rax *Register = NewRegister("rax", 0, QUADWORD)
	Rcx *Register = NewRegister("rcx", 1, QUADWORD)
	Rdx *Register = NewRegister("rdx", 2, QUADWORD)
	Rbx *Register = NewRegister("rbx", 3, QUADWORD)
	Rsp *Register = NewRegister("rsp", 4, QUADWORD)
	Rbp *Register = NewRegister("rbp", 5, QUADWORD)
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

	Xmm0  *Register = NewRegister("xmm0", 0, QUADDOUBLE)
	Xmm1  *Register = NewRegister("xmm1", 1, QUADDOUBLE)
	Xmm2  *Register = NewRegister("xmm2", 2, QUADDOUBLE)
	Xmm3  *Register = NewRegister("xmm3", 3, QUADDOUBLE)
	Xmm4  *Register = NewRegister("xmm4", 4, QUADDOUBLE)
	Xmm5  *Register = NewRegister("xmm5", 5, QUADDOUBLE)
	Xmm6  *Register = NewRegister("xmm6", 6, QUADDOUBLE)
	Xmm7  *Register = NewRegister("xmm7", 7, QUADDOUBLE)
	Xmm8  *Register = NewRegister("xmm8", 8, QUADDOUBLE)
	Xmm9  *Register = NewRegister("xmm9", 9, QUADDOUBLE)
	Xmm10 *Register = NewRegister("xmm10", 10, QUADDOUBLE)
	Xmm11 *Register = NewRegister("xmm11", 11, QUADDOUBLE)
	Xmm12 *Register = NewRegister("xmm12", 12, QUADDOUBLE)
	Xmm13 *Register = NewRegister("xmm13", 13, QUADDOUBLE)
	Xmm14 *Register = NewRegister("xmm14", 14, QUADDOUBLE)
	Xmm15 *Register = NewRegister("xmm15", 15, QUADDOUBLE)

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
)
