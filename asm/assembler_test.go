package asm

import "testing"

func Test_EncodeREXPrefix(t *testing.T) {
	unit := NewREXPrefix(true, true, true, true).Encode()
	expected := uint8(79)
	if unit != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
	unit = NewREXPrefix(false, false, false, false).Encode()
	expected = uint8(64)
	if unit != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}

	unit = NewREXPrefix(false, false, false, true).Encode()
	expected = uint8(65)
	if unit != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
	unit = NewREXPrefix(false, false, true, true).Encode()
	expected = uint8(67)
	if unit != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
	unit = NewREXPrefix(false, true, true, true).Encode()
	expected = uint8(71)
	if unit != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
}

func Test_EncodeModRM(t *testing.T) {
	unit := NewModRM(DirectRegisterMode, 0, 0).Encode()
	expected := uint8(192)
	if unit != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
	unit = NewModRM(DirectRegisterMode, 1, 0).Encode()
	expected = uint8(193)
	if unit != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
	unit = NewModRM(DirectRegisterMode, 1, 1).Encode()
	expected = uint8(201)
	if unit != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
	unit = NewModRM(IndirectRegisterMode, 0, 1).Encode()
	expected = uint8(8)
	if unit != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
	for i, reg := range []*Register{rax, rcx, rdx, rbx, rsp, rbp, rsi, rdi} {
		unit = NewModRM(DirectRegisterMode, reg.Encode(), 0).Encode()
		expected = uint8(192 + i)
		if unit != expected {
			t.Fatal("Expecting", expected, "got", unit)
		}
	}
	for i, reg := range []*Register{r8, r9, r10, r11, r12, r13, r14, r15} {
		unit = NewModRM(DirectRegisterMode, reg.Encode(), 0).Encode()
		expected = uint8(192 + i)
		if unit != expected {
			t.Fatal("Expecting", expected, "got", unit)
		}
	}
}

func Test_INC(t *testing.T) {
	unit, err := (&INC{rax}).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected := "  48 ff c0"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
	unit, err = (&INC{rcx}).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected = "  48 ff c1"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}

	unit, err = (&INC{r14}).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected = "  49 ff c6"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
}

func Test_DEC(t *testing.T) {
	unit, err := (&DEC{rax}).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected := "  48 ff c8"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
	unit, err = (&DEC{rcx}).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected = "  48 ff c9"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
	unit, err = (&DEC{r14}).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected = "  49 ff ce"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
}

func Test_MOV(t *testing.T) {
	unit, err := (&MOV{rax, rax}).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected := "  48 89 c0"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
	unit, err = (&MOV{rax, rcx}).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected = "  48 89 c1"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
	unit, err = (&MOV{rcx, rax}).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected = "  48 89 c8"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
	unit, err = (&MOV{rax, r14}).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected = "  49 89 c6"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
	unit, err = (&MOV{r14, rax}).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected = "  4c 89 f0"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
	unit, err = (&MOV{Uint64(0xffffffff), rax}).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected = "  48 b8 ff ff ff ff 00 \n  00 00 00"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
	unit, err = (&MOV{Uint64(0xffffffff), rcx}).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected = "  48 b9 ff ff ff ff 00 \n  00 00 00"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
	unit, err = (&MOV{Uint64(0xffffffff), r14}).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected = "  49 be ff ff ff ff 00 \n  00 00 00"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
}

func Test_JMP(t *testing.T) {
	unit, err := (&JMP{Uint8(3)}).Encode()
	if err != nil {
		t.Fatal(err)
	}
	expected := "  eb 03"
	if unit.String() != expected {
		t.Fatal("Expecting", expected, "got", unit)
	}
}
