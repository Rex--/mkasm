package main

import (
	"bytes"
	"fmt"
	"io"
	"sort"
)

type Memory map[int]int

// A P Object(.po) file is in the format used by pdpnasm.
// Each line represents either an origin address (prefixed with 0xF---)
// or an instruction to be placed in memory at the last specified origin + the offset (lines since)
func (m Memory) exportPObject(w io.Writer) {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	var lastAddr int = -1
	for _, addr := range keys {
		// Write change of address line
		if addr > lastAddr+1 {
			fmt.Fprintf(w, "%o\n", addr|0b1111000000000000)
		}
		// Write instruction word
		fmt.Fprintf(w, "%o\n", m[addr])
		// Update last address
		lastAddr = addr
	}
}

// The read in mode (RIM) format is a binary format originally used for paper tapes.
// It was the format used for the first bootstrapping programs on the PDP-8.
// The RIM loader was small enough to be keyed in manually using the switch
// register, using only 18 memory locations (16 instructions to key in) !
//
// The original format conveniently used 8-column tape so we use an 8-bit byte
// for the smallest block. The basic format is a 2-byte address followed by
// a 2-byte memory word. In both cases the actual address or data is in the
// lower 6-bits of each byte.
//
// The first byte of the address always has bit 7 set. The following byte
// representing the lower 6-bits of the address should not have bit 7 set.
//
// Both bytes of data should have bits 7 and 8 cleared. Same as the address,
// the first byte represents the MSB and the 2nd is the LSB.
//
// The program consists of these 4-bytes for every word to be programmed in
// memory. The program can be lead and trailed with zero or more of the
// leader/trailer byte value 0x80 or 1000 0000 in binary.
func (m Memory) exportRim(w io.Writer) {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	var p []byte

	p = append(p, 0o200, 0o200)
	for _, addr := range keys {
		// Add address
		p = append(p, byte((addr&0o7700)>>6)|0o100, byte(addr&0o77))
		inst := m[addr]
		p = append(p, byte((inst&0o7700)>>6), byte(inst&0o77))
	}
	p = append(p, 0o200, 0o200)
	_, err := w.Write(p)
	if err != nil {
		panic("Unable to write")
	}
}

// var urlBase = "http://localhost"

func (m Memory) exportURL(urlBase string) {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	link := urlBase + "?core="
	var lastAddr int = -1
	for i, addr := range keys {
		// Write change of address line
		if addr > lastAddr+1 {
			link += fmt.Sprintf("*0%o,", addr)
		}
		if i == len(keys)-1 {
			// Write last instruction word
			link += fmt.Sprintf("0%o", m[addr])
		} else {
			// Write instruction word
			link += fmt.Sprintf("0%o,", m[addr])
		}
		// Update last address
		lastAddr = addr
	}
	fmt.Print(link)
}

func (m Memory) exportListing(w io.Writer, lst map[int][]byte, labels map[int][]byte) {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	fmt.Fprintln(w, "Abs\tInst")
	fmt.Fprintln(w, "Addr\tData\tTag\t\tInstruction")
	fmt.Fprintln(w, "-----\t----\t--------\t-----------")
	lastAddr := keys[0]
	for _, addr := range keys {
		if addr-lastAddr > 1 {
			fmt.Fprintf(w, "    :\n")
		}
		inst := m[addr]
		var label []byte
		var line []byte
		if l, exists := labels[addr]; exists {
			// Add trailing comma
			label = append(bytes.TrimSpace(l), ',')

			// Cut label from instruction if on the same line
			line = bytes.TrimPrefix(lst[addr], label)

			// Add extra tab for alignment (allows for labels > 8 char long)
			if len(label) < 8 {
				label = append(label, '\t')
			}
		} else {
			line = lst[addr]
			label = []byte("\t")
		}
		// Trim leading/trailing whitespace and remove newlines (in case of errors)
		line = bytes.TrimSpace(bytes.ReplaceAll(line, []byte("\n"), []byte("")))

		before, after, found := bytes.Cut(line, []byte("/"))
		var comment []byte
		if found {
			line = bytes.TrimSpace(before)
			if len(line) < 8 {
				line = append(line, '\t', '\t')
			} else if len(line) < 16 {
				line = append(line, '\t')
			}
			comment = append([]byte("/ "), bytes.TrimSpace(after)...)
		} else {
			comment = []byte("")
		}

		// Check for string in line
		if bytes.ContainsAny(line, "\"") {
			// Only print the full string on the first line with the label
			if len(label) == 1 && label[0] == '\t' {
				line = []byte("\t\t")

				// Add comment with ascii character of memory location
				if len(comment) == 0 {
					char := byte(m[addr])
					if char == 0 {
						comment = []byte("/ NULL")
					} else if char == '\n' {
						comment = []byte("/ \"\\n\"")
					} else {
						comment = []byte{'/', ' ', '"', byte(m[addr]), '"'}
					}
				}
			} else {
				// Add comment of character if line is short enough
				if len(line) < 24 {
					char := byte(m[addr])
					if char == 0 {
						comment = []byte("/ NULL")
					} else if char == '\n' {
						comment = []byte("/ \"\\n\"")
					} else {
						comment = []byte{'/', ' ', '"', byte(m[addr]), '"'}
					}
				}

			}
		}
		fmt.Fprintf(w, "%.4o,\t%.4o\t%s\t%s\t%s\n", addr, inst, label, line, comment)
		lastAddr = addr
	}
	fmt.Fprintln(w, "$")
}

func (m Memory) exportSize() {
	wordsTotal := 0o7777
	wordsUsed := len(m)
	wordsFree := wordsTotal - wordsUsed
	fmt.Printf("used: %d  free: %d  total: %d (words)\n", wordsUsed, wordsFree, wordsTotal)
}
