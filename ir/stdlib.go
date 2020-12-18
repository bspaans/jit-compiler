package ir

const Stdlib = `
func Write(fid uint64, str []uint8, len uint64) int64 { 
	return syscall(1, fid, str, len) 
} 
func Open(filename []uint8, flags uint64, mode uint64) int64 { 
	return syscall(2, filename, flags, mode) 
} 
func Close(fid uint64) int64 { 
	return syscall(3, fid) 
} 
func Max(i int64, j int64) int64 {
	if i > j {
		return i
	} else {
		return j
	}
}
`
