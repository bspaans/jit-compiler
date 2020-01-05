package ir

import "fmt"

const Bleep = `

func Out(out []float64, n uint64, sampleRate uint64, pitch []float64) {
	stepSize = (2.0 * pitch) / sampleRate 
}
`

func init() {
	ir := MustParseIR(Bleep)
	fmt.Println(ir)
}
