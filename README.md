# jit-compiler

Compiles a simple intermediate representation into x86-64 assembly and then
machinecode, which can be executed in process.

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

## What's supported

Not an awful lot yet:

* MOV, MOVQ
* LEA
* PUSH and POP
* PUSHFQ
* ADD, SUB 
* ADDSD, SUBSD, MULSD and DIVSD
* INC and DEC
* CMP
* CVTSI2SS, CVTTSD2SI (convert int to and from float)
* SETE
* JMP and JMPE
* RET 
* SYSCALL

Immediate values, direct registers, displaced registers and RIP relative
addressing are mostly supported, but not every opcode will be able to handle
all. I'm looking into some code generation to make it easier to add more
complete coverage. 

Register allocation is really noddy and works until you run out of registers;
there is no allocating on the stack or heap yet.

In the higher level language the following constructs work:

* Unsigned 64bit integers
* 64bit floating point numbers
* Booleans
* Immutable byte arrays
* Assigning to variables
* If statements
* While loops
* Equals
* NOT logic expression
* Syscalls
* Return
* Int arithmetic `(+, -)`
* Float arithmetic `(+, -, *, /)`

I know, it's a bit much. Goodbye, Haskell.

## Next steps

* SIMD support 
* Ability to define and call functions
* Some sort of array (possibly only static size)
* Possibly a higher level language that compiles down into the IR
* A parser
* Try and use it from bleep
