package expr

import (
	"fmt"

	"github.com/bspaans/jit/asm"
	"github.com/bspaans/jit/asm/encoding"
	. "github.com/bspaans/jit/ir/shared"
	"github.com/bspaans/jit/lib"
)

type IR_Syscall struct {
	*BaseIRExpression
	Syscall uint
	Args    []IRExpression
}

func NewIR_Syscall(syscall_nr uint, args []IRExpression) *IR_Syscall {
	return &IR_Syscall{
		Syscall: syscall_nr,
		Args:    args,
	}
}
func (i *IR_Syscall) ReturnType(ctx *IR_Context) Type {
	return TUint64
}

func (i *IR_Syscall) String() string {
	return fmt.Sprintf("syscall(%v, %v)", i.Syscall, i.Args)
}

func (i *IR_Syscall) Encode(ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {

	result, _, clobbered, err := ABI_Call_Setup(ctx, i.Args, TUint64)
	if err != nil {
		return nil, err
	}

	instr := []lib.Instruction{
		&asm.MOV{encoding.Uint64(uint64(i.Syscall)), encoding.Rax},
		&asm.SYSCALL{},
		&asm.MOV{encoding.Rax, target},
	}
	for _, inst := range instr {
		result = append(result, inst)
		ctx.AddInstruction(inst)
	}
	restore := RestoreRegisters(ctx, clobbered)
	result = result.Add(restore)
	return result, nil
}

func (b *IR_Syscall) AddToDataSection(ctx *IR_Context) error {
	for _, arg := range b.Args {
		if err := arg.AddToDataSection(ctx); err != nil {
			return err
		}
	}
	return nil
}

type IR_Syscall_Linux uint

const (
	IR_Syscall_Linux_Read  IR_Syscall_Linux = 0
	IR_Syscall_Linux_Write IR_Syscall_Linux = 1
	IR_Syscall_Linux_Open  IR_Syscall_Linux = 2
	IR_Syscall_Linux_Close IR_Syscall_Linux = 3
)

func NewIR_LinuxWrite(fid IRExpression, b []uint8, size int) IRExpression {
	return NewIR_Syscall(uint(IR_Syscall_Linux_Write), []IRExpression{fid, NewIR_ByteArray(b), NewIR_Uint64(uint64(size))})
}
func NewIR_LinuxOpen(filename string, flags, mode int) IRExpression {
	return NewIR_Syscall(uint(IR_Syscall_Linux_Open), []IRExpression{NewIR_ByteArray([]uint8(filename + "\x00")), NewIR_Uint64(uint64(flags)), NewIR_Uint64(uint64(mode))})
}
func NewIR_LinuxClose(fid uint64) IRExpression {
	return NewIR_Syscall(uint(IR_Syscall_Linux_Close), []IRExpression{NewIR_Uint64(fid)})
}
