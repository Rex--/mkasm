package main

import (
	"fmt"
	"os"
	"strings"
)

func (l *Lexer) UnknownLexeme(lm *Lexeme, col int, msg string) {
	if col < 0 {
		col = lm.Col
	}
	fmt.Printf("\n****> Error: %s\n    |\n", msg)
	fmt.Printf("%3d | %s\n    | %*s\n\n", lm.Line, strings.TrimRight(string(l.line), "\n\r"), col, "^")
	os.Exit(1)
}

func (p *Parser) SyntaxError(lm *Lexeme, col int, msg string) {
	if col < 0 {
		col = lm.Col
	}
	var linestr string
	if strings.Contains(string(p.lex.line), string(lm.Bytes)) {
		linestr = strings.TrimRight(string(p.lex.line), "\n\r")
	} else if strings.Contains(string(p.lex.prevLine), string(lm.Bytes)) {
		linestr = strings.TrimRight(string(p.lex.prevLine), "\n\r")
	} else {
		linestr = strings.TrimRight(string(lm.Bytes), "\n\r")
	}
	fmt.Printf("\n****> Error: %s\n    |\n", msg)
	fmt.Printf("%3d | %s\n    | %*s\n\n", lm.Line, strings.TrimRight(linestr, "\n\r"), col, "^")
	os.Exit(2)
}

func (p *Parser) UndefinedSymbols() {
	fmt.Println("\n****> Error: undefined symbols\n    |")
	for _, l := range p.undef {
		fmt.Printf("%3d | %s\n    |\n", l.Line, string(l.Bytes))
		// fmt.Println("Undefined symbol", string(l.Bytes))
	}
	fmt.Println()
	os.Exit(3)
}
