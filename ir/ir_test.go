package ir

import (
	"fmt"
	"testing"

	. "github.com/bspaans/jit-compiler/ir/expr"
	. "github.com/bspaans/jit-compiler/ir/shared"
	. "github.com/bspaans/jit-compiler/ir/statements"
)

var ShouldRun = [][]IR{
	[]IR{NewIR_Assignment("a", NewIR_ByteArray([]uint8("test"))),
		NewIR_Return(NewIR_Variable("a")),
	},
	[]IR{NewIR_If(NewIR_Bool(true),
		NewIR_Assignment("f", NewIR_Uint64(53)),
		NewIR_Assignment("f", NewIR_Uint64(54)),
	),
		NewIR_Return(NewIR_Variable("f")),
	},
}

func Test_ShouldRun(t *testing.T) {
	for _, ir := range ShouldRun {
		b, err := Compile(ir)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(b)
		b.Execute()
	}
}

func Test_ParseExecute_Happy(t *testing.T) {
	units := []string{
		`f = 53`,
		`f = 51 + 2`,
		`f = 55 - 2`,
		`f = 3 + 25 * 2`,
		`f = (2 * 25) + 3`,
		`f = []uint64{42,52,53}[2]`,
		`f = ([]uint64{42,52,53})[2]`,
		`g = []uint64{42,52,53}; g[0] = g[0] + 11; f = g[0]`,
		`g = []float64{53.0}; h = uint64(g[0]) ; f = h`,
		`f = 0; while f != 53 { f = f + 1 }`,
		`if 15 == 15 { f = 53 } else { f = 100 }`,
		`b = struct{Field uint64}{53}; f = b.Field`,
		`b = struct{Field uint64
		            Field2 uint64}{51, 53}; f = b.Field2`,
		`b = func(i uint64) uint64 { return i - 2 }; f = b(55)`,
		`func b(i uint64) uint64 { return i - 2}; f = b(55)`,
	}
	for _, ir := range units {
		i, err := ParseIR(ir + "; return f")
		if err != nil {
			t.Fatal(err)
		}
		b, err := Compile([]IR{i})
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(b)
		value := b.Execute()
		if value != uint(53) {
			t.Fatal("Expecting 53 got", value, "in", ir)
		}
	}

}

func Test_Execute_Result(t *testing.T) {
	var units = [][]IR{
		[]IR{NewIR_Assignment("f", NewIR_Uint64(53))},
		[]IR{NewIR_If(NewIR_Bool(true),
			NewIR_Assignment("f", NewIR_Uint64(53)),
			NewIR_Assignment("f", NewIR_Uint64(54)),
		)},
		[]IR{NewIR_Assignment("f", NewIR_Uint64(3)),
			NewIR_If(NewIR_Bool(true),
				NewIR_Assignment("f", NewIR_Uint64(53)),
				NewIR_Assignment("f", NewIR_Uint64(54)),
			)},
		[]IR{NewIR_Assignment("f", NewIR_Uint64(53)),
			NewIR_Assignment("g", NewIR_Syscall(NewIR_Uint64(uint64(IR_Syscall_Linux_Write)), []IRExpression{NewIR_Uint64(1), NewIR_ByteArray([]uint8("hello world\n")), NewIR_Uint64(uint64(12))})),
		},
		[]IR{NewIR_Assignment("f", NewIR_Uint64(0)),
			NewIR_While(NewIR_Not(NewIR_Equals(NewIR_Variable("f"), NewIR_Uint64(53))),
				NewIR_Assignment("f", NewIR_Add(NewIR_Variable("f"), NewIR_Uint64(1))),
			),
		},
		[]IR{MustParseIR(`f = 0; while f != 53 { f = f + 1 }`)},
		[]IR{NewIR_Assignment("f", NewIR_Add(NewIR_Uint64(51), NewIR_Uint64(2)))},
		[]IR{NewIR_Assignment("f", NewIR_Cast(NewIR_Add(NewIR_Float64(51), NewIR_Float64(2)), TUint64))},
		[]IR{
			NewIR_Assignment("g", NewIR_Float64(53.343)),
			NewIR_Assignment("f", NewIR_Cast(NewIR_Variable("g"), TUint64)),
		},
		[]IR{
			NewIR_Assignment("g", NewIR_Float64(26.5)),
			NewIR_Assignment("h", NewIR_Float64(2.0)),
			NewIR_Assignment("i", NewIR_Mul(NewIR_Variable("g"), NewIR_Variable("h"))),
			NewIR_Assignment("f", NewIR_Cast(NewIR_Variable("i"), TUint64)),
		},
		[]IR{
			NewIR_Assignment("g", NewIR_Float64(106)),
			NewIR_Assignment("h", NewIR_Float64(2.0)),
			NewIR_Assignment("i", NewIR_Div(NewIR_Variable("g"), NewIR_Variable("h"))),
			NewIR_Assignment("f", NewIR_Cast(NewIR_Variable("i"), TUint64)),
		},
		[]IR{
			NewIR_Assignment("g", NewIR_Float64(55)),
			NewIR_Assignment("h", NewIR_Float64(2.0)),
			NewIR_Assignment("i", NewIR_Sub(NewIR_Variable("g"), NewIR_Variable("h"))),
			NewIR_Assignment("f", NewIR_Cast(NewIR_Variable("i"), TUint64)),
		},
		[]IR{
			NewIR_Assignment("a",
				NewIR_Function(&TFunction{TUint64, []Type{TUint64}, []string{"z"}},
					NewIR_Return(NewIR_Add(NewIR_Variable("z"), NewIR_Uint64(3))))),
			NewIR_Assignment("f", NewIR_Call("a", []IRExpression{NewIR_Uint64(50)})),
		},
		[]IR{
			NewIR_Assignment("a", NewIR_ByteArray([]uint8{50, 51, 52, 53})),
			NewIR_Assignment("f", NewIR_ArrayIndex(NewIR_Variable("a"), NewIR_Uint64(3))),
		},
		[]IR{
			NewIR_Assignment("b", NewIR_Struct(&TStruct{
				FieldTypes: []Type{TUint64},
				Fields:     []string{"first_field"},
			}, []IRExpression{NewIR_Uint64(53)})),
			NewIR_Assignment("f", NewIR_StructField(NewIR_Variable("b"), "first_field")),
		},
		[]IR{
			NewIR_Assignment("b", NewIR_Struct(&TStruct{
				FieldTypes: []Type{TUint64, TUint64},
				Fields:     []string{"first_field", "second_field"},
			}, []IRExpression{NewIR_Uint64(14), NewIR_Uint64(53)})),
			NewIR_Assignment("f", NewIR_StructField(NewIR_Variable("b"), "second_field")),
		},
	}
	for _, ir := range units {
		i := append(ir, NewIR_Return(NewIR_Variable("f")))
		b, err := Compile(i)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(b)
		value := b.Execute()
		if value != uint(53) {
			t.Fatal("Expecting 53 got", value, "in", ir)
		}
	}
}

func Test_IR_Length(t *testing.T) {

	ctx := NewIRContext()
	stmt := NewIR_Assignment("f", NewIR_Uint64(43))
	l, err := IR_Length(stmt, ctx)
	if err != nil {
		t.Fatal(err)
	}
	if l != 10 {
		t.Fatal("Expecting length 10 but got", l)
	}
}

func Test_IR_Length_does_not_affect_instruction_pointer(t *testing.T) {

	ctx := NewIRContext()
	stmt := NewIR_Assignment("f", NewIR_Uint64(43))
	rip := ctx.InstructionPointer
	_, err := IR_Length(stmt, ctx)
	if err != nil {
		t.Fatal(err)
	}
	if ctx.InstructionPointer != rip {
		t.Fatal("InstructionPointer changed")
	}
}
