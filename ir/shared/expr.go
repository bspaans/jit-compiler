package shared

import (
	"github.com/bspaans/jit-compiler/asm/encoding"
	"github.com/bspaans/jit-compiler/lib"
)

type IRExpressionType int
type IRExpression interface {
	Type() IRExpressionType
	ReturnType(ctx *IR_Context) Type
	AddToDataSection(ctx *IR_Context) error
	Encode(ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error)
	String() string
}

//go:generate stringer -type=IRExpressionType
const (
	Uint8       IRExpressionType = iota
	Uint16      IRExpressionType = iota
	Uint32      IRExpressionType = iota
	Uint64      IRExpressionType = iota
	Int8        IRExpressionType = iota
	Int16       IRExpressionType = iota
	Int32       IRExpressionType = iota
	Int64       IRExpressionType = iota
	Float64     IRExpressionType = iota
	ByteArray   IRExpressionType = iota
	StaticArray IRExpressionType = iota
	ArrayIndex  IRExpressionType = iota
	Bool        IRExpressionType = iota
	Struct      IRExpressionType = iota
	StructField IRExpressionType = iota
	Not         IRExpressionType = iota
	Add         IRExpressionType = iota
	Sub         IRExpressionType = iota
	Mul         IRExpressionType = iota
	Div         IRExpressionType = iota
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
