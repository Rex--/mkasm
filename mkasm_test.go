package main

import (
	"os"
	"testing"
)

func TestAssemble(t *testing.T) {
	// Assemble test.p8 and check its output
	testSrc, _ := os.Open("tests/test.p8")

	parser := NewParser(NewLexer(testSrc), &default_symbols)

	parser.parseP8Assembly()
}
