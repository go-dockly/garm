## STRING OPERATIONS
ARM64 doesn't have dedicated string instructions like x86's SCAS/CMPS/STOS
Instead, we can use:
1. Load/Store with post-increment
2. SIMD operations for bulk processing
```arm
strlen:
    mov     x2, x0                  // Save start address
    // Check alignment
    ands    x3, x0, #7              // Test if aligned
    b.eq    1f                      // Branch if aligned
    
    // Handle unaligned bytes until aligned
.align_loop:
    ldrb    w1, [x0], #1           // Load byte and increment
    cbz     w1, .done              // If zero, we're done
    ands    x3, x0, #7             // Check alignment
    b.ne    .align_loop           // Continue until aligned

1:  // Main loop - process 8 bytes at a time
    .align  4
.main_loop:
    ldr     x1, [x0], #8           // Load 8 bytes
    // Check for zero byte using clever bit manipulation
    sub     x3, x1, #0x0101010101010101
    and     x3, x3, #0x8080808080808080
    cbz     x3, .main_loop         // No zero byte, continue
    
    // Found zero - find exact position
    rev     x1, x1                 // Reverse for CLZ
    clz     x1, x1                 // Count leading zeros
    sub     x0, x0, #8             // Adjust pointer
    add     x0, x0, x1, lsr #3     // Add final offset
    
.done:
    sub     x0, x0, x2             // Calculate length
    ret
```

### SIMD String Operations Example
Process 16 bytes at once for string operations
```arm
.align 4
simd_strlen:
    dup     v0.16b, wzr            // Zero vector for comparison
    mov     x2, x0                 // Save start
    
1:  ld1     {v1.16b}, [x0], #16    // Load 16 bytes
    cmeq    v1.16b, v1.16b, v0.16b // Compare each byte with zero
    umaxv   b1, v1.16b             // Get maximum value
    fmov    w1, s1                 // Move to general register
    cbz     w1, 1b                 // If no zero found, continue
    
    // Find exact zero position...
    sub     x0, x0, #16
    // (Additional code to find exact zero position)
    sub     x0, x0, x2             // Calculate length
    ret
```

[NEXT -> macros](13_macros.md)

<div align="center">
  <img src="../img/argo-mascot.jpg" alt="Logo">
</div>
<p align="center">
    <img src="https://raw.githubusercontent.com/bornmay/bornmay/Update/svg/Bottom.svg" alt="Github Stats" />
</p>
<p align="right">(<a href="#top">back to top</a>)</p>
