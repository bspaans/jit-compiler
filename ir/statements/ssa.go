package statements

import (
	. "github.com/bspaans/jit-compiler/ir/shared"
)

func SSA_Rewrites_to_IR(rw SSA_Rewrites) IR {
	if len(rw) == 0 {
		return nil
	}
	var v IR
	v = NewIR_Assignment(rw[0].Variable, rw[0].Expr)
	if len(rw) == 1 {
		return v
	}
	for _, rewrite := range rw[1:] {
		v = NewIR_AndThen(v, NewIR_Assignment(rewrite.Variable, rewrite.Expr))
	}
	return v
}
