package asm

import (
	"fmt"

	"github.com/bspaans/jit/asm/encoding"
	"github.com/bspaans/jit/lib"
)

func init() {
	b, err := lib.CompileInstruction([]lib.Instruction{
		&MOV{encoding.Uint64(0), encoding.Rax},
		&MOV{encoding.Uint64(0), encoding.Rcx},
		&MOV{encoding.Uint64(0), encoding.Rdx},
		&MOV{encoding.Uint64(0), encoding.Rbx},
		&MOV{encoding.Uint64(0), encoding.Rbp},
		&MOV{encoding.Uint64(0), encoding.Rsi},
		&MOV{encoding.Uint64(0), encoding.Rdi},
		&MOV{encoding.Uint64(0), &encoding.DisplacedRegister{encoding.Rsp, 8}},
		&MOV{encoding.Uint64(0xffff), encoding.Rdi},
		&INC{encoding.Rax},
		&CMP{encoding.Rdi, encoding.Rax},
		&JNE{encoding.Uint8(0xf9)},
		&CMP{encoding.Rdi, encoding.Rax},
		&SETE{encoding.Al},
		&MOV{encoding.Uint64(123), encoding.Rcx},
		&ADD{encoding.Rcx, encoding.Rax},
		&ADD{encoding.Uint32(2), encoding.Rax},
		&PUSH{encoding.Rax},
		&POP{encoding.Rax},
		&PUSHFQ{},
		&POP{encoding.Rdx},
		&MOV{encoding.Rax, &encoding.DisplacedRegister{encoding.Rsp, 8}},
		&RET{},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(b)
	b.Execute()
}
