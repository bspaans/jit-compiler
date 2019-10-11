package shared

import (
	"github.com/bspaans/jit/asm"
)

type ABI interface {
	GetRegistersForArgs(args []Type) []*asm.Register
	ReturnTypeToOperand(ty Type) asm.Operand
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

func (a *ABI_AMDSystemV) ReturnTypeToOperand(arg Type) asm.Operand {
	if arg.Type() == T_Float64 {
		return asm.Xmm0
	}
	return asm.Rax
}

// returns instructions and clobbered registers
func PreserveRegisters(ctx *IR_Context, argTypes []Type, returnType Type) (asm.Instructions, map[asm.Operand]asm.Operand, []asm.Operand) {
	regs := ctx.ABI.GetRegistersForArgs(argTypes)
	returnOp := ctx.ABI.ReturnTypeToOperand(returnType)
	clobbered := []asm.Operand{}
	result := []asm.Instruction{}
	mapping := map[asm.Operand]asm.Operand{}
	push := &asm.PUSH{returnOp}
	result = append(result, push)
	clobbered = append(clobbered, returnOp)
	ctx.AddInstruction(push)
	var inUse bool
	for i, arg := range argTypes {
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
	for i, arg := range argTypes {
		reg := regs[i]
		if arg.Type() == T_Float64 {
			inUse = ctx.FloatRegisters[reg.Register]
		} else {
			inUse = ctx.Registers[reg.Register]
		}
		if inUse {
			mapping[reg] = &asm.DisplacedRegister{asm.Rsp, uint8((len(clobbered) - 2 - i) * 4)}
		}
	}
	return result, mapping, clobbered
}

func ABI_Call_Setup(ctx *IR_Context, args []IRExpression, returnType Type) (asm.Instructions, map[asm.Operand]asm.Operand, []asm.Operand, error) {
	argTypes := make([]Type, len(args))
	for i, arg := range args {
		argTypes[i] = arg.ReturnType(ctx)
	}
	result, mapping, clobbered := PreserveRegisters(ctx, argTypes, returnType)
	regs := ctx.ABI.GetRegistersForArgs(argTypes)

	ctx_ := ctx.Copy()
	for variable, location := range ctx_.VariableMap {
		if newLocation, found := mapping[location]; found {
			ctx_.VariableMap[variable] = newLocation
		}
	}
	for _, reg := range regs {
		if reg.Size == asm.QUADDOUBLE {
			ctx_.FloatRegisters[reg.Register] = true
			ctx_.FloatRegistersAllocated += 1
		} else {
			if !ctx_.Registers[reg.Register] {
				ctx_.Registers[reg.Register] = true
				ctx_.RegistersAllocated += 1
			}
		}
	}

	for i, arg := range args {
		instr, err := arg.Encode(ctx_, regs[i])
		if err != nil {
			return nil, nil, nil, err
		}
		ctx.AddInstructions(instr)
		result = result.Add(instr)
	}

	return result, mapping, clobbered, nil

}

func RestoreRegisters(ctx *IR_Context, clobbered []asm.Operand) asm.Instructions {
	// Pop in reverse order
	result := []asm.Instruction{}
	for j := len(clobbered) - 1; j >= 0; j-- {
		reg := clobbered[j]
		result = append(result, &asm.POP{reg})
		ctx.AddInstruction(&asm.POP{reg})
	}
	return result
}
