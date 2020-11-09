package asm

import (
	"testing"

	"github.com/bspaans/jit-compiler/asm/encoding"
	"github.com/bspaans/jit-compiler/lib"
)

func Test_INC(t *testing.T) {
	unit, err := (INC(encoding.Rax)).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected := "  48 ff c0"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
	unit, err = (INC(encoding.Rcx)).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected = "  48 ff c1"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}

	unit, err = (INC(encoding.R14)).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected = "  49 ff c6"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
}

func Test_DEC(t *testing.T) {
	unit, err := (DEC(encoding.Rax)).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected := "  48 ff c8"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
	unit, err = (DEC(encoding.Rcx)).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected = "  48 ff c9"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
	unit, err = (DEC(encoding.R14)).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected = "  49 ff ce"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
}

func Test_MOV(t *testing.T) {
	unit, err := (MOV(encoding.Rax, encoding.Rax)).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected := "  48 89 c0"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
	unit, err = (MOV(encoding.Rax, encoding.Rcx)).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected = "  48 89 c1"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
	unit, err = (MOV(encoding.Rcx, encoding.Rax)).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected = "  48 89 c8"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
	unit, err = (MOV(encoding.Rax, encoding.R14)).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected = "  49 89 c6"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
	unit, err = (MOV(encoding.R14, encoding.Rax)).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected = "  4c 89 f0"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
	unit, err = (MOV(encoding.Uint64(0xffffffff), encoding.Rax)).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected = "  48 b8 ff ff ff ff 00 00 \n  00 00"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
	unit, err = (MOV(encoding.Uint64(0xffffffff), encoding.Rcx)).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected = "  48 b9 ff ff ff ff 00 00 \n  00 00"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
	unit, err = (MOV(encoding.Uint64(0xffffffff), encoding.R14)).Encode()
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

func Test_SIB_Addressing(t *testing.T) {
	//unit, err := MOV(encoding.Rax, &encoding.SIBRegister{encoding.Rcx, encoding.Rax, encoding.Scale8}).Encode()
	table := [][]interface{}{
		[]interface{}{&encoding.SIBRegister{encoding.Rcx, encoding.Rax, encoding.Scale8}, encoding.Rax, "  48 8b 04 c1"},
		[]interface{}{&encoding.SIBRegister{encoding.Rcx, encoding.R9, encoding.Scale8}, encoding.Rax, "  4a 8b 04 c9"},
		[]interface{}{&encoding.SIBRegister{encoding.R9, encoding.Rax, encoding.Scale8}, encoding.Rax, "  49 8b 04 c1"},
		[]interface{}{&encoding.SIBRegister{encoding.R9, encoding.R10, encoding.Scale8}, encoding.Rax, "  4b 8b 04 d1"},
		// There is a special case for register 13, because the encoding
		// interferes with RIP relative encoding.  Need to use a 0 displacement
		[]interface{}{&encoding.SIBRegister{encoding.R13, encoding.R9, encoding.Scale8}, encoding.Rax, "  4b 8b 44 cd 00"},

		[]interface{}{encoding.Rax, &encoding.SIBRegister{encoding.Rcx, encoding.Rax, encoding.Scale8}, "  48 89 04 c1"},
		[]interface{}{encoding.Rax, &encoding.SIBRegister{encoding.Rcx, encoding.R9, encoding.Scale8}, "  4a 89 04 c9"},
		[]interface{}{encoding.Rax, &encoding.SIBRegister{encoding.R9, encoding.R9, encoding.Scale8}, "  4b 89 04 c9"},
		[]interface{}{encoding.Rax, &encoding.SIBRegister{encoding.R13, encoding.R9, encoding.Scale8}, "  4b 89 44 cd 00"},

		[]interface{}{encoding.Al, &encoding.SIBRegister{encoding.Rax, encoding.Rcx, encoding.Scale8}, "  40 88 04 c8"},
	}
	for _, testCase := range table {
		unit, err := MOV(testCase[0].(encoding.Operand), testCase[1].(encoding.Operand)).Encode()
		if err != nil {
			t.Fatal(err)
		}
		if unit.String() != testCase[2].(string) {
			t.Error("Expecting", testCase[2].(string), "got", unit, "in mov", testCase[0], testCase[1])
		}
	}
}

func Test_Execute(t *testing.T) {
	units := [][]lib.Instruction{
		[]lib.Instruction{
			MOV(encoding.Uint32(uint32(5)), &encoding.DisplacedRegister{encoding.Rsp, 8}),
			RETURN(),
		},
		[]lib.Instruction{
			MOV(encoding.Uint64(5), encoding.Rax),
			MOV(encoding.Rax, &encoding.DisplacedRegister{encoding.Rsp, 8}),
			RETURN(),
		},
		[]lib.Instruction{
			MOV(encoding.Uint64(2), encoding.Rax),
			MOV(encoding.Uint64(3), encoding.Rdi),
			ADD(encoding.Rdi, encoding.Rax),
			MOV(encoding.Rax, &encoding.DisplacedRegister{encoding.Rsp, 8}),
			RETURN(),
		},
		[]lib.Instruction{
			MOV(encoding.Uint64(2), encoding.R13),
			MOV(encoding.Uint64(3), encoding.R14),
			ADD(encoding.R13, encoding.R14),
			MOV(encoding.R14, &encoding.DisplacedRegister{encoding.Rsp, 8}),
			RETURN(),
		},
		[]lib.Instruction{
			MOV(encoding.Uint64(2), encoding.Rax),
			ADD(encoding.Uint32(3), encoding.Rax),
			MOV(encoding.Rax, &encoding.DisplacedRegister{encoding.Rsp, 8}),
			RETURN(),
		},
		[]lib.Instruction{
			MOV(encoding.Uint64(5), encoding.Rdi),
			CVTSI2SD(encoding.Rdi, encoding.Xmm4),
			MOV(encoding.Uint64(6), encoding.Rdi),
			CVTTSD2SI(encoding.Xmm4, encoding.Rax),
			MOV(encoding.Rax, &encoding.DisplacedRegister{encoding.Rsp, 8}),
			RETURN(),
		},
		[]lib.Instruction{
			MOV(encoding.Float64(5.0), encoding.Rdi),
			MOV(encoding.Rdi, encoding.Xmm4),
			CVTTSD2SI(encoding.Xmm4, encoding.Rax),
			MOV(encoding.Rax, &encoding.DisplacedRegister{encoding.Rsp, 8}),
			RETURN(),
		},
		[]lib.Instruction{
			MOV(encoding.Uint64(1), encoding.Rdi),
			CVTSI2SD(encoding.Rdi, encoding.Xmm5),
			MOV(encoding.Float64(4.0), encoding.Rdi),
			MOV(encoding.Rdi, encoding.Xmm4),
			ADD(encoding.Xmm5, encoding.Xmm4),
			CVTTSD2SI(encoding.Xmm4, encoding.Rax),
			MOV(encoding.Rax, &encoding.DisplacedRegister{encoding.Rsp, 8}),
			RETURN(),
		},
		[]lib.Instruction{
			MOV(encoding.Uint64(2), encoding.Rdi),
			CVTSI2SD(encoding.Rdi, encoding.Xmm5),
			MOV(encoding.Float64(7.0), encoding.Rdi),
			MOV(encoding.Rdi, encoding.Xmm4),
			SUB(encoding.Xmm5, encoding.Xmm4),
			CVTTSD2SI(encoding.Xmm4, encoding.Rax),
			MOV(encoding.Rax, &encoding.DisplacedRegister{encoding.Rsp, 8}),
			RETURN(),
		},
		[]lib.Instruction{
			MOV(encoding.Float64(2.0), encoding.Rcx),
			MOV(encoding.Float64(2.5), encoding.Rdi),
			MOV(encoding.Rcx, encoding.Xmm4),
			MOV(encoding.Rdi, encoding.Xmm5),
			MUL(encoding.Xmm5, encoding.Xmm4),
			CVTTSD2SI(encoding.Xmm4, encoding.Rax),
			MOV(encoding.Rax, &encoding.DisplacedRegister{encoding.Rsp, 8}),
			RETURN(),
		},
		[]lib.Instruction{
			MOV(encoding.Float64(2.0), encoding.Rcx),
			MOV(encoding.Float64(10.0), encoding.Rdi),
			MOV(encoding.Rdi, encoding.Xmm4),
			MOV(encoding.Rcx, encoding.Xmm5),
			DIV(encoding.Xmm5, encoding.Xmm4),
			CVTTSD2SI(encoding.Xmm4, encoding.Rax),
			MOV(encoding.Rax, &encoding.DisplacedRegister{encoding.Rsp, 8}),
			RETURN(),
		},
		[]lib.Instruction{
			MOV(encoding.Float64(2.0), encoding.Rcx),
			MOV(encoding.Float64(10.0), encoding.Rdi),
			MOV(encoding.Rdi, encoding.Xmm4),
			MOV(encoding.Rcx, encoding.Xmm5),
			MOV(encoding.Xmm4, encoding.Xmm0),
			MOV(encoding.Xmm5, encoding.Xmm1),
			DIV(encoding.Xmm1, encoding.Xmm0),
			CVTTSD2SI(encoding.Xmm0, encoding.Rax),
			MOV(encoding.Rax, &encoding.DisplacedRegister{encoding.Rsp, 8}),
			RETURN(),
		},
	}
	for _, unit := range units {
		debug := false
		b, err := lib.CompileInstruction(unit, debug)
		if err != nil {
			t.Fatal(err, "in", unit)
		}
		value := b.Execute(debug)
		if value != uint(5) {
			t.Fatal("Expecting 5 got", value, "in", unit, "\n", b)
		}
	}

}
