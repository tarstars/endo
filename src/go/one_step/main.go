package main

import (
	"fmt"
	"github.com/tarstars/endo/src/go/dna_processor"
	"io"
	"os"
)

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	flnmSource := os.Args[1]
	flnmDest := os.Args[2]

	fileSource, err := os.Open(flnmSource)
	handleError(err)
	defer fileSource.Close()

	fileDest, err := os.Create(flnmDest)
	handleError(err)
	defer fileDest.Close()

	data, err := io.ReadAll(fileSource)
	handleError(err)

	dna := dna_processor.NewSimpleDnaStorage(string(data))

	err = dna_processor.Step(dna, true)
	fmt.Print("error: ", err)

	_, err = fileDest.WriteString(dna.String())
	handleError(err)
}
