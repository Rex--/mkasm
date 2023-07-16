
*200

START,  JMS DV
B,      3
A,      100
        HLT
        JMP START
DV,     0
        CLA             / IGNORE AC
        TAD I DV        / GET ARGS
        ISZ DV
        DCA DA          / DIVIDEND
        TAD I DV
        ISZ DV
        / SNA             / B = 0 TEST
        / JMP ERROR       
        CLL CIA         / SUBRTACTION WILL BE DONE BY ADDING
        DCA DB          / -DIVISOR
        DCA REMAIN      / CLEAR MSB (DIVENDED)
        CLL CLA CMA RAL
DV1,    DCA QUOT
        TAD DA
        CLL RAL
        DCA DA
        TAD REMAIN
        RAL
        DCA REMAIN
        TAD REMAIN
        TAD DB
        SZL
        DCA REMAIN
        CLA
        TAD QUOT
        RAL
        SZL
        JMP DV1
        JMP I DV
DA,     0
REMAIN, 0
DB,     0
QUOT,   0
