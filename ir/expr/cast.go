package expr

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm"
	"github.com/bspaans/jit-compiler/asm/encoding"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

type IR_Cast struct {
	*BaseIRExpression
	Value      IRExpression
	CastToType Type
}

func NewIR_Cast(value IRExpression, typ Type) *IR_Cast {
	return &IR_Cast{
		BaseIRExpression: NewBaseIRExpression(Cast),
		Value:            value,
		CastToType:       typ,
	}
}

func (i *IR_Cast) ReturnType(ctx *IR_Context) Type {
	return i.CastToType
}

func (i *IR_Cast) String() string {
	return fmt.Sprintf("%s(%s)", i.CastToType.String(), i.Value.String())
}

func (i *IR_Cast) Encode(ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {
	result := []lib.Instruction{}
	valueType := i.Value.ReturnType(ctx)
	if valueType == nil {
		return nil, fmt.Errorf("nil return type in %s", i.Value.String())
	}
	if i.CastToType == TUint64 {
		if valueType == TUint64 {
			return i.Value.Encode(ctx, target)
		} else if valueType == TFloat64 {
			tmpReg := ctx.AllocateRegister(TFloat64)
			defer ctx.DeallocateRegister(tmpReg)

			expr, err := i.Value.Encode(ctx, tmpReg)
			if err != nil {
				return nil, err
			}
			for _, code := range expr {
				result = append(result, code)
			}
			cvt := asm.CVTTSD2SI(tmpReg, target)
			ctx.AddInstruction(cvt)
			result = append(result, cvt)
		} else {
			return nil, fmt.Errorf("Unsupport cast operation %s,(%s) in: %s", valueType.String(), i.CastToType.String(), i.String())
		}
	} else if i.CastToType == TFloat64 {
		if valueType == TFloat64 {
			return i.Value.Encode(ctx, target)
		} else if valueType == TUint64 {
			tmpReg := ctx.AllocateRegister(TUint64)
			defer ctx.DeallocateRegister(tmpReg)

			expr, err := i.Value.Encode(ctx, tmpReg)
			if err != nil {
				return nil, err
			}
			for _, code := range expr {
				result = append(result, code)
			}
			cvt := asm.CVTSI2SD(tmpReg, target)
			ctx.AddInstruction(cvt)
			result = append(result, cvt)
		} else {
			return nil, fmt.Errorf("Unsupport cast operation %s,(%s) in: %s",
				valueType.String(),
				i.CastToType.String(),
				i.String())
		}
	} else {
		return nil, fmt.Errorf("Unsupport cast operation %s,(%s) in: %s", valueType.String(), i.CastToType.String(), i.String())
	}
	return result, nil
}
