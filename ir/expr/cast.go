package expr

import (
	"fmt"

	"github.com/bspaans/jit/asm"
	. "github.com/bspaans/jit/ir/shared"
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
	return fmt.Sprintf("(%s).(%s)", i.Value.String(), i.CastToType.String())
}

func (i *IR_Cast) Encode(ctx *IR_Context, target *asm.Register) ([]asm.Instruction, error) {
	result := []asm.Instruction{}
	valueType := i.Value.ReturnType(ctx)
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
			cvt := &asm.CVTTSD2SI{tmpReg, target}
			ctx.AddInstruction(cvt)
			result = append(result, cvt)
		} else {
			return nil, fmt.Errorf("Unsupport cast operation: " + i.String())
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
			cvt := &asm.CVTSI2SD{tmpReg, target}
			ctx.AddInstruction(cvt)
			result = append(result, cvt)
		} else {
			return nil, fmt.Errorf("Unsupport cast operation: " + i.String())
		}
	} else {
		return nil, fmt.Errorf("Unsupport cast operation: " + i.String())
	}
	return result, nil
}
