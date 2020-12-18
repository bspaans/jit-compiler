package encoding

import "github.com/bspaans/jit-compiler/lib"

type Type uint8

const (
	T_Register Type = iota
	T_Uint12   Type = iota
)

type Operand interface {
	Type() Type
	String() string
	Width() lib.Size
}
