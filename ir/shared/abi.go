package shared

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm"
	"github.com/bspaans/jit-compiler/asm/encoding"
	"github.com/bspaans/jit-compiler/lib"
)

type ABI interface {
	GetRegistersForArgs(args []Type) []*encoding.Register
	ReturnTypeToOperand(ty Type) encoding.Operand
}

type ABI_AMDSystemV struct {
	intTargets   []*encoding.Register
	floatTargets []*encoding.Register
}

func NewABI_AMDSystemV() *ABI_AMDSystemV {
	return &ABI_AMDSystemV{
		intTargets:   []*encoding.Register{encoding.Rdi, encoding.Rsi, encoding.Rdx, encoding.R10, encoding.R8, encoding.R9},
		floatTargets: []*encoding.Register{encoding.Xmm0, encoding.Xmm1, encoding.Xmm2, encoding.Xmm3, encoding.Xmm4, encoding.Xmm5},
	}
}
func (a *ABI_AMDSystemV) GetRegistersForArgs(args []Type) []*encoding.Register {
	intRegisterIx := 0
	floatRegisterIx := 0
	result := []*encoding.Register{}

	var reg *encoding.Register
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

func (a *ABI_AMDSystemV) ReturnTypeToOperand(arg Type) encoding.Operand {
	if arg.Type() == T_Float64 {
		return encoding.Xmm0
	}
	return encoding.Rax
}

// returns instructions and clobbered registers
func PreserveRegisters(ctx *IR_Context, argTypes []Type, returnType Type) (lib.Instructions, map[encoding.Operand]encoding.Operand, []encoding.Operand) {
	clobbered := []encoding.Operand{}
	result := []lib.Instruction{}
	mapping := map[encoding.Operand]encoding.Operand{}

	// push the return register; TODO: check if in use?
	returnOp := ctx.ABI.ReturnTypeToOperand(returnType)
	push := asm.PUSH(returnOp)
	result = append(result, push)
	clobbered = append(clobbered, returnOp)
	ctx.AddInstruction(push)

	// Push registers that are already in use
	regs := ctx.ABI.GetRegistersForArgs(argTypes)
	var inUse bool
	for i, arg := range argTypes {
		reg := regs[i]
		if arg.Type() == T_Float64 {
			inUse = ctx.FloatRegisters[reg.Register]
		} else {
			inUse = ctx.Registers[reg.Register]
		}
		if inUse {
			result = append(result, asm.PUSH(reg))
			ctx.AddInstruction(asm.PUSH(reg))
			clobbered = append(clobbered, reg)
		}
	}
	// Build the register -> location on the stack mapping
	mappedClobbered := 1 // set to 1, to account for the return op TODO
	for i, arg := range argTypes {
		reg := regs[i]
		if arg.Type() == T_Float64 {
			inUse = ctx.FloatRegisters[reg.Register]
		} else {
			inUse = ctx.Registers[reg.Register]
		}
		if inUse {
			offset := (len(clobbered) - mappedClobbered) * int(arg.Width())
			mapping[reg] = &encoding.DisplacedRegister{encoding.Rsp, uint8(offset)}
			mappedClobbered += 1
		}
	}
	return result, mapping, clobbered
}

func ABI_Call_Setup(ctx *IR_Context, args []IRExpression, returnType Type) (lib.Instructions, map[encoding.Operand]encoding.Operand, []encoding.Operand, error) {
	argTypes := make([]Type, len(args))
	for i, arg := range args {
		argTypes[i] = arg.ReturnType(ctx)
		if argTypes[i] == nil {
			return nil, nil, nil, fmt.Errorf("Unknown type for value: %s", arg)
		}
	}
	result, mapping, clobbered := PreserveRegisters(ctx, argTypes, returnType)
	regs := ctx.ABI.GetRegistersForArgs(argTypes)

	ctx_ := ctx.Copy()
	for _, reg := range regs {
		if reg.Size == lib.QUADDOUBLE {
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
	for variable, location := range ctx_.VariableMap {
		if newLocation, found := mapping[location]; found {
			ctx_.VariableMap[variable] = newLocation
		}
	}

	return result, mapping, clobbered, nil

}

func RestoreRegisters(ctx *IR_Context, clobbered []encoding.Operand) lib.Instructions {
	// Pop in reverse order
	result := []lib.Instruction{}
	for j := len(clobbered) - 1; j >= 0; j-- {
		reg := clobbered[j]
		result = append(result, asm.POP(reg))
		ctx.AddInstruction(asm.POP(reg))
	}
	return result
}
