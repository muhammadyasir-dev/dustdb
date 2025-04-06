#include "textflag.h"

// Function to count trailing zeros
TEXT ·countTrailingZeros(SB), NOSPLIT, $0
    MOVQ x+0(FP), AX        // Load the input number into AX
    TZCNTQ AX, AX           // Count trailing zeros in AX
    MOVQ AX, ret+8(FP)      // Store the result
    RET
