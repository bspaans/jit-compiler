package x86_64

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm/x86_64"
	"github.com/bspaans/jit-compiler/asm/x86_64/encoding"
	"github.com/bspaans/jit-compiler/ir/expr"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

func encode_IR_Div(i *expr.IR_Div, ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {
	ctx.AddInstruction("operator " + encoding.Comment(i.String()))
	returnType1, returnType2 := i.Op1.ReturnType(ctx), i.Op2.ReturnType(ctx)
	if returnType1 != returnType2 {
		return nil, fmt.Errorf("Unsupported types (%s, %s) in / IR operation: %s", returnType1, returnType2, i.String())
	}
	if IsFloat(returnType1) {
		return encode_Operator(i.Op1, i.Op2, x86_64.IDIV2, i.String(), ctx, target)
	}
	if IsInteger(returnType1) {

		allocator := ctx.Allocator.(*X86_64_Allocator)
		raxInUse := allocator.Registers[0]

		shouldPreserveRdx := (returnType1.Width() != lib.BYTE) && allocator.Registers[2]
		shouldPreserveRax := target.(*encoding.Register).Register != 0 && raxInUse
		var tmpRdx *encoding.Register
		var tmpRax *encoding.Register

		result := lib.Instructions{}
		ctxCopy := ctx

		// Preserve the %rdx register
		if shouldPreserveRdx {
			tmpRdx = ctx.AllocateRegister(TUint64)
			defer ctx.DeallocateRegister(tmpRdx)
			preserveRdx := x86_64.MOV(encoding.Rdx, tmpRdx)
			result = append(result, preserveRdx)
			ctx.AddInstruction(result...)
			// Replace variables in the variablemap that point to rdx with the new register
			ctxCopy = ctxCopy.Copy()
			for v, vTarget := range ctxCopy.VariableMap {
				if r, ok := vTarget.(*encoding.Register); ok && r.Register == 2 {
					ctxCopy.VariableMap[v] = tmpRdx
				}
			}
		}
		// Preserve the %rax register
		if shouldPreserveRax {
			ctxCopy = ctxCopy.Copy()
			allocator = ctxCopy.Allocator.(*X86_64_Allocator)
			// Make sure we don't allocate %rdx
			if !allocator.Registers[2] && (returnType1.Width() != lib.BYTE) {
				allocator.Registers[2] = true
			}
			tmpRax = ctxCopy.AllocateRegister(TUint64)
			defer ctxCopy.DeallocateRegister(tmpRax)
			preserveRax := x86_64.MOV(encoding.Rax, tmpRax)
			result = append(result, preserveRax)
			ctx.AddInstruction(result...)
			// Replace variables in the variablemap that point to rax with the new register
			for v, vTarget := range ctxCopy.VariableMap {
				if r, ok := vTarget.(*encoding.Register); ok && r.Register == 0 {
					ctxCopy.VariableMap[v] = tmpRax
				}
			}
		}

		rax := encoding.Rax.ForOperandWidth(returnType1.Width())

		op1, err := encodeExpression(i.Op1, ctxCopy, rax)
		if err != nil {
			return nil, err
		}

		result = result.Add(op1)

		var reg encoding.Operand
		if i.Op2.Type() == Variable {
			variable := i.Op2.(*expr.IR_Variable).Value
			reg = ctxCopy.VariableMap[variable]
		} else {
			reg = ctxCopy.AllocateRegister(returnType2)
			defer ctxCopy.DeallocateRegister(reg.(*encoding.Register))

			expr, err := encodeExpression(i.Op2, ctxCopy, reg)
			if err != nil {
				return nil, err
			}
			result = lib.Instructions(result).Add(expr)
		}

		zeroRegisters := map[Type]*encoding.Register{
			TUint8:  encoding.Ah,
			TUint16: encoding.Dx,
			TUint32: encoding.Edx,
			TUint64: encoding.Rdx,
		}
		if IsSignedInteger(returnType1) {
			var instr lib.Instruction
			if returnType1 == TInt64 {
				instr = x86_64.CQO()
			} else if returnType1 == TInt32 {
				instr = x86_64.CDQ()
			} else if returnType1 == TInt16 {
				instr = x86_64.CWD()
			} else if returnType1 == TInt8 {
				instr = x86_64.CBW()
			}
			result = append(result, instr)
			ctx.AddInstruction(result...)
		} else {
			zero := zeroRegisters[returnType1]
			xor := x86_64.XOR(zero, zero)
			result = append(result, xor)
			ctx.AddInstruction(result...)
		}

		instr := x86_64.DIV(reg)
		if IsSignedInteger(returnType1) {
			instr = x86_64.IDIV1(reg)
		}
		ctx.AddInstruction(instr)
		result = append(result, instr)

		if rax != target {
			mov := x86_64.MOV(rax, target)
			ctx.AddInstruction(mov)
			result = append(result, mov)
		}

		// Restore %rax
		if shouldPreserveRax {
			restore := x86_64.MOV(tmpRax, encoding.Rax)
			ctx.AddInstruction(restore)
			result = append(result, restore)
		}
		// Restore %rdx
		if shouldPreserveRdx {
			restore := x86_64.MOV(tmpRdx, encoding.Rdx)
			ctx.AddInstruction(restore)
			result = append(result, restore)
		}
		return result, nil
	}
	return nil, fmt.Errorf("Unsupported types (%s, %s) in / IR operation: %s", returnType1, returnType2, i.String())
}
