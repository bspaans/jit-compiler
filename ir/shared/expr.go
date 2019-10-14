package shared

import (
	"github.com/bspaans/jit/asm/encoding"
	"github.com/bspaans/jit/lib"
)

type IRExpressionType int
type IRExpression interface {
	Type() IRExpressionType
	ReturnType(ctx *IR_Context) Type
	AddToDataSection(ctx *IR_Context) error
	Encode(ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error)
	String() string
}

const (
	Uint64      IRExpressionType = iota
	Float64     IRExpressionType = iota
	ByteArray   IRExpressionType = iota
	StaticArray IRExpressionType = iota
	ArrayIndex  IRExpressionType = iota
	Bool        IRExpressionType = iota
	Not         IRExpressionType = iota
	Add         IRExpressionType = iota
	Variable    IRExpressionType = iota
	Equals      IRExpressionType = iota
	Syscall     IRExpressionType = iota
	Cast        IRExpressionType = iota
	Function    IRExpressionType = iota
	Call        IRExpressionType = iota
)

type BaseIRExpression struct {
	typ IRExpressionType
}

func NewBaseIRExpression(typ IRExpressionType) *BaseIRExpression {
	return &BaseIRExpression{
		typ: typ,
	}
}

func (b *BaseIRExpression) Type() IRExpressionType {
	return b.typ
}
func (b *BaseIRExpression) AddToDataSection(ctx *IR_Context) error {
	return nil
}

func IREXpression_length(expr IRExpression, ctx *IR_Context, target encoding.Operand) (int, error) {
	commit := ctx.Commit
	ctx.Commit = false
	instr, err := expr.Encode(ctx, target)
	if err != nil {
		return 0, err
	}
	code, err := lib.Instructions(instr).Encode()
	if err != nil {
		return 0, err
	}
	ctx.Commit = commit
	return len(code), nil
}
