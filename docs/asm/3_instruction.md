<p align="center"> 
  Visitor count<br>
  <img src="https://profile-counter.glitch.me/sagar-viradiya/count.svg" />
</p>

# Instructions

ARMv8-A is a load/store architecture. Data processing instructions do not operate directly on data in memory as we find with the x86 architecture. The data is first loaded into registers, modified, and then stored back in memory or simply discarded once it’s no longer required. Most data processing instructions use one destination register and two source operands.

```asm
Xd, Xn, Operand2 // Instruction
```
Xd is the destination register. Xn is the register that is operated on. The use of X indicates a full 64-bit CPU register. Operand2 might be a register, a modified register, or an immediate value.

## Let's see some of the most common instructions we might encounter:
1. MOV – Move a value into a register
```asm 
MOV X0, #5 // Move the value 5 into register X0
```
2. ADD – Add two values
```asm 
ADD X0, X1, X2 // Add vals in X1 and X2, store result in X0
```
NOTE: register can hold a byte, halfword, word, or doubleword
3. SUB – Subtract one value from another
```asm 
SUB X0, X1, X2 // Subtract the value in X2 from X1, store result in X0
```
4. RSB – Reverse subtract
```asm 
RSB X0, X1, #0 // Subtract X1 from 0, store result in X0
```
5. MUL – Multiply two values
```asm 
MUL X0, X1, X2 // Multiply X1 by X2, store result in X0
```
6. DIV – Divide two values (on some ARM architectures, but not ARMv7)
```asm 
SDIV X0, X1, X2 // Divide X1 by X2, store result in X0
```
7. AND – Bitwise AND operation
```asm 
AND X0, X1, X2 // Perform a bitwise AND of X1 and X2, store the result in X0
```
8. ORR – Bitwise OR operation
```asm 
ORR X0, X1, X2 // Perform a bitwise OR of X1 and X2, store the result in X0
```
9. EOR – Bitwise Exclusive OR (XOR)
```asm 
XOR X0, X1, X2 // Perform bitwise XOR of X1 and X2, store result in X0
```
10. BIC – Bitwise AND with complement
```asm 
BIC X0, X1, X2 // Perform bitwise AND of X1 with complement of X2
```
11. CMP – Compare two values (sets flags based on the result)
```asm 
CMP X0, X1 // Compare X0 and X1 by subtracting X1 from X0 (affects flags)
```
12. CMN – Compare Negative (sets flags as if adding two values)
```asm 
CMN X0, X1 // Compare by adding X0 and X1 (affects flags)
```
13. TST – Test bits (bitwise AND, sets flags based on result)
```asm 
TST X0, X1 // Perform bitwise AND between X0 and X1, set flags
```
14. TEQ – Test Equivalence (bitwise XOR, sets flags based on result)
```asm 
TEQ X0, X1 // Perform bitwise XOR between X0 and X1, set flags
```
15. LDR – Load a value from memory into a register
```asm 
LDR X0, [X1] // Load value at address in X1 into X0
```
16. STR – Store a register value into memory
```asm 
STR X0, [X1] // Store value from X0 into memory address in X1
```
17. LDRB – Load a byte from memory
```asm 
LDRB X0, [X1] // Load byte at address in X1 into X0
```
18. STRB – Store a byte into memory
```asm 
STRB X0, [X1] // Store byte in X0 into address in X1
```
19. LDM – Load multiple registers from memory
```asm 
LDM X0, {X1, X2} // Load values at address X1 and X2 into X0
```
20. STM – Store multiple registers into memory
```asm 
STM X0, {X1, X2} // Store values in X1 and X2 into X0
```
21. B – Unconditional branch (jump to a label)
```asm 
B loop // Jump to label loop
```
22. BL – Branch with link (call a subroutine)
```asm 
BL func // Call subroutine at label func and save return address
```
23. B.EQ – Branch if equal (based on condition flags)
```asm 
B.EQ label // Branch to “label” if zero flag (Z) is set (equality)
```
24. B.NE – Branch if not equal
```asm 
B.NE label // Branch to “label” if zero flag is clear (inequality)
```
25. B.GT – Branch if greater than
```asm 
B.GT label // Branch if greater than (N and V flags match, and Z flag is clear)
```
26. B.LT – Branch if less than
```asm 
B.LT label // Branch if less than (N flag differs from V flag)
```
27. LSL – Logical shift left
```asm 
LSL X0, X1, #2 // Shift bits in X1 left by 2 places, store result in X0
```
28. LSR – Logical shift right
```asm 
LSR X0, X1, #2 // Shift bits in X1 right by 2 places, store result in X0
```
29. ROR – Rotate bits right
```asm 
ROR X0, X1, #2 // Rotate bits in X1 right by 2 places, store result in X0
```
30. Modulo (x0 = x0 % 37)
```asm 
mov     X1, 37
udiv    X2, X0, X1
msub    X0, X2, X1, X0
```
31. Ternary
```asm
subs    x1, x1, 1
csel    w0, w3, w4, eq          // w0 = (x1==0) ? w3 : w4
```
ARM assembly includes many more instructions. The total number varies depending on architecture (eg ARMv7, ARMv8) and instruction sets included (eg Thumb, NEON, SIMD...)

32. Return an Int
```asm
    mov	W0, 16
    ret
```
33. Return a Long
```asm
    mov	X0, 16
    ret
```
34. Return a Float
```asm
    fmov	s0, 16
    ret
```
35. Return a Double
```asm
    fmov	d0, 16
    ret
```
36. Bad (memory heavy)
```asm
str w0, [x1]        // Store to memory
ldr w0, [x1]        // Load right back
add w0, w0, #1      // Add
```
37. Good (register heavy)
```asm
mov w2, w0          // Keep in register
add w2, w2, #1      // Work with register
str w2, [x1]        // Store only when needed
```

## Instructions to avoid...

###  **Instructions to generally avoid or use carefully on ARM64:**

Less efficient instructions:
```asm
sdiv x0, x1, x2      // Division is very expensive (20-100 cycles)
udiv x0, x1, x2      // Unsigned division is also expensive

rbit x0, x1          // Reverse bits - can be expensive
rev  x0, x1          // Reverse bytes - consider if necessary

madd x0, x1, x2, x3  // Fused multiply-add can be slower than separate mul/add
                     // on some implementations
```
Better alternatives when possible:
```asm
lsr/asr/lsl          // Use shifts instead of division by powers of 2
mov/add              // Simple instructions are usually faster
```

### **Variable shifts can be expensive:**
Less efficient:
```asm
lsl x0, x1, x2       // Variable shift using register
```
More efficient if possible:
```asm
lsl x0, x1, #3       // Immediate shift by constant
```

### **Some memory access patterns to avoid:**
Less efficient:
```asm
ldp x0, x1, [x2, #8]!    // Pre-index can be slower
ldr x0, [x1, x2, lsl #3] // Complex addressing modes
```
More efficient:
```asm
ldp x0, x1, [x2], #8     // Post-index often better
ldr x0, [x1, #8]         // Simple offset
```

## ARM64-specific considerations:

### Unlike x86's CPUID, ARM64 has several ways to identify processor features:
```asm
mrs x0, MIDR_EL1         // Read Main ID Register
mrs x0, ID_AA64ISAR0_EL1 // Read Instruction Set Attributes
mrs x0, ID_AA64MMFR0_EL1 // Read Memory Model Features
```

### Different ARM implementations (Apple M1, Cortex-A76, etc.) have different performance characteristics. What's slow on one might be fast on another.

### General guidelines:
   - Prefer simple addressing modes
   - Avoid complex bit manipulation when possible
   - Use vector/SIMD instructions (NEON) for bulk operations
   - Avoid division operations when possible (use shifts for powers of 2)

### Modern considerations:
Avoid when possible:
```asm
dmb sy               // Full memory barrier - very expensive
dsb sy               // Data synchronization barrier
isb                  // Instruction synchronization barrier
```

Use more specific barriers if needed:
```asm
dmb ishld           // More specific barrier can be faster
dmb oshld           // Load-only barrier when appropriate
```
Understanding these characteristics can help in writing efficient code, especially for performance-critical sections.

[gARM - full instruction list](../../pkg/op/op.go)

[NEXT -> functions](4_function.md)

<div align="center">
  <img src="../img/argo-mascot.jpg" alt="Logo">
</div>
<p align="center">
		<img src="https://raw.githubusercontent.com/bornmay/bornmay/Update/svg/Bottom.svg" alt="Github Stats" />
</p>
<p align="right">(<a href="#top">back to top</a>)</p>