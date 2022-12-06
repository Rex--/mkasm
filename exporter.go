package main

import (
	"fmt"
	"io"
	"sort"
)

type Memory map[int]int

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

func (m Memory) print() {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, addr := range keys {
		inst := m[addr]
		fmt.Printf("%.4o | %.4o\t%.12b\n", addr, inst, inst)
	}
}
