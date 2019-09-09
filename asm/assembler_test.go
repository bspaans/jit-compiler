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
	for i, reg := range []*Register{rax, rcx, rdx, rbx, rsp, rbp, rsi, rdi} {
		unit = NewModRM(DirectRegisterMode, reg.Encode(), 0).Encode()
		expected = uint8(192 + i)
		if unit != expected {
			t.Fatal("Expecting", expected, "got", unit)
		}
	}
}
