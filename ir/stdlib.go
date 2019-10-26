package ir

import "fmt"

const Stdlib = `
func Write(fid uint64, str []uint8, len uint64) uint64 { 
	return syscall(1, fid, str, len) 
} 
func Open(filename []uint8, flags uint64, mode uint64) uint64 { 
	return syscall(2, filename, flags, mode) 
} 
func Close(fid uint64) uint64 { 
	return syscall(3, fid) 
} 

`

func init() {
	ir := MustParseIR(Stdlib)
	fmt.Println(ir)
}
