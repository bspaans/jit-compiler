# jit-compiler

[![Documentation](https://godoc.org/github.com/bspaans/jit-compiler?status.svg)](https://godoc.org/github.com/bspaans/jit-compiler) 


This is a Golang library containing an x86-64 assembler (see 'asm/') and a
higher level intermediate representation that compiles down into x86-64 (see
'ir/').

## Motivation

The original intent behind this project was to be able to compile complicated
Sequencer and Synthesizer definitions down to machine code (see my [bleep](https://github.com/bspaans/bleep))
project), but it's become a little bit more general purpose since, whilst still 
not achieving its original goal 👍 There's a very, very early prototype [here](https://github.com/bspaans/bleep-jit), 
but it's not doing much yet.


In `bleep`, as in many other synthesizers, we build complex sounds by combining
smaller building blocks (e.g. a sine wave, a delay filter, etc.) into bigger
instruments:

```
          +---- sine wave
         /
  delay <
         \
          +---- sqaure wave

```

We end up with a tree of sub-synthesizers and filters that together form the
final sound...

...This is nice, but can also be computationally expensive. Especially when 
multiple separate synthesizers are playing at the same time.

One of the reasons it is expensive is because the code is jumping around from
block to block, basically interpreting the tree. Wouldn't it be nice if we
could compile it all down into a single function on the fly? Maybe. This is a
slightly unnecessary experiment to find out, whilst at the same time learning
something about x86-64 and JIT compilation.

## What is JIT compilation?

Just In Time compilation is a method to convert, at runtime,  the execution of
datastructures into machine code. This can be a a lot faster than interpreting
the datastructures, as you are dealing directly with the processor and can
apply optimisations that aren't usually possible in the source language. It
does however mean you have to have some way to convert everything into binary;
hence projects like these.

## What's supported

### Assembler

The following x86-64 instructions are supported in the assembler. For a detailed 
overview see [`asm/x86_64/opcodes/`](https://github.com/bspaans/jit-compiler/tree/master/asm/x86_64/opcodes):

* MOV, MOVQ, MOVSD, MOVSX, MOVZX (moving things in and out of registers and memory)
* LEA (loading the address of memory locations into a register)
* PUSH and POP (stack em up)
* ADD, SUB, MUL, DIV, IMUL, IDIV (arithmetic)
* ADDSD, SUBSD, MULSD and DIVSD (float arithmetic)
* INC and DEC
* SHL and SHR (shift to the left and right)
* AND, OR and XOR (logic operations)
* CMP (compare numbers)
* CBW, CWD, CDQ, CQO (sign extend %al, %ax, %eax and %rax)
* CVTSI2SD, CVTTSD2SI (convert int to and from float)
* SETA, SETAE, SETB, SETBE, SETE, SETL, SETLE, SETG, SETGE, SETNE
* JMP, JA, JAE, JB, JBE, JE, JG, JGE, JL, JLE, JNA, JNAE, JNB, JNBE, JNE, JNG, JNGE, JNL, JNLE (jumps and conditional jumps)
* CALL and SYSCALL
* RET 
* PUSHFQ (push RFLAGS to the stack)
* Immediate values
* Addressing modes: direct and indirect registers, displaced registers, RIP relative, SIB

### Higher Level Language

The higher level language is kind of like a very stripped down Go/C like
language that makes it easier to generate code. It currently supports:

#### Data Types

* Unsigned 8bit, 16bit, 32bit and 64bit integers
* Signed 8bit, 16bit, 32bit and 64bit integers
* 64bit floating point numbers
* Booleans
* Static size arrays
* Structs 

#### Expressions

* Signed and unsigned integer arithmetic `(+, -, *, /)`
* Signed and unsigned integer comparisons `(==, !=, <, <=, >, >=)`
* Float arithmetic `(+, -, *, /)`
* Logic expressions `(&&, ||, !)`
* Array indexing
* Function calls
* Syscalls
* Casting types
* Equality testing
* Struct field indexing

#### Statements

* Assigning to variables
* Assigning to arrays
* If statements
* While loops
* Function definitions
* Return

#### Register allocation

Register allocation is really simple and works until you run out of registers;
there is no allocating on the stack or heap yet; preserving registers across
calls and syscalls is supported however.

## Examples

### Creating machine code 

```golang
import (
    "github.com/bspaans/jit-compiler/asm/x86_64"
    "github.com/bspaans/jit-compiler/asm/x86_64/encoding"
    "github.com/bspaans/jit-compiler/lib"
)


...

result := lib.Instructions
result = result.Add(x86_64.MOV(encoding.Rax, encoding.Rcx))
machineCode, err := result.Encode()
if err != nil {
    panic(err)
}
machineCode.Execute()

```

### Using the Intermediate Representation

```golang
import (
    "github.com/bspaans/jit-compiler/ir"
)

var code = `prev = 1; current = 1;
while current != 13 {
  tmp = current
  current = current + prev
  prev = tmp
}
return current
`

machineCode, err := ir.Compile(ir.MustParseIR(code))
if err != nil {
    panic(err)
}
fmt.Println(machineCode.Execute())
```

## Next steps

* SIMD support 
* Improve array support (SIMD, auto-vectorisation)
* Possibly a higher level language that compiles down into the IR
* Possible a WASM or ARM backend
* Test output against other assemblers
* Try and use it from bleep ([prototype](https://github.com/bspaans/bleep-jit))

## Contributing

Contributions are always welcome, but if you want to introduce a breaking
change please raise an issue first to discuss. For small additions and bug
fixes feel free to just create a PR.

## License

This package is licensed under a MIT License:

```
Copyright 2020, Bart Spaans

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

```
