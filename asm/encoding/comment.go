package encoding

import "github.com/bspaans/jit-compiler/lib"

type Comment string

func (c Comment) Encode() (lib.MachineCode, error) {
	return []uint8{}, nil
}

func (c Comment) String() string {
	return "# " + string(c)
}
