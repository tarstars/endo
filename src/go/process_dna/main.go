package main

import (
	"fmt"
	"github.com/tarstars/endo/src/go/dna_processor"
	"io"
	"os"
)

func main() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	dna := dna_processor.NewSimpleDnaStorage(string(data))

	for {
		dna_processor.Step(dna)
	}
}
