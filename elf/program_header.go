package elf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)

//go:generate stringer -type=PHType
type PHType uint32

const (
	// The array element is unused; other members’ values are undefined.  This
	// type lets the program header table have ignored entries
	PT_NULL PHType = 0
	// The array element specifies a loadable segment, described by Filesize
	// and Memsize.  The bytes from the file are mapped to the beginning of the
	// memory segment.  If the segment’s memorysize (Memsize) is larger than
	// the file size (Filesize), the "extra" bytes are defined to hold the
	// value 0 and to follow the segment’s initialized area.  The file size may
	// not be larger than the memorysize.  Loadable segment entries in the
	// program header table appear in ascending order, sorted on the
	// SegmentVirtualAddress member.
	PT_LOAD PHType = 1
	// The array element specifies dynamic linking information.
	PT_DYNAMIC PHType = 2
	// The array element specifies the location and size of a null-terminated
	// path name to invoke as an interpreter.  This segment type is meaningful
	// only for executable files (though it may occur for shared objects); it
	// may not occur more than once in a file.  If it is present, it must
	// precede any loadable segment entry. (e.g. /lib64/ld-linux-x86-64.so.2)
	PT_INTERP PHType = 3
	// The array element specifies the location and size of auxiliaryinformation.
	PT_NOTE PHType = 4
	// This segment type is reserved but has unspecified semantics
	PT_SHLIB PHType = 5
	// The array element, if present, specifies the location and size of
	// the program header table itself, both in the file and in the memory image
	// of the program.  This segment type may not occur more than once in a
	// file.  Moreover, it may occur only if the programheader table is part of
	// the memory image of the program.  If it is present, it must precede any
	// loadable segment entry.
	PT_PHDR PHType = 6
	// Values in this inclusive range are reserved for processor-specific semantics.
	PT_LOPROC PHType = 0x70000000
	// Values in this inclusive range are reserved for processor-specific semantics.
	PT_HIPROC PHType = 0x7fffffff
)

//go:generate stringer -type=PHFlags
type PHFlags uint32

const (
	// Execute
	PF_X PHFlags = 0x1
	// Write
	PF_W PHFlags = 0x2
	// Execute & Write
	PF_WX PHFlags = 0x3
	// Read
	PF_R PHFlags = 0x4
	// Read & Execute
	PF_RX PHFlags = 0x5
	// Read & Write
	PF_RW PHFlags = 0x6
	// Read & Write & Execute
	PF_RWX PHFlags = 0x7
	// Unspecified
	PF_MASKPROC PHFlags = 0xf0000000
)

const ProgramHeaderSize64 uint16 = 56

type ProgramHeader struct {
	// This member tells what kind of segment this array element describes or
	// how to interpret the array element’s information.
	Type PHType
	// This member gives flags relevant to the segment.
	Flags PHFlags
	// This member gives the offset from the beginning of the file at which the
	// first byte of the segment resides.
	Offset Elf64_Off
	// This member gives the virtual address at which the first byte of the
	// segment resides in memory
	SegmentVirtualAddress Elf64_Addr
	// On systems for which physical addressing is relevant, this member is
	// reserved for the segment’s physical address.  Because System V ignores
	// physical addressing for application programs,this member has unspecified
	// contents for executable files and shared objects
	SegmentPhysicalAddress Elf64_Addr
	// Segment size in file
	Filesize uint64
	// Segment size in memory
	Memsize uint64
	// Segment alignment, file & memory
	Align uint64

	Data []byte
}

func NewProgramHeader(typ PHType, flags PHFlags) *ProgramHeader {
	return &ProgramHeader{
		Type:  typ,
		Flags: flags,
	}
}

func ParseProgramHeader(header *ELFHeader, r *bytes.Reader) (*ProgramHeader, error) {
	result := &ProgramHeader{}
	byteOrder := header.GetByteOrder()
	if err := binary.Read(r, byteOrder, &result.Type); err != nil {
		return nil, err
	}
	if err := binary.Read(r, byteOrder, &result.Flags); err != nil {
		return nil, err
	}
	if err := binary.Read(r, byteOrder, &result.Offset); err != nil {
		return nil, err
	}
	if err := binary.Read(r, byteOrder, &result.SegmentVirtualAddress); err != nil {
		return nil, err
	}
	if err := binary.Read(r, byteOrder, &result.SegmentPhysicalAddress); err != nil {
		return nil, err
	}
	if err := binary.Read(r, byteOrder, &result.Filesize); err != nil {
		return nil, err
	}
	if err := binary.Read(r, byteOrder, &result.Memsize); err != nil {
		return nil, err
	}
	if err := binary.Read(r, byteOrder, &result.Align); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ProgramHeader) Encode(header *ELFHeader) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	byteOrder := header.GetByteOrder()
	if err := binary.Write(buffer, byteOrder, s.Type); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, byteOrder, s.Flags); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, byteOrder, s.Offset); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, byteOrder, s.SegmentVirtualAddress); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, byteOrder, s.SegmentPhysicalAddress); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, byteOrder, s.Filesize); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, byteOrder, s.Memsize); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, byteOrder, s.Align); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (s *ProgramHeader) String() string {
	table := [][]string{
		[]string{"Program Header:"},
		[]string{"  Type:             ", s.Type.String()},
		[]string{"  Flags:            ", s.Flags.String()},
		[]string{"  Offset:           ", fmt.Sprintf("0x%x", s.Offset)},
		[]string{"  Virtual Address:  ", fmt.Sprintf("0x%x", s.SegmentVirtualAddress)},
		[]string{"  Physical Address: ", fmt.Sprintf("0x%x", s.SegmentPhysicalAddress)},
		[]string{"  File Size:        ", fmt.Sprintf("0x%x (%v)", s.Filesize, s.Filesize)},
		[]string{"  Memory Size:      ", fmt.Sprintf("0x%x (%v)", s.Memsize, s.Memsize)},
		[]string{"  Address Align:    ", fmt.Sprintf("0x%x (%v)", s.Align, s.Align)},
	}
	result := []string{}
	for _, row := range table {
		result = append(result, strings.Join(row, "\t"))
	}
	return strings.Join(result, "\n")
}

func ParseProgramHeaderTable(header *ELFHeader, r *bytes.Reader) ([]*ProgramHeader, error) {

	result := []*ProgramHeader{}
	for i := uint16(0); i < header.ProgramHeaderNumberOfEntries; i++ {

		offset := int64(header.ProgramHeaderTableOffset)
		offset += int64(i) * int64(header.ProgramHeaderEntrySize)
		if _, err := r.Seek(offset, io.SeekStart); err != nil {
			return nil, fmt.Errorf("failed seek(1): %s", err.Error())
		}
		ph, err := ParseProgramHeader(header, r)
		if err != nil {
			return nil, fmt.Errorf("failed parse: %s", err.Error())
		}
		if _, err := r.Seek(int64(ph.Offset), io.SeekStart); err != nil {
			return nil, fmt.Errorf("failed seek(2): %s", err.Error())
		}
		tmpData := make([]byte, ph.Filesize)
		if _, err := r.Read(tmpData); err != nil {
			return nil, fmt.Errorf("failed read from %d for %d bytes: %s", ph.Offset, ph.Filesize, err.Error())
		}
		ph.Data = tmpData
		result = append(result, ph)
	}
	return result, nil
}

type ProgramHeaderTable []*ProgramHeader

func (p ProgramHeaderTable) Encode(header *ELFHeader) ([]byte, error) {
	result := make([]byte, int(header.ProgramHeaderEntrySize)*len(p))
	fmt.Println(len(result), header.ProgramHeaderEntrySize, len(p))
	for i, ph := range p {
		phBytes, err := ph.Encode(header)
		if err != nil {
			return nil, err
		}
		for j, b := range phBytes {
			result[i*int(header.ProgramHeaderEntrySize)+j] = b
		}
	}
	return result, nil
}
