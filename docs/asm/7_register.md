# Registers
**Register Availability**: The need for register packing is less crucial in ARM64 because:
   - 31 general-purpose registers (x0-x30) compared to x86's 8
   - Each register can be accessed as:
     - X0-X30 (64-bit)
     - W0-W30 (32-bit)
     - Various SIMD arrangements

Most **System Calls** only use general-purpose registers.
```bash
Wn	32-bit	GP 0-31
Xn	64-bit	GP 0-31
WZR	32-bit	Zero register
XZR	64-bit	Zero register
SP	64-bit	Stack pointer
```
The alignment of sp must be two times the size of a pointer (16 bytes)

# Condition Flags

ARM has a “process state” with condition flags that affect the behaviour of some instructions. Branch instructions can be used to change the flow of execution. Some of the data processing instructions allow setting the condition flags with the S suffix. e.g ANDS or ADDS. The flags are the Zero Flag (Z), the Carry Flag (C), the Negative Flag (N) and the is Overflow Flag (V).

- `N`	Bit 31. Set if the result of an operation is negative. Cleared if the result is positive or zero.
- `Z`	Bit 30. Set if the result of an operation is zero/equal. Cleared if non-zero/not equal.
- `C`	Bit 29. Set if an instruction results in a carry or overflow. Cleared if no carry.
- `V`	Bit 28. Set if an instruction results in an overflow. Cleared if no overflow.


# Register Allocation Convention

- X registers are 64-bit (e.g., X0, X1, X2...)
- W registers are 32-bit (e.g., W0, W1, W2...)
Each W register is actually the lower 32 bits of the corresponding X register

```bash
X0:  [63...............32][31................0]
                          ↑
W0:                       [31................0]
```

```asm
mov     x0, #0x1234567890ABCDEF    // 64-bit value
mov     w0, #0x12345678            // Only affects lower 32 bits
                                   // Upper 32 bits are zeroed!
```

```bash
X0 before:  0x1234567890ABCDEF
W0 before:  0x90ABCDEF            (lower 32 bits of X0)

X0 after:   0x0000000012345678    (upper bits zeroed)
W0 after:   0x12345678            (same as lower 32 bits of X0)
```

### Writing to a W register:

- Only affects lower 32 bits
- Automatically zeros the upper 32 bits of the X register
- Good for 32-bit arithmetic

### Writing to an X register:

- Affects all 64 bits
- Can access via W register for lower 32 bits

```asm
// For byte operations (8-bit)
ldrb    w0, [x1]        // Load a single byte into W0 (zeros upper bits)

// For 32-bit operations
add     w0, w1, w2      // 32-bit addition
mul     w0, w1, w2      // 32-bit multiplication

// For 64-bit operations
add     x0, x1, x2      // 64-bit addition
mul     x0, x1, x2      // 64-bit multiplication
```

### Use W registers when:

- Working with 32-bit or smaller values
- Doing 32-bit arithmetic
- Loading/storing bytes or 32-bit values
- Working with legacy 32-bit code


### Use X registers when:

- Working with addresses (always 64-bit)
- Doing 64-bit arithmetic
- Handling 64-bit values
- Working with pointers

```asm
// Loading different sizes
ldrb    w0, [x1]        // Load byte (8-bit)
ldrh    w0, [x1]        // Load halfword (16-bit)
ldr     w0, [x1]        // Load word (32-bit)
ldr     x0, [x1]        // Load doubleword (64-bit)

// Storing different sizes
strb    w0, [x1]        // Store byte (8-bit)
strh    w0, [x1]        // Store halfword (16-bit)
str     w0, [x1]        // Store word (32-bit)
str     x0, [x1]        // Store doubleword (64-bit)
```

## Parameter/Result Registers (Caller-saved)
- x0-x7: Parameter/scratch registers
- x0: First parameter and return value
- x1: Second parameter
- x2: Third parameter
- x3: Fourth parameter
- x4: Fifth parameter
- x5: Sixth parameter
- x6: Seventh parameter
- x7: Eighth parameter

## Caller-saved Registers (Temporary)
- x8: Indirect result location register
- x9-x15: Temporary registers
- Can be used without saving
- Must be saved by caller if needed across calls

## Callee-saved Registers
- x19-x28: Must be preserved across function calls
- Must save and restore if you use them
- Often used for local variables needed across function calls

## Special Purpose Registers
- x29 (fp): Frame pointer
- x30 (lr): Link register (return address)
- sp: Stack pointer
- x16, x17: Intra-procedure-call temporary registers
- x18: Platform register (reserved)

## Example Register Selection Decision Tree:
1. For parameters:
   - First 8 parameters go in x0-x7
   - Additional parameters go on stack

2. For return values:
   - Single values use x0
   - Pairs use x0 and x1
   - Larger returns use x8 (indirect return)

3. For local variables:
   - If needed across function calls: Use x19-x28
   - If temporary within function: Use x9-x15
   - If very temporary (few instructions): Use x0-x7

4. For string operations:
   - Source pointer often in x0
   - Destination pointer often in x1
   - Length/size often in x2
   - Loop counters often in x9-x15

## Key principles:

- Follow ABI convention for parameters (x0-x7)
- Use caller-saved registers (x9-x15) for temporary values
- Use callee-saved registers (x19-x28) for values needed across function calls
- Always preserve and restore any callee-saved registers you use
- Return values in x0 (and x1 if needed)

### Ensures:

- Predictable function interfaces
- Proper interoperability with other code
- Efficient register usage
- Correct behavior when functions call other functions

For generating efficient code with multiple function calls, we should use a register allocation strategy:

1. First, we analyze "live ranges" - how long values need to be preserved
2. Then allocate registers based on overlapping lifetimes
3. Use caller-saved (x0-x18) for short-lived values
4. Reserve callee-saved (x19-x28) for values that span function calls

Here's an example of the thinking process:

```asm
// Example Go code:
func main() {
    a := 1        // Short lived, goes directly to x0
    b := 2        // Short lived, goes directly to x1
    c := add(a,b) // Result needed later, store in x19
    
    d := 3        // Short lived, goes directly to x0
    e := 4        // Short lived, goes directly to x1
    f := add(d,e) // Result needed later, store in x20
    
    g := add(c,f) // Use x19,x20 as inputs via x0,x1
}

// Resulting assembly:
start:
    stp x29, x30, [sp, #-16]!
    mov x29, sp
    stp x19, x20, [sp, #-16]!   // Only save registers we'll actually use

    // First add - results needed later
    mov x0, #1                   // Immediate to param reg
    mov x1, #2                   // Immediate to param reg
    bl add
    mov x19, x0                  // Save result - spans function call

    // Second add - results needed later
    mov x0, #3                   // Reuse x0
    mov x1, #4                   // Reuse x1
    bl add
    mov x20, x0                  // Save result - spans function call

    // Final add - using preserved values
    mov x0, x19                  // Load preserved value
    mov x1, x20                  // Load preserved value
    bl add                       // Result in x0, not saved as it's final

    ldp x19, x20, [sp], #16     // Restore what we saved
    ldp x29, x30, [sp], #16
    ret

```

Key strategies for a code generator:

1. Build a graph of value lifetimes:
```
Value   Live Range
a       birth -> first add
b       birth -> first add
c       first add -> final add
d       birth -> second add
e       birth -> second add
f       second add -> final add
```

2. Register allocation rules:
```
- Values only needed for immediate function call:
  → Use parameter registers (x0-x7)

- Values needed across function calls:
  → Allocate callee-saved register (x19-x28)
  → Save/restore at function boundaries

- Registers can be reused when values are no longer needed
  → x0,x1 reused between calls
  → x19,x20 only used for cross-call values
```

3. Code generation phases:
```
a) Analyze lifetimes
b) Count max simultaneous live values
c) Allocate callee-saved registers for long-lived values
d) Use caller-saved registers for temporaries
e) Generate prologue to save used callee-saved regs
f) Generate epilogue to restore them
```

4. Optimization opportunities:
```
- Immediate → parameter register when possible
- Reuse parameter registers between calls
- Only save/restore registers actually used
- Keep values in caller-saved regs if lifetime permits
```

This strategy minimizes:
- Number of registers used
- Stack operations
- Register moves

The key is tracking value lifetimes and using the appropriate register class (caller vs callee saved) based on whether values need to survive function calls.

```asm
func:
    // Save FP/LR
    stp x29, x30, [sp, #-16]!
    mov x29, sp
    
    // If using callee-saved, save them
    stp x19, x20, [sp, #-16]!
    
    // Can freely use x0-x7, x9-x15 for temporary work
    // Use x19-x28 for values needed across calls
    
    // Restore in reverse order
    ldp x19, x20, [sp], #16
    ldp x29, x30, [sp], #16
    ret
```

# Register instructions

### Addressing modes
memory access instructions commonly support three addressing modes:

#### Offset addressing 
An offset is applied to an address from a base register and the result is used to perform memory access `[rN, offset]`

#### Pre-indexed addressing 
An offset is applied to an address from a base register, the result is used to perform memory access and written back to the base register. It looks like `[rN, offset]!`
The exclamation mark “!” implies adding the offset after the load/store.
```asm
// load a byte from x1 plus 1, then advance x1 by 1
ldrb   w0, [x1, 1]!

// load a half-word from x1 plus 2, then advance x1 by 2
ldrh   w0, [x1, 2]!

// load a word from x1 plus 4, then advance x1 by 4
ldr    w0, [x1, 4]!

// load a doubleword from x1 plus 8, then advance x1 by 8
ldr    x0, [x1, 8]!
```
#### Post-indexed addressing 
An address is used as-is from a base register for memory access. The offset is applied and the result is stored back to the base register. It looks like `[rN], offset`
This mode accesses the value first and then adds the offset to base.
```asm
  // load a byte from x1, then advance x1 by 1
  ldrb   w0, [x1], 1

  // load a half-word from x1, then advance x1 by 2
  ldrh   w0, [x1], 2

  // load a word from x1, then advance x1 by 4
  ldr    w0, [x1], 4

  // load a doubleword from x1, then advance x1 by 8
  ldr    x0, [x1], 8
```
### Saving Registers
19 registers are free to use without having to preserve them for the caller. Compare to x86 where only 3 registers are available or 5 on AMD64. ARM stores the return address in the Link Register (LR) which is an alias for X30 register. A callee is expected to save LR/X30 if it calls a subroutine
```asm
    // push {x0}
    // [base - 16] = x0
    // base = base - 16
    str    x0, [sp, -16]!

    // pop {x0}
    // x0 = [base]
    // base = base + 16
    ldr    x0, [sp], 16

    // push {x0, x1}
    stp    x0, x1, [sp, -16]!

    // pop {x0, x1}
    ldp    x0, x1, [sp], 16
```
The 16 in stp x29, x30, [sp, #-16]! is about stack space:

- pushes two 8-byte registers (x29 and x30)
- 8 bytes × 2 registers = 16 bytes total
- stack must be aligned as unaligned access can cause exceptions
- push and pop instructions have been deprecated in favour of load and store, use STR/STP to save and LDR/LDP to store. Here’s how you can save/restore registers using the stack

### Copying Registers
```asm
    // Move x1 to x0
    mov     x0, x1

    // Extract bits 0-63 from x1 and store in x0 zero extended.
    ubfx   x0, x1, 0, 63

    // x0 = (x1 & ~0)
    bic    x0, x1, xzr

    // x0 = x1 >> 0
    lsr    x0, x1, 0

    // Use a circular shift (rotate) to move x1 to x0
    ror    x0, x1, 0
    
    // Extract bits 0-63 from x1 and insert into x0
    bfxil  x0, x1, 0, 63
```
###  Init register to zero
Initialize a counter “i = 0” or pass NULL/0 to a system call like so
```asm
    // Move an immediate value of zero into the register.
    mov    x0, 0

    // Copy the zero register.
    mov    x0, xzr

    // Exclusive-OR the register with itself.
    eor    x0, x0, x0

    // Subtract the register from itself.
    sub    x0, x0, x0

    // Mask the register with zero register using a bitwise AND.
    // An immediate value of zero will work here too.
    and    x0, x0, xzr

    // Multiply the register by the zero register.
    mul    x0, x0, xzr

    // Extract 64 bits from xzr and place in x0.
    bfxil  x0, xzr, 0, 63
    
    // Circular shift (rotate) right.
    ror    x0, xzr, 0

    // Logical shift right.
    lsr    x0, xzr, 0
    
    // Reverse bytes of zero register.
    rev    x0, xzr
```
### Init register to 1.
Rarely starts a counter at 1, but it’s common enough
```asm
    // Move 1 into x0.
    mov     x0, 1

    // Compare x0 with x0 and set x0 if equal.
    cmp     x0, x0
    cset    x0, eq

    // Bitwise NOT the zero register and store in x0. Negate x0.
    mvn     x0, xzr
    neg     x0, x0
```
### Init register to -1.
Some sys calls require this
```asm
    // move -1 into register
    mov     x0, -1

    // copy the zero register inverted
    mvn     x0, xzr

    // x0 = ~(x0 ^ x0)
    eon     x0, x0, x0

    // x0 = (x0 | ~xzr)
    orn     x0, x0, xzr

    // x0 = (int)0xFF
    mov     w0, 255
    sxtb    x0, w0

    // x0 = (x0 == x0) ? -1 : x0
    cmp     x0, x0
    csetm   x0, eq
```
### Init register to 0x80000000.
might seem vague, but crypto/X25519 uses this value for reduction step
```asm
    mov     w0, 0x80000000

    // Set bit 31 of w0.
    mov     w0, 1
    mov     w0, w0, lsl 31

    // Set bit 31 of w0.
    mov     w0, 1
    ror     w0, w0, 1

    // Set bit 31 of w0.
    mov     w0, 1
    rbit    w0, w0

    // Set bit 31 of w0.
    eon     w0, w0, w0
    lsr     w0, w0, 1
    add     w0, w0, 1
    
    // Set bit 31 of w0.
    mov     w0, -1
    extr    w0, w0, wzr, 1
```
### Check for 1/TRUE.
some ways to test for equality
```asm
    // Compare x0 with 1, branch if equal.
    cmp     x0, 1
    beq     true

    // Compare x0 with zero register, branch if not equal.
    cmp     x0, xzr
    bne     true
    
    // Subtract 1 from x0 and set flags. Branch if equal. (Z flag is set)
    subs    x0, x0, 1
    beq     true

    // Negate x0 and set flags. Branch if x0 is negative.
    negs    x0, x0
    bmi     true

    // Conditional branch if x0 is not zero.
    cbnz    x0, true

    // Test bit 0 and branch if not zero.
    tbnz    x0, 0, true
```

### Check for 0/FALSE.
```asm
    // x0 == 0
    cmp     x0, 0
    beq     false

    // x0 == 0
    cmp     x0, xzr
    beq     false

    ands    x0, x0, x0
    beq     false

    // same as ANDS, but discards result
    tst     x0, x0
    beq     false

    // x0 == -0
    negs    x0
    beq     false

    // (x0 - 1) == -1
    subs    x0, x0, 1
    bmi     false

    // if (!x0) goto false
    cbz     x0, false

    // if (!x0) goto false
    tbz     x0, 0, false
```
### Check for -1
Some functions will return a negative number like -1 to indicate failure. CMN is used in the first example. This behaves exactly like CMP, except it is adding the source value (register or immediate) to the destination register, setting the flags and discarding the result.
```asm
    // w0 == -1
    cmn     w0, 1
    beq     failed

    // w0 == 0
    cmn     w0, wzr
    bmi     failed

    // negative?
    ands    w0, w0, w0
    bmi     failed

    // same as AND, but discards result
    tst     w0, w0
    bmi     failed

    // w0 & 0x80000000
    tbz     w0, 31, failed
```
```asm
// Save registers (using stack rather than SIMD since ARM64 has plenty of GP registers)
stp x29, x30, [sp, #-16]!    // Save frame pointer and link register
stp x19, x20, [sp, #-16]!    // Save callee-saved registers
stp x21, x22, [sp, #-16]!    // Save more if needed

// Your code here using registers
// x0-x18 are caller-saved, free to use
// x19-x28 are callee-saved, we saved them above
// x29 (frame pointer) and x30 (link register) were saved

// Restore registers
ldp x21, x22, [sp], #16      // Restore in reverse order
ldp x19, x20, [sp], #16
ldp x29, x30, [sp], #16
ret
```

[NEXT -> immediates](8_immediate.md)

<div align="center">
    <img src="../img/argo-mascot.jpg" alt="Logo">
</div>
<p align="center">
    <img src="https://raw.githubusercontent.com/bornmay/bornmay/Update/svg/Bottom.svg" alt="Github Stats" />
</p>
<p align="right">(<a href="#top">back to top</a>)</p>