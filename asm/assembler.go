package asm

import (
	"encoding/hex"
	"fmt"
	"syscall"
	"unsafe"
)

type Size uint8

const (
	BYTE     Size = 1
	WORD     Size = 2
	DOUBLE   Size = 4
	QUADWORD Size = 8
)

type MachineCode []uint8

func (m MachineCode) String() string {
	h := hex.EncodeToString(m)
	result := []rune{' ', ' '}
	for i, c := range h {
		result = append(result, c)
		if i%2 == 1 && i+1 < len(h) {
			result = append(result, ' ')
		}
		if i%14 == 13 && i+1 < len(h) {
			result = append(result, '\n', ' ', ' ')
		}
	}
	return string(result)
}

func (m MachineCode) Execute() {
	mmapFunc, err := syscall.Mmap(
		-1,
		0,
		len(m),
		syscall.PROT_READ|syscall.PROT_WRITE|syscall.PROT_EXEC, syscall.MAP_PRIVATE|syscall.MAP_ANONYMOUS,
	)
	if err != nil {
		fmt.Printf("mmap err: %v", err)
	}
	for i, b := range m {
		mmapFunc[i] = b
	}
	type execFunc func() uint
	unsafeFunc := (uintptr)(unsafe.Pointer(&mmapFunc))
	f := *(*execFunc)(unsafe.Pointer(&unsafeFunc))
	fmt.Println("Result:", f())
}

type Instruction interface {
	Encode() (MachineCode, error)
	String() string
}

type Instructions []Instruction

func (i Instructions) Encode() (MachineCode, error) {
	result := []uint8{}
	for _, instr := range i {
		b, err := instr.Encode()
		if err != nil {
			return nil, err
		}
		for _, code := range b {
			result = append(result, code)
		}
	}
	return result, nil
}

func CompileInstruction(instr []Instruction) (MachineCode, error) {
	result := []uint8{}
	address := 0
	for _, i := range instr {
		b, err := i.Encode()
		if err != nil {
			return nil, err
		}
		fmt.Printf("0x%x: %s\n", address, i.String())
		address += len(b)
		fmt.Println(MachineCode(b))
		for _, code := range b {
			result = append(result, code)
		}
	}
	return result, nil
}

func init() {
	b, err := CompileInstruction([]Instruction{
		&MOV{Uint64(0), Rax},
		&MOV{Uint64(0), Rcx},
		&MOV{Uint64(0), Rdx},
		&MOV{Uint64(0), Rbx},
		&MOV{Uint64(0), Rbp},
		&MOV{Uint64(0), Rsi},
		&MOV{Uint64(0), Rdi},
		&MOV{Uint64(0), &DisplacedRegister{Rsp, 8}},
		&MOV{Uint64(0xffff), Rdi},
		&INC{Rax},
		&CMP{Rdi, Rax},
		&JNE{Uint8(0xf9)},
		&CMP{Rdi, Rax},
		&SETE{Al},
		&MOV{Uint64(123), Rcx},
		&ADD{Rcx, Rax},
		&ADD{Uint32(2), Rax},
		&PUSH{Rax},
		&POP{Rax},
		&PUSHFQ{},
		&POP{Rdx},
		&MOV{Rax, &DisplacedRegister{Rsp, 8}},
		&RET{},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(b)
	b.Execute()
}
