package encoding

import "github.com/bspaans/jit-compiler/lib"

//go:generate stringer -type=Type
type Type uint8

const (
	T_Register          Type = iota // e.g. %rax
	T_IndirectRegister  Type = iota // e.g. (%rax)
	T_RIPRelative       Type = iota // e.g. -$0x18(%rip)
	T_SIBRegister       Type = iota // e.g. (%rax, %rcx, 8)  (the address of %rxc * 8 + %rax)
	T_DisplacedRegister Type = iota // e.g. 0x9(%rax)
	// TODO
	T_DisplacedSIBRegister Type = iota // e.g. 0x9(%rax, %rcx, 8) the address of %rcx * 8 + %rax + 9)
	T_Uint8                Type = iota
	T_Uint16               Type = iota
	T_Uint32               Type = iota
	T_Uint64               Type = iota
	T_Int32                Type = iota
	T_Float32              Type = iota
	T_Float64              Type = iota
)

type Operand interface {
	Type() Type
	String() string
	Width() lib.Size
}

func IsRegister(op Operand) bool {
	t := op.Type()
	return t == T_Register || t == T_IndirectRegister || t == T_DisplacedRegister || t == T_RIPRelative || t == T_SIBRegister
}

func IsInt(op Operand) bool {
	t := op.Type()
	return t == T_Uint8 || t == T_Uint16 || t == T_Uint32 || t == T_Uint64
}
