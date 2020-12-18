package elf

/*

An  executable  file using the ELF file format consists of an ELF header,
followed by a program header table or a section header table, or both.  The
ELF header is always at offset zero of the file.  The program header table
and the section header table's offset in the file are defined in the ELF
header.  The two tables describe the rest of the particularities of the file.


*/

type Elf64_Addr uint64
type Elf64_Off uint64

type ELFClass uint8

const (
	ELFCLASSNONE ELFClass = 0
	ELFCLASS32   ELFClass = 1
	ELFCLASS64   ELFClass = 2
	ELFCLASSNUM  ELFClass = 3
)

type ELFData uint8

const (
	ELFDATANONE ELFData = 0
	ELFDATA2LSB ELFData = 1
	ELFDATA2MSB ELFData = 2
)

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

type ELFVersion uint32

const (
	EV_NONE    ELFVersion = 0 // Invalid version
	EV_CURRENT ELFVersion = 1 // Current version
)

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
		ABI_Version: 0, // TODO?
		Type:        ET_EXEC,
		Machine:     EM_X86_64,
		Version:     EV_CURRENT,
	}
}

func (e *ELFHeader) Encode() []uint8 {
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
	return e_ident
}

type ELF struct {
	*ELFHeader

	// A program header table, if present, tells the system how to create a
	// process image. Files used to build a process image (execute a program)
	// must have a program header table; relocatable files do not need one.
	//ProgramHeaderTable
	//Sections
}
