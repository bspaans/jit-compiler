package ir

import "testing"

func Test_Parser_Happy(t *testing.T) {
	shouldParse := []string{
		"a123 = 1234",
		"a123 = 1234 + 123",
		"a123 = 1234 + z + b",
		"a123 = 1234 + z + b; z = a + 2",
		"a123 = 1234 + z + b; z = a + 2",
		"a123 = 1234 + z + b; z = a + 2; b = 3; return b",
		"if a + 3 { b = 3 } else { z = 300 }",
		"while a != 3 { a = a * 1 }",
		"while a != 3 { a = a * 1; b = 3.1415 }",
		"a123 = 1234 + b[3]",
		"a123 = 1234 + b[3 + z]",
		"a123 = []uint64{1,2,3,4,5}",
		"a123 = []uint64{1,2,3,4,5}[i]",
		"a123 = []float64{1.0,2.0,3.0,4.0,5.0}[i]",
		"a123 = []float64{1.0, 2.0, 3.0, 4.0, 5.0}[i]",
		"a123 = []float64{1.0, 2.0, 3.0, 4.0, 5.0 }[i]",
		"a123 = []float64{ 1.0, 2.0, 3.0, 4.0, 5.0 }[i]",
		"a123 = func(a uint64) uint64 { b = a * 2 }",
		"a123 = func(a uint64, c float64) uint64 { b = a * 2 }",

		`f = 2.5`,
		`f = freq * 0.00027210884353741496; currentIndex = 0`,
		`f[0] = 12`,
		"a123 = b() + c()",
		"a123 = b() + c(1)",
		"a123 = b() + c(1, 2)",
		"a123 = b() + c(1,2)",
		"a123 = b() + c(1,2,3 )",
		"a123 = b() + c(  1  ,  2,  3 )",
		"a123 = uint64(1)",
		"a123 = float64(1)",
		`a = 1; b = 2`,
		`a = 1 
         b = 2`,
		`a = struct { 
			Field uint64
			AnotherField uint64
		 }{
			 53,
			 53,
		 }`,
		`a = b.Field`,
		`a = (5 + 4) * 6`,
		`a = ([]uint64{1,2,3})[2]`,
	}
	for _, p := range shouldParse {
		_, err := ParseIR(p)
		if err != nil {
			t.Fatalf("Failed to parse %v: %v", p, err)
		}
	}
}

func Test_Parser_Sad(t *testing.T) {
	shouldParse := []string{
		"a123 = uint64(1, 2)",
		"a123 = float64(1, 2)",
	}
	for _, p := range shouldParse {
		_, err := ParseIR(p)
		if err == nil {
			t.Fatalf("Parsing of '%v' succeeded, but should have failed.", p)
		}
	}
}
