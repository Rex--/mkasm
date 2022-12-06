package main

import (
	"fmt"
	"os"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Usage:", os.Args[0], "<src_file> <out_file>")
		os.Exit(1)
	}

	// Open file
	srcFile, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer srcFile.Close()

	lexer := NewLexer(srcFile)
	parser := NewParser(lexer, &default_symbols)
	parser.parseP8Assembly()

	parser.mem.print()

	// Open out file
	outFile, err := os.Create(os.Args[2])
	if err != nil {
		panic(err)
	}
	defer outFile.Close()
	// Write object file
	parser.mem.exportPObject(outFile)
}
