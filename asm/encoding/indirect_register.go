package encoding

import "fmt"

type IndirectRegister struct {
	*Register
}

func (t *IndirectRegister) Type() Type {
	return T_IndirectRegister
}

func (t *IndirectRegister) String() string {
	return fmt.Sprintf("(%s)", t.Register.String())
}
