package expr

import (
	"fmt"

	. "github.com/bspaans/jit-compiler/ir/shared"
)

type IR_Syscall struct {
	*BaseIRExpression
	Syscall IRExpression
	Args    []IRExpression
}

func NewIR_Syscall(syscall IRExpression, args []IRExpression) *IR_Syscall {
	return &IR_Syscall{
		BaseIRExpression: NewBaseIRExpression(Syscall),
		Syscall:          syscall,
		Args:             args,
	}
}
func (i *IR_Syscall) ReturnType(ctx *IR_Context) Type {
	return TUint64
}

func (i *IR_Syscall) String() string {
	return fmt.Sprintf("syscall(%v, %v)", i.Syscall, i.Args)
}

func (b *IR_Syscall) AddToDataSection(ctx *IR_Context) error {
	for _, arg := range b.Args {
		if err := arg.AddToDataSection(ctx); err != nil {
			return err
		}
	}
	return nil
}
func (b *IR_Syscall) SSA_Transform(ctx *SSA_Context) (SSA_Rewrites, IRExpression) {
	rewrites, expr := b.Syscall.SSA_Transform(ctx)
	newArgs := make([]IRExpression, len(b.Args))
	for i, arg := range b.Args {
		if IsLiteralOrVariable(arg) {
			newArgs[i] = arg
		} else {
			rw, expr := arg.SSA_Transform(ctx)
			for _, rewrite := range rw {
				rewrites = append(rewrites, rewrite)
			}
			v := ctx.GenerateVariable()
			rewrites = append(rewrites, NewSSA_Rewrite(v, expr))
			newArgs[i] = NewIR_Variable(v)
		}
	}
	if IsLiteralOrVariable(b.Syscall) {
		return rewrites, NewIR_Syscall(b.Syscall, newArgs)
	}
	v := ctx.GenerateVariable()
	rewrites = append(rewrites, NewSSA_Rewrite(v, expr))
	return rewrites, NewIR_Syscall(NewIR_Variable(v), newArgs)
}

type IR_Syscall_Linux uint

const (
	IR_Syscall_Linux_Read  IR_Syscall_Linux = 0
	IR_Syscall_Linux_Write IR_Syscall_Linux = 1
	IR_Syscall_Linux_Open  IR_Syscall_Linux = 2
	IR_Syscall_Linux_Close IR_Syscall_Linux = 3
)

func NewIR_LinuxWrite(fid IRExpression, b []uint8, size int) IRExpression {
	return NewIR_Syscall(NewIR_Uint64(uint64(IR_Syscall_Linux_Write)), []IRExpression{fid, NewIR_ByteArray(b), NewIR_Uint64(uint64(size))})
}
func NewIR_LinuxOpen(filename string, flags, mode int) IRExpression {
	return NewIR_Syscall(NewIR_Uint64(uint64(IR_Syscall_Linux_Open)), []IRExpression{NewIR_ByteArray([]uint8(filename + "\x00")), NewIR_Uint64(uint64(flags)), NewIR_Uint64(uint64(mode))})
}
func NewIR_LinuxClose(fid uint64) IRExpression {
	return NewIR_Syscall(NewIR_Uint64(uint64(IR_Syscall_Linux_Close)), []IRExpression{NewIR_Uint64(fid)})
}
