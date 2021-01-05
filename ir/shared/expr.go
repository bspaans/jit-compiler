package shared

import (
	"github.com/bspaans/jit-compiler/asm/x86_64/encoding"
	"github.com/bspaans/jit-compiler/lib"
)

type IRExpressionType int
type IRExpression interface {
	Type() IRExpressionType
	ReturnType(ctx *IR_Context) Type
	AddToDataSection(ctx *IR_Context) error
	String() string
	SSA_Transform(*SSA_Context) (SSA_Rewrites, IRExpression)
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
	And         IRExpressionType = iota
	Or          IRExpressionType = iota
	Not         IRExpressionType = iota
	Add         IRExpressionType = iota
	Sub         IRExpressionType = iota
	Mul         IRExpressionType = iota
	Div         IRExpressionType = iota
	Variable    IRExpressionType = iota
	Equals      IRExpressionType = iota
	LT          IRExpressionType = iota
	LTE         IRExpressionType = iota
	GT          IRExpressionType = iota
	GTE         IRExpressionType = iota
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

func IsLiteral(e IRExpression) bool {
	t := e.Type()
	return t == Uint8 || t == Uint16 || t == Uint32 || t == Uint64 ||
		t == Int8 || t == Int16 || t == Int32 || t == Int64 ||
		t == Float64 || t == ByteArray || t == StaticArray || t == Bool
}

func IsVariable(e IRExpression) bool {
	return e.Type() == Variable
}

func IsLiteralOrVariable(e IRExpression) bool {
	return IsVariable(e) || IsLiteral(e)
}

func IREXpression_length(expr IRExpression, ctx *IR_Context, target encoding.Operand) (int, error) {
	commit := ctx.Commit
	ctx.Commit = false
	instr, err := ctx.Architecture.EncodeExpression(expr, ctx, target)
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
