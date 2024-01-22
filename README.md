# mkasm
An assembler for the PDP-8 using modified PAL syntax.

## Features
`mkasm` supports PAL-III syntax with a few additional useful features.

### Memory Reference Instructions

All standard PDP-8 Memory Reference instructions are implemented. Reference
specifiers `I` and `Z` are available for specifying indirect and page
addressing, respectively.


### OPR Micro-Instructions
Currently only group 1 and 2 operate instructions are supported.

<!-- **Group 1 Microinstructions** -->

<!-- **Group 2 Microinstructions** -->


### IOT Instructions
Standard IOT instructions for a teletype are built in.


### Additional Features
mkasm includes some features not found in the PAL assemblers. These have to be enabled with the `-D` flag.
 
**Single Quotes** - Single quotes can be used to get ASCII codes of characters.
This can be useful for storing character values into memory. The trailing
quote is optional. \
*Example:* `DASH,  '-'`

**Double Quotes** - Double quotes can be used to store C-style ASCII strings.
Each character in the string is stored in consecutive memory locations
followed by a NUL (0). \
*Example:* `HELLO, "Hello, World!\n"`

## Usage
Compile programs written in PAL Assembly into several different formats.
The most basic usage is `mkasm example.pa` which producs a Pobj binary
`example.po`.

Current supported output formats are:

* **Pobj**: Human readable format produced by pdpnasm. Each instruction is
encoded as 4 ASCII digits that represent a 12-bit octal number. Addresses are
prefixed with `17`.

* **RIM**: Format used on the OG PDP-8. This format was used to bootstrap the
PDP-8 as a simple RIM loader could be keyed in manually. This binary format
encodes each instruction into 4 bytes. The first two bytes are the address
followed by two bytes containing the instruction to store at the address. The
first byte in the address always has bit 7 set.

* **URL**: Format used for [mkweb](https://pdp8.mckinnon.ninja).

```
Usage: mkasm [options] <src_file> [out_file]

Options:
  -D    Support additional PAL-D syntax
  -dump
        Dump program listing to stdout
  -err-ctx int
        Lines of context surrounding errors
  -help
        Print this message and exit
  -list
        Generate program listing file
  -pobj
        Output in PObject (.po) format
  -rim
        Output in RIM format
  -size
        Print program size information
  -url
        Output in URL format
  -url-base string
        Base URL to use for URL format.
```


Build
-----
Building this project requires `go`, this can be downloaded online.

To build the assembler, first clone this repository:

    git clone https://github.com/Rex--/mkasm.git

Then build the project:

    cd mkasm
    go build .

This should produce the binary `mkasm` in the directory.


Copying
-------
Copyright (c) 2024 Rex McKinnon \
This software is available for free under the permissive University of
Illinois/NCSA Open Source License. See the LICENSE file for full details.
