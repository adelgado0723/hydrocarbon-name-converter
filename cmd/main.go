package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/adelgado0723/hw-adelgado0723/pkg/converter"
)

const helpMessage = `Usage: go run cmd/main.go [SMILES string]`

func handleInput() (string, error) {
	if len(os.Args) < 2 {
		return "", fmt.Errorf("Error: Please provide a SMILES string")
	}

	help := flag.Bool("help", false, helpMessage)
	flag.Parse()

	if *help {
		fmt.Println(helpMessage)
		os.Exit(0)
	}

	return converter.Convert(flag.Arg(0))
}

func main() {
	smiles, err := handleInput()
	if err != nil {
		fmt.Println(err)
		fmt.Println(helpMessage)
		os.Exit(1)
	}
	fmt.Println(smiles)
	os.Exit(0)
}
