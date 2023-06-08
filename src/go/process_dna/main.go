package main

import (
	"fmt"
	"github.com/tarstars/endo/src/go/dna_processor"
	"io"
	"os"
	"time"
)

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	data, err := io.ReadAll(os.Stdin)
	handleError(err)

	dna := dna_processor.NewSimpleDnaStorage(string(data))
	meter := 0
	startTime := time.Now()

	go func() {
		ticker := time.NewTicker(1 * time.Second) // every second
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				elapsed := time.Since(startTime)
				ips := float64(meter) / elapsed.Seconds() // iterations per second
				fmt.Fprintf(os.Stderr, "Iterations per second: %.2f\n", ips)
			}
		}
	}()

	for {
		err = dna_processor.Step(dna, meter, false)
		if err == dna_processor.Finish {
			break
		}
		meter += 1
	}
}
