
### Arithmetic

x0 == -1?
```asm
  cmn     x0, 1
  beq     minus_one
```

x0 == 0
```asm
  cmp     x0, 0
  beq     zero

allocate 32 bytes of stack
```asm
  sub     sp, sp, 32
```

x0 = x0 % 37
```asm
  mov     x1, 37
  udiv    x2, x0, x1
  msub    x0, x2, x1, x0
```

x0 = 0
```asm
  sub     x0, x0, x0
```

### Logical

Multiplication can be performed using logical shift left LSL. Division can be performed using logical shift right LSR. Modulo operations can be performed using bitwise AND. The only condition is that the multiplier and divisor be a power of two. The first three examples shown here demonstrate those operations.

x1 = x0 / 8
```asm
  lsr     x1, x0, 3
```

x1 = x0 * 4
```asm
  lsl     x1, x0, 2
```

x1 = x0 % 16
```asm
  and     x1, x0, 15
```

x0 == 0?
```asm
  tst     x0, x0
  beq     zero
```

x0 = 0
```asm
  eor     x0, x0, x0
```

### Bit Manipulation
ARM64 provides several instructions for bit manipulation:
- CLZ - Count Leading Zeros
- CTZ/RBIT+CLZ - Count Trailing Zeros
- CLS - Count Leading Sign bits

1. Find highest set bit (equivalent to BSR)
```asm
find_highest_set_bit:
    // Input in x0, result in x0
    clz     x1, x0              // Count leading zeros
    mov     x2, #63             // For 64-bit (use 31 for 32-bit)
    sub     x0, x2, x1          // Convert to bit position
    // x0 now contains the position of highest set bit (0-63)
    ret
```
2. Find highest power of 2 less than or equal to input
```asm
highest_power_of_2:
    // Input in x0, result in x0
    clz     x1, x0              // Count leading zeros
    mov     x2, #63             // Maximum bit position
    sub     x1, x2, x1          // Get highest bit position
    mov     x0, #1              // Start with 1
    lsl     x0, x0, x1          // Shift to create power of 2
    ret
```
3. Round up to next power of 2
```asm
round_up_power_2:
    // Input in x0, result in x0
    sub     x1, x0, #1          // Handle case when input is already power of 2
    clz     x2, x1              // Count leading zeros
    mov     x3, #63             // For 64-bit
    sub     x2, x3, x2          // Get position
    mov     x0, #1              // Start with 1
    add     x2, x2, #1          // Add 1 for next power
    lsl     x0, x0, x2          // Shift to create power of 2
    ret
```
4. Check if number is power of 2
```asm
is_power_of_2:
    // Input in x0, result in x0 (1 if power of 2, 0 if not)
    sub     x1, x0, #1          // Subtract 1
    and     x0, x0, x1          // AND with original
    cmp     x0, #0              // If result is 0, was power of 2
    cset    x0, eq              // Set result based on comparison
    ret
```
5. Find nearest power of 2 (rounding to nearest)
```asm
nearest_power_of_2:
    // Input in x0, result in x0
    clz     x1, x0              // Count leading zeros
    mov     x2, #63             // For 64-bit
    sub     x1, x2, x1          // Get position of highest bit
    mov     x3, #1
    lsl     x3, x3, x1          // Higher power of 2
    lsr     x4, x3, #1          // Lower power of 2

    // Find which is closer
    sub     x5, x0, x4          // Distance to lower
    sub     x6, x3, x0          // Distance to higher
    cmp     x5, x6
    csel    x0, x3, x4, gt      // Select closer power
    ret
```
6. Efficient byte reverse using bit scanning
```asm
byte_reverse:
    // Input in x0, result in x0
    rbit    x0, x0              // Reverse all bits
    lsr     x0, x0, #56         // Shift right to align byte
    ret
```
7. Find log base 2 of a number (floor)
```asm
log2_floor:
    // Input in x0, result in x0
    clz     x1, x0              // Count leading zeros
    mov     x2, #63             // For 64-bit
    sub     x0, x2, x1          // Convert to log2
    ret
```
Key advantages of ARM64's CLZ compared to x86's BSR:

1. **Performance**:
   - CLZ is typically single-cycle on ARM64
   - More predictable than x86's BSR
   - No flags affected unless explicitly requested

2. **Flexibility**:
   - Works on both 32-bit and 64-bit values
   - Can be combined with other bit manipulation instructions
   - RBIT (Reverse Bits) provides additional functionality

3. **Additional Features**:
   - CLS (Count Leading Sign bits) for signed values
   - RBIT for full bit reversal
   - Conditional selection (CSEL) for branching optimization

4. **Common Use Cases**:
   - Computing log base 2
   - Finding nearest power of 2
   - Implementing binary search
   - Bit field manipulation
   - Memory allocation (finding appropriate block sizes)

### Division

1. Basic unsigned division
```asm
unsigned_division:
    // Input: dividend in x0, divisor in x1
    udiv    x0, x0, x1          // Single instruction, no clearing needed!
    ret
```
2. Basic signed division
```asm
signed_division:
    // Input: dividend in x0, divisor in x1
    sdiv    x0, x0, x1          // Single instruction for signed division
    ret
```
3. Division with remainder (unsigned)
```asm
unsigned_div_with_remainder:
    // Input: dividend in x0, divisor in x1
    // Output: quotient in x0, remainder in x1
    udiv    x2, x0, x1          // Get quotient
    msub    x1, x2, x1, x0      // remainder = dividend - (quotient * divisor)
    mov     x0, x2              // Move quotient to x0
    ret
```
4. Division with remainder (signed)
```asm
signed_div_with_remainder:
    // Input: dividend in x0, divisor in x1
    // Output: quotient in x0, remainder in x1
    sdiv    x2, x0, x1          // Get quotient
    msub    x1, x2, x1, x0      // remainder = dividend - (quotient * divisor)
    mov     x0, x2              // Move quotient to x0
    ret
```
5. Optimized division by power of 2 (unsigned)
```asm
divide_by_power2_unsigned:
    // Input: dividend in x0, power in x1
    lsr     x0, x0, x1          // Simple right shift
    ret
```
6. Optimized division by power of 2 (signed)
```asm
divide_by_power2_signed:
    // Input: dividend in x0, power in x1
    asr     x0, x0, x1          // Arithmetic right shift
    ret
```
7. Division by constant (optimized using multiplication)
```asm
divide_by_constant:
    // Example: Division by 7
    // Use multiplication by magic number (2^64)/7 + 1
    movz    x1, #0x2492, lsl #48
    movk    x1, #0x4925, lsl #32
    movk    x1, #0x4925, lsl #16
    movk    x1, #0x4925
    umulh   x0, x0, x1          // Upper 64 bits of multiply
    ret
```
8. Checking for division by zero
```asm
safe_unsigned_division:
    // Input: dividend in x0, divisor in x1
    cmp     x1, #0              // Check for zero divisor
    b.eq    division_by_zero    // Handle division by zero
    udiv    x0, x0, x1          // Safe to divide
    ret
division_by_zero:
    // Handle error condition
    mov     x0, #-1             // Or appropriate error value
    ret
```
9. Batch processing multiple divisions using SIMD
```asm
simd_unsigned_division:
    // Process 4 32-bit divisions at once
    // Assuming inputs in v0 and v1
    udiv    v2.4s, v0.4s, v1.4s  // Divide 4 numbers at once
    // Result in v2
    ret
```

Key differences from x86:

1. **Separate Instructions**:
   - ARM64 has distinct `UDIV` and `SDIV` instructions
   - No need to clear registers like x86's `XOR EDX,EDX`
   - No equivalent to x86's `CDQ`/`IDIV` complexity

2. **Performance Characteristics**:
   - Division operations are typically 4-12 cycles
   - No additional setup cycles needed
   - Unsigned division (`UDIV`) is often slightly faster than signed (`SDIV`)

3. **Optimization Options**:
   - For powers of 2, use shifts (`LSR`/`ASR`)
   - For constants, use multiplication by reciprocal
   - SIMD operations available for batch processing

4. **Key Advantages**:
   - Simpler code (no register clearing needed)
   - More predictable performance
   - Better error handling options
   - SIMD support for parallel operations

#### Best practices

1. Use unsigned division when possible (`UDIV`)
2. For powers of 2, use shifts
3. For known constants, use multiplication by reciprocal
4. Consider SIMD for bulk operations
5. Always handle division by zero

### Multiply to Divide
For division by constant, we can use multiplication by reciprocal
```asm
divide_by_7:
    // Multiply by scaled reciprocal of 7
    // (2^32 + 6) / 7 = 0x24924925
    movz    x1, #0x2492, lsl #16
    movk    x1, #0x4925
    umulh   x0, x0, x1    // Upper 64 bits of multiply
    ret

// OPTIMIZED DIVISION BY CONSTANTS
divide_by_10:
    // Multiply by magic number: (2^32 + 9) / 10 = 0x1999999A
    movz    x1, #0x1999, lsl #16
    movk    x1, #0x999A
    smulh   x2, x0, x1    // Signed multiply high
    add     x2, x2, x0    // Add correction
    asr     x2, x2, #3    // Shift right by 3
    // Handle negative numbers
    cmp     x0, #0
    cneg    x2, x2, lt    // Conditional negate if negative
    mov     x0, x2
    ret
```
**UMULH/SMULH** instructions are specifically for getting the upper 64 bits of multiplication. 
It makes division by constant more efficient than x86.

Technique works well because:

- Multiplication is typically 3-4 cycles
- Division can be 12+ cycles

### Lanes
NEON cannot do 128-bit math. The reason it has space this large is because you can put data into “lanes” in order to do parallel processing.

A 128-bit register can have:
```
16 8-bit lanes
8 16-bit lanes
4 32-bit lanes
2 64-bit lanes
```

Native 64-bit Operations in ARM64

```asm
// For 64-bit addition
ldr x0, [x1]        // Load 64-bit value
ldr x2, [x3]        // Load another 64-bit value
add x4, x0, x2      // 64-bit addition in one instruction

// For 128-bit addition (similar to ADC usage in x86)
ldp x0, x1, [x2]    // Load 128-bit value (low, high)
ldp x3, x4, [x5]    // Load second 128-bit value
adds x6, x0, x3     // Add low 64 bits, set carry
adc  x7, x1, x4     // Add high 64 bits with carry
```

Using SIMD/NEON (equivalent to MMX approach)

```asm
// Using SIMD for large integer operations
ld1 {v0.2d}, [x0]   // Load 128 bits into vector register
ld1 {v1.2d}, [x1]   // Load another 128 bits
add v2.2d, v0.2d, v1.2d  // Add as 2x64-bit integers

// For larger numbers (256-bit example)
ld1 {v0.2d-v1.2d}, [x0]  // Load 256 bits
ld1 {v2.2d-v3.2d}, [x1]  // Load another 256 bits
add v4.2d, v0.2d, v2.2d  // Add low 128 bits
add v5.2d, v1.2d, v3.2d  // Add high 128 bits
```
NEON (ARM's SIMD) is more powerful than MMX


Optimized Multiple-Precision Arithmetic

```asm
// Adding large numbers (e.g., 256-bit)
ldp x0, x1, [x10]      // Load first 128 bits
ldp x2, x3, [x10, #16] // Load second 128 bits
ldp x4, x5, [x11]      // Load third 128 bits
ldp x6, x7, [x11, #16] // Load fourth 128 bits

adds x0, x0, x4        // Add low 64 bits
adcs x1, x1, x5        // Add with carry
adcs x2, x2, x6        // Add with carry
adc  x3, x3, x7        // Add final bits with carry
```

Modern optimization tips for ARM64:
```asm
// Prefer paired loads/stores for better memory access
ldp x0, x1, [x2]       // Load pair
stp x0, x1, [x2]       // Store pair

// Use NEON for bulk operations
ld1 {v0.4s}, [x0]      // Load 128 bits
add v1.4s, v0.4s, v2.4s // Parallel add

// For very large numbers, consider using the crypto extensions
// if available (they include specialized big number operations)
```

#### Performance considerations:
ARM64 has a more complex pipeline and different instruction latencies. 
For example, most arithmetic operations like ADD have a latency of 1 cycle, but some operations like floating-point multiply can take 3-4 cycles and division 12+ cycles.

- ARM64's NEON instructions are typically faster than scalar operations for large numbers
- Unlike P4's ADC/SBB issue, ARM64's carries are efficient
- Modern ARM processors can often do multiple vector operations in parallel
- Consider using the crypto extensions for big number arithmetic if avail

ARM64 can often execute multiple instructions in parallel if they're independent, so spreading out dependent instructions gives the processor more opportunities for instruction-level parallelism.
We want to avoid back-to-back instructions where one instruction depends on the result of the previous one. Here's the equivalent example:
```asm
// Less optimal:
    mul x0, x1, x2      // multiply takes multiple cycles
    add x0, x0, #1      // depends on previous result
    cmp x0, #100        // depends on previous result
    b.gt large_number

// Better:
    mul x0, x1, x2      // multiply takes multiple cycles
    add x3, x4, #5      // independent operation
    add x5, x6, #7      // independent operation
    add x0, x0, #1      // by now, mul result should be ready
    cmp x0, #100
    b.gt large_number
```
It is beneficial to help hardware scheduler by organizing code to minimize obvious dependencies, especially around high-latency instructions such as multiply, divide, or memory ops.

[NEXT -> strings](12_strings.md)

<div align="center">
  <img src="../img/argo-mascot.jpg" alt="Logo">
</div>
<p align="center">
    <img src="https://raw.githubusercontent.com/bornmay/bornmay/Update/svg/Bottom.svg" alt="Github Stats" />
</p>
<p align="right">(<a href="#top">back to top</a>)</p>
