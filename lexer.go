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
	STRING
	CHAR
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
	// Th previous lexeme that was parsed
	Prev Lexeme
	// The current lexeme being parsed
	This Lexeme
	// The next lexeme to be parsed
	Next Lexeme

	// File
	f *os.File
	// Command line arguments
	args *CLIArgs
	// Lexer line scanner
	s *bufio.Scanner

	// Number of the line of file
	lineNum int
	// Buffer which holds the raw line bytes
	line []byte
	// Scan position in line
	pos int

	// Buffer to hold the previous line for easy error reporting
	prevLine []byte
}

func NewLexer(f *os.File, args *CLIArgs) (l *Lexer) {
	l = new(Lexer)

	// Save file object for reference
	l.f = f
	l.args = args

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
	l.Prev = l.This
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

	if l.args.LangPalD { // PAL-D doesn't actually support this

		// Check for double quoted strings
		if l.line[l.pos] == '"' {
			l.Next.Type = STRING
			start := l.pos
			l.pos++
			for c := l.line[l.pos]; c != '"'; {
				l.pos++
				if l.pos >= len(l.line) {
					break
				} else {
					c = l.line[l.pos]
				}
			}

			if l.pos < len(l.line) && l.line[l.pos] == '"' { // Include trailing "
				l.pos++
				l.Next.Bytes = l.line[start:l.pos]
			} else {
				l.Next.Bytes = l.line[start:l.pos]
				l.UnknownLexeme(&l.Next, -1, "unterminated string")
				// panic("unterminated string")
			}
			return // Bail
		}

		// Check for single quoted characters
		if l.line[l.pos] == '\'' {
			l.Next.Type = CHAR
			start := l.pos
			l.pos++                    // Skip leading '
			if l.line[l.pos] == '\\' { // Catch characters escaped with \
				l.pos++
				// valid := isEscaped(l.line[l.pos])
				// if !valid {
				// 	l.UnknownLexeme(&l.Next, l.pos+1, "unknown escaped character")
				// 	// panic("unknown char")
				// }
			}

			l.pos++ // Skip the character

			// fmt.Println(string(l.line))
			// fmt.Println(string(l.line[l.pos]))

			// Check for ending ' (not required)
			if l.line[l.pos] == '\'' {
				l.Next.Bytes = l.line[start : l.pos+1]
			} else if isWhitespace(l.line[l.pos]) {
				l.Next.Bytes = l.line[start:l.pos]
			} else {
				l.UnknownLexeme(&l.Next, l.pos, "unknown character")
			}
			l.pos++
			return // Bail

		}
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

		var matchFunc func(c byte) bool = isDigit
		if c := l.line[l.pos]; c == 'b' || c == 'o' || c == 'd' {
			// Number base that only includes digits 0-9
			l.pos++
		} else if c == 'x' {
			// Hex numbers can contain some letters as digits
			matchFunc = isHexDigit
			l.pos++
		}
		for c := l.line[l.pos]; matchFunc(c); {
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
		l.UnknownLexeme(&l.This, l.pos+1, "unknown character")
		// l.Next.Bytes = l.line[l.pos : l.pos+1]
		// l.pos++
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
	l.prevLine = l.line
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
	for c := l.line[l.pos]; isWhitespace(c); {
		l.pos++
		if l.pos == leng {
			break
		} else {
			c = l.line[l.pos]
		}
	}
}

func isWhitespace(c byte) bool {
	if c == ' ' || c == '\t' {
		return true
	}
	return false
}

func isDigit(c byte) bool {
	if c >= '0' && c <= '9' {
		return true
	}
	return false
}

func isHexDigit(c byte) bool {
	if isDigit(c) ||
		(c == 'a' || c == 'A') ||
		(c == 'b' || c == 'B') ||
		(c == 'c' || c == 'C') ||
		(c == 'd' || c == 'D') ||
		(c == 'e' || c == 'E') ||
		(c == 'f' || c == 'F') {
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

// func isEscaped(c byte) bool {
// 	if c == 'n' || c == 'r' || c == 't' || c == '\\' {
// 		return true
// 	} else {
// 		return false
// 	}
// }
