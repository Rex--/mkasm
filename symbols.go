package main

type SymType int

const (
	PSEUDO SymType = iota
	MRI
	SI
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
	// if lastVal, exists := (*st)[symbol]; exists && val != lastVal.Val {
	// 	// fmt.Printf("Redefined existing symbol: %s = %o\n", symbol, val)
	// } else if !exists {
	// 	// fmt.Printf("Defined new symbol: %s = %o\n", symbol, val)
	// }
	(*st)[symbol] = Symbol{SI, val}
}

func (st *SymbolTable) Label(symbol string, val int) {
	// if lastVal, exists := (*st)[symbol]; exists && val != lastVal.Val {
	// fmt.Printf("Redefined existing label: %s = %o\n", symbol, val)
	// } else if !exists {
	// fmt.Printf("Defined new symbol: %s = %o\n", symbol, val)
	// }
	(*st)[symbol] = Symbol{LABEL, val}
}

var default_symbols SymbolTable = SymbolTable{
	// Pseudo instructions
	"EXPUNGE": Symbol{PSEUDO, -1},
	"FIXTAB":  Symbol{PSEUDO, -1},
	"FIXMRI":  Symbol{PSEUDO, -1},

	// Memory reference instructions
	"AND": Symbol{MRI, 0},
	"TAD": Symbol{MRI, 0o1000},
	"ISZ": Symbol{MRI, 0o2000},
	"DCA": Symbol{MRI, 0o3000},
	"JMS": Symbol{MRI, 0o4000},
	"JMP": Symbol{MRI, 0o5000},

	// Group 1 operate instructions
	"NOP": Symbol{SI, 0o7000},
	"IAC": Symbol{SI, 0o7001},
	"RAL": Symbol{SI, 0o7004},
	"RTL": Symbol{SI, 0o7006},
	"RAR": Symbol{SI, 0o7010},
	"RTR": Symbol{SI, 0o7012},
	"CML": Symbol{SI, 0o7020},
	"CMA": Symbol{SI, 0o7040},
	"CIA": Symbol{SI, 0o7041},
	"CLL": Symbol{SI, 0o7100},
	"STL": Symbol{SI, 0o7120},
	"CLA": Symbol{SI, 0o7200}, // This does the same thing as 0o7600 - pick your favorite
	"GLK": Symbol{SI, 0o7204},
	"STA": Symbol{SI, 0o7240},

	// Group 2 operate instructions
	"HLT": Symbol{SI, 0o7402},
	"OSR": Symbol{SI, 0o7404},
	"SKP": Symbol{SI, 0o7410},
	"SNL": Symbol{SI, 0o7420},
	"SZL": Symbol{SI, 0o7430},
	"SZA": Symbol{SI, 0o7440},
	"SNA": Symbol{SI, 0o7450},
	"SMA": Symbol{SI, 0o7500},
	"SPA": Symbol{SI, 0o7510},
	// "CLA": Symbol{SI, 0o7600}, // This does the same thing as 0o7200 - pick your favorite
	"LAS": Symbol{SI, 0o7604},

	// IOT - Program Interrupt
	"ION": Symbol{SI, 0o6001},
	"IOF": Symbol{SI, 0o6002},

	// IOT - High Speed Perforated Tape Reader
	"RSF": Symbol{SI, 0o6011},
	"RRB": Symbol{SI, 0o6012},
	"RFC": Symbol{SI, 0o6014},

	// IOT - High Speed Perforated Tape Punch
	"PSF": Symbol{SI, 0o6021},
	"PCF": Symbol{SI, 0o6022},
	"PPC": Symbol{SI, 0o6024},
	"PLS": Symbol{SI, 0o6026},

	// IOT - Teletype Keyboard/Reader
	"KSF": Symbol{SI, 0o6031},
	"KCC": Symbol{SI, 0o6032},
	"KRS": Symbol{SI, 0o6034},
	"KRB": Symbol{SI, 0o6036},

	// IOT - Teletype Teleprinter/Punch
	"TSF": Symbol{SI, 0o6041},
	"TCF": Symbol{SI, 0o6042},
	"TPC": Symbol{SI, 0o6044},
	"TLS": Symbol{SI, 0o6046},
}
