package aarch64

/* The 64-bit ARM (AArch64) calling convention allocates the 31 general-purpose registers as:

    x31 (SP): Stack pointer or a zero register, depending on context.
    x30 (LR): Procedure link register, used to return from subroutines.
    x29 (FP): Frame pointer.
    x19 to x29: Callee-saved.
    x18 (PR): Platform register. Used for some operating-system-specific special purpose, or an additional caller-saved register.
    x16 (IP0) and x17 (IP1): Intra-Procedure-call scratch registers.
    x9 to x15: Local variables, caller saved.
    x8 (XR): Indirect return value address.
    x0 to x7: Argument values passed to and results returned from a subroutine.

The 32 floating-point registers are allocated as:

    v0 to v7: Argument values passed to and results returned from a subroutine.
    v8 to v15: callee-saved, but only the bottom 64 bits need to be preserved.
    v16 to v31: Local variables, caller saved.
*/

type ABI_AArch64 struct {
}
