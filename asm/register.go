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

func Get64BitRegisterByIndex(ix uint8) *Register {
	for _, reg := range Registers64 {
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

	Ax   *Register = NewRegister("eax", 0, WORD)
	Cx   *Register = NewRegister("ecx", 1, WORD)
	Dx   *Register = NewRegister("edx", 2, WORD)
	Bx   *Register = NewRegister("ebx", 3, WORD)
	Sp   *Register = NewRegister("esp", 4, WORD)
	Bp   *Register = NewRegister("ebp", 5, WORD)
	Si   *Register = NewRegister("esi", 6, WORD)
	Di   *Register = NewRegister("edi", 7, WORD)
	R8w  *Register = NewRegister("r8d", 8, WORD)
	R9w  *Register = NewRegister("r9d", 9, WORD)
	R10w *Register = NewRegister("r10d", 10, WORD)
	R11w *Register = NewRegister("r11d", 11, WORD)
	R12w *Register = NewRegister("r12d", 12, WORD)
	R13w *Register = NewRegister("r13d", 13, WORD)
	R14w *Register = NewRegister("r14d", 14, WORD)
	R15w *Register = NewRegister("r15d", 15, WORD)

	Al   *Register = NewRegister("eax", 0, BYTE)
	Cl   *Register = NewRegister("ecx", 1, BYTE)
	Dl   *Register = NewRegister("edx", 2, BYTE)
	Bl   *Register = NewRegister("ebx", 3, BYTE)
	Spl  *Register = NewRegister("esp", 4, BYTE)
	Bpl  *Register = NewRegister("ebp", 5, BYTE)
	Sil  *Register = NewRegister("esi", 6, BYTE)
	Dil  *Register = NewRegister("edi", 7, BYTE)
	R8b  *Register = NewRegister("r8d", 8, BYTE)
	R9b  *Register = NewRegister("r9d", 9, BYTE)
	R10b *Register = NewRegister("r10d", 10, BYTE)
	R11b *Register = NewRegister("r11d", 11, BYTE)
	R12b *Register = NewRegister("r12d", 12, BYTE)
	R13b *Register = NewRegister("r13d", 13, BYTE)
	R14b *Register = NewRegister("r14d", 14, BYTE)
	R15b *Register = NewRegister("r15d", 15, BYTE)
)
