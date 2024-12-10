
### Memory

Load, Store and Addressing Modes

load a byte from x1
```asm
  ldrb    w0, [x1]
```
load a signed byte from x1
```asm
  ldrsb   w0, [x1]
```
store a 32-bit word to address in x1
```asm
  str     w0, [x1]
```
load two 32-bit words from stack, advance sp by 8
```asm
  ldp     w0, w1, [sp], 8
```
store two 64-bit words at [sp-96] and subtract 96 from sp
```asm
  stp     x0, x1, [sp, -96]!
```
load 32-bit immediate from literal pool
```asm
  ldr     w0, =0x12345678
```
#### Base register only

load a byte from x1
```asm
  ldrb   w0, [x1]
```
load a half-word from x1
```asm
  ldrh   w0, [x1]
```
load a word from x1
```asm
  ldr    w0, [x1]
```
load a doubleword from x1
```asm
  ldr    x0, [x1]
```
#### Base register plus offset
load a byte from x1 plus 1
```asm
  ldrb   w0, [x1, 1]
```
load a half-word from x1 plus 2
```asm
  ldrh   w0, [x1, 2]
```
load a word from x1 plus 4
```asm
  ldr    w0, [x1, 4]
```
load a doubleword from x1 plus 8
```asm
  ldr    x0, [x1, 8]
```
load a doubleword from x1 using x2 as index
w2 is multiplied by 8
```asm
  ldr    x0, [x1, x2, lsl 3]
```
load a doubleword from x1 using w2 as index
w2 is zero-extended and multiplied by 8
```asm
  ldr    x0, [x1, w2, uxtw 3]
```
#### Pre-index
The exclamation mark “!” implies adding the offset after the load or store.

load a byte from x1 plus 1, then advance x1 by 1
```asm
  ldrb   w0, [x1, 1]!
```
load a half-word from x1 plus 2, then advance x1 by 2
```asm
  ldrh   w0, [x1, 2]!
```
load a word from x1 plus 4, then advance x1 by 4
```asm
  ldr    w0, [x1, 4]!
```
load a doubleword from x1 plus 8, then advance x1 by 8
```asm
  ldr    x0, [x1, 8]!
```
#### Post-index
This mode accesses the value first and then adds the offset to base.

load a byte from x1, then advance x1 by 1
```asm
  ldrb   w0, [x1], 1
```
load a half-word from x1, then advance x1 by 2
```asm
  ldrh   w0, [x1], 2
```
load a word from x1, then advance x1 by 4
```asm
  ldr    w0, [x1], 4
```
load a doubleword from x1, then advance x1 by 8
```asm
  ldr    x0, [x1], 8
```
Literal (PC-relative)
These instructions work similar to RIP-relative addressing on AMD64.

load address of label
```asm
  adr    x0, label
```
load address of label
```asm
  adrp   x0, label
```

### Bit Manipulation

load 32-bit immediate value, we do the following.
```asm
  movl    w0, 0x12345678
```
Move 0x12345678 into w0
```asm
    mov     w0, 0x5678
    mov     w1, 0x1234
    bfi     w0, w1, 16, 16
```
Extract 8-bits from x1 into the x0 register at position 0
```asm
    ubfx    x0, x1, 8, 8  // If x1 is 0x12345678, 0x00000056 is placed in x0
```
Extract 8-bits from x1 and insert with zeros into the x0 register at position 8
```asm
    ubfiz   x0, x1, 8, 8  // If x1 is 0x12345678, 0x00005600 is placed in x0
```
Extract 8-bits from x1 and insert into x0 at position 0
```asm
    bfxil   x0, x1, 0, 8  // if x1 is 0x12345678 and x0 is 0x09ABCDEF. x0 after execution has 0x09ABCD78
```
Clear lower 8 bits
```asm
    bfxil   x0, xzr, 0, 8
```
Zero-extend 8-bits
```asm
    uxtb    x0, x0
```

[NEXT -> registers](7_register.md)

<div align="center">
  <img src="../img/argo-mascot.jpg" alt="Logo">
</div>
<p align="center">
    <img src="https://raw.githubusercontent.com/bornmay/bornmay/Update/svg/Bottom.svg" alt="Github Stats" />
</p>
<p align="right">(<a href="#top">back to top</a>)</p>
