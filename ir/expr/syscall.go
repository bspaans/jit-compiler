package expr

import (
	"fmt"

	"github.com/bspaans/jit/asm"
	. "github.com/bspaans/jit/ir/shared"
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

func (i *IR_Syscall) Encode(ctx *IR_Context, target asm.Operand) ([]asm.Instruction, error) {

	result := []asm.Instruction{}
	if ctx.Registers[0] {
		push := &asm.PUSH{asm.Rax}
		result = append(result, push)
		ctx.AddInstruction(push)
	}

	// TODO result, clobbered, err := ABI_Call_Setup(ctx, args, i.ReturnType(ctx))
	targets := []*asm.Register{asm.Rdi, asm.Rsi, asm.Rdx, asm.R10, asm.R8, asm.R9}
	targetRegisterIndices := []uint8{7, 6, 2, 10, 8, 9}
	clobbered := 0
	for j, argTarget := range targets {
		if j < len(i.Args) {
			// Push registers on the stack if they are in use
			registerIndex := targetRegisterIndices[j]
			if ctx.Registers[registerIndex] {
				reg := asm.Get64BitRegisterByIndex(registerIndex)
				result = append(result, &asm.PUSH{reg})
				ctx.AddInstruction(&asm.PUSH{reg})
				clobbered += 1
			}
			instr, err := i.Args[j].Encode(ctx, argTarget)
			if err != nil {
				return nil, err
			}
			for _, code := range instr {
				result = append(result, code)
			}
		}
	}

	instr := []asm.Instruction{
		&asm.MOV{asm.Uint64(uint64(i.Syscall)), asm.Rax},
		&asm.SYSCALL{},
		&asm.MOV{asm.Rax, target},
	}
	// Restore registers from the stack
	for j := clobbered; j > 0; j-- {
		registerIndex := targetRegisterIndices[j-1]
		reg := asm.Get64BitRegisterByIndex(registerIndex)
		instr = append(instr, &asm.POP{reg})
	}
	// restore rax
	if ctx.Registers[0] {
		instr = append(instr, &asm.POP{asm.Rax})
	}
	for _, inst := range instr {
		result = append(result, inst)
		ctx.AddInstruction(inst)
	}
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
