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

var (
	rax *Register = NewRegister("rax", 0, QUADWORD)
	rcx *Register = NewRegister("rcx", 1, QUADWORD)
	rdx *Register = NewRegister("rdx", 2, QUADWORD)
	rbx *Register = NewRegister("rbx", 3, QUADWORD)
	rsp *Register = NewRegister("rsp", 4, QUADWORD)
	rbp *Register = NewRegister("rbp", 5, QUADWORD)
	rsi *Register = NewRegister("rsi", 6, QUADWORD)
	rdi *Register = NewRegister("rdi", 7, QUADWORD)
	r8  *Register = NewRegister("r8", 8, QUADWORD)
	r9  *Register = NewRegister("r9", 9, QUADWORD)
	r10 *Register = NewRegister("r10", 10, QUADWORD)
	r11 *Register = NewRegister("r11", 11, QUADWORD)
	r12 *Register = NewRegister("r12", 12, QUADWORD)
	r13 *Register = NewRegister("r13", 13, QUADWORD)
	r14 *Register = NewRegister("r14", 14, QUADWORD)
	r15 *Register = NewRegister("r15", 15, QUADWORD)

	eax  *Register = NewRegister("eax", 0, DOUBLE)
	ecx  *Register = NewRegister("ecx", 1, DOUBLE)
	edx  *Register = NewRegister("edx", 2, DOUBLE)
	ebx  *Register = NewRegister("ebx", 3, DOUBLE)
	esp  *Register = NewRegister("esp", 4, DOUBLE)
	ebp  *Register = NewRegister("ebp", 5, DOUBLE)
	esi  *Register = NewRegister("esi", 6, DOUBLE)
	edi  *Register = NewRegister("edi", 7, DOUBLE)
	r8d  *Register = NewRegister("r8d", 8, DOUBLE)
	r9d  *Register = NewRegister("r9d", 9, DOUBLE)
	r10d *Register = NewRegister("r10d", 10, DOUBLE)
	r11d *Register = NewRegister("r11d", 11, DOUBLE)
	r12d *Register = NewRegister("r12d", 12, DOUBLE)
	r13d *Register = NewRegister("r13d", 13, DOUBLE)
	r14d *Register = NewRegister("r14d", 14, DOUBLE)
	r15d *Register = NewRegister("r15d", 15, DOUBLE)
)

/*

Registers:

64 bit

rax // syscall number / return
rcx // used to pass fourth argument to functions
rdx // used to pass third argument to functions
rbx
rsp // stack pointer
rbp
rsi // pointer used to pass 2nd argument to functions
rdi // used to pass first argument to functions
r8 // used to pass fifth arg
r9 // ,,   ,,  ,,  sixth ,,
r10
r11
r12
r13
r14
r15

lower 32 bits:

eax
ebx
ecx
edx
esi
edi
ebp
esp
r8d
r9d
r10d
r11d
r12d
r13d
r14d
r15d

lower 16 bits:

ax
bx
cx
dx
si
di
bp
sp
r8w
r9w
r10w
r11w
r12w
r13w
r14w
r15w

lower 8 bits:

al
bl
cl
dl
sil
dil
bpl
spl
r8b
r9b
r10b
r11b
r12b
r13b
r14b
r15b

*/
