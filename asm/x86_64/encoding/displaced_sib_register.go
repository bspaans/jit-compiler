package encoding

import (
	"fmt"

	"github.com/bspaans/jit-compiler/lib"
)

type DisplacedSIBRegister struct {
	Scale Scale
	Index *Register
	Base  *Register
	// TODO: also support the 16 bit form
	Displacement uint8
}

func (t *DisplacedSIBRegister) Type() lib.Type {
	return lib.T_DisplacedSIBRegister
}

func (t *DisplacedSIBRegister) String() string {
	return fmt.Sprintf("0x%x(%s, %s, %s)", t.Displacement, t.Base.String(), t.Index.String(), t.Scale.String())
}
