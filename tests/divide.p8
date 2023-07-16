*200
DIV,    CLA
        DCA QUOT        / Clear Quotient register
        TAD DIVA        / Load Divendend
        DCA DIVD        / Store working divended
        TAD DIVB        / Load divisor
        CIA             / -divisor
        DCA DIVR        / Store -divisor
        TAD DIVD        / Load initial dividend

DLOOP,  TAD DIVR        / Add -divisor
        DCA DIVD        / Store new dividend
        TAD QUOT        / Load quotient
        IAC             / Increment quotient
        DCA QUOT        / Store new quotient
        TAD DIVD        / Load intermediate dividend
        SZA             / Skip if AC == 0
        JMP DLOOP       / Jump to divide loop
        TAD QUOT        / Load answer into AC
        HLT             / Halt
        JMP DIV         / Jump to beginning



QUOT,   0
DIVD,   0
DIVR,   0
DIVA,   0d144
DIVB,   0d12