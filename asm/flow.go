package asm

type RET struct {
}

func (i *RET) Encode() (MachineCode, error) {
	return []uint8{0xc3}, nil
}

func (i *RET) String() string {
	return "ret"
}

type SYSCALL struct {
}

func (i *SYSCALL) Encode() (MachineCode, error) {
	return []uint8{0x0f, 0x05}, nil
}

func (i *SYSCALL) String() string {
	return "syscall"
}
