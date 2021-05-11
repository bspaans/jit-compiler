package lib

import (
	"encoding/hex"
	"fmt"
	"runtime"
	"sync"
	"syscall"
	"unsafe"
)

const mmapProtFlags = syscall.PROT_READ | syscall.PROT_WRITE | syscall.PROT_EXEC

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
	mmapFunc, err := syscall.Mmap(-1, 0, len(m), mmapProtFlags, syscall.MAP_PRIVATE|mmapFlags)
	if err != nil {
		panic(fmt.Sprintf("mmap err: %v", err))
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

type execFunc func() int

func (m MachineCode) Compile() *CompiledCode {
	mem, err := syscall.Mmap(-1, 0, len(m), mmapProtFlags, syscall.MAP_PRIVATE|mmapFlags)
	if err != nil {
		panic(fmt.Sprintf("mmap err: %v", err))
	}

	cc := &CompiledCode{mem: mem, mc: m}
	runtime.SetFinalizer(cc, func(cc *CompiledCode) {
		cc.release()
	})
	return cc
}

func (m MachineCode) Add(m2 MachineCode) MachineCode {
	m = append(m, m2...)
	return m
}

type CompiledCode struct {
	mem []byte
	mc  MachineCode
	mux sync.Mutex
}

func (cc *CompiledCode) Execute() int {
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if cc.mem == nil {
		panic("used after munmap")
	}
	copy(cc.mem, cc.mc)
	unsafeFunc := (uintptr)(unsafe.Pointer(&cc.mem))
	f := *(*execFunc)(unsafe.Pointer(&unsafeFunc))
	return f()
}

func (cc *CompiledCode) release() {
	if cc.mem != nil {
		syscall.Munmap(cc.mem)
		cc.mem = nil
		cc.mc = nil
	}
}
