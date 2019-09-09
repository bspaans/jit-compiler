package asm

import (
	"encoding/hex"
	"errors"
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
	return hex.EncodeToString(m)
}

func (m MachineCode) Execute() {
	code := []uint8{}
	for _, c := range m {
		code = append(code, c)
	}
	// Use rax register as return arg:
	// movq %rax, 0x8(%rsp);
	// ret
	// e.g. move the return value on the stack.
	for _, c := range []uint8{0x48, 0x89, 0x44, 0x24, 0x08, 0xc3} {
		code = append(code, c)
	}
	mmapFunc, err := syscall.Mmap(
		-1,
		0,
		len(code),
		syscall.PROT_READ|syscall.PROT_WRITE|syscall.PROT_EXEC, syscall.MAP_PRIVATE|syscall.MAP_ANONYMOUS,
	)
	if err != nil {
		fmt.Printf("mmap err: %v", err)
	}
	for i, b := range code {
		mmapFunc[i] = b
	}
	type execFunc func() int
	unsafeFunc := (uintptr)(unsafe.Pointer(&mmapFunc))
	f := *(*execFunc)(unsafe.Pointer(&unsafeFunc))
	fmt.Println("Result:", f())
}

func EncodeModRM(mod, reg, rm uint8) uint8 {
	return rm
}

type Instruction interface {
	Encode() (MachineCode, error)
	String() string
}

type INC struct {
	Register *Register
}

func (i *INC) Encode() (MachineCode, error) {
	if i.Register == nil {
		return nil, errors.New("Missing register")
	}
	if i.Register.Size == QUADWORD {
		rexB := i.Register.Register > 7
		rex := NewREXPrefix(true, false, false, rexB).Encode()
		modrm := NewModRM(DirectRegisterMode, i.Register.Encode(), 0).Encode()
		return []uint8{rex, 0xff, modrm}, nil
	}
	return nil, errors.New("Unsupported register size")
}

func (i *INC) String() string {
	return "inc " + i.Register.String()
}

type DEC struct {
	Register *Register
}

func (i *DEC) Encode() (MachineCode, error) {
	if i.Register == nil {
		return nil, errors.New("Missing register")
	}
	if i.Register.Size == QUADWORD {
		rexB := i.Register.Register > 7
		rex := NewREXPrefix(true, false, false, rexB).Encode()
		modrm := NewModRM(DirectRegisterMode, i.Register.Encode(), 1).Encode()
		return []uint8{rex, 0xff, modrm}, nil
	}
	return nil, errors.New("Unsupported register size")
}

func (i *DEC) String() string {
	return "dec " + i.Register.String()
}

type RET struct {
}

func (i *RET) Encode() (MachineCode, error) {
	return []uint8{0xc3}, nil
}

func (i *RET) String() string {
	return "ret"
}

type Type uint8

func CompileInstruction(instr []Instruction) (MachineCode, error) {
	result := []uint8{}
	for _, i := range instr {
		b, err := i.Encode()
		if err != nil {
			return nil, err
		}
		fmt.Println(i)
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
		&INC{rax},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(b)
	b.Execute()
}
