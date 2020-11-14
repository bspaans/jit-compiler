package expr

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm"
	"github.com/bspaans/jit-compiler/asm/encoding"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

type IR_Div struct {
	*BaseIRExpression
	Op1 IRExpression
	Op2 IRExpression
}

func NewIR_Div(op1, op2 IRExpression) *IR_Div {
	return &IR_Div{
		BaseIRExpression: NewBaseIRExpression(Div),
		Op1:              op1,
		Op2:              op2,
	}
}

func (i *IR_Div) ReturnType(ctx *IR_Context) Type {
	return i.Op1.ReturnType(ctx)
}

func (i *IR_Div) Encode(ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {
	ctx.AddInstruction("operator " + encoding.Comment(i.String()))
	returnType1, returnType2 := i.Op1.ReturnType(ctx), i.Op2.ReturnType(ctx)
	if returnType1 != returnType2 {
		return nil, fmt.Errorf("Unsupported types (%s, %s) in / IR operation: %s", returnType1, returnType2, i.String())
	}
	if returnType1 == TFloat64 {
		return NewIR_Operator(asm.IDIV, "*", i.Op1, i.Op2).Encode(ctx, target)
	}
	if returnType1 == TUint64 || returnType1 == TUint32 || returnType1 == TUint8 {

		raxInUse := ctx.Registers[0]

		shouldPreserveRdx := (returnType1 == TUint64 || returnType1 == TUint32) && ctx.Registers[2]
		shouldPreserveRax := target.(*encoding.Register).Register != 0 && raxInUse
		var tmpRdx *encoding.Register
		var tmpRax *encoding.Register

		result := lib.Instructions{}
		ctxCopy := ctx

		// Preserve the %rdx register
		if shouldPreserveRdx {
			tmpRdx = ctx.AllocateRegister(TUint64)
			defer ctx.DeallocateRegister(tmpRdx)
			preserveRdx := asm.MOV(encoding.Rdx, tmpRdx)
			result = append(result, preserveRdx)
			ctx.AddInstructions(result)
			// Replace variables in the variablemap that point to rax with the new register
			ctxCopy = ctxCopy.Copy()
			for v, vTarget := range ctxCopy.VariableMap {
				if r, ok := vTarget.(*encoding.Register); ok && r.Register == 2 {
					ctxCopy.VariableMap[v] = tmpRax
				}
			}
		}
		// Set %ah to 0 in 8bit mode
		if returnType1 == TUint8 {
			mov0 := asm.MOV_immediate(0, encoding.Ah) // TODO use xor
			result = append(result, mov0)
			ctx.AddInstructions(result)

			// Set %edx to 0 in 32 bit mode
		} else if returnType1 == TUint32 {
			mov0 := asm.MOV_immediate(0, encoding.Edx) // TODO use xor
			result = append(result, mov0)
			ctx.AddInstructions(result)

			// Set %rdi to 0 in 64bit mode
		} else if returnType1 == TUint64 {
			mov0 := asm.MOV_immediate(0, encoding.Rdx) // TODO use xor
			result = append(result, mov0)
			ctx.AddInstructions(result)
		}

		// Preserve the %rax register
		if shouldPreserveRax {
			ctxCopy = ctxCopy.Copy()
			// Make sure we don't allocate %rdx
			if !ctxCopy.Registers[2] && (returnType1 == TUint64 || returnType1 == TUint32) {
				ctxCopy.Registers[2] = true
			}
			tmpRax = ctxCopy.AllocateRegister(TUint64)
			defer ctxCopy.DeallocateRegister(tmpRax)
			preserveRax := asm.MOV(encoding.Rax, tmpRax)
			result = append(result, preserveRax)
			ctx.AddInstructions(result)
			// Replace variables in the variablemap that point to rax with the new register
			for v, vTarget := range ctxCopy.VariableMap {
				if r, ok := vTarget.(*encoding.Register); ok && r.Register == 0 {
					ctxCopy.VariableMap[v] = tmpRax
				}
			}
		}

		rax := encoding.Rax.ForOperandWidth(returnType1.Width())

		op1, err := i.Op1.Encode(ctxCopy, rax)
		if err != nil {
			return nil, err
		}

		result = result.Add(op1)

		var reg encoding.Operand
		if i.Op2.Type() == Variable {
			variable := i.Op2.(*IR_Variable).Value
			reg = ctxCopy.VariableMap[variable]
		} else {
			reg = ctxCopy.AllocateRegister(returnType2)
			defer ctxCopy.DeallocateRegister(reg.(*encoding.Register))

			expr, err := i.Op2.Encode(ctxCopy, reg)
			if err != nil {
				return nil, err
			}
			result = lib.Instructions(result).Add(expr)
		}
		instr := asm.DIV(reg)
		ctx.AddInstruction(instr)
		result = append(result, instr)

		mov := asm.MOV(rax, target)
		ctx.AddInstruction(mov)
		result = append(result, mov)

		// Restore %rax
		if shouldPreserveRax {
			restore := asm.MOV(tmpRax, encoding.Rax)
			ctx.AddInstruction(restore)
			result = append(result, restore)
		}
		// Restore %rdx
		if shouldPreserveRdx {
			restore := asm.MOV(tmpRdx, encoding.Rdx)
			ctx.AddInstruction(restore)
			result = append(result, restore)
		}
		return result, nil
	}
	return nil, fmt.Errorf("Unsupported types (%s, %s) in / IR operation: %s", returnType1, returnType2, i.String())
}

func (i *IR_Div) String() string {
	return fmt.Sprintf("%s / %s", i.Op1.String(), i.Op2.String())
}

func (b *IR_Div) AddToDataSection(ctx *IR_Context) error {
	if err := b.Op1.AddToDataSection(ctx); err != nil {
		return err
	}
	return b.Op2.AddToDataSection(ctx)
}
