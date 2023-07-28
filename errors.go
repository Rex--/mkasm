package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var ErrorLexemes []*Lexeme
var ErrorStrings []string

func (l *Lexer) UnknownLexeme(lm *Lexeme, col int, msg string) {
	if col < 0 {
		col = lm.Col
	}
	fmt.Print(formatErrorMsg("unknown lexeme: " + msg))
	fmt.Printf("%3d | %s\n    | %*s\n\n", lm.Line, strings.TrimRight(string(l.line), "\n\r"), col, "^")
	os.Exit(1)
}

func (p *Parser) SyntaxError(lm *Lexeme, msg string) {
	ErrorLexemes = append(ErrorLexemes, lm)
	ErrorStrings = append(ErrorStrings, "syntax error: "+msg)
}

func (p *Parser) IllegalReferenceError(lm *Lexeme, msg string) {
	ErrorLexemes = append(ErrorLexemes, lm)
	ErrorStrings = append(ErrorStrings, "illegal reference: "+msg)
}

func (p *Parser) UndefinedSymbolError(lm *Lexeme, msg string) {
	ErrorLexemes = append(ErrorLexemes, lm)
	ErrorStrings = append(ErrorStrings, "undefined symbol: "+msg)
}

func (p *Parser) UndefinedSymbols() {
	for _, l := range p.undef {
		ErrorLexemes = append(ErrorLexemes, &l)
		ErrorStrings = append(ErrorStrings, "undefined symbol")
	}
}

func (p *Parser) ResetErrors() {
	ErrorLexemes = make([]*Lexeme, 0)
	ErrorStrings = make([]string, 0)
}

func (p *Parser) HasErrors() bool {
	return len(ErrorLexemes) > 0
}

func (p *Parser) PrintErrors() {
	for i := 0; i < len(ErrorLexemes); i++ {
		lexemeStr := string(ErrorLexemes[i].Bytes)
		fmt.Print(formatErrorMsg(ErrorStrings[i] + ": '" + lexemeStr + "'"))
		printLine(p.lex.ferr, ErrorLexemes[i], p.lex.args.ErrCtx)
	}
}

func formatErrorMsg(msg string) string {
	return fmt.Sprintf("****> Error: %s\n", msg)
}

func printLine(f *os.File, lm *Lexeme, ctx int) {
	f.Seek(0, 0)
	lineReader := bufio.NewReader(f)
	for i := 1; i < lm.Line; i++ {
		pl, err := lineReader.ReadString('\n')
		if err != nil && err != io.EOF {
			panic(err)
		}
		if i >= lm.Line-ctx { // Print surrounding context
			fmt.Printf("%3d | %s\n", lm.Line-(lm.Line-i), strings.TrimRight(pl, "\n\r"))
		}
	}
	errLine, err := lineReader.ReadString('\n')
	if err != nil && err != io.EOF {
		panic(err)
	}
	fmt.Printf("%3d | %s\n      %*s%s\n", lm.Line, strings.TrimRight(errLine, "\n\r"), lm.Col, "^", strings.Repeat("~", len(lm.Bytes)-1))

	for i := 1; i <= ctx; i++ {
		pl, err := lineReader.ReadString('\n')
		if err != nil && err != io.EOF {
			panic(err)
		}
		// Print surrounding context
		fmt.Printf("%3d | %s\n", lm.Line+i, strings.TrimRight(pl, "\n\r"))
	}
	fmt.Println()
}
