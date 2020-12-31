package encoding

import "fmt"

type SIBRegister struct {
	*Register
	Index *Register
	Scale Scale
}

func (t *SIBRegister) Type() Type {
	return T_SIBRegister
}

func (t *SIBRegister) String() string {
	return fmt.Sprintf("(%s, %s, %s)", t.Register.String(), t.Index.String(), t.Scale.String())
}
