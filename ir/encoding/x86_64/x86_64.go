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

func (x *X86_64) EncodeExpression(expr IRExpression, ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {
	return encodeExpression(expr, ctx, target)
}

func (x *X86_64) EncodeStatement(stmt IR, ctx *IR_Context) ([]lib.Instruction, error) {
	return encodeStatement(stmt, ctx)
}

func (x *X86_64) EncodeDataSection(stmts []IR, ctx *IR_Context) (readonly, rw []uint8, err error) {
	for _, stmt := range stmts {
		readonly_, rw_, err := encodeDataSection(stmt, ctx, readonly, rw)
		if err != nil {
			return nil, nil, err
		}
		readonly = readonly_
		rw = rw_
	}
	return readonly, rw, nil
}

func encodeExpression(e IRExpression, ctx *IR_Context, target encoding.Operand) ([]lib.Instruction, error) {
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

func encodeDataSection(i IR, ctx *IR_Context, readonly, rw []uint8) ([]uint8, []uint8, error) {
	switch v := i.(type) {
	case *statements.IR_AndThen:
		readonly_, rw_, err := encodeDataSection(v.Stmt1, ctx, readonly, rw)
		if err != nil {
			return nil, nil, err
		}
		return encodeDataSection(v.Stmt2, ctx, readonly_, rw_)
	case *statements.IR_ArrayAssignment:
		readonly_, rw_, err := encodeExpressionForDataSection(v.Index, ctx, readonly, rw)
		if err != nil {
			return nil, nil, err
		}
		return encodeExpressionForDataSection(v.Expr, ctx, readonly_, rw_)
	case *statements.IR_Assignment:
		return encodeExpressionForDataSection(v.Expr, ctx, readonly, rw)
	case *statements.IR_FunctionDef:
		return encodeExpressionForDataSection(v.Expr, ctx, readonly, rw)
	case *statements.IR_If:
	case *statements.IR_Return:
	case *statements.IR_While:
	default:
		return nil, nil, fmt.Errorf("Unsupported '%s' statement in x86_64 data section encoder", i.String())
	}
	return readonly, rw, nil
}

func encodeExpressionForDataSection(i IRExpression, ctx *IR_Context, readonly, rw []uint8) ([]uint8, []uint8, error) {
	encodeOperators := func(op1, op2 IRExpression, readonly, rw []uint8) ([]uint8, []uint8, error) {
		readonly_, rw_, err := encodeExpressionForDataSection(op1, ctx, readonly, rw)
		if err != nil {
			return nil, nil, err
		}
		return encodeExpressionForDataSection(op2, ctx, readonly_, rw_)
	}
	switch v := i.(type) {
	case *expr.IR_ByteArray:
		v.Address = len(rw) + 2
		rw = append(rw, v.Value...)
		return readonly, rw, nil
	case *expr.IR_Add:
		return encodeOperators(v.Op1, v.Op2, readonly, rw)
	case *expr.IR_And:
		return encodeOperators(v.Op1, v.Op2, readonly, rw)
	case *expr.IR_ArrayIndex:
		return encodeOperators(v.Array, v.Index, readonly, rw)
	case *expr.IR_Call:
		for _, arg := range v.Args {
			readonly_, rw_, err := encodeExpressionForDataSection(arg, ctx, readonly, rw)
			if err != nil {
				return nil, nil, err
			}
			readonly = readonly_
			rw = rw_
		}
		return readonly, rw, nil
	case *expr.IR_Div:
		return encodeOperators(v.Op1, v.Op2, readonly, rw)
	case *expr.IR_Equals:
		return encodeOperators(v.Op1, v.Op2, readonly, rw)
	case *expr.IR_Function:
		readonly_, rw_, err := encodeDataSection(v.Body, ctx, readonly, rw)
		if err != nil {
			return nil, nil, err
		}
		readonly_, rw_, err = encode_IR_Function_for_DataSection(v, ctx, readonly_, rw_)

		if err != nil {
			return nil, nil, err
		}
		return readonly_, rw_, nil
	case *expr.IR_GT:
		return encodeOperators(v.Op1, v.Op2, readonly, rw)
	case *expr.IR_GTE:
		return encodeOperators(v.Op1, v.Op2, readonly, rw)
	case *expr.IR_LT:
		return encodeOperators(v.Op1, v.Op2, readonly, rw)
	case *expr.IR_LTE:
		return encodeOperators(v.Op1, v.Op2, readonly, rw)
	case *expr.IR_Mul:
		return encodeOperators(v.Op1, v.Op2, readonly, rw)
	case *expr.IR_Not:
		return encodeExpressionForDataSection(v.Op1, ctx, readonly, rw)
	case *expr.IR_Or:
		return encodeOperators(v.Op1, v.Op2, readonly, rw)
	case *expr.IR_StaticArray:
		return encode_IR_StaticArray_for_DataSection(v, ctx, readonly, rw)
	case *expr.IR_Struct:
		return encode_IR_Struct_for_DataSection(v, ctx, readonly, rw)
	case *expr.IR_StructField:
		return encodeExpressionForDataSection(v.Struct, ctx, readonly, rw)
	case *expr.IR_Syscall:
		for _, arg := range v.Args {
			readonly_, rw_, err := encodeExpressionForDataSection(arg, ctx, readonly, rw)
			if err != nil {
				return nil, nil, err
			}
			readonly = readonly_
			rw = rw_
		}
		return readonly, rw, nil
	case *expr.IR_Sub:
		return encodeOperators(v.Op1, v.Op2, readonly, rw)
	case *expr.IR_Bool, *expr.IR_Cast, *expr.IR_Variable, *expr.IR_Float64,
		*expr.IR_Uint8, *expr.IR_Uint16, *expr.IR_Uint32, *expr.IR_Uint64,
		*expr.IR_Int8, *expr.IR_Int16, *expr.IR_Int32, *expr.IR_Int64:
		return readonly, rw, nil
	default:
		return nil, nil, fmt.Errorf("Unsupported '%s' expr in x86_64 data section encoder", i.String())
	}
	return readonly, rw, nil
}
