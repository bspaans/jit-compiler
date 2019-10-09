package main

import (
	"github.com/bspaans/jit/ir"
	"github.com/bspaans/jit/ir/shared"
)

func main() {
	stmt := []shared.IR{}
	ir.CompileIR(stmt)
}
