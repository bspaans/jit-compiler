package asm

import (
	"fmt"
	"testing"

	"github.com/bspaans/jit/asm/encoding"
	"github.com/bspaans/jit/lib"
)

func Test_INC(t *testing.T) {
	unit, err := (&INC{encoding.Rax}).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected := "  48 ff c0"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
	unit, err = (&INC{encoding.Rcx}).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected = "  48 ff c1"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}

	unit, err = (&INC{encoding.R14}).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected = "  49 ff c6"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
}

func Test_DEC(t *testing.T) {
	unit, err := (&DEC{encoding.Rax}).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected := "  48 ff c8"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
	unit, err = (&DEC{encoding.Rcx}).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected = "  48 ff c9"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
	unit, err = (&DEC{encoding.R14}).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected = "  49 ff ce"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
}

func Test_MOV(t *testing.T) {
	unit, err := (&MOV{encoding.Rax, encoding.Rax}).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected := "  48 89 c0"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
	unit, err = (&MOV{encoding.Rax, encoding.Rcx}).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected = "  48 89 c1"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
	unit, err = (&MOV{encoding.Rcx, encoding.Rax}).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected = "  48 89 c8"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
	unit, err = (&MOV{encoding.Rax, encoding.R14}).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected = "  49 89 c6"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
	unit, err = (&MOV{encoding.R14, encoding.Rax}).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected = "  4c 89 f0"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
	unit, err = (&MOV{encoding.Uint64(0xffffffff), encoding.Rax}).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected = "  48 b8 ff ff ff ff 00 00 \n  00 00"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
	unit, err = (&MOV{encoding.Uint64(0xffffffff), encoding.Rcx}).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected = "  48 b9 ff ff ff ff 00 00 \n  00 00"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
	unit, err = (&MOV{encoding.Uint64(0xffffffff), encoding.R14}).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected = "  49 be ff ff ff ff 00 00 \n  00 00"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
}

func Test_JMP(t *testing.T) {
	unit, err := JMP(encoding.Uint8(3)).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected := "  eb 03"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
}

func Test_Execute(t *testing.T) {
	units := [][]lib.Instruction{
		[]lib.Instruction{
			&MOV{encoding.Uint32(uint32(5)), &encoding.DisplacedRegister{encoding.Rsp, 8}},
			&RET{},
		},
		[]lib.Instruction{
			&MOV{encoding.Uint64(5), encoding.Rax},
			&MOV{encoding.Rax, &encoding.DisplacedRegister{encoding.Rsp, 8}},
			&RET{},
		},
		[]lib.Instruction{
			&MOV{encoding.Uint64(2), encoding.Rax},
			&MOV{encoding.Uint64(3), encoding.Rdi},
			&ADD{encoding.Rdi, encoding.Rax},
			&MOV{encoding.Rax, &encoding.DisplacedRegister{encoding.Rsp, 8}},
			&RET{},
		},
		[]lib.Instruction{
			&MOV{encoding.Uint64(2), encoding.R13},
			&MOV{encoding.Uint64(3), encoding.R14},
			&ADD{encoding.R13, encoding.R14},
			&MOV{encoding.R14, &encoding.DisplacedRegister{encoding.Rsp, 8}},
			&RET{},
		},
		[]lib.Instruction{
			&MOV{encoding.Uint64(2), encoding.Rax},
			&ADD{encoding.Uint32(3), encoding.Rax},
			&MOV{encoding.Rax, &encoding.DisplacedRegister{encoding.Rsp, 8}},
			&RET{},
		},
		[]lib.Instruction{
			&MOV{encoding.Uint64(5), encoding.Rdi},
			&CVTSI2SD{encoding.Rdi, encoding.Xmm4},
			&MOV{encoding.Uint64(6), encoding.Rdi},
			&CVTTSD2SI{encoding.Xmm4, encoding.Rax},
			&MOV{encoding.Rax, &encoding.DisplacedRegister{encoding.Rsp, 8}},
			&RET{},
		},
		[]lib.Instruction{
			&MOV{encoding.Float64(5.0), encoding.Rdi},
			&MOVQ{encoding.Rdi, encoding.Xmm4},
			&CVTTSD2SI{encoding.Xmm4, encoding.Rax},
			&MOV{encoding.Rax, &encoding.DisplacedRegister{encoding.Rsp, 8}},
			&RET{},
		},
		[]lib.Instruction{
			&MOV{encoding.Uint64(1), encoding.Rdi},
			&CVTSI2SD{encoding.Rdi, encoding.Xmm5},
			&MOV{encoding.Float64(4.0), encoding.Rdi},
			&MOVQ{encoding.Rdi, encoding.Xmm4},
			&ADD{encoding.Xmm5, encoding.Xmm4},
			&CVTTSD2SI{encoding.Xmm4, encoding.Rax},
			&MOV{encoding.Rax, &encoding.DisplacedRegister{encoding.Rsp, 8}},
			&RET{},
		},
		[]lib.Instruction{
			&MOV{encoding.Uint64(2), encoding.Rdi},
			&CVTSI2SD{encoding.Rdi, encoding.Xmm5},
			&MOV{encoding.Float64(7.0), encoding.Rdi},
			&MOVQ{encoding.Rdi, encoding.Xmm4},
			&SUB{encoding.Xmm5, encoding.Xmm4},
			&CVTTSD2SI{encoding.Xmm4, encoding.Rax},
			&MOV{encoding.Rax, &encoding.DisplacedRegister{encoding.Rsp, 8}},
			&RET{},
		},
		[]lib.Instruction{
			&MOV{encoding.Float64(2.0), encoding.Rcx},
			&MOV{encoding.Float64(2.5), encoding.Rdi},
			&MOVQ{encoding.Rcx, encoding.Xmm4},
			&MOVQ{encoding.Rdi, encoding.Xmm5},
			&MUL{encoding.Xmm5, encoding.Xmm4},
			&CVTTSD2SI{encoding.Xmm4, encoding.Rax},
			&MOV{encoding.Rax, &encoding.DisplacedRegister{encoding.Rsp, 8}},
			&RET{},
		},
		[]lib.Instruction{
			&MOV{encoding.Float64(2.0), encoding.Rcx},
			&MOV{encoding.Float64(10.0), encoding.Rdi},
			&MOVQ{encoding.Rdi, encoding.Xmm4},
			&MOVQ{encoding.Rcx, encoding.Xmm5},
			&DIV{encoding.Xmm5, encoding.Xmm4},
			&CVTTSD2SI{encoding.Xmm4, encoding.Rax},
			&MOV{encoding.Rax, &encoding.DisplacedRegister{encoding.Rsp, 8}},
			&RET{},
		},
		[]lib.Instruction{
			&MOV{encoding.Float64(2.0), encoding.Rcx},
			&MOV{encoding.Float64(10.0), encoding.Rdi},
			&MOVQ{encoding.Rdi, encoding.Xmm4},
			&MOVQ{encoding.Rcx, encoding.Xmm5},
			&MOVSD{encoding.Xmm4, encoding.Xmm0},
			&MOVSD{encoding.Xmm5, encoding.Xmm1},
			&DIV{encoding.Xmm1, encoding.Xmm0},
			&CVTTSD2SI{encoding.Xmm0, encoding.Rax},
			&MOV{encoding.Rax, &encoding.DisplacedRegister{encoding.Rsp, 8}},
			&RET{},
		},
	}
	for _, unit := range units {
		b, err := lib.CompileInstruction(unit)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(b)
		value := b.Execute()
		if value != uint(5) {
			t.Fatal("Expecting 5 got", value, "in", unit)
		}
	}

}
