package encoding

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
