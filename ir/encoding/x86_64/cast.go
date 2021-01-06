package x86_64

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm/x86_64"
	"github.com/bspaans/jit-compiler/asm/x86_64/encoding"
	"github.com/bspaans/jit-compiler/ir/expr"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

func encode_IR_Cast(i *expr.IR_Cast, ctx *IR_Context, target lib.Operand) ([]lib.Instruction, error) {
	ctx.AddInstruction("cast " + encoding.Comment(i.String()))
	result := []lib.Instruction{}
	valueType := i.Value.ReturnType(ctx)
	if valueType == nil {
		return nil, fmt.Errorf("nil return type in %s", i.Value.String())
	}
	// TODO: use movsx and movzx
	if i.CastToType == TUint64 {
		if valueType == TUint64 {
			return encodeExpression(i.Value, ctx, target)
		} else if valueType == TUint32 {
			tmpReg := ctx.AllocateRegister(valueType)
			defer ctx.DeallocateRegister(tmpReg)
			expr, err := encodeExpression(i.Value, ctx, tmpReg)
			if err != nil {
				return nil, err
			}
			for _, code := range expr {
				result = append(result, code)
			}
			mov := x86_64.MOV(tmpReg.(*encoding.Register).Get64BitRegister(), target)
			ctx.AddInstruction(mov)
			result = append(result, mov)
			return result, nil
		} else if valueType == TUint16 || valueType == TUint8 {
			tmpReg := ctx.AllocateRegister(valueType)
			defer ctx.DeallocateRegister(tmpReg)
			expr, err := encodeExpression(i.Value, ctx, tmpReg)
			if err != nil {
				return nil, err
			}
			for _, code := range expr {
				result = append(result, code)
			}
			mov := x86_64.MOVZX(tmpReg, target)
			ctx.AddInstruction(mov)
			result = append(result, mov)
			return result, nil
		} else if valueType == TFloat64 {
			tmpReg := ctx.AllocateRegister(TFloat64)
			defer ctx.DeallocateRegister(tmpReg)

			expr, err := encodeExpression(i.Value, ctx, tmpReg)
			if err != nil {
				return nil, err
			}
			for _, code := range expr {
				result = append(result, code)
			}
			cvt := x86_64.CVTTSD2SI(tmpReg, target)
			ctx.AddInstruction(cvt)
			result = append(result, cvt)
			return result, nil
		}
	} else if i.CastToType == TUint8 {
		if valueType == TUint64 || valueType == TUint32 || valueType == TUint16 || valueType == TUint8 {
			result, err := encodeExpression(i.Value, ctx, target)
			if err != nil {
				return nil, err
			}
			return result, nil
		}
	} else if i.CastToType == TUint16 {
		if valueType == TUint64 || valueType == TUint32 || valueType == TUint16 || valueType == TUint8 {
			result, err := encodeExpression(i.Value, ctx, target)
			if err != nil {
				return nil, err
			}
			return result, nil
		}
	} else if i.CastToType == TUint32 {
		if valueType == TUint64 || valueType == TUint32 || valueType == TUint16 || valueType == TUint8 {
			result, err := encodeExpression(i.Value, ctx, target)
			if err != nil {
				return nil, err
			}
			return result, nil
		}
	} else if i.CastToType == TFloat64 {
		if valueType == TFloat64 {
			return encodeExpression(i.Value, ctx, target)
		} else if valueType == TUint64 {
			tmpReg := ctx.AllocateRegister(TUint64)
			defer ctx.DeallocateRegister(tmpReg)

			expr, err := encodeExpression(i.Value, ctx, tmpReg)
			if err != nil {
				return nil, err
			}
			for _, code := range expr {
				result = append(result, code)
			}
			cvt := x86_64.CVTSI2SD(tmpReg, target)
			ctx.AddInstruction(cvt)
			result = append(result, cvt)
			return result, nil
		}
	}
	return nil, fmt.Errorf("Unsupported cast operation %s -> (%s) in: %s", valueType.String(), i.CastToType.String(), i.String())
}
