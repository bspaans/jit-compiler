package shared

import (
	"github.com/bspaans/jit-compiler/asm/x86_64/encoding"
	"github.com/bspaans/jit-compiler/lib"
)

// TODO: should only contain call setup and teardown; this makes no sense.
type ABI interface {
	GetRegistersForArgs(args []Type) []*encoding.Register
	ReturnTypeToOperand(ty Type) lib.Operand
}
