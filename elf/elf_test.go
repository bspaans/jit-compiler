package elf

import (
	"fmt"
	"testing"
)

func Test_ELF(t *testing.T) {

	fmt.Println(NewELFHeader().Encode())
}
