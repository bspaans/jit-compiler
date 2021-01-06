package shared

import (
	"github.com/bspaans/jit-compiler/asm/x86_64/encoding"
)

type ABI interface {
	GetRegistersForArgs(args []Type) []*encoding.Register
	ReturnTypeToOperand(ty Type) encoding.Operand
}
