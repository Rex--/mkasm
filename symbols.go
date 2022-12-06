package main

type SymType int

const (
	SI SymType = iota
	MRI
	LABEL
)

type Symbol struct {
	Type SymType
	Val  int
}

type SymbolTable map[string]Symbol

func (st *SymbolTable) Get(symbol string) *Symbol {
	if sym, exists := (*st)[symbol]; exists {
		// fmt.Printf("Found symbol: %v\n", sym)
		return &sym
	}
	// fmt.Println("Symbol not found:", symbol)
	return nil
}

func (st *SymbolTable) Set(symbol string, val int) {
	if lastVal, exists := (*st)[symbol]; exists && val != lastVal.Val {
		// fmt.Printf("Redefined existing symbol: %s = %o\n", symbol, val)
	} else if !exists {
		// fmt.Printf("Defined new symbol: %s = %o\n", symbol, val)
	}
	(*st)[symbol] = Symbol{SI, val}
}

func (st *SymbolTable) Label(symbol string, val int) {
	if lastVal, exists := (*st)[symbol]; exists && val != lastVal.Val {
		// fmt.Printf("Redefined existing label: %s = %o\n", symbol, val)
	} else if !exists {
		// fmt.Printf("Defined new symbol: %s = %o\n", symbol, val)
	}
	(*st)[symbol] = Symbol{LABEL, val}
}

var default_symbols SymbolTable = SymbolTable{
	// Memory reference instructions
	"AND": Symbol{MRI, 0},
	"TAD": Symbol{MRI, 0o1000},
	"SZA": Symbol{MRI, 0o2000},
	"DCA": Symbol{MRI, 0o3000},
	"JMS": Symbol{MRI, 0o4000},
	"JMP": Symbol{MRI, 0o5000},

	// Group 1 operate instructions
	"CLL": Symbol{SI, 0o7100},
	"CLA": Symbol{SI, 0o7200},

	// Group 2 operate instructions
	"HLT": Symbol{SI, 0o7402},
	"SNA": Symbol{SI, 0o7450},
}
