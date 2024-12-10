### Immediates
```asm
.data
    number: .word 42    // integer constant

.text
.align 2
.global main
```
Creates a label called number
- Reserves 4 bytes of memory
- Initializes 4 bytes with the value 42
- Stores bytes in little-endian format (most modern systems)

.align 2 is about instruction alignment in memory

- means align to 2^2 = 4 bytes
- required for ARM instructions which must be 4-byte aligned


```bash
.byte   // 1  byte  (8 bits)   Range: -128 to 127 or 0 to 255
// half
.hword  // 2  bytes (16 bits)  Range: -32,768 to 32,767 or 0 to 65,535
.word   // 4  bytes (32 bits)  Range: -2^31 to 2^31-1 or 0 to 2^32-1
// double
.dword  // 8  bytes (64 bits)  Range: -2^63 to 2^63-1 or 0 to 2^64-1
// quad
.qword  // 16 bytes (128 bits) Range: -2^127 to 2^127-1 or 0 to 2^128-1
```
Looking at the actual bytes in memory, 42 would be stored as:

```asm
2A 00 00 00
```
In arm64, memory addresses are 64-bit values. Typically loaded in two parts because:

- A single instruction can only hold a limited size immediate value
- Programs use relative addressing for position-independent code (PIC)

```asm
adrp    x0, number               // Loads upper bits (page address)
add     x0, x0, :lo12:number     // Adds lower 12 bits offset
```
adrp loads the page address (think of it as the "neighborhood")

- A page is 4KB (4096 bytes), so it aligns to 4KB boundaries
- It loads the address with the bottom 12 bits set to zero

:lo12: is a relocation operator that

- Takes the full address of `number`
- Extracts just the bottom 12 bits (the `house number` in our neighborhood)
- These bits represent the offset within the 4KB page

If our data was at address 0x1234A678

- adrp would load 0x12340000 (page address)
- :lo12: would give 0x678 (offset)
- add combines them to the full address

This instruction sequence is necessary because ARM can't load a full 64-bit address in a single instruction.

Functions typically share the same .data and .text sections within a single assembly file. The sections are file-level (or program-level) constructs, not function-level.

```asm
.data
    number1: .word 42    // Shared data section
    number2: .word 100
    
.text                    // Shared text section
.align 2
.global fn1
fn1:
    stp x29, x30, [sp, #-16]!
    // load number1 ...
    ldp x29, x30, [sp], #16
    ret

.align 2
.global fn2
fn2:
    stp x29, x30, [sp, #-16]!
    // load number2 ...
    ldp x29, x30, [sp], #16
    ret
```

[NEXT -> conditionals](9_conditional.md)

<div align="center">
	<img src="../img/argo-mascot.jpg" alt="Logo">
</div>
<p align="center">
	<img src="https://raw.githubusercontent.com/bornmay/bornmay/Update/svg/Bottom.svg" alt="Github Stats" />
</p>
<p align="right">(<a href="#top">back to top</a>)</p>