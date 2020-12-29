package elf

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/bspaans/jit-compiler/lib"
)

/*

An  executable  file using the ELF file format consists of an ELF header,
followed by a program header table or a section header table, or both.  The
ELF header is always at offset zero of the file.  The program header table
and the section header table's offset in the file are defined in the ELF
header.  The two tables describe the rest of the particularities of the file.


*/

type ELF struct {
	*ELFHeader

	// A program header table, if present, tells the system how to create a
	// process image. Files used to build a process image (execute a program)
	// must have a program header table; relocatable files do not need one.
	ProgramHeaders ProgramHeaderTable

	Sections []*Section
}

func NewELF() *ELF {
	return &ELF{}
}

func ParseELF(r *bytes.Reader) (*ELF, error) {
	header, err := ParseELFHeader(r)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse ELF header: %s", err.Error())
	}

	table, err := ParseSectionHeaderTable(header, r)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse section header table: %s", err.Error())
	}
	sections, err := ParseSections(header, table, r)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse sections: %s", err.Error())
	}
	programHeaders, err := ParseProgramHeaderTable(header, r)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse program header table: %s", err.Error())
	}
	return &ELF{
		ELFHeader:      header,
		Sections:       sections,
		ProgramHeaders: programHeaders,
	}, nil
}

func ParseELFFile(path string) (*ELF, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseELF(bytes.NewReader(f))
}

// Encodes the ELF header, program header table and section header table (TODO).
// Does not copy .Data!
func (e *ELF) EncodeHeaders() ([]byte, error) {
	headerSize := 64
	e.ELFHeader.ProgramHeaderEntrySize = ProgramHeaderSize64
	e.ELFHeader.SectionHeaderEntrySize = 64
	e.ELFHeader.ProgramHeaderNumberOfEntries = uint16(len(e.ProgramHeaders))
	e.ELFHeader.SectionHeaderNumberOfEntries = uint16(len(e.Sections))
	e.ELFHeader.ProgramHeaderTableOffset = Elf64_Off(headerSize)
	programHeaders, err := e.ProgramHeaders.Encode(e.ELFHeader)
	if err != nil {
		return nil, err
	}
	e.ELFHeader.SectionHeaderTableOffset = 0 // Elf64_Off(headerSize + len(programHeaders))
	header, err := e.ELFHeader.Encode()
	if err != nil {
		return nil, err
	}
	if len(header) != headerSize {
		panic("header size mismatch")
	}
	result := append(header, programHeaders...)
	return result, nil
}

func (e *ELF) String() string {
	result := e.ELFHeader.String()
	for _, header := range e.ProgramHeaders {
		result += "\n" + header.String()
	}
	for _, section := range e.Sections {
		result += "\n" + section.String()
	}
	return result
}

func (e *ELF) GetSection(name string) *Section {
	for _, s := range e.Sections {
		if s.Name == name {
			return s
		}
	}
	return nil
}

// Demo function.
// Don't use this in anything serious: it puts data in an executable segment.
func CreateTinyBinary(m lib.MachineCode, path string) error {
	elf := NewELF()
	elf.ELFHeader = NewELFHeader()

	programHeaderEntrySize := ProgramHeaderSize64
	headerSize := uint16(64)
	addr := Elf64_Addr(0x400000)
	size := uint64(headerSize+programHeaderEntrySize) + uint64(len(m))

	elf.Entry = addr + Elf64_Addr(headerSize+programHeaderEntrySize)

	ph := NewProgramHeader(PT_LOAD, PF_RWX)
	ph.Offset = Elf64_Off(0)
	ph.SegmentVirtualAddress = addr
	ph.SegmentPhysicalAddress = addr
	ph.Filesize = size
	ph.Memsize = size
	ph.Align = size
	elf.ProgramHeaders = ProgramHeaderTable{ph}
	result, err := elf.EncodeHeaders()
	if err != nil {
		return err
	}
	result = append(result, m...)
	fmt.Println(len(result), size)

	elfReloaded, err := ParseELF(bytes.NewReader(result))
	if err != nil {
		return err
	}
	fmt.Println(elfReloaded.String())

	return ioutil.WriteFile(path, result, 0755)
}
