package encoding

import (
	"fmt"

	"github.com/bspaans/jit-compiler/lib"
)

type IndirectRegister struct {
	*Register
}

func (t *IndirectRegister) Type() lib.Type {
	return lib.T_IndirectRegister
}

func (t *IndirectRegister) String() string {
	return fmt.Sprintf("(%s)", t.Register.String())
}
