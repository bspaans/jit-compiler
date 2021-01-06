package aarch64

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm/aarch64/encoding"
	"github.com/bspaans/jit-compiler/ir/expr"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/ir/statements"
	"github.com/bspaans/jit-compiler/lib"
)

type AArch64 struct {
}

func (x *AArch64) EncodeExpression(expr IRExpression, ctx *IR_Context, target lib.Operand) ([]lib.Instruction, error) {
	return encodeExpression(expr, ctx, target)
}

func (x *AArch64) EncodeStatement(stmt IR, ctx *IR_Context) ([]lib.Instruction, error) {
	return encodeStatement(stmt, ctx)
}

func (x *AArch64) EncodeDataSection(stmts []IR, ctx *IR_Context) (*Segments, error) {
	segments := NewSegments()
	for _, stmt := range stmts {
		if err := encodeDataSection(stmt, ctx, segments); err != nil {
			return nil, err
		}
	}
	return segments, nil
}

func encodeExpression(e IRExpression, ctx *IR_Context, target lib.Operand) ([]lib.Instruction, error) {
	switch v := e.(type) {
	case *expr.IR_Variable:
		return encode_IR_Variable(v, ctx, target)
	case *expr.IR_Int64:
		return encode_IR_Int64(v, ctx, target)
	case *expr.IR_Add:
	case *expr.IR_And:
	case *expr.IR_ArrayIndex:
	case *expr.IR_Bool:
	case *expr.IR_ByteArray:
	case *expr.IR_Call:
	case *expr.IR_Cast:
	case *expr.IR_Div:
	case *expr.IR_Equals:
	case *expr.IR_Float64:
	case *expr.IR_Function:
	case *expr.IR_GT:
	case *expr.IR_GTE:
	case *expr.IR_Int8:
	case *expr.IR_Int16:
	case *expr.IR_Int32:
	case *expr.IR_LT:
	case *expr.IR_LTE:
	case *expr.IR_Mul:
	case *expr.IR_Not:
	case *expr.IR_Or:
	case *expr.IR_StaticArray:
	case *expr.IR_Struct:
	case *expr.IR_StructField:
	case *expr.IR_Syscall:
	case *expr.IR_Sub:
	case *expr.IR_Uint8:
	case *expr.IR_Uint16:
	case *expr.IR_Uint32:
	case *expr.IR_Uint64:
	}
	return nil, fmt.Errorf("Unsupported '%s' :: %s expression in x86_64 encoder", e.String(), e.Type().String())
}

func encodeStatement(stmt IR, ctx *IR_Context) ([]lib.Instruction, error) {
	switch v := stmt.(type) {
	case *statements.IR_Assignment:
		return encode_IR_Assignment(v, ctx)
	case *statements.IR_AndThen:
		return encode_IR_AndThen(v, ctx)
	case *statements.IR_ArrayAssignment:
	case *statements.IR_FunctionDef:
	case *statements.IR_If:
	case *statements.IR_Return:
	case *statements.IR_While:
	}
	return nil, fmt.Errorf("Unsupported '%s' statement in x86_64 encoder", stmt.String())
}

func encodeDataSection(i IR, ctx *IR_Context, segments *Segments) error {
	switch v := i.(type) {
	case *statements.IR_AndThen:
		if err := encodeDataSection(v.Stmt1, ctx, segments); err != nil {
			return err
		}
		return encodeDataSection(v.Stmt2, ctx, segments)
	case *statements.IR_ArrayAssignment:
		if err := encodeExpressionForDataSection(v.Index, ctx, segments); err != nil {
			return err
		}
		return encodeExpressionForDataSection(v.Expr, ctx, segments)
	case *statements.IR_Assignment:
		return encodeExpressionForDataSection(v.Expr, ctx, segments)
	case *statements.IR_FunctionDef:
		return encodeExpressionForDataSection(v.Expr, ctx, segments)
	case *statements.IR_If:
		// TODO
	case *statements.IR_Return:
		// TODO
	case *statements.IR_While:
	default:
		return fmt.Errorf("Unsupported '%s' statement in x86_64 data section encoder", i.String())
	}
	return nil
}

func encodeExpressionForDataSection(i IRExpression, ctx *IR_Context, segments *Segments) error {
	encodeOperators := func(op1, op2 IRExpression) error {
		if err := encodeExpressionForDataSection(op1, ctx, segments); err != nil {
			return err
		}
		return encodeExpressionForDataSection(op2, ctx, segments)
	}
	switch v := i.(type) {
	case *expr.IR_ByteArray:
		v.Address = segments.Add(ReadWrite, v.Value...)
		return nil
	case *expr.IR_Add:
		return encodeOperators(v.Op1, v.Op2)
	case *expr.IR_And:
		return encodeOperators(v.Op1, v.Op2)
	case *expr.IR_ArrayIndex:
		return encodeOperators(v.Array, v.Index)
	case *expr.IR_Call:
		for _, arg := range v.Args {
			if err := encodeExpressionForDataSection(arg, ctx, segments); err != nil {
				return err
			}
		}
		return nil
	case *expr.IR_Div:
		return encodeOperators(v.Op1, v.Op2)
	case *expr.IR_Equals:
		return encodeOperators(v.Op1, v.Op2)
	case *expr.IR_Function:
		if err := encodeDataSection(v.Body, ctx, segments); err != nil {
			return err
		}
		// TODO: return encode_IR_Function_for_DataSection(v, ctx, segments)
		return nil
	case *expr.IR_GT:
		return encodeOperators(v.Op1, v.Op2)
	case *expr.IR_GTE:
		return encodeOperators(v.Op1, v.Op2)
	case *expr.IR_LT:
		return encodeOperators(v.Op1, v.Op2)
	case *expr.IR_LTE:
		return encodeOperators(v.Op1, v.Op2)
	case *expr.IR_Mul:
		return encodeOperators(v.Op1, v.Op2)
	case *expr.IR_Not:
		return encodeExpressionForDataSection(v.Op1, ctx, segments)
	case *expr.IR_Or:
		return encodeOperators(v.Op1, v.Op2)
	case *expr.IR_StaticArray:
		// TODO: return encode_IR_StaticArray_for_DataSection(v, segments)
		return nil
	case *expr.IR_Struct:
		// TODO: return encode_IR_Struct_for_DataSection(v, ctx, segments)
		return nil
	case *expr.IR_StructField:
		return encodeExpressionForDataSection(v.Struct, ctx, segments)
	case *expr.IR_Syscall:
		for _, arg := range v.Args {
			if err := encodeExpressionForDataSection(arg, ctx, segments); err != nil {
				return err
			}
		}
		return nil
	case *expr.IR_Sub:
		return encodeOperators(v.Op1, v.Op2)
	case *expr.IR_Bool, *expr.IR_Cast, *expr.IR_Variable, *expr.IR_Float64,
		*expr.IR_Uint8, *expr.IR_Uint16, *expr.IR_Uint32, *expr.IR_Uint64,
		*expr.IR_Int8, *expr.IR_Int16, *expr.IR_Int32, *expr.IR_Int64:
		return nil
	default:
		return fmt.Errorf("Unsupported '%s' expr in x86_64 data section encoder", i.String())
	}
	return nil
}

func (x *AArch64) GetAllocator() Allocator {
	return NewAArch64_Allocator()
}

type AArch64_Allocator struct {
	Registers               []bool
	RegistersAllocated      uint8
	FloatRegisters          []bool
	FloatRegistersAllocated uint8
}

func NewAArch64_Allocator() *AArch64_Allocator {
	x := &AArch64_Allocator{
		Registers:               make([]bool, 32),
		RegistersAllocated:      0,
		FloatRegisters:          make([]bool, 32), // TODO?
		FloatRegistersAllocated: 0,
	}
	x.Registers[31] = true // stack pointer
	x.RegistersAllocated = 1
	return x
}

func (i *AArch64_Allocator) AllocateRegister(typ Type) lib.Operand {
	if typ == TFloat64 {
		return encoding.GetFloatingPointRegisterByIndex(i.allocateFloatRegister())
	}
	return encoding.Get64BitRegisterByIndex(i.allocateRegister()).ForOperandWidth(typ.Width())
}

func (i *AArch64_Allocator) DeallocateRegister(op lib.Operand) {
	reg, ok := op.(*encoding.Register)
	if !ok {
		return
	}
	if reg.Size == lib.QUADDOUBLE {
		i.deallocateFloatRegister(reg.Register)
		return
	}
	i.deallocateRegister(reg.Register)
}

func (i *AArch64_Allocator) allocateRegister() uint8 {
	if i.RegistersAllocated >= 16 {
		panic("Register allocation limit. Needs stack handling")
	}
	for j := 0; j < len(i.Registers); j++ {
		if !i.Registers[j] {
			i.Registers[j] = true
			i.RegistersAllocated += 1
			return uint8(j)
		}
	}
	panic("Register allocation limit reached with incorrect allocation counter. Needs stack handling")
}

func (i *AArch64_Allocator) deallocateRegister(reg uint8) {
	i.Registers[reg] = false
	i.RegistersAllocated -= 1
}

func (i *AArch64_Allocator) allocateFloatRegister() uint8 {
	if i.FloatRegistersAllocated >= 16 {
		panic("FloatRegister allocation limit. Needs stack handling")
	}
	for j := 0; j < len(i.FloatRegisters); j++ {
		if !i.FloatRegisters[j] {
			i.FloatRegisters[j] = true
			i.FloatRegistersAllocated += 1
			return uint8(j)
		}
	}
	panic("FloatRegister allocation limit reached with incorrect allocation counter. Needs stack handling")
}

func (i *AArch64_Allocator) deallocateFloatRegister(reg uint8) {
	i.FloatRegisters[reg] = false
	i.FloatRegistersAllocated -= 1
}

func (i *AArch64_Allocator) Copy() Allocator {
	regs := make([]bool, 16)
	floatRegs := make([]bool, 16)
	for j := 0; j < 16; j++ {
		regs[j] = i.Registers[j]
		floatRegs[j] = i.FloatRegisters[j]
	}
	return &AArch64_Allocator{
		Registers:               regs,
		RegistersAllocated:      i.RegistersAllocated,
		FloatRegisters:          floatRegs,
		FloatRegistersAllocated: i.FloatRegistersAllocated,
	}
}
