package ir

import (
	"runtime"
	"testing"

	"github.com/bspaans/jit-compiler/ir/encoding/x86_64"
	. "github.com/bspaans/jit-compiler/ir/expr"
	. "github.com/bspaans/jit-compiler/ir/shared"
	. "github.com/bspaans/jit-compiler/ir/statements"
	"github.com/bspaans/jit-compiler/lib"
)

var (
	TargetArch = &x86_64.X86_64{}
	TargetABI  = x86_64.NewABI_AMDSystemV()
)

var ShouldRun = [][]IR{
	{
		NewIR_Assignment("a", NewIR_ByteArray([]uint8("test"))),
		NewIR_Return(NewIR_Variable("a")),
	},
	{
		NewIR_If(NewIR_Bool(true),
			NewIR_Assignment("f", NewIR_Uint64(53)),
			NewIR_Assignment("f", NewIR_Uint64(54)),
		),
		NewIR_Return(NewIR_Variable("f")),
	},
}

func Test_ShouldRun(t *testing.T) {
	for _, ir := range ShouldRun {
		debug := false
		b, err := Compile(TargetArch, TargetABI, ir, false)
		if err != nil {
			t.Fatal(err)
		}
		b.Execute(debug)
	}
}

func Test_ParseExecute_Happy(t *testing.T) {
	parseExecute_Happy(t)
}

func Benchmark_Parse(b *testing.B) {
	b.Run("No-SSA/No-Cache", func(b *testing.B) {
		benchParseCompileExecute(b, false, false)
	})
	b.Run("SSA/Cache", func(b *testing.B) {
		benchParseCompileExecute(b, true, true)
	})
	b.Run("SSA/No-Cache", func(b *testing.B) {
		benchParseCompileExecute(b, true, false)
	})
	b.Run("No-SSA/Cache", func(b *testing.B) {
		benchParseCompileExecute(b, false, true)
	})
}

func benchParseCompileExecute(b *testing.B, ssa, compile bool) {
	b.Helper()
	var compiled []*lib.CompiledCode
	if compile {
		compiled = make([]*lib.CompiledCode, 0, len(all53))
		for _, ir := range all53 {
			i, err := ParseIR(ir + "; return f")
			if err != nil {
				b.Fatal(err, "in", ir)
			}
			if ssa {
				i = i.SSA_Transform(NewSSA_Context())
			}
			mc, err := Compile(TargetArch, TargetABI, []IR{i}, false)
			if err != nil {
				b.Fatal(err, "in", ir)
			}
			compiled = append(compiled, mc.Compile())
		}
	}

	runtime.GC()
	runtime.GC()
	runtime.GC()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if compile {
			for i, fn := range compiled {
				value := fn.Execute()
				if value != int(53) {
					b.Fatal("Expecting 53 got", value, "in", all53[i], "\n")
				}
			}
		} else {
			for _, ir := range all53 {
				i, err := ParseIR(ir + "; return f")
				if err != nil {
					b.Fatal(err, "in", ir)
				}
				if ssa {
					i = i.SSA_Transform(NewSSA_Context())
				}
				mc, err := Compile(TargetArch, TargetABI, []IR{i}, false)
				if err != nil {
					b.Fatal(err, "in", ir)
				}
				value := mc.Execute(false)
				if value != int(53) {
					mc.Execute(true)
					b.Fatal("Expecting 53 got", value, "in", ir, "\n", mc)
				}
			}
		}
	}
}

func parseExecute_Happy(tb testing.TB) {
	tb.Helper()

	for _, ir := range all53 {
		i, err := ParseIR(ir + "; return f")
		if err != nil {
			tb.Fatal(err, "in", ir)
		}
		debug := false
		b, err := Compile(TargetArch, TargetABI, []IR{i}, false)
		if err != nil {
			if !debug {
				Compile(TargetArch, TargetABI, []IR{i}, false)
			}
			tb.Fatal(err, "in", ir)
		}
		value := b.Execute(false)
		if value != int(53) {
			if !debug {
				Compile(TargetArch, TargetABI, []IR{i}, false)
				b.Execute(false)
			}
			tb.Fatal("Expecting 53 got", value, "in", ir, "\n", b)
		}

		transformed := i.SSA_Transform(NewSSA_Context())
		b2, err := Compile(TargetArch, TargetABI, []IR{transformed}, false)
		if err != nil {
			tb.Fatal(err)
		}
		value = b2.Execute(false)
		if value != int(53) {
			if !debug {
				Compile(TargetArch, TargetABI, []IR{transformed}, false)
			}
			tb.Fatal("Expecting 53 got", value, "in", ir, " after SSA transform\n", transformed)
		}

	}
}

func Test_Execute_Result(t *testing.T) {
	units := [][]IR{
		{NewIR_Assignment("f", NewIR_Uint64(53))},
		{NewIR_If(NewIR_Bool(true),
			NewIR_Assignment("f", NewIR_Uint64(53)),
			NewIR_Assignment("f", NewIR_Uint64(54)),
		)},
		{
			NewIR_Assignment("f", NewIR_Uint64(3)),
			NewIR_If(NewIR_Bool(true),
				NewIR_Assignment("f", NewIR_Uint64(53)),
				NewIR_Assignment("f", NewIR_Uint64(54)),
			),
		},
		{
			NewIR_Assignment("f", NewIR_Uint64(53)),
			NewIR_Assignment("g", NewIR_Syscall(NewIR_Uint64(uint64(IR_Syscall_Linux_Write)), []IRExpression{NewIR_Uint64(1), NewIR_ByteArray([]uint8("hello world\n")), NewIR_Uint64(uint64(12))})),
		},
		{
			NewIR_Assignment("f", NewIR_Uint64(0)),
			NewIR_While(NewIR_Not(NewIR_Equals(NewIR_Variable("f"), NewIR_Uint64(53))),
				NewIR_Assignment("f", NewIR_Add(NewIR_Variable("f"), NewIR_Uint64(1))),
			),
		},
		{MustParseIR(`f = 0; while f != 53 { f = f + 1 }`)},
		{NewIR_Assignment("f", NewIR_Add(NewIR_Uint64(51), NewIR_Uint64(2)))},
		{NewIR_Assignment("f", NewIR_Cast(NewIR_Add(NewIR_Float64(51), NewIR_Float64(2)), TUint64))},
		{
			NewIR_Assignment("g", NewIR_Float64(53.343)),
			NewIR_Assignment("f", NewIR_Cast(NewIR_Variable("g"), TUint64)),
		},
		{
			NewIR_Assignment("g", NewIR_Float64(26.5)),
			NewIR_Assignment("h", NewIR_Float64(2.0)),
			NewIR_Assignment("i", NewIR_Mul(NewIR_Variable("g"), NewIR_Variable("h"))),
			NewIR_Assignment("f", NewIR_Cast(NewIR_Variable("i"), TUint64)),
		},
		{
			NewIR_Assignment("g", NewIR_Float64(106)),
			NewIR_Assignment("h", NewIR_Float64(2.0)),
			NewIR_Assignment("i", NewIR_Div(NewIR_Variable("g"), NewIR_Variable("h"))),
			NewIR_Assignment("f", NewIR_Cast(NewIR_Variable("i"), TUint64)),
		},
		{
			NewIR_Assignment("g", NewIR_Float64(55)),
			NewIR_Assignment("h", NewIR_Float64(2.0)),
			NewIR_Assignment("i", NewIR_Sub(NewIR_Variable("g"), NewIR_Variable("h"))),
			NewIR_Assignment("f", NewIR_Cast(NewIR_Variable("i"), TUint64)),
		},
		{
			NewIR_Assignment("a",
				NewIR_Function(&TFunction{TUint64, []Type{TUint64}, []string{"z"}},
					NewIR_Return(NewIR_Add(NewIR_Variable("z"), NewIR_Uint64(3))))),
			NewIR_Assignment("f", NewIR_Call("a", []IRExpression{NewIR_Uint64(50)})),
		},
		{
			NewIR_Assignment("a", NewIR_ByteArray([]uint8{50, 51, 52, 53})),
			NewIR_Assignment("f", NewIR_ArrayIndex(NewIR_Variable("a"), NewIR_Uint64(3))),
		},
		{
			NewIR_Assignment("b", NewIR_Struct(&TStruct{
				FieldTypes: []Type{TUint64},
				Fields:     []string{"first_field"},
			}, []IRExpression{NewIR_Uint64(53)})),
			NewIR_Assignment("f", NewIR_StructField(NewIR_Variable("b"), "first_field")),
		},
		{
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
		b, err := Compile(TargetArch, TargetABI, i, false)
		if err != nil {
			t.Fatal(err)
		}
		value := b.Execute(debug)
		if value != int(53) {
			t.Fatal("Expecting 53 got", value, "in", ir, "\n", b)
		}
	}
}

func Test_IR_Length(t *testing.T) {
	ctx := NewIRContext(TargetArch, TargetABI)
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
	ctx := NewIRContext(TargetArch, TargetABI)
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

var all53 = []string{

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

	// float64
	`f = uint64(53.0)`,
	`f = uint64(51.0 + 2.0)`,
	`f = uint64(55.0 - 2.0)`,
	`f = uint64(3.0 + 25.0 * 2.0)`,
	`f = uint64((2.0 * 25.0) + 3.0)`,
	`f = uint64(3.0 + (2.0 * 25.0))`,
	`f = uint64((100.0/ 2.0) + 3.0)`,
	`f = uint64(3.0 + (100.0 / 2.0))`,
	`h = 2.0; f = uint64((h * 25.0) + 3.0)`,
	`h = 25.0; f = uint64((2.0 * h) + 3.0)`,
	`h = 3.0; f = uint64((2.0 * 25.0) + h)`,
	`f = uint64(-53.0 * -1.0)`,
	`f = uint64(-53.0 / -1.0)`,

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
	// `g = []uint64{42,52,53}; g[0] = uint64(11) + g[0]; f = g[0]`,
	//`g = []uint64{42,52,53}; g[0] = uint64(1) + g[1]; f = g[0]`,
	// `g = []uint64{42,52,53}; g[0] = g[0] + uint64(11); f = g[0]`,
	//`g = []uint64{42,52,53}; g[1] = g[1] + uint64(1); f = g[1]`,

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
	`i = 0; while i < 53 { i = i + 1} ; f = i`,
	`i = 0; while i <= 52 { i = i + 1} ; f = i`,
	`i = 100; while i > 53 { i = i - 1} ; f = i`,
	`i = 100; while i >= 54 { i = i - 1} ; f = i`,
	`k = 2;j = 1; i = 0; while i != 53 { i = i + 1} ; f = i`,
	`k = 2;j = 1; i = 0; while i != 53 { i = i + 1} ; f = 53`,
	`f = 0; while f != 53 { f = f + 1 }`,
	`j = 1; i = 5; while i == 5 { j = j + 1; if j == 5 { i = 53 } else { i = 5 }}; f = i`,

	// if statements with int64
	`if 15 == 15 { f = 53 } else { f = 100 }`,
	`k = 21; j = 1; if 15 == 15 { f = 53 } else { f = 100 }`,
	`if 13 != 15 { f = 53 } else { f = 100 }`,
	`if 13 < 15 { f = 53 } else { f = 100 }`,
	`if 14 <= 15 { f = 53 } else { f = 100 }`,
	`if 15 <= 15 { f = 53 } else { f = 100 }`,
	`if 13 == 15 { f = 100 } else { f = 53 }`,
	`if 13 > 15 { f = 100 } else { f = 53 }`,
	`if 13 >= 15 { f = 100 } else { f = 53 }`,
	`if 16 > 15 { f = 53 } else { f = 100 }`,
	`if 15 >= 15 { f = 53 } else { f = 100 }`,
	`if (15 == 15) && (17 == 17) { f = 53 } else { f = 100 }`,
	`if (14 < 15) && (14 <= 17) { f = 53 } else { f = 100 }`,
	`if (16 > 15) && (19 >= 17) { f = 53 } else { f = 100 }`,
	`if (15 == 15) && (17 == 14) { f = 100 } else { f = 53 }`,
	`if (15 == 14) && (17 == 16) { f = 100 } else { f = 53 }`,
	`if (15 == 14) && (17 == 17) { f = 100 } else { f = 53 }`,
	`if (15 == 14) || (17 == 17) { f = 53 } else { f = 100 }`,
	`if (15 == 15) || (17 == 14) { f = 53 } else { f = 100 }`,
	`if (15 == 14) || (17 == 14) { f = 100} else { f = 53 }`,

	// boolean variables and if
	`b = true; if b { f = 53 } else { f = 100 }`,
	`b = !false; if b { f = 53 } else { f = 100 }`,
	`b = false; if !b { f = 53 } else { f = 100 }`,
	`b = true || true; if b { f = 53 } else { f = 100 }`,
	`b = true || false; if b { f = 53 } else { f = 100 }`,
	`b = false || true; if b { f = 53 } else { f = 100 }`,
	`b = false || false; if !b { f = 53 } else { f = 100 }`,
	`b = true && true; if b { f = 53 } else { f = 100 }`,
	`b = true && false; if !b { f = 53 } else { f = 100 }`,
	`b = false && true; if !b { f = 53 } else { f = 100 }`,
	`b = false && false; if !b { f = 53 } else { f = 100 }`,
	`b = 10 > 9; if b { f = 53 } else { f = 100 }`,
	`b = 10 >= 9; if b { f = 53 } else { f = 100 }`,
	`b = 10 < 9; if !b { f = 53 } else { f = 100 }`,
	`b = 10 <= 9; if !b { f = 53 } else { f = 100 }`,
	`b = int8(15) < int8(-1); if !b { f = 53 } else { f = 100 }`,
	`c = int8(127) <= int8(-127); if !c { f = 53 } else { f = 100 }`,
	`b = int8(15) < int8(-1); c = int8(127) <= int8(-127); if (!b) && (!c) { f = 53 } else { f = 100 }`,
	`b = int8(15) < int8(-1) ; c = !b ; d = int8(127) <= int8(-127) ; e = !d ; if c && e { f = 53 } else { f = 100 }`,

	// if statements with uint8
	`if uint8(13) < uint8(15) { f = 53 } else { f = 100 }`,
	`if uint8(15) <= uint8(15) { f = 53 } else { f = 100 }`,
	`if uint8(14) <= uint8(15) { f = 53 } else { f = 100 }`,
	`if (uint8(13) < uint8(15)) && (uint8(15) <= uint8(15)) { f = 53 } else { f = 100 }`,

	// if statements with uint16
	`if uint16(13) < uint16(15) { f = 53 } else { f = 100 }`,
	`if uint16(13) <= uint16(15) { f = 53 } else { f = 100 }`,
	`if uint16(15) <= uint16(15) { f = 53 } else { f = 100 }`,

	// if statements with uint32
	`if uint32(13) < uint32(15) { f = 53 } else { f = 100 }`,
	`if uint32(15) <= uint32(15) { f = 53 } else { f = 100 }`,
	`if uint32(13) <= uint32(15) { f = 53 } else { f = 100 }`,

	// if statements with int8
	`if int8(13) < int8(15) { f = 53 } else { f = 100 }`,
	`if int8(-1) < int8(15) { f = 53 } else { f = 100 }`,
	`if int8(-1) <= int8(15) { f = 53 } else { f = 100 }`,
	`if int8(15) <= int8(15) { f = 53 } else { f = 100 }`,
	`if (int8(-1) < int8(15)) && (int8(-125) <= int8(15)) { f = 53 } else { f = 100 }`,
	`if (int8(15) > int8(-1)) && (int8(127) >= int8(-127)) { f = 53 } else { f = 100 }`,
	`if (!(int8(15) < int8(-1))) && (!(int8(127) <= int8(-127))) { f = 53 } else { f = 100 }`,

	// if statements with int16
	`if int16(13) < int16(15) { f = 53 } else { f = 100 }`,
	`if int16(-1) < int16(15) { f = 53 } else { f = 100 }`,
	`if int16(-1) <= int16(15) { f = 53 } else { f = 100 }`,
	`if int16(15) <= int16(15) { f = 53 } else { f = 100 }`,
	// TODO: fix operator precedence
	`if (int16(-1) < int16(15)) && (int16(-127) <= int16(15)) { f = 53 } else { f = 100 }`,
	`if (int16(15) > int16(-1)) && (int16(127) >= int16(-127)) { f = 53 } else { f = 100 }`,
	`if (!(int16(15) < int16(-1))) && (!(int16(127) <= int16(-127))) { f = 53 } else { f = 100 }`,
	`if (!(int16(-2) > int16(-1))) && (!(int16(17) >= int16(18))) { f = 53 } else { f = 100 }`,

	// if statements with int32
	`if int32(13) < int32(15) { f = 53 } else { f = 100 }`,
	`if int32(-1) < int32(15) { f = 53 } else { f = 100 }`,
	`if int32(15) <= int32(15) { f = 53 } else { f = 100 }`,
	`if int32(-1) <= int32(15) { f = 53 } else { f = 100 }`,
	`if (int32(-1) < int32(15)) && (int32(-127) <= int32(15)) { f = 53 } else { f = 100 }`,
	`if (int32(15) > int32(-1)) && (int32(127) >= int32(-127)) { f = 53 } else { f = 100 }`,
	`if (!(int32(15) < int32(-1))) && (!(int32(127) <= int32(-127))) { f = 53 } else { f = 100 }`,

	// structs
	`b = struct{Field int64}{53}; f = b.Field`,
	`b = struct{Field int64
				Field2 int64}{51, 53}; f = b.Field2`,

	// functions
	`b = func(i uint64) uint64 { return i - uint64(2) }; f = b(55)`,
	`func b(i uint64) uint64 { return i - uint64(2)}; f = b(55)`,
}
