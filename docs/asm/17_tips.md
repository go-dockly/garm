The following is based on [Assembly Optimization Tips by Mark Larson](https://masm32.com/masmcode/marklarson/index.htm)
and translated to arm64 equivalent

## Minimize memory
ARM64 is a load-store architecture, you almost always need to:

- Load from memory into a register
- Perform operation
- Store back to memory if needed

The principle of "minimize memory operands" is very important in ARM64 

1. Minimize memory accesses where possible (keep values in registers)
2. Use conditional instructions where available (although ARM's predication system is different)
3. Align code and data appropriately (ARM64 has different alignment requirements but the principle holds)
4. Consider instruction latency and throughput
5. Avoid data dependencies when possible for better pipelining

```asm
// Bad (memory heavy)
str w0, [x1]        // Store to memory
ldr w0, [x1]        // Load right back
add w0, w0, #1      // Add

// Good (register heavy)
mov w2, w0          // Keep in register
add w2, w2, #1      // Work with register
str w2, [x1]        // Store only when needed

// Instead of complex memory operations, use simple loads/stores
// Bad (if it existed):
// complex_memory_copy    x0, x1, #16    // Hypothetical complex instruction

// Good:
ldp x2, x3, [x1], #16     // Load pair and update pointer
stp x2, x3, [x0], #16     // Store pair and update pointer

// Use modern alternatives:
// Instead of rep movsb equivalent, use:
// - ldp/stp for known small sizes
// - memcpy for larger copies
// - Consider using SIMD (vector) registers for bulk data movement
```

## Register Preservation

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

### Use paired instructions

```asm
// Better than single loads/stores
ldp x1, x2, [x0]      // Load pair
stp x1, x2, [x0]      // Store pair
```

### Take advantage of zero register initialisation

```asm
// xzr is always zero
add x0, x1, xzr       // Move register (clearer than mov)
```
### Use post-index and pre-index addressing

```asm
ldr x0, [x1], #8      // Load and increment
ldr x0, [x1, #8]!     // Increment and load
```
### Consider NEON/SIMD instructions for data-parallel operations

```asm
ld1 {v0.4s}, [x0]     // Load 4 32-bit values
add v1.4s, v0.4s, v0.4s // Parallel add
```

### Rotating in ARM64
```asm
// Prefer fixed immediate rotates
ror w0, w0, #1        // Good: fixed immediate
ror w0, w0, w1        // Less efficient: variable rotate

// ROR - Rotate right
mov     x0, #0x1234
ror     x0, x0, #8     // Rotate right by 8 bits
```
Note: Left rotation can be achieved with ROR using (32/64 - n)

### Utilize Flags
```asm
// Less efficient
sub w0, w0, #1        // Decrement
cmp w0, #0            // Unnecessary compare
b.ne loop             // Branch if not zero

// More efficient
subs w0, w0, #1       // Decrement and set flags
b.ne loop             // Branch using flags from subs
```

### Multiply Instead of Divide
```asm
// For power-of-2 divisions, use shift
lsr w0, w0, #3        // Divide by 8

// For other constants, use multiply-high
// Example dividing by 3:
movz w1, #0xAAAB      
mul  w0, w0, w1       
lsr  w0, w0, #1       // Result in w0
```

Zero/Sign Extension in ARM64
ARM64 has dedicated instructions:

```asm
// Zero extend
uxtb w0, w0           // Byte to word
uxth w0, w0           // Halfword to word

// Sign extend
sxtb w0, w0           // Byte to word
sxth w0, w0           // Halfword to word
```

### Zeroing Registers in ARM64

```asm
// Best way to zero a register
mov w0, wzr           // Using zero register
// or
eor w0, w0, w0        // Using XOR (might be useful in specific cases)
```

### Unsigned Division in ARM64

```asm
// Signed division
sdiv w0, w0, w1

// Unsigned division (typically faster)
udiv w0, w0, w1
```

### Avoiding Dependencies in ARM64

```asm
// Poor scheduling (dependency)
add w0, w0, #1
add w1, w1, #1
cmp w0, #1            // Depends on first add

// Better scheduling
add w0, w0, #1
mov w2, #1            // Independent instruction
add w1, w1, #1
cmp w0, w2            // Dependencies better distributed
```

### Use Combined Operations

```asm
// Add/subtract with shift
add w0, w1, w2, lsl #2    // w0 = w1 + (w2 << 2)
```

### Take Advantage of Conditional Instructions

```asm
cmp w0, #0
csel w0, w1, w2, eq       // w0 = (w0 == 0) ? w1 : w2
```

### Use Post/Pre-Index for Efficient Memory Access

```asm
ldr w0, [x1], #4          // Load and increment
ldr w0, [x1, #4]!         // Increment and load
```

### Utilize Compare and Branch Instructions

```asm
cbz w0, label             // Compare and branch if zero
cbnz w0, label            // Compare and branch if not zero
```

The famous LEA (Load Effective Address) optimization from x86. In ARM64, while we don't have an exact equivalent of LEA, we have several ways to achieve similar optimizations:

### Combined Arithmetic Operations in ARM64

```asm
// ARM64 offers combined add/shift operations:
add x1, x2, x2, lsl #2    // x1 = x2 + (x2 << 2)  // Similar to LEA [reg*4]
add x1, x2, x2, lsl #1    // x1 = x2 + (x2 << 1)  // Similar to LEA [reg*3]

// For your specific example (multiply by 4 and add 3):
loop:
    subs w0, w0, #1        // Decrement and set flags
    add  w1, w2, w2, lsl #2  // w1 = w2 * 5
    add  w1, w1, #3         // Add 3
    b.ne loop              // Branch if not zero
```

### Flag-Independent Operations
ARM64 has explicit flag-setting versions of instructions, so you can control when flags are affected:

```asm
// Regular add (doesn't affect flags)
add w0, w0, w1

// Flag-setting add (affects flags)
adds w0, w0, w1

// Example combining both:
subs w0, w0, #1           // Sets flags
add  w1, w2, w2, lsl #2   // Doesn't affect flags
b.ne loop                 // Uses flags from subs
```

### More Complex Address Calculations
ARM64 can do some complex addressing in one instruction:
```asm
add x0, x1, x2, lsl #3    // x0 = x1 + (x2 << 3)
add x0, x1, x2, lsr #2    // x0 = x1 + (x2 >> 2)

// Can combine with extend operations:
add x0, x1, w2, sxtw #2   // Sign-extend w2 to 64 bits, shift left by 2, then add
```

### Alternative Optimizations

ARM64 offers other efficient ways to do math:
```asm
madd w0, w1, w2, w3      // w0 = w1 * w2 + w3
msub w0, w1, w2, w3      // w0 = w1 * w2 - w3

// For array indexing (similar to common LEA use):
add x0, x1, x2, lsl #2   // For 4-byte elements
add x0, x1, x2, lsl #3   // For 8-byte elements
```

Key differences from x86's LEA:

- ARM64 separates memory addressing from arithmetic operations
- You often need to chain multiple instructions for complex calculations
- Flag behavior is more explicit and predictable
- ARM64 has more registers, reducing the need for complex addressing modes
- Clearer separation of concerns
- More predictable behavior
- Explicit flag control
- Rich set of combined arithmetic operations

## Built-in byte swap instructions:
There are several options for handling byte swapping and bit manipulation analogous to x86's BSWAP, ROL, and ROR
1. REV - Reverse bytes across entire register
2. REV16 - Reverse bytes in each halfword
3. REV32 - Reverse bytes in each word
```asm
// Cool trick to switch from Big Endian to Little Endian, using REV instruction:
    mov     w0, #0          // Clear the 32-bit register
    movz    w0, #234        // Set lower 16 bits to 234
    rev     w0, w0          // Reverse bytes (similar to BSWAP)
    movz    w1, #345        // Store 345 in another register temporarily
    bfi     w0, w1, #0, #16 // (Bit Field Insert) 345 into lower 16 bits
    rev     w0, w0          // Swap back to access first value
    add     w1, w0, #5      // Add 5 to the first value (now in lower 16 bits)
    rev     w0, w0          // Swap to access second value
    add     w1, w0, #7      // Add 7 to second value

// Alternative approach using bit field operations:
    movz    x0, #234            // Store first value
    movk    x0, #345, lsl #16   // Store second value in upper halfword
    add     w1, w0, #5          // Add to lower halfword
    add     w2, w0, lsr #16, #7 // Add to upper halfword
    bfi     x0, x1, #0, #16     // (Bit Field Insert) result back
    bfi     x0, x2, #16, #16    // (Bit Field Insert) second result

// Counting set bits (population count):
// ARM64 provides dedicated instructions:
    cnt     v0.8b, v0.8b   // Count set bits (SIMD)
    fmov    x0, d0         // Move result to general register
```
Note: 
- UBFX/SBFX (Bit Field Extract)
- CLZ (Count Leading Zeros) for zero-byte detection

#### **Better Native Support**: Unlike x86, ARM64 has dedicated byte-reverse instructions (`REV`, `REV16`, `REV32`) that are generally more efficient than the x86 `BSWAP`.

#### **Better Bit Manipulation**: Instead of using rotates for packing, ARM64 offers:
   - `BFI` (Bit Field Insert)
   - `UBFX`/`SBFX` (Bit Field Extract)
   - `MOVK` for moving 16-bit immediates to specific halfwords

#### **Performance Characteristics**: Unlike the P4 where rotates were faster than `BSWAP`, in modern ARM64:
   - `REV` instructions are highly optimized
   - Bit field instructions often provide better alternative
   - SIMD operations can parallelize many of these ops

### Zero-extension 
is handled differently than x86. Several instructions have zero-extension built in:
1. LDRB (Load Register Byte) with zero-extension
```asm
zero_extend_byte:
    // This automatically zero-extends to 32 bits
    ldrb    w0, [x1]        // Load byte and zero-extend to 32 bits
    // w0 now contains the zero-extended value, no extra instructions needed
```
2. LDRH (Load Register Halfword) with zero-extension
```asm
zero_extend_halfword:
    // This automatically zero-extends to 32 bits
    ldrh    w0, [x1]        // Load halfword and zero-extend to 32 bits
```
3. Using dedicated zero-extend instructions
```asm
explicit_zero_extend:
    // UXTB - Extract unsigned byte
    uxtb    w0, w1          // Zero-extend byte to 32 bits
    // UXTH - Extract unsigned halfword
    uxth    w0, w1          // Zero-extend halfword to 32 bits
```
4. Working with 64-bit registers
```asm
zero_extend_to_64:
    // When loading to X registers, use zero-extending variants
    ldrb    x0, [x1]        // Zero-extends to full 64 bits
    
    // Or use explicit extension
    uxtb    x0, x1          // Zero-extend byte to 64 bits
```

5. Efficient handling of immediate values
```asm
load_small_immediate:
    // MOVZ automatically zero-extends
    movz    w0, #123        // Loads 123, upper bits are zeroed
```
6. Combining operations
```asm
combine_operations:
    // Load byte and add in one sequence
    ldrb    w0, [x1]        // Load and zero-extend
    add     w0, w0, #1      // Add to zero-extended value
    
    // Or with halfword
    ldrh    w0, [x1]        // Load and zero-extend
    add     w0, w0, #1      // Add to zero-extended value
```
7. Bitfield operations (alternative approach)
```asm
bitfield_zero_extend:
    // Using UBFX (Unsigned Bit Field Extract)
    ubfx    w0, w1, #0, #8  // Extract and zero-extend 8 bits
    ubfx    w0, w1, #0, #16 // Extract and zero-extend 16 bits
```
8. Working with arrays/buffers
```asm
process_byte_array:
    // Process multiple bytes efficiently
.loop:
    ldrb    w0, [x1], #1    // Load byte, post-increment
    // w0 is automatically zero-extended
    // ... process byte ...
    subs    x2, x2, #1      // Decrement counter
    b.ne    .loop           // Continue if not done
```
### Automatic Zero Extension:

- ARM64 load instructions (LDRB, LDRH) automatically zero-extend when loading to W registers
- No need for separate zero-extension operations in many cases like using W registers (32-bit), upper 32 bits of X register are automatically zeroed
- No partial register stalls like in x86 makes it much more flexible

Useful when the value is already in a register:

- UXTB (Unsigned Extend Byte)
- UXTH (Unsigned Extend Halfword)

Performance Characteristics:

- Zero-extension operations are typically single-cycle and can often be combined with other operations

## Byte Extraction Optimizations
Original C code
```c
        unsigned char c = ((the_array[i])>>(Pass<<3)) & 0xFF;

; I got rid of the "pass" variable by unrolling the loop 4 times.
; So I had 4 of these each one seperated by lots of C code.
        unsigned char c = (the_array[i])>>0) & 0xFF;
        unsigned char c = (the_array[i])>>8) & 0xFF;
        unsigned char c = (the_array[i])>>16) & 0xFF;
        unsigned char c = (the_array[i])>>24) & 0xFF;
```
Less optimal version using shifts and masks:
```asm
original_version:
    ldr     w0, [x1]            // Load the dword
    lsr     w0, w0, #16         // Shift right by 16
    and     w0, w0, #0xFF       // Mask to get byte
```
Better approaches:
1. Using byte-aligned loads
Most efficient - direct byte access:
```asm 
efficient_version:
    ldrb    w0, [x1, #2]        // Directly load third byte
    // No additional instructions needed!
```
2. Using UBFX (Unsigned Bit Field Extract)
Alternative when value is already in register:
```asm
bitfield_version:
    ldr     w0, [x1]            // Load the dword
    ubfx    w0, w0, #16, #8     // Extract bits 16-23
```
3. Unrolled version for all 4 bytes
Most efficient approach when needing all bytes:
```asm
unrolled_version:
    ldr     w0, [x1]            // Load the dword once
    // Then either:
    
    // Option A: Using byte-aligned loads
    ldrb    w1, [x1, #0]        // Get byte 0
    ldrb    w2, [x1, #1]        // Get byte 1
    ldrb    w3, [x1, #2]        // Get byte 2
    ldrb    w4, [x1, #3]        // Get byte 3
    
    // Option B: Using UBFX
    ubfx    w1, w0, #0, #8      // Extract byte 0
    ubfx    w2, w0, #8, #8      // Extract byte 1
    ubfx    w3, w0, #16, #8     // Extract byte 2
    ubfx    w4, w0, #24, #8     // Extract byte 3
    
    // Option C: Using LSR and AND (least efficient)
    and     w1, w0, #0xFF       // Get byte 0
    lsr     w2, w0, #8          // Shift for byte 1
    and     w2, w2, #0xFF       // Mask byte 1
    lsr     w3, w0, #16         // Shift for byte 2
    and     w3, w3, #0xFF       // Mask byte 2
    lsr     w4, w0, #24         // Shift for byte 3
    // No AND needed for byte 3 as upper bits are already 0
```
4. SIMD approach for processing multiple dwords
```asm
simd_version:
    ld1     {v0.4s}, [x1]           // Load 4 dwords
    uzp1    v1.16b, v0.16b, v0.16b  // Unzip bytes
    // Now v1 contains all bytes unpacked
```

Consider yourself an ARM64 assembly PRO now. Congratulations 🎉 hero!
It's time to go and conquer the world...

<div align="center">
  <img src="../img/argo-mascot.jpg" alt="Logo">
</div>
<p align="center">
    <img src="https://raw.githubusercontent.com/bornmay/bornmay/Update/svg/Bottom.svg" alt="Github Stats" />
</p>
<p align="right">(<a href="#top">back to top</a>)</p>
