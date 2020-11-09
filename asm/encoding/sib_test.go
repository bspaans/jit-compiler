package encoding

import "testing"

func Test_EncodeSIB(t *testing.T) {

	got := NewSIB(Scale0, 0, 0).Encode()
	expected := uint8(0)
	if got != expected {
		t.Fatal("Expecting", expected, "got", got)
	}

	got = NewSIB(Scale0, 0, 1).Encode()
	expected = uint8(1)
	if got != expected {
		t.Fatal("Expecting", expected, "got", got)
	}
	got = NewSIB(Scale8, 7, 1).Encode()
	expected = uint8(0xf9)
	if got != expected {
		t.Fatal("Expecting", expected, "got", got)
	}

}

func Test_ScaleString(t *testing.T) {
	if Scale8.String() != "8" {
		t.Fatal("Expecting 8, got", Scale8.String())
	}
}
