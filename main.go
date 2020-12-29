package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bspaans/jit-compiler/ir"
	"github.com/bspaans/jit-compiler/ir/shared"
)

func REPL() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		statements, err := ir.ParseIR(ir.Stdlib + text)
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

func CompileFiles() {
	source := ""
	for _, file := range os.Args[1:] {
		text, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}
		source += string(text) + "\n"
	}
	statements, err := ir.ParseIR(ir.Stdlib + source)
	if err != nil {
		panic(err)
	}

	debug := true
	statements = statements.SSA_Transform(shared.NewSSA_Context())
	if err := ir.CompileToBinary([]shared.IR{statements}, debug, "test.bin"); err != nil {
		panic(err)

	}
}

func main() {
	if len(os.Args) == 1 {
		REPL()
	} else {
		CompileFiles()
	}
}
