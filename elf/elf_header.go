package elf

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Elf64_Addr uint64
type Elf64_Off uint64

//go:generate stringer -type=ELFClass
type ELFClass uint8

const (
	ELFCLASSNONE ELFClass = 0
	ELFCLASS32   ELFClass = 1
	ELFCLASS64   ELFClass = 2
	ELFCLASSNUM  ELFClass = 3
)

//go:generate stringer -type=ELFData
type ELFData uint8

const (
	ELFDATANONE ELFData = 0
	ELFDATA2LSB ELFData = 1
	ELFDATA2MSB ELFData = 2
)

//go:generate stringer -type=ELFType
type ELFType uint16

const (
	ET_NONE   ELFType = 0x000  // An unknown type.
	ET_REL    ELFType = 0x0001 // A relocatable file.
	ET_EXEC   ELFType = 0x0002 // An executable file.
	ET_DYN    ELFType = 0x0003 // A shared object.
	ET_CORE   ELFType = 0x0004 // A core file.
	ET_LOPROC ELFType = 0xff00
	ET_HIPROC ELFType = 0xffff
)

//go:generate stringer -type=ELFMachine
type ELFMachine uint16

const (
	EM_NONE   ELFMachine = 0x0000 // An unknown machine
	EM_M32    ELFMachine = 0x0001 // AT&T WE 32100
	EM_SPARC  ELFMachine = 0x0002 // Sun Microsystems SPARC
	EM_386    ELFMachine = 0x0003 // Intel 80386
	EM_68K    ELFMachine = 0x0004 // Motorola 68000
	EM_88K    ELFMachine = 0x0005 // Motorola 88000
	EM_860    ELFMachine = 0x0007 // Intel 80860
	EM_MIPS   ELFMachine = 0x0008 // MIPS RS3000 (big-endian only)
	EM_PPC    ELFMachine = 0x0014 // PowerPC
	EM_PPC64  ELFMachine = 0x0015 // PowerPC 64-bit
	EM_ARM    ELFMachine = 0x0028 // Advanced RISC Machines
	EM_X86_64 ELFMachine = 0x003e // AMD x86-64
)

//go:generate stringer -type=ELFVersion
type ELFVersion uint32

const (
	EV_NONE    ELFVersion = 0 // Invalid version
	EV_CURRENT ELFVersion = 1 // Current version
)

//go:generate stringer -type=ELFOS_ABI
type ELFOS_ABI uint8

const (
	ELF_OS_ABI_NONE  ELFOS_ABI = 0
	ELF_OS_ABI_LINUX ELFOS_ABI = 3
)

type ELFHeader struct {
	// ELF Identification
	Class ELFClass
	// ELF Identification
	Data ELFData
	// ELF Identification
	OS_ABI ELFOS_ABI
	// ELF Identification
	ABI_Version uint8

	Type    ELFType
	Machine ELFMachine
	Version ELFVersion

	// This member gives the virtual address to which the system first
	// transfers control, thus starting the process. If the file has no
	// associated entry point, this member holds zero.
	Entry Elf64_Addr

	// This member holds the program header table's file offset in bytes. If
	// the file has no program header table, this member holds zero.
	ProgramHeaderTableOffset Elf64_Off

	// This member holds the section header table's file offset in bytes. If
	// the file has no section header table, this member holds zero.
	SectionHeaderTableOffset Elf64_Off

	Flags                        uint32
	HeaderSize                   uint16
	ProgramHeaderEntrySize       uint16
	ProgramHeaderNumberOfEntries uint16
	SectionHeaderEntrySize       uint16
	SectionHeaderNumberOfEntries uint16

	// This member holds the section header table index of the entry associated
	// with the section name string table. If the file has no section name
	// string table, this member holds the value SHN_UNDEF.
	Shstrndx uint16
}

func NewELFHeader() *ELFHeader {
	return &ELFHeader{
		Class:       ELFCLASS64,
		Data:        ELFDATA2LSB,
		OS_ABI:      ELF_OS_ABI_LINUX,
		ABI_Version: 0,
		Type:        ET_EXEC,
		Machine:     EM_X86_64,
		Version:     EV_CURRENT,
	}
}

func ParseELFHeader(r *bytes.Reader) (*ELFHeader, error) {
	if r.Size() < 16 {
		return nil, errors.New("Buffer too small to contain ELF header")
	}
	header := make([]byte, 16)
	n, err := r.Read(header)

	if n != 16 {
		return nil, errors.New("Buffer too small to contain ELF header")
	}
	if err != nil {
		return nil, err
	}
	if header[0] != 0x7f || header[1] != 'E' || header[2] != 'L' || header[3] != 'F' {
		return nil, errors.New("Buffer does not start with header ELF number")
	}
	result := &ELFHeader{}
	result.Class = ELFClass(header[4])
	result.Data = ELFData(header[5])
	result.OS_ABI = ELFOS_ABI(header[7])
	result.ABI_Version = header[8]
	byteOrder := result.GetByteOrder()
	if err := binary.Read(r, byteOrder, &result.Type); err != nil {
		return nil, err
	}
	if err := binary.Read(r, byteOrder, &result.Machine); err != nil {
		return nil, err
	}
	if err := binary.Read(r, byteOrder, &result.Version); err != nil {
		return nil, err
	}
	if err := binary.Read(r, byteOrder, &result.Entry); err != nil {
		return nil, err
	}
	if err := binary.Read(r, byteOrder, &result.ProgramHeaderTableOffset); err != nil {
		return nil, err
	}
	if err := binary.Read(r, byteOrder, &result.SectionHeaderTableOffset); err != nil {
		return nil, err
	}
	if err := binary.Read(r, byteOrder, &result.Flags); err != nil {
		return nil, err
	}
	if err := binary.Read(r, byteOrder, &result.HeaderSize); err != nil {
		return nil, err
	}
	if err := binary.Read(r, byteOrder, &result.ProgramHeaderEntrySize); err != nil {
		return nil, err
	}
	if err := binary.Read(r, byteOrder, &result.ProgramHeaderNumberOfEntries); err != nil {
		return nil, err
	}
	if err := binary.Read(r, byteOrder, &result.SectionHeaderEntrySize); err != nil {
		return nil, err
	}
	if err := binary.Read(r, byteOrder, &result.SectionHeaderNumberOfEntries); err != nil {
		return nil, err
	}
	if err := binary.Read(r, byteOrder, &result.Shstrndx); err != nil {
		return nil, err
	}
	return result, nil
}

func (e *ELFHeader) GetByteOrder() binary.ByteOrder {
	if e.Data == ELFDATA2LSB {
		return binary.LittleEndian
	}
	return binary.BigEndian
}

func (e *ELFHeader) Encode() ([]uint8, error) {
	e_ident := make([]uint8, 16)
	e_ident[0] = 0x7f
	e_ident[1] = 'E'
	e_ident[2] = 'L'
	e_ident[3] = 'F'
	e_ident[4] = uint8(e.Class)
	e_ident[5] = uint8(e.Data)
	e_ident[6] = uint8(EV_CURRENT)
	e_ident[7] = uint8(e.OS_ABI)
	e_ident[8] = uint8(e.ABI_Version)
	byteOrder := e.GetByteOrder()
	buffer := bytes.NewBuffer(e_ident)
	if err := binary.Write(buffer, byteOrder, e.Type); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, byteOrder, e.Machine); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, byteOrder, e.Version); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, byteOrder, e.Entry); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, byteOrder, e.ProgramHeaderTableOffset); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, byteOrder, e.SectionHeaderTableOffset); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, byteOrder, e.Flags); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, byteOrder, e.HeaderSize); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, byteOrder, e.ProgramHeaderEntrySize); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, byteOrder, e.ProgramHeaderNumberOfEntries); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, byteOrder, e.SectionHeaderEntrySize); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, byteOrder, e.SectionHeaderNumberOfEntries); err != nil {
		return nil, err
	}
	if err := binary.Write(buffer, byteOrder, e.Shstrndx); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (e *ELFHeader) String() string {
	table := [][]string{
		[]string{"ELF Header:"},
		[]string{"  Class:       ", e.Class.String()},
		[]string{"  Data:        ", e.Data.String()},
		[]string{"  OS ABI:      ", e.OS_ABI.String()},
		[]string{"  ABI Version: ", fmt.Sprintf("%v", e.ABI_Version)},
		[]string{"  Type:        ", e.Type.String()},
		[]string{"  Machine:     ", e.Machine.String()},
		[]string{"  Version:     ", e.Version.String(), strconv.Itoa(int(e.Version))},
		[]string{"  Entry:       ", fmt.Sprintf("0x%x", e.Entry)},
		[]string{"  Program Header Table:"},
		[]string{"    Offset:     ", fmt.Sprintf("0x%x", e.ProgramHeaderTableOffset)},
		[]string{"    Entry size: ", fmt.Sprintf("%v", e.ProgramHeaderEntrySize)},
		[]string{"    Num:        ", fmt.Sprintf("%v", e.ProgramHeaderNumberOfEntries)},
		[]string{"  Section Header Table:"},
		[]string{"    Offset:     ", fmt.Sprintf("0x%x", e.SectionHeaderTableOffset)},
		[]string{"    Entry size: ", fmt.Sprintf("%v", e.SectionHeaderEntrySize)},
		[]string{"    Num:        ", fmt.Sprintf("%v", e.SectionHeaderNumberOfEntries)},
		[]string{"  Flags:        ", fmt.Sprintf("0x%x", e.Flags)},
		[]string{"  Shstrndx:     ", fmt.Sprintf("0x%x", e.Shstrndx)},
	}
	result := []string{}
	for _, row := range table {
		result = append(result, strings.Join(row, "\t"))
	}
	return strings.Join(result, "\n")
}
