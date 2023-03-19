// rd = rs1 | imm  = 0x0A | 0x05 = 0x0F

RVector: @
    lw x1, 0x28(x0)     // Load x1 with the contents of: 0x28 BA = 0x0A WA
    andi x2, x1, 0x05
    ebreak              // Stop
