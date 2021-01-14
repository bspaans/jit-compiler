package elf

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Section struct {
	Name      string
	Type      SHType
	Flags     SHFlags
	Size      uint32
	Addr      Elf64_Addr
	Offset    Elf64_Off
	Link      uint32
	Info      uint32
	AddrAlign uint32
	EntSize   uint32
	Data      []byte

	header *ELFHeader
}

func NewSection(name string, typ SHType, flags SHFlags) *Section {
	return &Section{
		Name:  name,
		Type:  typ,
		Flags: flags,
	}
}

func (s *Section) String() string {

	stringIf := func(condition bool, result []string) []string {
		if !condition {
			return []string{}
		}
		return result
	}

	table := [][]string{
		[]string{"Section:"},
		[]string{"  Name:          ", s.Name},
		[]string{"  Type:          ", s.Type.String()},
		[]string{"  Flags:         ", s.Flags.String()},
		stringIf(s.Addr != 0, []string{"  Address:       ", fmt.Sprintf("0x%x", s.Addr)}),
		stringIf(s.Offset != 0, []string{"  Offset:        ", fmt.Sprintf("0x%x", s.Offset)}),
		stringIf(s.Size != 0, []string{"  Size:          ", fmt.Sprintf("%v", s.Size)}),
		stringIf(s.Link != 0, []string{"  Link:          ", fmt.Sprintf("%v", s.Link)}),
		stringIf(s.Info != 0, []string{"  Info:          ", fmt.Sprintf("%v", s.Info)}),
		stringIf(s.AddrAlign != 0, []string{"  Address Align: ", fmt.Sprintf("0x%x", s.AddrAlign)}),
		stringIf(s.EntSize != 0, []string{"  Entry Size:    ", fmt.Sprintf("0x%x", s.EntSize)}),
	}
	result := []string{}
	for _, row := range table {
		if len(row) != 0 {
			result = append(result, strings.Join(row, "\t"))
		}
	}
	return strings.Join(result, "\n")
}

func (s *Section) GetStringTable() *StringTable {
	return NewStringTable(s.Data)
}
func (s *Section) GetSymbolTable(strTable *StringTable) (*SymbolTable, error) {
	return ParseSymbolTable(s.header, strTable, bytes.NewReader(s.Data))
}

// This section holds uninitialized data that contribute to the program’s memory
// image.  By definition, the system initializes the data with zeros when the
// program begins to run.  The section occupies no file space, as indicated by
// the section type, SHT_NOBITS
func NewBSSSection() *Section {
	return NewSection(".bss", SHT_NOBITS, SHF_ALLOC&SHF_WRITE)
}

// This section holds version control information.
func NewCommentSection() *Section {
	return NewSection(".comment", SHT_PROGBITS, SHF_NULL)
}

// These sections hold initialized data that contribute to the program’s memory image
func NewDataSection() *Section {
	return NewSection(".data", SHT_PROGBITS, SHF_ALLOC&SHF_WRITE)
}

// These sections hold read-only data that typically contribute to a
// non-writable segment in the process image.
func NewReadOnlyDataSection() *Section {
	return NewSection(".rodata", SHT_PROGBITS, SHF_ALLOC)
}

// This section holds section names.
func NewSectionHeaderStringSection() *Section {
	return NewSection(".shstrtab", SHT_STRTAB, SHF_NULL)
}

// This section holds the "text", or executable instructions, of a program
func NewTextSection() *Section {
	return NewSection(".text", SHT_PROGBITS, SHF_ALLOC)
}

func ParseSections(header *ELFHeader, shTable []*SectionHeader, r *bytes.Reader) ([]*Section, error) {
	if header.SectionHeaderNumberOfEntries == 0 {
		return []*Section{}, nil
	}
	if header.Shstrndx == SHN_UNDEF {
		return nil, errors.New("ELF files without a section header string table are unsupported")
	}
	if int(header.Shstrndx) >= len(shTable) {
		return nil, errors.New("Shstrndx out of bounds")
	}

	strTable, err := ParseStringTable(header, shTable[header.Shstrndx], r)
	if err != nil {
		return nil, err
	}

	result := []*Section{}
	for _, sectionHeader := range shTable {
		name, err := strTable.GetString(int(sectionHeader.Name))
		if err != nil {
			return nil, err
		}
		data := []byte{}
		if sectionHeader.Type != SHT_NULL {
			_, err := r.Seek(int64(sectionHeader.Offset), io.SeekStart)
			if err != nil {
				return nil, err
			}
			tmpData := make([]byte, sectionHeader.Size)
			n, err := r.Read(tmpData)
			if err != nil {
				return nil, err
			}
			if n != int(sectionHeader.Size) {
				return nil, errors.New("read bytes != section header size")
			}
			data = tmpData
		}
		section := &Section{
			Name:      name,
			Type:      sectionHeader.Type,
			Flags:     sectionHeader.Flags,
			Size:      sectionHeader.Size,
			Addr:      sectionHeader.Addr,
			Offset:    sectionHeader.Offset,
			Link:      sectionHeader.Link,
			Info:      sectionHeader.Info,
			AddrAlign: sectionHeader.AddrAlign,
			EntSize:   sectionHeader.EntSize,
			Data:      data,
			header:    header,
		}
		result = append(result, section)

	}
	return result, nil

}
