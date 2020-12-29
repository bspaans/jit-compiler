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

func Test_ELF_encode(t *testing.T) {

	elf := NewELFHeader()
	elf.Class = ELFCLASS32
	elf.Data = ELFDATA2LSB
	elf.OS_ABI = ELF_OS_ABI_LINUX
	elf.ABI_Version = 9
	elf.Type = ET_EXEC
	elf.Machine = EM_88K
	elf.Version = EV_CURRENT
	elf.Entry = 0xff
	elf.ProgramHeaderTableOffset = 0x1234

	headerBytes, err := elf.Encode()
	if err != nil {
		t.Fatal(err)
	}
	reader := bytes.NewReader(headerBytes)

	parsedHeader, err := ParseELFHeader(reader)
	if err != nil {
		t.Fatal(err)
	}
	if parsedHeader.Class != ELFCLASS32 {
		t.Errorf("Wrong class %v", parsedHeader.Class)
	}
	if parsedHeader.Data != ELFDATA2LSB {
		t.Errorf("Wrong data %v", parsedHeader.Data)
	}
	if parsedHeader.OS_ABI != ELF_OS_ABI_LINUX {
		t.Errorf("Wrong OS ABI %v", parsedHeader.OS_ABI)
	}
	if parsedHeader.ABI_Version != 9 {
		t.Errorf("Wrong ABI version %v", parsedHeader.ABI_Version)
	}
	if parsedHeader.Type != ET_EXEC {
		t.Errorf("Wrong type %v", parsedHeader.Type)
	}
	if parsedHeader.Machine != EM_88K {
		t.Errorf("Wrong machine %v", parsedHeader.Machine)
	}
	if parsedHeader.Version != EV_CURRENT {
		t.Errorf("Wrong version %v", parsedHeader.Version)
	}
	if parsedHeader.Entry != 0xff {
		t.Errorf("Wrong entry %v", parsedHeader.Entry)
	}
	if parsedHeader.ProgramHeaderTableOffset != 0x1234 {
		t.Errorf("Wrong ph offset %v", parsedHeader.ProgramHeaderTableOffset)
	}
}
