package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/bspaans/jit-compiler/ir"
	"github.com/bspaans/jit-compiler/ir/shared"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		statements, err := ir.ParseIR(text)
		if err != nil {
			fmt.Println("Parse error: ", err.Error())
			continue

		}
		debug := true
		statements = statements.SSA_Transform(shared.NewSSA_Context())
		instr, err := ir.Compile([]shared.IR{statements}, debug)
		if err != nil {
			fmt.Println("Compile error: ", err.Error())
			continue

		}
		fmt.Println(instr.Execute(debug))
	}
}
