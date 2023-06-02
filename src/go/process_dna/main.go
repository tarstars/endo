package main

import (
	"fmt"
	"github.com/tarstars/endo/src/go/dna_processor"
	"io"
	"log"
	"os"
)

func main() {
	fmt.Println("Start")
	file, err := os.Open("/home/tarstars/prj/endo/doc/endo.dna")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	dna := dna_processor.NewSimpleDnaStorage(string(data))

	for {
		dna_processor.Step(dna, true)
	}
}
