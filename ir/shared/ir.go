package shared

import (
	"github.com/bspaans/jit/asm"
)

type IRType int

const (
	Assignment IRType = iota
	If         IRType = iota
	While      IRType = iota
	Return     IRType = iota
	AndThen    IRType = iota
)

type IR interface {
	Type() IRType
	String() string
	AddToDataSection(ctx *IR_Context) error
	Encode(*IR_Context) ([]asm.Instruction, error)
}

type BaseIR struct {
	typ IRType
}

func NewBaseIR(typ IRType) *BaseIR {
	return &BaseIR{
		typ: typ,
	}
}
func (b *BaseIR) Type() IRType {
	return b.typ
}
func (b *BaseIR) AddToDataSection(ctx *IR_Context) error {
	return nil
}

func IR_Length(stmt IR, ctx *IR_Context) (int, error) {
	commit := ctx.Commit
	ctx.Commit = false
	instr, err := stmt.Encode(ctx)
	if err != nil {
		return 0, err
	}
	code, err := asm.Instructions(instr).Encode()
	if err != nil {
		return 0, err
	}
	ctx.Commit = commit
	return len(code), nil
}
