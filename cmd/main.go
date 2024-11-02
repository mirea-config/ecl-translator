package main

import (
	"ecl-translator/internal/translate"
	"flag"
	"fmt"
	"io"
	"os"
)

var flagOutput string

func init() {
	flag.StringVar(&flagOutput, "o", "./main.ecl", "output .ecl file path")
}

func main() {
	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("ecl: 0 arguments provided")
		return
	}

	jsonPath := args[0]
	jsonFile, err := os.Open(jsonPath)
	if err != nil {
		fmt.Printf("ecl: %s\n", err.Error())
		return
	}
	defer jsonFile.Close()

	src, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Printf("ecl: %s\n", err.Error())
		return
	}

	if err = translate.Translate(src, flagOutput); err != nil {
		fmt.Printf("ecl: %s\n", err.Error())
		return
	}
}
