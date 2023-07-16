/ This simple program shifts a number (A) by a number (B) of bits left (<<).


*200
SHFT,   CLA
        TAD A
        DCA SHFTA
        TAD B
        DCA SHFTB
SHFTL,  CLA / Clear AC to start shift loop
        TAD SHFTA       / Load A into AC
        CLL RAL         / Shift left
        DCA SHFTA       / Store intermediate A
        STA             / Generate a -1 in AC
        TAD SHFTB       / Decrement B
        DCA SHFTB       / Save B
        TAD SHFTB       / Load B into AC
        SZA         / Skip if AC == 0 (Shifting is done)
        JMP SHFTL       / Jump to beginning of loop
        CLA
        TAD SHFTA
        HLT
        JMP SHFT
SHFTA,  0
SHFTB,  0
A,      1
B,      0d11
$
