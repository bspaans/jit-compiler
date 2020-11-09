package encoding

import "github.com/bspaans/jit-compiler/lib"

//go:generate stringer -type=Type
type Type uint8

const (
	T_Register          Type = iota
	T_IndirectRegister  Type = iota
	T_DisplacedRegister Type = iota
	T_RIPRelative       Type = iota
	T_Uint8             Type = iota
	T_Uint16            Type = iota
	T_Uint32            Type = iota
	T_Uint64            Type = iota
	T_Int32             Type = iota
	T_Float32           Type = iota
	T_Float64           Type = iota
)

type Operand interface {
	Type() Type
	String() string
	Width() lib.Size
}

func IsRegister(op Operand) bool {
	t := op.Type()
	return t == T_Register || t == T_IndirectRegister || t == T_DisplacedRegister || t == T_RIPRelative
}

func IsInt(op Operand) bool {
	t := op.Type()
	return t == T_Uint8 || t == T_Uint16 || t == T_Uint32 || t == T_Uint64
}
