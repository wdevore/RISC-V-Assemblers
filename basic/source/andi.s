// rd = rs1 & imm  = 0x0A & 0x05 = 0x00

RVector: @
    @: Main            // Reset vector

Main: @
    lw x1, @Data(x0)     // Load x1 with the contents of Data+0
    andi x2, x1, 0x05
    ebreak              // Stop

Data: @
    d: 0000000A
    @: Data
