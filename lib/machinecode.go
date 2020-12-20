package lib

import (
	"encoding/hex"
	"fmt"
	"syscall"
	"unsafe"
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
		if i%16 == 15 && i+1 < len(h) {
			result = append(result, '\n', ' ', ' ')
		}
	}
	return string(result)
}

func (m MachineCode) Execute(debug bool) int {
	mmapFunc, err := syscall.Mmap(
		-1,
		0,
		len(m),
		syscall.PROT_READ|syscall.PROT_WRITE|syscall.PROT_EXEC, syscall.MAP_PRIVATE|mmapFlags,
	)
	if err != nil {
		fmt.Printf("mmap err: %v", err)
	}
	for i, b := range m {
		mmapFunc[i] = b
	}
	type execFunc func() int
	unsafeFunc := (uintptr)(unsafe.Pointer(&mmapFunc))
	f := *(*execFunc)(unsafe.Pointer(&unsafeFunc))
	value := f()
	if debug {
		fmt.Println("\nResult :", value)
		fmt.Printf("Hex    : %x\n", value)
		fmt.Printf("Size   : %d bytes\n\n", len(m))
	}
	return value
}
func (m MachineCode) Add(m2 MachineCode) MachineCode {
	for _, code := range m2 {
		m = append(m, code)
	}
	return m
}
