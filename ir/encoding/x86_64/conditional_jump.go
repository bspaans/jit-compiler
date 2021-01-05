package x86_64

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm/x86_64"
	"github.com/bspaans/jit-compiler/asm/x86_64/encoding"
	"github.com/bspaans/jit-compiler/ir/expr"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

func conditionalJump(ctx *IR_Context, condition IRExpression, stmtLen int) ([]lib.Instruction, error) {

	jmpSize := 2
	reg := ctx.AllocateRegister(TBool)
	defer ctx.DeallocateRegister(reg)

	var result []lib.Instruction
	var instr []lib.Instruction
	var err error
	switch condition.(type) {
	case *expr.IR_Equals:
		c := condition.(*expr.IR_Equals)
		result, err = encode_IR_Equals(c, ctx, reg, false)
		instr = []lib.Instruction{
			x86_64.JNE(encoding.Uint8(stmtLen + jmpSize)),
		}
	case *expr.IR_LT:
		c := condition.(*expr.IR_LT)
		result, err = encode_IR_LT(c, ctx, reg, false)
		if IsSignedInteger(c.Op1.ReturnType(ctx)) {
			instr = []lib.Instruction{
				x86_64.JNL(encoding.Uint8(stmtLen + jmpSize)),
			}
		} else {
			instr = []lib.Instruction{
				x86_64.JNB(encoding.Uint8(stmtLen + jmpSize)),
			}
		}
	case *expr.IR_LTE:
		c := condition.(*expr.IR_LTE)
		result, err = encode_IR_LTE(c, ctx, reg, false)
		if IsSignedInteger(c.Op1.ReturnType(ctx)) {
			instr = []lib.Instruction{
				x86_64.JNLE(encoding.Uint8(stmtLen + jmpSize)),
			}
		} else {
			instr = []lib.Instruction{
				x86_64.JNBE(encoding.Uint8(stmtLen + jmpSize)),
			}
		}
	case *expr.IR_GT:
		c := condition.(*expr.IR_GT)
		result, err = encode_IR_GT(c, ctx, reg, false)
		if IsSignedInteger(c.Op1.ReturnType(ctx)) {
			instr = []lib.Instruction{
				x86_64.JNG(encoding.Uint8(stmtLen + jmpSize)),
			}
		} else {
			instr = []lib.Instruction{
				x86_64.JNA(encoding.Uint8(stmtLen + jmpSize)),
			}
		}
	case *expr.IR_GTE:
		c := condition.(*expr.IR_GTE)
		result, err = encode_IR_GTE(c, ctx, reg, false)
		if IsSignedInteger(c.Op1.ReturnType(ctx)) {
			instr = []lib.Instruction{
				x86_64.JNGE(encoding.Uint8(stmtLen + jmpSize)),
			}
		} else {
			instr = []lib.Instruction{
				x86_64.JNAE(encoding.Uint8(stmtLen + jmpSize)),
			}
		}
	case *expr.IR_Not:
		result, err = encode_IR_Not(condition.(*expr.IR_Not), ctx, reg, false)
		instr = []lib.Instruction{
			x86_64.JE(encoding.Uint8(stmtLen + jmpSize)),
		}
	case *expr.IR_And:
		result, err = encode_IR_And(condition.(*expr.IR_And), ctx, reg)
		instr = []lib.Instruction{
			x86_64.JE(encoding.Uint8(stmtLen + jmpSize)),
		}
	case *expr.IR_Or:
		result, err = encode_IR_Or(condition.(*expr.IR_Or), ctx, reg)
		instr = []lib.Instruction{
			x86_64.JE(encoding.Uint8(stmtLen + jmpSize)),
		}
	case *expr.IR_Bool:
		result, err = encodeExpression(condition, ctx, reg)
		instr = []lib.Instruction{
			x86_64.CMP_immediate(1, reg),
			x86_64.JNE(encoding.Uint8(stmtLen + jmpSize)),
		}
	default:
		return nil, fmt.Errorf("Unsupported condition %s (type: %v)", condition.String(), condition.Type())
	}
	if err != nil {
		return nil, err
	}
	for _, inst := range instr {
		ctx.AddInstruction(inst)
		result = append(result, inst)
	}
	return result, nil
}
