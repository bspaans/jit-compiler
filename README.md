# jit-compiler

This is a Golang library containing an x86-64 assembler (see 'asm/') and a
higher level intermediate representation that compiles down into x86-64 (see
'ir/').

## Why Though?

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

The following x86-64 instructions are supported in the assembler. For a detailed 
overview see `asm/opcodes/opcodes.go` and `asm/opcodes/opcode_groups.go`:

* MOV, MOVQ, MOVSD, MOVZX (moving things in and out of registers and memory)
* LEA (loading the address of memory locations into a register)
* PUSH and POP (stack em up)
* ADD, SUB, MUL, DIV (arithmetic)
* ADDSD, SUBSD, MULSD and DIVSD (float arithmetic)
* INC and DEC
* SHL and SHR (shift to the left and right)
* XOR
* CMP
* CVTSI2SD, CVTTSD2SI (convert int to and from float)
* SETA, SETAE, SETB, SETBE, SETE, SETNE
* JMP, JNE, JMPE
* CALL and SYSCALL
* RET 
* PUSHFQ (push RFLAGS to the stack)
* Immediate values
* Addressing modes: direct and indirect registers, displaced registers, RIP relative, SIB

In the higher level language the following constructs work:

* Unsigned 8bit, 16bit, 32bit and 64bit integers
* 64bit floating point numbers
* Booleans
* Static size arrays
* Array indexing
* Assigning to variables
* Assigning to arrays
* If statements
* While loops
* Defining and calling functions
* Equals
* NOT logic expression
* Syscalls
* Return
* Int arithmetic `(+, -, *)`
* Float arithmetic `(+, -, *, /)`
* Parsing

Register allocation is really noddy and works until you run out of registers;
there is no allocating on the stack or heap yet; preserving registers across
calls and syscalls is supported however.

## Examples

### Creating machine code 

```golang
import (
    "github.com/bspaans/jit-compiler/asm"
    "github.com/bspaans/jit-compiler/asm/encoding"
    "github.com/bspaans/jit-compiler/lib"
)


...

result := lib.Instructions
result = result.Add(asm.MOV(encoding.Rax, encoding.Rcx))
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
