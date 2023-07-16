*10
SPTR,   STACK

*50
PUSH=DCA SPTR      / Helper macro for pushing AC onto the stack
POP,    0               / Subroutine POP pops the top of the stack into AC
        CLA             / Clear AC
        TAD SPTR        / Add top stack addr into AC
        DCA SPA         / Store address
        STA             / Generate a -1 in AC
        TAD SPA         / Addr + (-1)
        DCA SPTR        / Store decremented stack pointer
        TAD I SPA       / Indirect load the stored addr into AC
SPA,    0

*100
STACK,  0
        0
        0
        0
        0
        0
        0
        0
        0
        0

*200
$