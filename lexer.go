package main

import (
	"bufio"
	"bytes"
	"os"
)

type LexType int

const (
	SYMBOL LexType = iota
	PUNCTUATION
	NUMBER
	COMMENT
	EOL
	EOF
	UNKNOWN
)

type Lexeme struct {
	// The type of the lexeme
	Type LexType

	// The raw lexeme bytes
	Bytes []byte

	// Location of lexeme in file
	Line int
	Col  int
}

type Lexer struct {
	// The current lexeme being parsed
	This Lexeme
	// The next lexeme to be parsed
	Next Lexeme

	// File
	f *os.File
	// Lexer line scanner
	s *bufio.Scanner

	// Number of the line of file
	lineNum int
	// Buffer which holds the raw line bytes
	line []byte
	// Scan position in line
	pos int
}

func NewLexer(f *os.File) (l *Lexer) {
	l = new(Lexer)

	// Save file object for reference
	l.f = f

	// Create a new scanner on our reader and set our custom splitLine function
	l.s = bufio.NewScanner(f)
	l.s.Split(scanLines)

	// Read the first line into line buffer
	l.readLine()

	// Read the first lexeme into Next. A successive call to Advance will place
	// this lexeme into This, and scan a new one into Next.
	l.Advance()

	return
}

func (l *Lexer) Reset() {
	l.lineNum = 0
	l.pos = 0

	// Seek file to beginning
	_, err := l.f.Seek(0, 0)
	if err != nil {
		panic("seek error")
	}

	// Create a new scanner on our reader because I couldn't figure out a
	// reliable way to reset the scanner.
	l.s = bufio.NewScanner(l.f)
	l.s.Split(scanLines)
	l.readLine()
	l.Advance()
}

// Advance the current lexeme by one position, moving next -> this and reading
// a new lexeme into next
func (l *Lexer) Advance() {
	l.This = l.Next
	l.Next.Type = UNKNOWN
	l.Next.Bytes = nil
	l.Next.Line = l.lineNum

	// fmt.Println("Scanning line:", l.line)
	// fmt.Printf("%d, %d\t[%d]\t%s\n", l.This.Line, l.This.Col, l.This.Type, strings.TrimSpace(string(l.This.Bytes)))

	// Skip Whitespace
	l.skipWhitespace()

	// Set column of next lexeme
	l.Next.Col = l.pos + 1

	// Check if we're at EOF
	if l.pos == -1 || l.line[l.pos] == 0 || l.line[l.pos] == '$' {
		l.Next.Type = EOF
		l.Next.Bytes = []byte{0}
		return
	}

	// Check if we're at EOL
	if l.line[l.pos] == '\n' || l.line[l.pos] == ';' {
		l.Next.Type = EOL
		l.Next.Bytes = []byte{'\n'} // Should we emit the actual line ending? (either ';' or '\n')
		// If its a colon delimited line, go ahead and increment pos. This will
		// catch colons at the end of line and any empty lines afterwards.
		if l.line[l.pos] == ';' {
			l.pos++
		}
		// Skip blank lines or lines that only contain whitespace
		for l.line[l.pos] == '\n' {
			l.readLine()
			if l.pos == -1 {
				break
			}
			l.skipWhitespace()
		}
		return // Bail
	}

	// Check for comment and read the rest of the line as a comment lexeme
	if l.line[l.pos] == '/' {
		l.Next.Type = COMMENT
		l.Next.Bytes = l.line[l.pos : len(l.line)-1]
		l.pos = len(l.line) - 1
		return // Bail early
	}

	// Check for valid punctuation lexemes
	if p := l.line[l.pos]; p == '=' || p == '*' || p == ',' || p == '.' || p == '-' || p == '+' {
		l.Next.Type = PUNCTUATION
		l.Next.Bytes = l.line[l.pos : l.pos+1]
		l.pos++
		return // Bail early
	}

	// Check for symbols, numbers, and unknown lexemes
	if isLetter(l.line[l.pos]) {
		// fmt.Println("Found symbol:", string(l.line[l.pos:]))
		// Symbols are alphanumeric and start with a letter
		start := l.pos
		l.pos++
		for c := l.line[l.pos]; isAlphaNum(c); {
			l.pos++
			if l.pos == len(l.line) {
				break
			} else {
				c = l.line[l.pos]
			}
		}
		l.Next.Type = SYMBOL
		l.Next.Bytes = l.line[start:l.pos]

	} else if isDigit(l.line[l.pos]) {
		// fmt.Println("Found number:", string(l.line[l.pos:]))
		//Numbers contain digits
		start := l.pos
		l.pos++
		for c := l.line[l.pos]; isDigit(c); {
			l.pos++
			if l.pos == len(l.line) {
				break
			}
			c = l.line[l.pos]
		}
		l.Next.Type = NUMBER
		l.Next.Bytes = l.line[start:l.pos]

	} else {
		// Invalid character
		l.Next.Bytes = l.line[l.pos : l.pos+1]
		l.pos++
	}
}

// Custom scanLine function. Lines keep their trailing \n, and a NULL byte is
// appended as the EOF character
func scanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, data[0 : i+1], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it with a terminating NULL
	if atEOF {
		return len(data), append(data, 0), nil
	}
	// Request more data.
	return 0, nil, nil
}

// Reads the current line into buffer
func (l *Lexer) readLine() {
	if l.s.Scan() {
		l.line = l.s.Bytes()
		l.pos = 0
	} else {
		l.line = nil
		l.pos = -1
	}
	l.lineNum++
}

func (l *Lexer) skipWhitespace() {
	leng := len(l.line)
	if l.pos == leng || l.pos == -1 {
		return // Bail early if pointer is at end of line or file
	}
	for c := l.line[l.pos]; c == ' ' || c == '\t'; {
		l.pos++
		if l.pos == leng {
			break
		} else {
			c = l.line[l.pos]
		}
	}
}

func isDigit(c byte) bool {
	if (c >= '0' && c <= '9') || c == 'x' || c == 'b' || c == 'o' || c == 'd' {
		return true
	}
	return false
}

func isLetter(c byte) bool {
	if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') {
		return true
	}
	return false
}

func isAlphaNum(c byte) bool {
	if isDigit(c) || isLetter(c) {
		return true
	}
	return false
}
