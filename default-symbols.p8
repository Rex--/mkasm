/ This file includes the built-in symbol definitions.
/ It does not include memory reference instructions.

/ Group 1 operate instructions
NOP=7000
IAC=7001
RAL=7004
RTL=7006
RAR=7010
RTR=7012
CML=7020
CMA=7040
CIA=7041
CLL=7100
STL=7120
CLA=7200 / This does the same thing as 7600 - pick your favorite
GLK=7204
STA=7240

/ Group 2 operate instructions
HLT=7402
OSR=7404
SKP=7410
SNL=7420
SZL=7430
SZA=7440
SNA=7450
SMA=7500
SPA=7510
/ CLA=7600 / This does the same thing as 7200 - pick your favorite
LAS=7604

/ IOT - Program Interrupt
ION=6001
IOF=6002

/ IOT - High Speed Perforated Tape Reader
RSF=6011
RRB=6012
RFC=6014

/ IOT - High Speed Perforated Tape Punch
PSF=6021
PCF=6022
PPC=6024
PLS=6026

/ IOT - Teletype Keyboard/Reader
KSF=6031
KCC=6032
KRS=6034
KRB=6036

/ IOT - Teletype Teleprinter/Punch
TSF=6041
TCF=6042
TPC=6044
TLS=6046