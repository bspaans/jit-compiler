package elf

import (
	"bytes"
	"fmt"
	"testing"
)

func Test_ELF_parse_header_sad(t *testing.T) {
	cases := [][]uint8{
		[]uint8{0x7f, 'E', 'L', 'K'},
	}
	for _, c := range cases {
		r := bytes.NewReader(c)
		_, err := ParseELFHeader(r)
		if err == nil {
			t.Errorf("Expecting error parsing: %v", c)
		}
	}
}

func Test_ELF_parse_header_happy(t *testing.T) {

	cases := [][]uint8{
		[]uint8{0x7f, 'E', 'L', 'F',
			uint8(ELFCLASS64), uint8(ELFDATA2MSB),
			uint8(EV_CURRENT), uint8(ELF_OS_ABI_LINUX),
			0x0,                // version
			0x0, 0x0, 0x0, 0x0, // padding
			0x0, 0x0, 0x0, // padding
			0x0, 0x02, // type
			0x0, 0x3e, // machine
			0x0, 0x0, 0x0, 0x1, // version
			0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, // entry
			0x0, 0x0, 0x0, 0xf, 0x0, 0x0, 0x0, 0x0, // ph offset
			0x0, 0x0, 0x0, 0xf, 0x0, 0x0, 0x0, 0x0, // sh offset
			0x0, 0x0, 0x0, 0x0, // flags
			0x0, 0x0, // header size
			0x0, 0x0, // program header entry size
			0x0, 0x0, // program header number of entries
			0x0, 0x0, // section header entry size
			0x0, 0x0, // section header number of entries
			0x0, 0x0, // section header table index
		},
	}
	for _, c := range cases {
		r := bytes.NewReader(c)
		header, err := ParseELFHeader(r)
		if err != nil {
			t.Errorf("Unexpected error in ELF header %v: %v", c, err.Error())
		}
		if header != nil {
			fmt.Println(header.String())
		}
	}
}
