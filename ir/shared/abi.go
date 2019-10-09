package shared

import (
	"github.com/bspaans/jit/asm"
)

type ABI interface {
	GetRegistersForArgs(args []Type) []*asm.Register
}

type ABI_AMDSystemV struct {
	intTargets   []*asm.Register
	floatTargets []*asm.Register
}

func NewABI_AMDSystemV() *ABI_AMDSystemV {
	return &ABI_AMDSystemV{
		intTargets:   []*asm.Register{asm.Rdi, asm.Rsi, asm.Rdx, asm.R10, asm.R8, asm.R9},
		floatTargets: []*asm.Register{asm.Xmm0, asm.Xmm1, asm.Xmm2, asm.Xmm3, asm.Xmm4, asm.Xmm5},
	}
}
func (a *ABI_AMDSystemV) GetRegistersForArgs(args []Type) []*asm.Register {
	intRegisterIx := 0
	floatRegisterIx := 0
	result := []*asm.Register{}

	var reg *asm.Register
	for _, arg := range args {
		if arg.Type() == T_Float64 {
			reg = a.floatTargets[floatRegisterIx]
			floatRegisterIx += 1
		} else {
			reg = a.intTargets[intRegisterIx]
			intRegisterIx += 1
		}
		result = append(result, reg)
	}
	return result
}

// returns instructions and clobbered registers
func PreserveRegisters(abi ABI, ctx *IR_Context, args []Type, returnType Type) (asm.Instructions, []*asm.Register) {
	regs := abi.GetRegistersForArgs(args)
	clobbered := []*asm.Register{}
	result := []asm.Instruction{}
	if returnType.Type() == T_Float64 {
		push := &asm.PUSH{asm.Xmm0}
		result = append(result, push)
		clobbered = append(clobbered, asm.Xmm0)
		ctx.AddInstruction(push)
	} else {
		if ctx.Registers[0] {
			push := &asm.PUSH{asm.Rax}
			result = append(result, push)
			clobbered = append(clobbered, asm.Rax)
			ctx.AddInstruction(push)
		}
	}
	var inUse bool
	for i, arg := range args {
		reg := regs[i]
		if arg.Type() == T_Float64 {
			inUse = ctx.FloatRegisters[reg.Register]
		} else {
			inUse = ctx.Registers[reg.Register]
		}
		if inUse {
			result = append(result, &asm.PUSH{reg})
			ctx.AddInstruction(&asm.PUSH{reg})
			clobbered = append(clobbered, reg)
		}
	}
	return result, clobbered

}

func RestoreRegisters(ctx *IR_Context, clobbered []*asm.Register) asm.Instructions {
	// Pop in reverse order
	result := []asm.Instruction{}
	for j := len(clobbered) - 1; j >= 0; j-- {
		reg := clobbered[j]
		result = append(result, &asm.POP{reg})
		ctx.AddInstruction(&asm.POP{reg})
	}
	return result
}
