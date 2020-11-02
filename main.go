package main

import (
	"github.com/bspaans/jit-compiler/ir"
	"github.com/bspaans/jit-compiler/ir/shared"
)

func main() {
	stmt := []shared.IR{}
	ir.CompileIR(stmt)
}
