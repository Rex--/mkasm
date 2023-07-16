*200
/ START,  CLA CLL
/         TAD FA
/         DCA MUL-2
/         TAD FB
/         DCA MUL-1
/         JMS MUL
/         HLT
/         JMP START
/ FA,     0d1024
/ FB,     0d2

/ MUL,	0		/ Multiplication subroutine
MUL,    CLA
        DCA MULDW       / Set working register to 0
        TAD MULB        / Load multiplier into AC
        DCA MULRW       / Save multiplier into working register
MLOOP,	CLA		/ Clear AC
        TAD MULDW       / Load intermediate value
        TAD MULA	/ Add multiplicand
        DCA MULDW	/ Store intermediate value
        CMA		/ Complement AC (-1)
        TAD MULRW	/ Decrement multiplier
        DCA MULRW	/ Save multiplier
        TAD MULRW	/ Load multiplier into AC
        SZA		/ Skip if AC == 0
        JMP MLOOP	/ Jump back to start of multiply loop
        TAD MULDW	/ Load answer into AC
        HLT
        JMP MUL
        / JMP I MUL	/ Return from subroutine
MULA,	0d150
MULB,	0d3
MULRW,  0
MULDW,	0
$
