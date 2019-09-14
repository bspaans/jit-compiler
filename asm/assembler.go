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

type RET struct {
}

func (i *RET) Encode() (MachineCode, error) {
	return []uint8{0xc3}, nil
}

func (i *RET) String() string {
	return "ret"
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
		&MOV{Uint64(0), rax},
		&MOV{Uint64(0), rcx},
		&MOV{Uint64(0), rdx},
		&MOV{Uint64(0), rbx},
		&MOV{Uint64(0), rbp},
		&MOV{Uint64(0), rsi},
		&MOV{Uint64(0), rdi},
		&MOV{Uint64(0), &DisplacedRegister{rsp, 8}},
		&MOV{Uint64(0xffff), rdi},
		&INC{rax},
		&CMP{rdi, rax},
		&JNE{Uint8(0xf9)},
		&MOV{Uint64(123), rcx},
		&ADD{rcx, rax},
		&ADD{Uint32(2), rax},
		&MOV{rax, &DisplacedRegister{rsp, 8}},
		&RET{},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(b)
	b.Execute()
}
