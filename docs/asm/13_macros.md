## Macros

#### load the 32-bit number 0x12345678 into register w0
```asm
  movz    w0, 0x5678
  movk    w0, 0x1234, lsl 16
```

#### load a 64-bit immediate using MOV
```asm
  .macro movq Xn, imm
      movz    \Xn,  \imm & 0xFFFF
      movk    \Xn, (\imm >> 16) & 0xFFFF, lsl 16
      movk    \Xn, (\imm >> 32) & 0xFFFF, lsl 32
      movk    \Xn, (\imm >> 48) & 0xFFFF, lsl 48
  .endm
```

#### load a 32-bit immediate using MOV
```asm
  .macro movl Wn, imm
      movz    \Wn,  \imm & 0xFFFF
      movk    \Wn, (\imm >> 16) & 0xFFFF, lsl 16
  .endm
```

The first instruction MOVZ loads 0x5678 into w0, zero extending to 32-bits. MOVK loads 0x1234 into the upper 16-bits using a shift, while preserving the lower 16-bits. Some assemblers provide a pseudo-instruction called MOVL that expands into the two instructions above. However, the GNU Assembler doesn’t recognize it, so here are two macros for GAS that can load a 32-bit or 64-bit immediate value into a general purpose register.

#### load a 64-bit immediate using MOV
```asm
  .macro movq Xn, imm
      movz    \Xn,  \imm & 0xFFFF
      movk    \Xn, (\imm >> 16) & 0xFFFF, lsl 16
      movk    \Xn, (\imm >> 32) & 0xFFFF, lsl 32
      movk    \Xn, (\imm >> 48) & 0xFFFF, lsl 48
  .endm
```

#### load a 32-bit immediate using MOV
```asm
  .macro movl Wn, imm
      movz    \Wn,  \imm & 0xFFFF
      movk    \Wn, (\imm >> 16) & 0xFFFF, lsl 16
  .endm
```
Then if we need to load a 32-bit immediate value, we do the following
```asm
  movl    w0, 0x12345678
```

#### The famous ternary 

```asm
cmp w0, #0
csel w0, w1, w2, eq       // w0 = (w0 == 0) ? w1 : w2
```

turn it into a macro

```asm
.macro SELECT_IF_ZERO dst, test, true_val, false_val
    cmp     \test, #0
    csel    \dst, \true_val, \false_val, eq
.endm
```
and use it like so
```asm
SELECT_IF_ZERO w0, w0, w1, w2    // w0 = (w0 == 0) ? w1 : w2
```
This macro takes 4 parameters: destination register, test register, value if true, value if false. Makes the comparison with zero and
performs the conditional select based on equality.

Now let's make it more flexible and compare against any value, not just zero to make it feel more like a C-style ternary operator (condition ? true_val : false_val)

```asm
.macro TERNARY dst_cond, question, true_val, false_val
    cmp     \dst_cond, \question    // condition ?
    csel    \dst_cond, \true_val, \false_val, eq   // true_val : false_val
.endm
```
use it like so

```asm
TERNARY w0, #0, w1, w2    // w0 = (w0 == 0) ? w1 : w2
TERNARY w0, w1, w2, w3    // w0 = (w0 == w1) ? w2 : w3
```
Much cleaner! The first parameter acts as the value being tested and the destination for result, just like a typical ternary expression.

[NEXT -> syscalls](14_syscall.md)

<div align="center">
  <img src="../img/argo-mascot.jpg" alt="Logo">
</div>
<p align="center">
    <img src="https://raw.githubusercontent.com/bornmay/bornmay/Update/svg/Bottom.svg" alt="Github Stats" />
</p>
<p align="right">(<a href="#top">back to top</a>)</p>
