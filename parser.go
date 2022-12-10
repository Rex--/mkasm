package main

import (
	"fmt"
	"strconv"
)

type InstType int

const ()

type Instruction struct {
	// Location in memory of word
	Loc int
	// Optional label that corresponds with this memory location
	Label *Symbol
	// Assembled instruction
	Data int
	// Raw Instruction
	Raw []byte
}

type Parser struct {
	lex    *Lexer
	symtab *SymbolTable
	lc     int
	mem    Memory
	apass  bool // Another Pass?
	pdepth int  // Parsed depth
	mdepth int  // Max depth
}

func NewParser(l *Lexer, st *SymbolTable) *Parser {
	// Create our parser
	return &Parser{
		lex:    l,
		symtab: st,
		lc:     0o200,
		mem:    make(Memory),
		mdepth: 10,
	}
}

func (p *Parser) parseP8Assembly() {
	p.pdepth++
	fmt.Printf("Making pass #%d\n", p.pdepth)

	// Create our parser
	// p := Parser{
	// 	lex:    l,
	// 	symtab: make(SymbolTable),
	// 	mem:    make(map[int16]InstType),
	// }
	// Add defaults to table
	// for i, sym := range *defaults {
	// 	p.symtab[i] = sym
	// }
loop:
	for {
		p.lex.Advance()
		// fmt.Printf("%d, %d\t[%d]\t%s\n", l.This.Line, l.This.Col, l.This.Type, strings.TrimSpace(string(l.This.Bytes)))

		switch p.lex.This.Type {

		case PUNCTUATION:
			switch p.lex.This.Bytes[0] {

			case '*': // Asterisk means Next sets the location counter
				p.lex.Advance()
				var str string
				p.lc, str = p.parseExpression()
				if str != "" {
					panic("Unknown symbol: " + str)
				}
				// fmt.Printf("Setting location counter: %o\n", p.lc)

			case '.':
				fallthrough
			case '-':
				fallthrough
			case '+':
				inst, expr := p.parseExpression()
				if expr != "" {
					fmt.Println("Another pass required:", expr)
					p.apass = true
				} else {
					p.mem[p.lc] = inst
				}
				p.lc++
			}

		case SYMBOL:
			if p.lex.Next.Type == PUNCTUATION && p.lex.Next.Bytes[0] != '.' { // Symbol definition
				switch p.lex.Next.Bytes[0] {
				case '=':
					p.parseSymbolDefinition()
				case ',':
					p.parseLabel()

				case '-':
					fallthrough
				case '+':
					inst, expr := p.parseExpression()
					if expr != "" {
						fmt.Println("Another pass required:", expr)
						p.apass = true
					} else {
						p.mem[p.lc] = inst
					}
					p.lc++
				}

			} else { // Were using the symbol
				// Lookup symbol
				sym := p.symtab.Get(string(p.lex.This.Bytes))
				if sym != nil && sym.Type == MRI {
					// Memory reference instruction
					// symStr := string(p.lex.This.Bytes)
					p.lex.Advance()
					var oprStr string
					var indirect, zeroPage bool
					// Check for (I)ndirect flag
					if p.lex.This.Type == SYMBOL && p.lex.This.Bytes[0] == 'I' {
						oprStr += "I"
						indirect = true
						p.lex.Advance()
					}
					// Check for (Z)ero page flag
					if p.lex.This.Type == SYMBOL && p.lex.This.Bytes[0] == 'Z' {
						oprStr += " Z"
						zeroPage = true
						p.lex.Advance()
					}
					// Parse expression of address operand
					result, expr := p.parseExpression()
					if expr != "" {
						fmt.Println("Another pass required:", expr)
						p.apass = true
					} else {
						if indirect {
							result |= 0b000100000000
						}
						if zeroPage {
							result &= 0b111101111111
						} // else {
						// 	result |= 0b000010000000
						// }
						result |= sym.Val
						p.mem[p.lc] = result
						// fmt.Printf("MRI: %s %s %o %b\n", symStr, oprStr, result, result)
					}
				} else {
					inst, expr := p.parseExpression()
					if expr != "" {
						fmt.Println("Another pass required:", expr)
						p.apass = true
					} else {
						p.mem[p.lc] = inst
					}
				}
				p.lc++
			}

		case NUMBER:
			inst := p.parseNumber()
			p.mem[p.lc] = inst
			p.lc++

		case EOF:
			break loop
		}
	}

	if p.apass && p.pdepth < p.mdepth {
		p.pdepth++
		fmt.Printf("Making pass #%d\n", p.pdepth)
		// Reset lexer to beginning of file
		p.lex.Reset()
		// Reset parser state
		p.lc = 0
		p.apass = false
		goto loop

	} else if p.pdepth >= p.mdepth {
		fmt.Println("Parsing failed")
	}
}

func (p *Parser) parseNumber() int {
	var err error
	var i64 int64

	if len(p.lex.This.Bytes) > 2 && isLetter(p.lex.This.Bytes[1]) {
		// Number base explicitly set with '0<x|o|b|d>' prefix

		if p.lex.This.Bytes[1] == 'd' {
			// Parse decimal number
			// fmt.Println(string(p.lex.This.Bytes[2:]))
			i64, err = strconv.ParseInt(string(p.lex.This.Bytes[2:]), 10, 16)
		} else {
			// ParseInt supports hex, bin, and octal automatically when passed 0 base
			i64, err = strconv.ParseInt(string(p.lex.This.Bytes), 0, 16)
		}

	} else {
		// Default to parsing number as octal
		i64, err = strconv.ParseInt(string(p.lex.This.Bytes), 8, 16)
	}

	if err != nil {
		panic("Number error (Too large?)")
	}
	// fmt.Println("Parsed number:", string(p.lex.This.Bytes), "->", strconv.Itoa(int(i64)))
	// fmt.Printf("NUM: %o\t%s ->\t\t%o\n", p.lc, string(p.lex.This.Bytes), int(i64))
	return int(i64)
}

func (p *Parser) parseExpression() (int, string) {
	// fmt.Print("Parsing expression: ")
	var start string = string(p.lex.This.Bytes)
	var sign, operand string

	if p.lex.This.Type == PUNCTUATION { // (<+|->A) OR (. [<+|-> B]) formatted expression

		if p.lex.This.Bytes[0] == '.' { // (. [<+|-> B]) formatted expression
			if p.lex.Next.Type == PUNCTUATION {
				var a int = p.lc
				var b int
				p.lex.Advance()
				sign = string(p.lex.This.Bytes)
				p.lex.Advance()

				operand = string(p.lex.This.Bytes)
				if isLetter(p.lex.This.Bytes[0]) { // Lookup symbol
					sym := p.symtab.Get(operand)
					if sym != nil {
						b = sym.Val
					} else {
						return -1, operand
					}
				} else if isDigit(p.lex.This.Bytes[0]) { // Parse number
					b = p.parseNumber()
				} else {
					panic("unknown expression operand")
				}

				var ans int
				switch sign {
				case "-":
					ans = a - b
				case "+":
					ans = a + b
				default:
					panic("unknown operation")
				}
				// fmt.Printf("OPR: %o\t%s%s%s\t%o\n", p.lc, start, sign, operand, ans)
				return ans, ""

			} else if p.lex.Next.Type == COMMENT || p.lex.Next.Type == EOL {
				return p.lc, ""
			}
		} else { // <+|->A formatted expression
			start = ""
			sign = string(p.lex.This.Bytes)
			p.lex.Advance()
			var a int
			operand = string(p.lex.This.Bytes)
			if isLetter(p.lex.This.Bytes[0]) { // Lookup symbol
				sym := p.symtab.Get(operand)
				if sym != nil {
					a = sym.Val
				} else {
					return -1, operand
				}
			} else if isDigit(p.lex.This.Bytes[0]) { // Parse number
				a = p.parseNumber()
			} else {
				panic("unknown expression operand")
			}

			var ans int
			switch sign {
			case "-":
				ans = a * -1
			case "+":
				ans = a
			default:
				panic("unknown operation")
			}

			return ans, ""
		}
	} else if p.lex.Next.Type == PUNCTUATION { // (A <+|-> B) formatted expression
		// a := string(l.This.Bytes)
		var a, b int
		if isLetter(p.lex.This.Bytes[0]) { // Lookup symbol
			sym := p.symtab.Get(start)
			if sym != nil {
				a = sym.Val
			} else {
				// Skip to end of expression
				for p.lex.This.Type != EOL {
					p.lex.Advance()
				}
				return -1, start
			}
		} else if isDigit(p.lex.This.Bytes[0]) { // Parse number
			a = p.parseNumber()
		} else {
			panic("unknown expression operand")
		}
		p.lex.Advance()

		sign = string(p.lex.This.Bytes)
		p.lex.Advance()

		operand = string(p.lex.This.Bytes)
		if isLetter(p.lex.This.Bytes[0]) {
			osym := p.symtab.Get(operand)
			if osym != nil {
				b = osym.Val
			} else {
				// Skip to end of expression
				for p.lex.This.Type != EOL {
					p.lex.Advance()
				}
				return -1, start
			}
		} else if isDigit(p.lex.This.Bytes[0]) {
			b = p.parseNumber()
		} else {
			panic("unknown expression operand")
		}

		var answer int
		switch sign {
		case "-":
			answer = a - b
		case "+":
			answer = a + b
		default:
			panic("unsupported operation in expression")
		}
		// fmt.Printf("OPR: %o\t%s%s%s\t%o\n", p.lc, start, sign, operand, answer)
		return answer, ""

	} else if p.lex.Next.Type == SYMBOL { // (A B) formatted expression (AND)
		start = string(p.lex.This.Bytes)
		p.lex.Advance()
		operand = string(p.lex.This.Bytes)
		sSym := p.symtab.Get(start)
		eSym := p.symtab.Get(operand)
		if sSym != nil && eSym != nil {
			// fmt.Printf("AND: %o\t%s|%s\t%o\n", p.lc, start, operand, sSym.Val|eSym.Val)
			return (sSym.Val | eSym.Val), ""
		}
	} else if p.lex.Next.Type == COMMENT || p.lex.Next.Type == EOL { // (A) formatted expression
		if isLetter(p.lex.This.Bytes[0]) {
			sym := p.symtab.Get(start)
			if sym != nil {
				// fmt.Printf("EXP: %o\t%s ->\t\t%o\n", p.lc, start, sym.Val)
				return sym.Val, ""
			}
		} else if isDigit(p.lex.This.Bytes[0]) {
			return p.parseNumber(), ""
		} else {
			panic("unknown expression operand")
		}
		return -1, start
	} else {
		panic("error: unknown syntax")
	}

	return -1, "error"
	// fmt.Printf("%s%s%s\n", start, sign, operand)
}

func (p *Parser) parseSymbolDefinition() {
	symbol := string(p.lex.This.Bytes)
	p.lex.Advance() // Symbol to define
	p.lex.Advance() // Equal sign '='
	value, str := p.parseExpression()
	if str == "" {
		p.symtab.Set(symbol, int(value))
	} else {
		fmt.Printf("Another pass required: %s (%s)\n", str, symbol)
		p.apass = true
	}
}

func (p *Parser) parseLabel() {
	symbol := string(p.lex.This.Bytes)
	p.lex.Advance() // Comma ','
	p.symtab.Label(symbol, p.lc)
}
