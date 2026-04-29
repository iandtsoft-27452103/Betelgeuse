//go:build amd64 && !purego

#include "textflag.h"

// func PextAsm(src, mask uint64) uint64
TEXT ·PextAsm(SB), NOSPLIT, $0-24
    MOVQ src+0(FP), AX
    MOVQ mask+8(FP), CX
    PEXTQ CX, AX, AX   // BMI2 instruction
    MOVQ AX, ret+16(FP)
    RET
