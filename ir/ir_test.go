package ir

import (
	"fmt"
	"testing"
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
			NewIR_Assignment("g", NewIR_Syscall(uint(IR_Syscall_Linux_Write), []IRExpression{NewIR_Uint64(1), NewIR_ByteArray([]uint8("hello world\n")), NewIR_Uint64(uint64(12))})),
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
