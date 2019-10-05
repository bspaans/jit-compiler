package asm

type Type uint8

const (
	T_Register          Type = 0
	T_DisplacedRegister Type = iota
	T_RIPRelative       Type = iota
	T_Uint8             Type = iota
	T_Uint16            Type = iota
	T_Uint32            Type = iota
	T_Uint64            Type = iota
	T_Int32             Type = iota
)

type Operand interface {
	Type() Type
	String() string
}
