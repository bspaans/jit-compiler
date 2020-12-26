package elf

import (
	"bytes"
	"io/ioutil"
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
	//ProgramHeaderTable

	Sections []*Section
}

func ParseELF(r *bytes.Reader) (*ELF, error) {
	header, err := ParseELFHeader(r)
	if err != nil {
		return nil, err
	}

	table, err := ParseSectionHeaderTable(header, r)
	if err != nil {
		return nil, err
	}
	sections, err := ParseSections(header, table, r)
	if err != nil {
		return nil, err
	}
	return &ELF{
		ELFHeader: header,
		Sections:  sections,
	}, nil
}

func ParseELFFile(path string) (*ELF, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseELF(bytes.NewReader(f))
}

func (e *ELF) String() string {
	result := e.ELFHeader.String()
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

func init() {
	/*
		elf, err := ParseELFFile("../jit-compiler")
		if err != nil {
			panic(err)
		}
		fmt.Println(elf)
		fmt.Println(string(elf.GetSection(".note.go.buildid").Data))
	*/
}
