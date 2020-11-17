package statements

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm"
	"github.com/bspaans/jit-compiler/asm/encoding"
	"github.com/bspaans/jit-compiler/ir/expr"
	"github.com/bspaans/jit-compiler/ir/shared"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/lib"
)

func ConditionalJump(ctx *IR_Context, condition IRExpression, stmtLen int) ([]lib.Instruction, error) {

	jmpSize := 2
	reg := ctx.AllocateRegister(TBool)
	defer ctx.DeallocateRegister(reg)

	var result []lib.Instruction
	var instr []lib.Instruction
	var err error
	switch condition.(type) {
	case *expr.IR_Equals:
		result, err = condition.(*expr.IR_Equals).EncodeWithoutSETE(ctx, reg)
		instr = []lib.Instruction{
			asm.JNE(encoding.Uint8(stmtLen + jmpSize)),
		}
	case *expr.IR_LT:
		c := condition.(*expr.IR_LT)
		result, err = c.EncodeWithoutSETE(ctx, reg)
		if shared.IsSignedInteger(c.Op1.ReturnType(ctx)) {
			instr = []lib.Instruction{
				asm.JNL(encoding.Uint8(stmtLen + jmpSize)),
			}
		} else {
			instr = []lib.Instruction{
				asm.JNB(encoding.Uint8(stmtLen + jmpSize)),
			}
		}
	case *expr.IR_LTE:
		c := condition.(*expr.IR_LTE)
		result, err = c.EncodeWithoutSETE(ctx, reg)
		if shared.IsSignedInteger(c.Op1.ReturnType(ctx)) {
			instr = []lib.Instruction{
				asm.JNLE(encoding.Uint8(stmtLen + jmpSize)),
			}
		} else {
			instr = []lib.Instruction{
				asm.JNBE(encoding.Uint8(stmtLen + jmpSize)),
			}
		}
	case *expr.IR_GT:
		c := condition.(*expr.IR_GT)
		result, err = c.EncodeWithoutSETE(ctx, reg)
		if shared.IsSignedInteger(c.Op1.ReturnType(ctx)) {
			instr = []lib.Instruction{
				asm.JNG(encoding.Uint8(stmtLen + jmpSize)),
			}
		} else {
			instr = []lib.Instruction{
				asm.JNA(encoding.Uint8(stmtLen + jmpSize)),
			}
		}
	case *expr.IR_GTE:
		c := condition.(*expr.IR_GTE)
		result, err = c.EncodeWithoutSETE(ctx, reg)
		if shared.IsSignedInteger(c.Op1.ReturnType(ctx)) {
			instr = []lib.Instruction{
				asm.JNGE(encoding.Uint8(stmtLen + jmpSize)),
			}
		} else {
			instr = []lib.Instruction{
				asm.JNAE(encoding.Uint8(stmtLen + jmpSize)),
			}
		}
	case *expr.IR_Not:
		result, err = condition.(*expr.IR_Not).EncodeWithoutSETE(ctx, reg)
		instr = []lib.Instruction{
			asm.JE(encoding.Uint8(stmtLen + jmpSize)),
		}
	case *expr.IR_Bool:
		result, err = condition.Encode(ctx, reg)
		instr = []lib.Instruction{
			asm.CMP_immediate(1, reg),
			asm.JNE(encoding.Uint8(stmtLen + jmpSize)),
		}
	default:
		return nil, fmt.Errorf("Unsupported condition %s", condition.String())
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
