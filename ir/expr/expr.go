package expr

import (
	"github.com/bspaans/jit/asm"
	. "github.com/bspaans/jit/ir/shared"
)

type IRExpressionType int
type IRExpression interface {
	Type() IRExpressionType
	ReturnType(ctx *IR_Context) Type
	AddToDataSection(ctx *IR_Context)
	Encode(ctx *IR_Context, target *asm.Register) ([]asm.Instruction, error)
	String() string
}

const (
	Uint64    IRExpressionType = iota
	Float64   IRExpressionType = iota
	ByteArray IRExpressionType = iota
	Bool      IRExpressionType = iota
	Not       IRExpressionType = iota
	Add       IRExpressionType = iota
	Variable  IRExpressionType = iota
	Equals    IRExpressionType = iota
	Syscall   IRExpressionType = iota
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
func (b *BaseIRExpression) AddToDataSection(ctx *IR_Context) {}

func IREXpression_length(expr IRExpression, ctx *IR_Context, target *asm.Register) (int, error) {
	commit := ctx.Commit
	ctx.Commit = false
	instr, err := expr.Encode(ctx, target)
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
