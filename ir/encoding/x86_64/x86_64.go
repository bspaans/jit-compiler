package x86_64

import (
	"fmt"

	"github.com/bspaans/jit-compiler/asm/x86_64/encoding"
	"github.com/bspaans/jit-compiler/ir/expr"
	. "github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/ir/statements"
	"github.com/bspaans/jit-compiler/lib"
)

type X86_64 struct {
}

func (x *X86_64) EncodeExpression(expr IRExpression, ctx *IR_Context, target lib.Operand) ([]lib.Instruction, error) {
	return encodeExpression(expr, ctx, target)
}

func (x *X86_64) EncodeStatement(stmt IR, ctx *IR_Context) ([]lib.Instruction, error) {
	return encodeStatement(stmt, ctx)
}

func (x *X86_64) EncodeDataSection(stmts []IR, ctx *IR_Context) (*Segments, error) {
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
	case *expr.IR_Add:
		return encode_IR_Add(v, ctx, target)
	case *expr.IR_And:
		return encode_IR_And(v, ctx, target)
	case *expr.IR_ArrayIndex:
		return encode_IR_ArrayIndex(v, ctx, target)
	case *expr.IR_Bool:
		return encode_IR_Bool(v, ctx, target)
	case *expr.IR_ByteArray:
		return encode_IR_ByteArray(v, ctx, target)
	case *expr.IR_Call:
		return encode_IR_Call(v, ctx, target)
	case *expr.IR_Cast:
		return encode_IR_Cast(v, ctx, target)
	case *expr.IR_Div:
		return encode_IR_Div(v, ctx, target)
	case *expr.IR_Equals:
		return encode_IR_Equals(v, ctx, target, true)
	case *expr.IR_Float64:
		return encode_IR_Float64(v, ctx, target)
	case *expr.IR_Function:
		return encode_IR_Function(v, ctx, target)
	case *expr.IR_GT:
		return encode_IR_GT(v, ctx, target, true)
	case *expr.IR_GTE:
		return encode_IR_GTE(v, ctx, target, true)
	case *expr.IR_Int8:
		return encode_IR_Int8(v, ctx, target)
	case *expr.IR_Int16:
		return encode_IR_Int16(v, ctx, target)
	case *expr.IR_Int32:
		return encode_IR_Int32(v, ctx, target)
	case *expr.IR_Int64:
		return encode_IR_Int64(v, ctx, target)
	case *expr.IR_LT:
		return encode_IR_LT(v, ctx, target, true)
	case *expr.IR_LTE:
		return encode_IR_LTE(v, ctx, target, true)
	case *expr.IR_Mul:
		return encode_IR_Mul(v, ctx, target)
	case *expr.IR_Not:
		return encode_IR_Not(v, ctx, target, true)
	case *expr.IR_Or:
		return encode_IR_Or(v, ctx, target)
	case *expr.IR_StaticArray:
		return encode_IR_StaticArray(v, ctx, target)
	case *expr.IR_Struct:
		return encode_IR_Struct(v, ctx, target)
	case *expr.IR_StructField:
		return encode_IR_StructField(v, ctx, target)
	case *expr.IR_Syscall:
		return encode_IR_Syscall(v, ctx, target)
	case *expr.IR_Sub:
		return encode_IR_Sub(v, ctx, target)
	case *expr.IR_Uint8:
		return encode_IR_Uint8(v, ctx, target)
	case *expr.IR_Uint16:
		return encode_IR_Uint16(v, ctx, target)
	case *expr.IR_Uint32:
		return encode_IR_Uint32(v, ctx, target)
	case *expr.IR_Uint64:
		return encode_IR_Uint64(v, ctx, target)
	case *expr.IR_Variable:
		return encode_IR_Variable(v, ctx, target)
	default:
		return nil, fmt.Errorf("Unsupported '%s' expression in x86_64 encoder", e.String())
	}
}

func encodeStatement(stmt IR, ctx *IR_Context) ([]lib.Instruction, error) {
	switch v := stmt.(type) {
	case *statements.IR_AndThen:
		return encode_IR_AndThen(v, ctx)
	case *statements.IR_ArrayAssignment:
		return encode_IR_ArrayAssignment(v, ctx)
	case *statements.IR_Assignment:
		return encode_IR_Assignment(v, ctx)
	case *statements.IR_FunctionDef:
		return encode_IR_FunctionDef(v, ctx)
	case *statements.IR_If:
		return encode_IR_If(v, ctx)
	case *statements.IR_Return:
		return encode_IR_Return(v, ctx)
	case *statements.IR_While:
		return encode_IR_While(v, ctx)
	default:
		return nil, fmt.Errorf("Unsupported '%s' statement in x86_64 encoder", stmt.String())
	}
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
		return encode_IR_Function_for_DataSection(v, ctx, segments)
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
		return encode_IR_StaticArray_for_DataSection(v, segments)
	case *expr.IR_Struct:
		return encode_IR_Struct_for_DataSection(v, ctx, segments)
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

func (x *X86_64) GetAllocator() Allocator {
	return NewX86_64_Allocator()
}

type X86_64_Allocator struct {
	Registers               []bool
	RegistersAllocated      uint8
	FloatRegisters          []bool
	FloatRegistersAllocated uint8
}

func NewX86_64_Allocator() *X86_64_Allocator {
	x := &X86_64_Allocator{
		Registers:               make([]bool, 16),
		RegistersAllocated:      0,
		FloatRegisters:          make([]bool, 16),
		FloatRegistersAllocated: 0,
	}
	// Always allocate the stack and frame pointer so that they don't
	// get overwritten. We could be smarter here, but meh.
	x.Registers[4] = true // stack pointer
	x.Registers[5] = true // frame pointer
	x.RegistersAllocated = 2
	return x
}

func (i *X86_64_Allocator) AllocateRegister(typ Type) lib.Operand {
	if typ == TFloat64 {
		return encoding.GetFloatingPointRegisterByIndex(i.allocateFloatRegister())
	}
	return encoding.Get64BitRegisterByIndex(i.allocateRegister()).ForOperandWidth(typ.Width())
}

func (i *X86_64_Allocator) DeallocateRegister(op lib.Operand) {
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

func (i *X86_64_Allocator) allocateRegister() uint8 {
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

func (i *X86_64_Allocator) deallocateRegister(reg uint8) {
	i.Registers[reg] = false
	i.RegistersAllocated -= 1
}

func (i *X86_64_Allocator) allocateFloatRegister() uint8 {
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

func (i *X86_64_Allocator) deallocateFloatRegister(reg uint8) {
	i.FloatRegisters[reg] = false
	i.FloatRegistersAllocated -= 1
}

func (i *X86_64_Allocator) Copy() Allocator {
	regs := make([]bool, 16)
	floatRegs := make([]bool, 16)
	for j := 0; j < 16; j++ {
		regs[j] = i.Registers[j]
		floatRegs[j] = i.FloatRegisters[j]
	}
	return &X86_64_Allocator{
		Registers:               regs,
		RegistersAllocated:      i.RegistersAllocated,
		FloatRegisters:          floatRegs,
		FloatRegistersAllocated: i.FloatRegistersAllocated,
	}
}
