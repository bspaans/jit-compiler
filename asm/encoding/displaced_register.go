package encoding

import "fmt"

type DisplacedRegister struct {
	*Register
	// TODO: also support the 16 bit form
	Displacement uint8
}

func (t *DisplacedRegister) Type() Type {
	return T_DisplacedRegister
}

func (t *DisplacedRegister) String() string {
	return fmt.Sprintf("0x%x(%s)", t.Displacement, t.Register.String())
}
