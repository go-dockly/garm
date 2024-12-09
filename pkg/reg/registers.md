Parameter/Result Registers (Caller-saved):
x0  - First parameter and return value
x1  - Second parameter
x2  - Third parameter
x3  - Fourth parameter
x4  - Fifth parameter
x5  - Sixth parameter
x6  - Seventh parameter
x7  - Eighth parameter
x8  - Indirect result location register / syscall number

Temporary Registers (Caller-saved):
x9  - Temporary register 1
x10 - Temporary register 2
x11 - Temporary register 3
x12 - Temporary register 4
x13 - Temporary register 5
x14 - Temporary register 6
x15 - Temporary register 7
x16 - IP0: Intra-procedure-call temporary register 1
x17 - IP1: Intra-procedure-call temporary register 2
x18 - Platform register (reserved)

Callee-saved Registers:
x19 - Callee-saved register 1
x20 - Callee-saved register 2
x21 - Callee-saved register 3
x22 - Callee-saved register 4
x23 - Callee-saved register 5
x24 - Callee-saved register 6
x25 - Callee-saved register 7
x26 - Callee-saved register 8
x27 - Callee-saved register 9
x28 - Callee-saved register 10

Special Purpose:
x29 - Frame pointer (FP)
x30 - Link register (LR)
sp  - Stack pointer (x31 when used as source/destination)

Usage Guidelines:
- x0-x7:   Use for function parameters and short-lived values
- x8:      Special for indirect returns and syscalls
- x9-x15:  Use for temporary calculations within a function
- x16-x17: Reserved for linker/dynamic linking
- x18:     Reserved for platform use
- x19-x28: Use for values that must survive function calls
- x29-sp:  Must be handled according to ABI rules