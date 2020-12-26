package elf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
	"strings"
)

//go:generate stringer -type=SHType
type SHType uint32

const (
	// This value marks the section header as inactive; it does not have an
	// associated section.  Other members of the section header have undefined
	// values.
	SHT_NULL SHType = 0
	// The section holds information defined by the program, whose format and
	// meaning are determined solely by the program.
	SHT_PROGBITS SHType = 1
	// These sections hold a symbol table.  Currently, an object file may have
	// only one section of each type, but this restriction may be relaxed in the
	// future.  Typically, SHT_SYMTAB provides symbols for link editing,
	// though it may also be used for dynamic linking.  As a complete symbol
	// table, it may contain many symbols unnecessary for dynamic linking.
	// Consequently, an object file may also contain a SHT_DYNSYM
	// section, which holds a minimal set of dynamic linking symbols, to
	// save space.
	SHT_SYMTAB SHType = 2
	// The section holds a string table.  An object file may have multiple
	// string table sections.
	SHT_STRTAB SHType = 3
	// The section holds relocation entries with explicit addends, such as type
	// Elf32_Rela afor the 32-bit class of object files.  An object file
	// may have multiple relocation sections.
	SHT_RELA SHType = 4
	// The section holds a symbol hash table.  All objects participating in
	// dynamic linking must contain a symbol hash table.Currently, an object
	// file may have only one hash table, but thisrestriction may be relaxed in
	// the future.
	SHT_HASH SHType = 5
	// The section holds information for dynamic linking. Currently an object
	// file may have only one dynamic section, but this restriction may be
	// relaxed in the future.
	SHT_DYNAMIC SHType = 6
	// The section holds information that marks the file in some way. (no shit)
	SHT_NOTE SHType = 7
	// A section of this type occupies no space in the file but other-wise
	// resembles SHT_PROGBITS.  Although this section contains no
	// bytes, the sh_offset member contains the conceptual file offset.
	SHT_NOBITS SHType = 8
	// The section holds relocation entries without explicit addends,such as
	// type Elf32_Rel for the 32-bit class of object files.  Anobject
	// file may have multiple relocation sections.
	SHT_REL SHType = 9
	// This section type is reserved but has unspecified semantics.
	SHT_SHLIB SHType = 10
	// These sections hold a symbol table.  Currently, an object file may have
	// only one section of each type, but this restriction may be relaxed in the
	// future.  Typically, SHT_SYMTAB provides symbols for link editing,
	// though it may also be used for dynamic linking.  As a complete symbol
	// table, it may contain many symbols unnecessary for dynamic linking.
	// Consequently, an object file may also contain a SHT_DYNSYM
	// section, which holds a minimal set of dynamic linking symbols, to
	// save space.
	SHT_DYNSYM SHType = 11
	SHT_NUM    SHType = 12
	// Values in this inclusive range are reserved for processor-specific
	// semantics.  If meanings are specified, the processorsupplement explains
	// them.
	SHT_LOPROC SHType = 0x70000000
	// Values in this inclusive range are reserved for processor-specific
	// semantics.  If meanings are specified, the processorsupplement explains
	// them.
	SHT_HIPROC SHType = 0x7fffffff
	// This value specifies the lower bound of the range of indexesreserved for application programs
	SHT_LOUSER SHType = 0x80000000
	// This value specifies the upper bound of the range of indexesreserved for application programs.
	SHT_HIUSER SHType = 0xffffffff
)

const SHN_UNDEF = 0x0

//go:generate stringer -type=SHFlags
type SHFlags uint64

const (
	// No flags
	SHF_NULL SHFlags = 0x0
	// The section contains data that should be writable during process execution.
	SHF_WRITE SHFlags = 0x1
	// The section occupies memory during process execution.Some control
	// sections do not reside in the memory image of an object file; this
	// attribute is off for those sections.
	SHF_ALLOC SHFlags = 0x2
	// The section contains executable machine instructions.
	SHF_EXECINSTR      SHFlags = 0x4
	SHF_RELA_LIVEPATCH SHFlags = 0x00100000
	SHF_RO_AFTER_INIT  SHFlags = 0x00200000
	// All bits included in this mask are reserved for processor-specific
	// semantics.  If meanings are specified, the processorsupplement explains
	// them.
	SHF_MASKPROC SHFlags = 0xf0000000
)

type SectionHeader struct {
	// This member specifies the name of the section.  Its value is an index
	// into the section header string table section, giving the location of a
	// null-terminated string.
	Name uint32
	// This member categorizes the section’s contents and semantics.
	Type SHType
	// Sections support 1-bit flags that describe miscellaneous attributes.
	Flags SHFlags
	// If the section will appear in the memory image of a process, this member
	// gives the address at which the section’s first byte should reside.
	// Otherwise, the member contains 0
	Addr Elf64_Addr
	// This member’s value gives the byte offset from the beginning of the file
	// to the first byte in the section. One section type, SHT_NOBITS,
	// occupies no space in the file, and its sh_offset
	// member locates the conceptual placement in the file.
	Offset Elf64_Off
	// This member gives the section’s size in bytes.  Unless the section type is
	// SHT_NOBITS, the section occupiess h_size bytes in the file.
	// A section of type SHT_NOBITS may have a non-zerosize, but it occupies no space in the file.
	Size uint32
	// This member holds a section header table index link, whoseinterpretation depends on the section type.
	Link uint32
	// This member holds extra information, whose interpretationdepends on the section type.
	Info uint32
	// Some sections have address alignment constraints.  For example, if a
	// section holds a doubleword, the system must ensure doubleword alignment
	// for the entire section.  That is, the value of sh_addr must be
	// congruent to 0, modulo the value of sh_addralign.  Currently,
	// only 0 and positive integral powersof two are allowed.  Values 0 and 1
	// mean the section has noalignment constraints.
	AddrAlign uint32
	// Some sections hold a table of fixed-size entries, such as a symbol
	// table.  For such a section, this member gives the size in bytes of each
	// entry.  The member contains 0 if the section doesnot hold a table of
	// fixed-size entries
	EntSize uint32
}

func ParseSectionHeader(header *ELFHeader, r *bytes.Reader) (*SectionHeader, error) {
	result := &SectionHeader{}
	byteOrder := header.GetByteOrder()
	if err := binary.Read(r, byteOrder, &result.Name); err != nil {
		return nil, err
	}
	if err := binary.Read(r, byteOrder, &result.Type); err != nil {
		return nil, err
	}
	if err := binary.Read(r, byteOrder, &result.Flags); err != nil {
		return nil, err
	}
	if err := binary.Read(r, byteOrder, &result.Addr); err != nil {
		return nil, err
	}
	if err := binary.Read(r, byteOrder, &result.Offset); err != nil {
		return nil, err
	}
	if err := binary.Read(r, byteOrder, &result.Size); err != nil {
		return nil, err
	}
	if err := binary.Read(r, byteOrder, &result.Link); err != nil {
		return nil, err
	}
	if err := binary.Read(r, byteOrder, &result.Info); err != nil {
		return nil, err
	}
	if err := binary.Read(r, byteOrder, &result.AddrAlign); err != nil {
		return nil, err
	}
	if err := binary.Read(r, byteOrder, &result.EntSize); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *SectionHeader) String() string {
	table := [][]string{
		[]string{"Section Header:"},
		[]string{"  Name:          ", strconv.Itoa(int(s.Name))},
		[]string{"  Type:          ", s.Type.String()},
		[]string{"  Flags:         ", s.Flags.String()},
		[]string{"  Address:       ", fmt.Sprintf("0x%x", s.Addr)},
		[]string{"  Offset:        ", fmt.Sprintf("0x%x", s.Offset)},
		[]string{"  Size:          ", fmt.Sprintf("%v", s.Size)},
		[]string{"  Link:          ", fmt.Sprintf("%v", s.Link)},
		[]string{"  Info:          ", fmt.Sprintf("%v", s.Info)},
		[]string{"  Address Align: ", fmt.Sprintf("0x%x", s.AddrAlign)},
		[]string{"  Entry Size:    ", fmt.Sprintf("0x%x", s.EntSize)},
	}
	result := []string{}
	for _, row := range table {
		result = append(result, strings.Join(row, "\t"))
	}
	return strings.Join(result, "\n")
}

func ParseSectionHeaderTable(header *ELFHeader, r *bytes.Reader) ([]*SectionHeader, error) {
	result := []*SectionHeader{}
	for i := uint16(0); i < header.SectionHeaderNumberOfEntries; i++ {

		offset := int64(header.SectionHeaderTableOffset)
		offset += int64(i) * int64(header.SectionHeaderEntrySize)
		if _, err := r.Seek(offset, io.SeekStart); err != nil {
			return nil, err
		}
		sh, err := ParseSectionHeader(header, r)
		if err != nil {
			return nil, err
		}
		result = append(result, sh)

	}
	return result, nil
}
