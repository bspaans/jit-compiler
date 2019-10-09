package main

import (
	"github.com/bspaans/jit/ir"
	"github.com/bspaans/jit/ir/statements"
)

func main() {
	stmt := []statements.IR{}
	ir.CompileIR(stmt)
}
