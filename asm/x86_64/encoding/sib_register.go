package encoding

import (
	"fmt"

	"github.com/bspaans/jit-compiler/lib"
)

type SIBRegister struct {
	*Register
	Index *Register
	Scale Scale
}

func (t *SIBRegister) Type() lib.Type {
	return lib.T_SIBRegister
}

func (t *SIBRegister) String() string {
	return fmt.Sprintf("(%s, %s, %s)", t.Register.String(), t.Index.String(), t.Scale.String())
}
