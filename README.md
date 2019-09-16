# jit-compiler

A better name will come one day.

Compiles a simple intermediate representation into x86-64 assembly and then
machinecode, which can be executed in process (and later written out as a
standalone binary, maybe)

## Why Though?

The original intent behind this project was to be able to compile complicated
Sequencer and Synthesizer definitions down to machine code (see my `bleep`
project), but it's become a little bit more general purpose since, whilst still 
not achieving its original goal 👍 

## What's supported

Not an awful lot yet:

* MOV
* PUSH and POP
* PUSHFQ
* ADD
* INC and DEC
* CMP
* SETE
* JMP and JMPE
* RET 
* SYSCALL

Immediate values, direct registers and displaced registers are mostly
supported, but not every opcode will be able to handle all. I'm looking into
some code generation to make it easier to add more complete coverage. 

Register allocation is really noddy and works until you run out of registers;
there is no allocating on the stack or heap yet.

In the higher level language the following constructs work:

* Unsigned 64bit integers
* Booleans
* Assigning to variables
* If statements
* Return

I know, it's a bit much.

## Next steps

* Adding float support
* Adding more arithmetic 
* Ability to define and call functions
* Possibly a higher level language that compiles down into the IR
* A parser
* Try and use it from bleep
