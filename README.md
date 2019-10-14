# jit-compiler

This is a Golang library containing an x86-64 assembler (see 'asm/') and a
higher level intermediate representation that compiles down into x86-64 (see
'ir/').

## Why Though?

The original intent behind this project was to be able to compile complicated
Sequencer and Synthesizer definitions down to machine code (see my `bleep`
project), but it's become a little bit more general purpose since, whilst still 
not achieving its original goal 👍 

In `bleep`, as in many other synthesizers, we build complex sounds by combining
smaller building blocks (e.g. a sine wave, a delay filter, etc.) into bigger
blocks:

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

* MOV, MOVQ, MOVSD
* LEA
* PUSH and POP
* PUSHFQ
* ADD, SUB, DIV
* ADDSD, SUBSD, MULSD and DIVSD
* INC and DEC
* CMP
* CVTSI2SD, CVTTSD2SI (convert int to and from float)
* SETA, SETAE, SETB, SETBE, SETE, SETNE
* JMP and JMPE
* CALL and SYSCALL
* RET 

Immediate values, direct and indirect registers, displaced registers and RIP
relative addressing are supported.

In the higher level language the following constructs work:

* Unsigned 64bit integers
* 64bit floating point numbers
* Booleans
* Static size arrays
* Array indexing
* Assigning to variables
* If statements
* While loops
* Defining and calling functions
* Equals
* NOT logic expression
* Syscalls
* Return
* Int arithmetic `(+, -, *)`
* Float arithmetic `(+, -, *, /)`

Register allocation is really noddy and works until you run out of registers;
there is no allocating on the stack or heap yet; preserving registers across
calls and syscalls is supported however.

## Next steps

* SIMD support 
* Improve array support (SIMD, auto-vectorisation)
* Possibly a higher level language that compiles down into the IR
* A parser
* Try and use it from bleep
