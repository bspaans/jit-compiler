package asm

import "fmt"

type Float32 float32
type Float64 float64

/*

Float or double arguments in %xmm0 - %xmm7
Float or double results in %xmm0
%xmm8 - %xmm15 are temporaries (caller saved)
Registers are arranged in a stack (fadd vs. faddp)

addss <- adds single
addsd <- adds double
addsx, subsx, mulsx, divsx, ...

addsd %xmm0, %xmm1


Conversion:
- cvtsx2sx  source, dest
- cvttsx2sx souce, dest
x is either s, d or i
with i add an extra extension for l or q

convert long to double
cvtsi2sdq %rdi, %xmm0

convert float to int
cvtts2sil %xmm0, %eax


SIMD instructions work on pairs, because registers are actually 128 bits wide:

addpx source, dest
subpx, mulpx, divpx

Add two pairs of doubles:
addpd %xmm0, %xmm1

Multiply four pairs of floats:
mulps %xmm0, %xmm1
*/

func (f Float32) Type() Type {
	return T_Float32
}
func (f Float32) String() string {
	return fmt.Sprintf("$%f", f)
}
func (f Float32) Encode() []uint8 {
	return []uint8{}
}

func (f Float64) Type() Type {
	return T_Float64
}
func (f Float64) String() string {
	return fmt.Sprintf("$%f", f)
}
func (f Float64) Encode() []uint8 {
	return []uint8{}
}
