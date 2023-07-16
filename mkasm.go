package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
)

type CLIArgs struct {
	ProgName  string
	InFile    string
	OutFile   string
	CustomExt bool

	Pobj bool
	Ihex bool
	Rim  bool
	Bin  bool
	URL  bool

	Dump bool
}

func printUsage() {
	fmt.Println("Usage:", os.Args[0], "[options] <src_file> [out_file]")
	fmt.Printf("\nOptions:\n")
	flag.PrintDefaults()
}

func parseArgs() CLIArgs {

	args := CLIArgs{}

	// Set program name
	args.ProgName = os.Args[0]

	flag.Usage = printUsage

	// Add flags
	flag.BoolVar(&args.Pobj, "pobj", false, "Output in PObject (.po) format")
	flag.BoolVar(&args.Rim, "rim", false, "Output in RIM format")
	flag.BoolVar(&args.URL, "url", false, "Output as encoded url")
	flag.BoolVar(&args.Dump, "dump", false, "Dump assembled program to stdout")
	help := flag.Bool("help", false, "Print this message and exit")

	// Parse
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	// Get remaining positional arguments (infile [outfile])
	if len(flag.Args()) == 1 {
		args.InFile = flag.Arg(0)
		// Get outfile based on in file
		args.OutFile = strings.TrimSuffix(flag.Arg(0), path.Ext(flag.Arg(0)))
	} else if len(flag.Args()) == 2 {
		args.InFile = flag.Arg(0)
		// Get the extension of the outfile and output in that format if known
		ext := path.Ext(flag.Arg(1))
		switch ext {
		case ".rim":
			fallthrough
		case ".rm":
			fallthrough
		case ".RIM":
			fallthrough
		case ".RM":
			args.Rim = true
			args.OutFile = strings.TrimSuffix(flag.Arg(1), ext)

		case ".pobj":
			fallthrough
		case ".po":
			fallthrough
		case ".PO":
			args.Pobj = true
			args.OutFile = strings.TrimSuffix(flag.Arg(1), ext)

		default:
			// Save the extension if we don't recognize it
			args.CustomExt = true
			args.OutFile = flag.Arg(1)
		}
	} else {
		flag.Usage()
		os.Exit(1)
	}

	// Set a default output format if we couldn't deduce one
	if !args.Pobj && !args.Rim && !args.URL {
		// Default currently is pobj because it's human readable
		args.Pobj = true
	}

	return args
}

func main() {

	args := parseArgs()

	// Open file
	srcFile, err := os.Open(args.InFile)
	if err != nil {
		panic(err)
	}
	defer srcFile.Close()

	lexer := NewLexer(srcFile)
	parser := NewParser(lexer, &default_symbols)
	parser.parseP8Assembly()

	if args.Dump {
		parser.mem.print()
	}

	// Open out file
	// outFile, err := os.Create(args.OutFile)
	// if err != nil {
	// 	panic(err)
	// }
	// defer outFile.Close()

	// Write output file in specified format(s)
	if args.Pobj {
		outPath := args.OutFile
		if !args.CustomExt {
			outPath += ".po"
		}
		outFile, err := os.Create(outPath)
		if err != nil {
			panic(err)
		}
		fmt.Println("Writing PObj output file:", outPath)
		parser.mem.exportPObject(outFile)
		outFile.Close()
	}

	if args.Rim {
		outPath := args.OutFile
		if !args.CustomExt {
			outPath += ".rim"
		}
		outFile, err := os.Create(outPath)
		if err != nil {
			panic(err)
		}
		fmt.Println("Writing RIM output file:", outPath)
		parser.mem.exportRim(outFile)
		outFile.Close()
	}

	if args.URL {
		fmt.Println("Output URL:")
		parser.mem.exportURL()
	}
}
