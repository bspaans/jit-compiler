package encoding

import (
	"fmt"

	"github.com/bspaans/jit-compiler/lib"
)

// Get address relative to instruction pointer
type RIPRelative struct {
	Displacement Int32
}

func (t *RIPRelative) Type() Type {
	return T_RIPRelative
}

func (t *RIPRelative) String() string {
	if t.Displacement < 0 {
		return fmt.Sprintf("-$0x%x(%%rip)", int(t.Displacement)*-1)
	} else {
		return fmt.Sprintf("$0x%x(%%rip)", t.Displacement)
	}
}
func (t *RIPRelative) Width() lib.Size {
	return lib.QUADWORD
}
