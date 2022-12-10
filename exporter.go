package main

import (
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

// func (m Memory) exportHexfile(w io.Writer) {
// 	keys := make([]int, 0, len(m))
// 	for k := range m {
// 		keys = append(keys, k)
// 	}
// 	sort.Ints(keys)

// }

func (m Memory) print() {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, addr := range keys {
		inst := m[addr]
		fmt.Printf("%.12b\t%.4o | %.4o\t%.12b\n", addr, addr, inst, inst)
	}
}
