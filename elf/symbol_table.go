package elf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)

//go:generate stringer -type=SymbolBinding
type SymbolBinding uint8

const (
	// Local symbols are not visible outside the object file containing their
	// definition. Local symbols of the same name may exist in multiple files
	// without interfering with each other.
	SB_LOCAL SymbolBinding = 0

	// Global symbols are visible to all object files being combined. One
	// file’s definition of a global symbol will satisfy another file’s
	// undefined reference to the same global symbol.
	SB_GLOBAL SymbolBinding = 1

	// Weak symbols resemble global symbols, but their definitions have lower
	// precedence.
	SB_WEAK SymbolBinding = 2

	// Processor-specific
	SB_LOPROC SymbolBinding = 13

	// Processor-specific
	SB_HIPROC SymbolBinding = 15
)

//go:generate stringer -type=SymbolType
type SymbolType uint8

const (
	// The symbol’s type is not specified.
	STT_NOTYPE SymbolType = 0

	// The symbol is associated with a data object, such as a variable, an
	// array, and so on.
	STT_OBJECT SymbolType = 1

	// The symbol is associated with a function or other executable code.
	STT_FUNC SymbolType = 2

	// The symbol is associated with a section. Symbol table entries of this
	// type exist primarily for relocation and normally have STB_LOCAL
	// binding.
	STT_SECTION SymbolType = 3

	// Conventionally, the symbol’s name gives the name of the source file
	// associated with the object file. A file symbol has STB_LOCAL binding,
	// its section index is SHN_ABS, and it precedes the other STB_LOCAL
	// symbols for the file, if it is present.
	STT_FILE SymbolType = 4

	// Processor-specific
	STT_LOPROC SymbolType = 13

	// Processor-specific
	STT_HIPROC SymbolType = 15
)

type Symbol struct {
	// This member holds an index into the object file’s symbol string table,
	// which holds the character representations of the symbol names. If the
	// value is non-zero, it represents a string table index that gives the
	// symbol name. Otherwise, the symbol table entry has no name.
	Name uint32
	name string

	// This member specifies the symbol’s binding attribute.
	Binding SymbolBinding

	// This member specifies the symbol’s type attribute.
	Type SymbolType

	// This member currently holds 0 and has no defined meaning.
	Other byte

	// Every symbol table entry is ‘‘defined’’ in relation to some section;
	// this member holds the relevant section header table index.
	Shndx uint16

	// This member gives the value of the associated symbol. Depend- ing on the
	// context, this may be an absolute value, an address, and so on
	Value Elf64_Addr

	// Many symbols have associated sizes. For example, a data object’s size is
	// the number of bytes contained in the object. This member holds 0 if the
	// symbol has no size or an unknown size.
	Size uint64
}

func ParseSymbol(header *ELFHeader, stringTable *StringTable, r *bytes.Reader) (*Symbol, error) {
	result := &Symbol{}
	byteOrder := header.GetByteOrder()
	if err := binary.Read(r, byteOrder, &result.Name); err != nil {
		return nil, err
	}
	if stringTable != nil {
		name, err := stringTable.GetString(int(result.Name))
		if err != nil {
			return nil, fmt.Errorf("Failed to resolve symbol name (%d) from stringtable: %s", result.Name, err)
		}
		result.name = name
	}
	info := uint8(0)
	if err := binary.Read(r, byteOrder, &info); err != nil {
		return nil, err
	}
	result.Binding = SymbolBinding(info >> 4)
	result.Type = SymbolType(info & 0xf)
	if err := binary.Read(r, byteOrder, &result.Other); err != nil {
		return nil, err
	}
	if err := binary.Read(r, byteOrder, &result.Shndx); err != nil {
		return nil, err
	}
	if err := binary.Read(r, byteOrder, &result.Value); err != nil {
		return nil, err
	}
	if err := binary.Read(r, byteOrder, &result.Size); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Symbol) String() string {
	table := [][]string{
		[]string{"Symbol:"},
		[]string{"  Name:          ", s.name, fmt.Sprintf("(%d)", s.Name)},
		[]string{"  Value:         ", fmt.Sprintf("%v", s.Value)},
		[]string{"  Size:          ", fmt.Sprintf("%v", s.Size)},
		[]string{"  Binding:       ", s.Binding.String()},
		[]string{"  Type:          ", s.Type.String()},
		[]string{"  Shndx:         ", fmt.Sprintf("%v", s.Shndx)},
	}
	result := []string{}
	for _, row := range table {
		if len(row) != 0 {
			result = append(result, strings.Join(row, "\t"))
		}
	}
	return strings.Join(result, "\n")
}

type SymbolTable struct {
	Symbols []*Symbol
	lookup  map[string]*Symbol
}

func (s *SymbolTable) GetSymbol(name string) *Symbol {
	return s.lookup[name]
}

func ParseSymbolTable(header *ELFHeader, stringTable *StringTable, r *bytes.Reader) (*SymbolTable, error) {
	result := []*Symbol{}
	lookup := map[string]*Symbol{}
	var err error
	for err != io.EOF {
		sym, err_ := ParseSymbol(header, stringTable, r)
		if err_ != nil {
			if err_ == io.EOF {
				err = err_
				continue
			}
			return nil, err_
		}
		result = append(result, sym)
		lookup[sym.name] = sym
	}
	return &SymbolTable{result, lookup}, nil
}
