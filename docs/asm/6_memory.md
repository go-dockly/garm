
### Memory

Load, Store and Addressing Modes

load a byte from x1
```arm
  ldrb    w0, [x1]
```
load a signed byte from x1
```arm
  ldrsb   w0, [x1]
```
store a 32-bit word to address in x1
```arm
  str     w0, [x1]
```
load two 32-bit words from stack, advance sp by 8
```arm
  ldp     w0, w1, [sp], 8
```
store two 64-bit words at [sp-96] and subtract 96 from sp
```arm
  stp     x0, x1, [sp, -96]!
```
load 32-bit immediate from literal pool
```arm
  ldr     w0, =0x12345678
```
#### Base register only

load a byte from x1
```arm
  ldrb   w0, [x1]
```
load a half-word from x1
```arm
  ldrh   w0, [x1]
```
load a word from x1
```arm
  ldr    w0, [x1]
```
load a doubleword from x1
```arm
  ldr    x0, [x1]
```
#### Base register plus offset
load a byte from x1 plus 1
```arm
  ldrb   w0, [x1, 1]
```
load a half-word from x1 plus 2
```arm
  ldrh   w0, [x1, 2]
```
load a word from x1 plus 4
```arm
  ldr    w0, [x1, 4]
```
load a doubleword from x1 plus 8
```arm
  ldr    x0, [x1, 8]
```
load a doubleword from x1 using x2 as index
w2 is multiplied by 8
```arm
  ldr    x0, [x1, x2, lsl 3]
```
load a doubleword from x1 using w2 as index
w2 is zero-extended and multiplied by 8
```arm
  ldr    x0, [x1, w2, uxtw 3]
```
#### Pre-index
The exclamation mark “!” implies adding the offset after the load or store.

load a byte from x1 plus 1, then advance x1 by 1
```arm
  ldrb   w0, [x1, 1]!
```
load a half-word from x1 plus 2, then advance x1 by 2
```arm
  ldrh   w0, [x1, 2]!
```
load a word from x1 plus 4, then advance x1 by 4
```arm
  ldr    w0, [x1, 4]!
```
load a doubleword from x1 plus 8, then advance x1 by 8
```arm
  ldr    x0, [x1, 8]!
```
#### Post-index
This mode accesses the value first and then adds the offset to base.

load a byte from x1, then advance x1 by 1
```arm
  ldrb   w0, [x1], 1
```
load a half-word from x1, then advance x1 by 2
```arm
  ldrh   w0, [x1], 2
```
load a word from x1, then advance x1 by 4
```arm
  ldr    w0, [x1], 4
```
load a doubleword from x1, then advance x1 by 8
```arm
  ldr    x0, [x1], 8
```
Literal (PC-relative)
These instructions work similar to RIP-relative addressing on AMD64.

load address of label
```arm
  adr    x0, label
```
load address of label
```arm
  adrp   x0, label
```

### Bit Manipulation

load 32-bit immediate value, we do the following.
```arm
  movl    w0, 0x12345678
```
Move 0x12345678 into w0
```arm
    mov     w0, 0x5678
    mov     w1, 0x1234
    bfi     w0, w1, 16, 16
```
Extract 8-bits from x1 into the x0 register at position 0
```arm
    ubfx    x0, x1, 8, 8 @ If x1 is 0x12345678, 0x00000056 is placed in x0
```
Extract 8-bits from x1 and insert with zeros into the x0 register at position 8
```arm
    ubfiz   x0, x1, 8, 8 @ If x1 is 0x12345678, 0x00005600 is placed in x0
```
Extract 8-bits from x1 and insert into x0 at position 0
```arm
    bfxil   x0, x1, 0, 8 @ if x1 is 0x12345678 and x0 is 0x09ABCDEF. x0 after execution has 0x09ABCD78
```
Clear lower 8 bits
```arm
    bfxil   x0, xzr, 0, 8
```
Zero-extend 8-bits
```arm
    uxtb    x0, x0
```
#### Macros

load the 32-bit number 0x12345678 into register w0
```arm
  movz    w0, 0x5678
  movk    w0, 0x1234, lsl 16
```

load a 64-bit immediate using MOV
```arm
  .macro movq Xn, imm
      movz    \Xn,  \imm & 0xFFFF
      movk    \Xn, (\imm >> 16) & 0xFFFF, lsl 16
      movk    \Xn, (\imm >> 32) & 0xFFFF, lsl 32
      movk    \Xn, (\imm >> 48) & 0xFFFF, lsl 48
  .endm
```

load a 32-bit immediate using MOV
```arm
  .macro movl Wn, imm
      movz    \Wn,  \imm & 0xFFFF
      movk    \Wn, (\imm >> 16) & 0xFFFF, lsl 16
  .endm
```

The first instruction MOVZ loads 0x5678 into w0, zero extending to 32-bits. MOVK loads 0x1234 into the upper 16-bits using a shift, while preserving the lower 16-bits. Some assemblers provide a pseudo-instruction called MOVL that expands into the two instructions above. However, the GNU Assembler doesn’t recognize it, so here are two macros for GAS that can load a 32-bit or 64-bit immediate value into a general purpose register.

load a 64-bit immediate using MOV
```arm
  .macro movq Xn, imm
      movz    \Xn,  \imm & 0xFFFF
      movk    \Xn, (\imm >> 16) & 0xFFFF, lsl 16
      movk    \Xn, (\imm >> 32) & 0xFFFF, lsl 32
      movk    \Xn, (\imm >> 48) & 0xFFFF, lsl 48
  .endm
```

load a 32-bit immediate using MOV
```arm
  .macro movl Wn, imm
      movz    \Wn,  \imm & 0xFFFF
      movk    \Wn, (\imm >> 16) & 0xFFFF, lsl 16
  .endm
```
Then if we need to load a 32-bit immediate value, we do the following
```arm
  movl    w0, 0x12345678
```
[src: AS Compiler walk through](https://modexp.wordpress.com/2018/10/30/arm64-assembly/)

[NEXT -> registers](register.md)

<div align="center">
  <img src="../img/argo-mascot.jpg" alt="Logo">
</div>
<p align="center">
    <img src="https://raw.githubusercontent.com/bornmay/bornmay/Update/svg/Bottom.svg" alt="Github Stats" />
</p>
<p align="right">(<a href="#top">back to top</a>)</p>
