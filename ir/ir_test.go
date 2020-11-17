package ir

import (
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
		debug := false
		b, err := Compile(ir, debug)
		if err != nil {
			t.Fatal(err)
		}
		b.Execute(debug)
	}
}

func Test_ParseExecute_Happy(t *testing.T) {
	units := []string{

		// int64 (default)
		`f = 53`,
		`f = 51 + 2`,
		`f = 55 - 2`,
		`f = 3 + 25 * 2`,
		`f = (2 * 25) + 3`,
		`f = 3 + (2 * 25)`,
		`f = (100 / 2) + 3`,
		`f = 3 + (100 / 2)`,
		`h = 2; f = (h * 25) + 3`,
		`h = 25; f = (2 * h) + 3`,
		`h = 3; f = (2 * 25) + h`,
		`f = -53 * -1`,
		`f = -53 / -1`,

		// uint8
		`f = uint8(51) + uint8(2)`,
		`f = uint8(55) - uint8(2)`,
		`f = (uint8(2) * uint8(25)) + uint8(3)`,
		`f = uint8(3) + (uint8(2) * uint8(25))`,
		`f = uint8(3) + (uint8(100) / uint8(2))`,
		`f = (uint8(100) / uint8(2)) + uint8(3)`,

		// int8
		`f = int8(54) + int8(-1)`,
		`f = int8(106) / int8(2)`,
		`f = int8(-53) / int8(-1)`,
		`f = int8(-53) / int8(-1)`,
		`f = int8(-53) * int8(-1)`,
		`g= 255552;f = int8(-53) / int8(-1)`,

		// uint16
		`f = uint16(51) + uint16(2)`,
		`f = uint16(55) - uint16(2)`,
		`f = (uint16(2) * uint16(25)) + uint16(3)`,
		`f = uint16(3) + (uint16(2) * uint16(25))`,
		`f = uint16(3) + (uint16(100) / uint16(2))`,
		`f = (uint16(100) / uint16(2)) + uint16(3)`,

		// int16
		`f = int16(54) + int16(-1)`,
		`f = int16(-53) * int16(-1)`,
		`g= 25555212213;f = int16(-53) / int16(-1)`,

		// uint32
		`f = uint32(51) + uint32(2)`,
		`f = uint32(55) - uint32(2)`,
		`f = (uint32(2) * uint32(25)) + uint32(3)`,
		`f = uint32(3) + (uint32(2) * uint32(25))`,
		`f = uint32(3) + (uint32(100) / uint32(2))`,
		`f = (uint32(100) / uint32(2)) + uint32(3)`,

		// int32
		`f = int32(54) + int32(-1)`,
		`f = int32(-53) * int32(-1)`,
		`f = int32(-53) / int32(-1)`,

		// []uint64
		`f = []uint64{53}[0]`,
		`f = []uint64{42,52,53}[2]`,
		`f = ([]uint64{42,52,53})[2]`,
		`g = []uint64{42,52,53}; f = g[2]`,
		`g = []uint64{42,52,53}; h = []uint64{42,52,53}; f = g[2]`,
		`g = []uint64{42,52,53}; h = []uint64{42,52,53}; f = h[2]`,
		`g = []uint64{13,13}; i = 0; while i != 53 { i = i + 1} ; f = i`,
		`g = []uint64{13,13}; i = 0; while i != 53 { h = 2; i = i + 1} ; f = 53`,
		`g = []uint64{42,52,53}; f = g[2]`,
		`g = []uint64{42,52,53}; g[2] = g[2]; f = g[2]`,
		`g = []uint64{42,52,53}; f = g[0] + uint64(11)`,
		`g = []uint64{42,52,53}; g[0] = 53; f = g[0]`, // TODO this shouldn't actually work without auto casting
		`g = []uint64{42,52,53}; g[1] = 53; f = g[1]`,
		`g = []uint64{42,52,33}; g[2] = 53; f = g[2]`,
		`g = []uint64{42,52,53}; g[0] = 42 + 11; f = g[0]`,
		`g = []uint64{42,52,53}; g[0] = uint64(11) + g[0]; f = g[0]`,
		`g = []uint64{42,52,53}; g[0] = uint64(1) + g[1]; f = g[0]`,
		`g = []uint64{42,52,53}; g[0] = g[0] + uint64(11); f = g[0]`,
		`g = []uint64{42,52,53}; g[1] = g[1] + uint64(1); f = g[1]`,

		// []uint8
		`g = []uint8{42,52,53}; g[0] = uint8(42) + uint8(11); f = g[0]`,
		`g = []uint8{42,52,53}; g[0] = g[0] + uint8(11); f = g[0]`,
		`g = []uint8{53} ; f = uint64(g[0])`,
		`g = []uint8{52,53} ; f = uint64(g[1])`,
		`g = []uint8{51} ; g[0] = g[0] + uint8(2); f = uint64(g[0])`,
		`g = []uint8{51} ; f = uint64(2) + uint64(g[0])`,
		// TODO `g = []uint8{51} ; f = 2 + uint64(g[0])`,

		// []uint16
		`g = []uint16{42,52,53}; g[0] = uint16(42) + uint16(11); f = g[0]`,
		`g = []uint16{42,52,53}; g[0] = g[0] + uint16(11); f = g[0]`,
		`g = []uint16{53} ; f = uint64(g[0])`,
		`g = []uint16{52,53} ; f = uint64(g[1])`,
		`g = []uint16{51} ; g[0] = g[0] + uint16(2); f = uint64(g[0])`,
		`g = []uint16{51} ; f = uint64(2) + uint64(g[0])`,

		// []uint32
		`g = []uint32{42,52,53}; g[0] = uint32(42) + uint32(11); f = g[0]`,
		`g = []uint32{42,52,53}; g[0] = g[0] + uint32(11); f = g[0]`,
		`g = []uint32{53} ; f = uint64(g[0])`,
		`g = []uint32{52,53} ; f = uint64(g[1])`,
		`g = []uint32{51} ; g[0] = g[0] + uint32(2); f = uint64(g[0])`,
		`g = []uint32{51} ; f = uint64(2) + uint64(g[0])`,

		// []float64
		`g = []float64{53.0}; h = uint64(g[0]) ; f = h`,
		`g = []float64{53.0}; h =g[0] ; f = uint64(h)`,
		`g = []float64{51.0}; g[0] = g[0] + 2.0 ; f = uint64(g[0])`,

		// while loops with int64
		`i = 0; while i != 53 { i = i + 1} ; f = i`,
		`k = 2;j = 1; i = 0; while i != 53 { i = i + 1} ; f = i`,
		`k = 2;j = 1; i = 0; while i != 53 { i = i + 1} ; f = 53`,
		`f = 0; while f != 53 { f = f + 1 }`,
		`j = 1; i = 5; while i == 5 { j = j + 1; if j == 5 { i = 53 } else { i = 5 }}; f = i`,

		// if statements with int64
		`if 15 == 15 { f = 53 } else { f = 100 }`,
		`k = 21; j = 1; if 15 == 15 { f = 53 } else { f = 100 }`,
		`if 13 != 15 { f = 53 } else { f = 100 }`,
		`if 13 == 15 { f = 100 } else { f = 53 }`,

		// structs
		`b = struct{Field int64}{53}; f = b.Field`,
		`b = struct{Field int64
		            Field2 int64}{51, 53}; f = b.Field2`,

		// functions
		`b = func(i uint64) uint64 { return i - uint64(2) }; f = b(55)`,
		`func b(i uint64) uint64 { return i - uint64(2)}; f = b(55)`,
	}
	for _, ir := range units {
		i, err := ParseIR(ir + "; return f")
		if err != nil {
			t.Fatal(err, "in", ir)
		}
		debug := false
		b, err := Compile([]IR{i}, debug)
		if err != nil {
			if !debug {
				Compile([]IR{i}, true)
			}
			t.Fatal(err, "in", ir)
		}
		value := b.Execute(debug)
		if value != uint(53) {
			if !debug {
				Compile([]IR{i}, true)
				b.Execute(true)
			}
			t.Fatal("Expecting 53 got", value, "in", ir, "\n", b)
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
		debug := false
		b, err := Compile(i, debug)
		if err != nil {
			t.Fatal(err)
		}
		value := b.Execute(debug)
		if value != uint(53) {
			t.Fatal("Expecting 53 got", value, "in", ir, "\n", b)
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
	if l != 7 {
		t.Fatal("Expecting length 7 but got", l)
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
