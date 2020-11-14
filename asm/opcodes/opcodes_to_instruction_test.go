package opcodes

import (
	"testing"

	"github.com/bspaans/jit-compiler/asm/encoding"
)

func Test_OpcodeToInstruction_happy(t *testing.T) {
	result := OpcodeToInstruction("mov", MOV_rm8_r8, 2, encoding.Al, encoding.Cl)

	expected := "mov %cl, %al"
	if result.String() != expected {
		t.Fatal("expecting", expected, "got", result.String())
	}
}
