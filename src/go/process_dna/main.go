package main

import (
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
	//flnmSource := os.Args[1]
	//flnmDest := os.Args[2]
	//
	//fileSource, err := os.Open(flnmSource)
	//handleError(err)
	//defer fileSource.Close()
	//
	//fileDest, err := os.Create(flnmDest)
	//handleError(err)
	//defer fileDest.Close()
	//
	//data, err := io.ReadAll(fileSource)
	data, err := io.ReadAll(os.Stdin)
	handleError(err)

	dna := dna_processor.NewSimpleDnaStorage(string(data))
	meter := 0

	for {
		err = dna_processor.Step(dna, meter, false)
		if err == dna_processor.Finish {
			break
		}
		meter += 1
	}

	//_, err = fileDest.WriteString(dna.String())
	//handleError(err)
}
