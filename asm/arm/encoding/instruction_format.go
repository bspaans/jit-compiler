/*

The Arm architecture is a Reduced Instruction Set Computer (RISC) architecture
with the following RISC architecture features:

* A large uniform register file.

* A load/store architecture, where data-processing operations only operate on
  register contents, not directly on memory contents.

* Simple addressing modes, with all load/store addresses determined from
  register contents and instruction fields only.

AArch64 is the 64-bit Execution state, meaning addresses are held in 64-bit registers, and instructions in the base instruction set can use 64-bit registers for their processing.

*/
package encoding

import "github.com/bspaans/jit-compiler/lib"

type InstructionFormat struct {
	LeadingBits uint8 // 3 bits
	Op0         uint8 // 4 bits
	Data        uint8 // Remainging 25 bits
}

func (i *InstructionFormat) Encode() lib.MachineCode {
	return []uint8{}

}
