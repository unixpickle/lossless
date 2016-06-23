package main

import (
	"fmt"
	"os"

	"github.com/unixpickle/lossless"
)

func main() {
	if len(os.Args) < 4 {
		dieUsage()
	}

	subCmd := os.Args[1]
	model := os.Args[2]

	predictor := lossless.GetPredictor(model)
	if predictor == nil {
		fmt.Fprintln(os.Stderr, "invalid model:", model)
		dieUsage()
	}

	input, err := os.Open(os.Args[3])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer input.Close()

	if len(os.Args) == 4 {
		if subCmd == "crossentropy" {
			res, err := lossless.CrossEntropy(predictor, input)
			fmt.Println("Cross entropy:", res)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error during computation:", err)
				os.Exit(1)
			}
			return
		}
		dieUsage()
	}

	output, err := os.Create(os.Args[4])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer output.Close()

	if subCmd == "compress" {
		err = lossless.Compress(predictor, input, output)
	} else if subCmd == "decompress" {
		err = lossless.Decompress(predictor, input, output)
	} else {
		fmt.Fprintln(os.Stderr, "invalid sub-command:", subCmd)
		dieUsage()
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func dieUsage() {
	fmt.Fprintln(os.Stderr, "Usage: cmd compress|decompress <model> <input.txt> <output.txt>\n"+
		"       cmd crossentropy <model> <input.txt>\n\n"+
		"Available models:\n")

	for _, model := range lossless.PredictorIDs() {
		fmt.Fprintln(os.Stderr, " * "+model)
	}
	fmt.Fprintln(os.Stderr)

	os.Exit(1)
}
