// Generational GC Constants
.equ GEN0_SIZE,     4 * 1024 * 1024  // 4MB Eden/Young space
.equ GEN1_SIZE,     16 * 1024 * 1024 // 16MB Survivor space  
.equ GEN2_SIZE,     44 * 1024 * 1024 // 44MB Tenured/Old space
.equ AGE_THRESHOLD,  15               // Objects survive this many GCs before promotion

// Generation Headers (32 bytes each)
.struct 0
gen_start:    .space 8     // Start address
gen_end:      .space 8     // End address
gen_top:      .space 8     // Allocation pointer
gen_survived: .space 8     // Bytes survived last GC
.end

// Object Header Extensions (4 bytes)
.struct 0 
obj_age:     .space 1      // Survival count
obj_gen:     .space 1      // Current generation
obj_forward: .space 2      // Forwarding pointer offset
.end

// GC State Extensions
.struct 0
gc_gen0:     .space 32     // Eden space
gc_gen1:     .space 32     // Survivor space
gc_gen2:     .space 32     // Tenured space
gc_alc_gen:  .space 8      // Current allocation generation
.end

// Fast Path Allocation
.macro GEN_ALLOC size, type
    // Try eden allocation first
    adr     x19, gc_state
    ldr     x20, [x19, #gc_gen0 + gen_top]
    add     x21, x20, \size
    ldr     x22, [x19, #gc_gen0 + gen_end]
    cmp     x21, x22
    b.gt    1f              // Eden full - trigger minor GC
    
    // Fast path eden allocation
    str     x21, [x19, #gc_gen0 + gen_top]
    mov     x0, x20
    
    // Initialize header
    str     \size, [x0, #obj_size]
    str     \type, [x0, #obj_type] 
    mov     w1, #0
    strb    w1, [x0, #obj_age]    // Age 0
    strb    w1, [x0, #obj_gen]    // Generation 0
    b       2f

1:  // Minor collection
    bl      minor_gc
    
    // Retry allocation
    GEN_ALLOC \size, \type

2:  // Done
.endm

// Minor (Young Generation) Collection
minor_gc:
    INIT_FRAME 6
    
    // Save registers
    stp     x19, x20, [sp, #-16]!
    stp     x21, x22, [sp, #-16]!
    stp     x23, x24, [sp, #-16]!
    
    // Copy survivors from eden to survivor space
    adr     x19, gc_state
    ldr     x20, [x19, #gc_gen0 + gen_start]
    ldr     x21, [x19, #gc_gen0 + gen_top]
    
1:  cmp     x20, x21
    b.ge    2f                      // Done with eden
    
    // Check if object survives
    bl      is_live_object
    cbz     x0, 3f                  // Dead object
    
    // Copy to survivor space
    bl      copy_to_survivor
    
3:  // Next object
    ldr     x22, [x20, #obj_size]
    add     x20, x20, x22
    b       1b
    
2:  // Reset eden
    ldr     x20, [x19, #gc_gen0 + gen_start]
    str     x20, [x19, #gc_gen0 + gen_top]
    
    // Process survivors
    bl      age_survivors
    
    // Restore registers
    ldp     x23, x24, [sp], #16
    ldp     x21, x22, [sp], #16
    ldp     x19, x20, [sp], #16
    
    SAFE_STACK_RET 6

// Copy Object to Survivor Space
copy_to_survivor:
    // Get survivor space allocation point
    ldr     x22, [x19, #gc_gen1 + gen_top]
    
    // Copy object
    ldr     x23, [x20, #obj_size]
    mov     x0, x22                 // Destination
    mov     x1, x20                 // Source
    mov     x2, x23                 // Size
    bl      memcpy
    
    // Update age
    ldrb    w23, [x20, #obj_age]
    add     w23, w23, #1
    strb    w23, [x22, #obj_age]
    
    // Set forwarding pointer
    sub     x23, x22, x20          // Offset
    strh    w23, [x20, #obj_forward]
    
    // Update survivor top
    add     x22, x22, x23
    str     x22, [x19, #gc_gen1 + gen_top]
    
    ret

// Age and Promote Survivors
age_survivors:
    // For each object in survivor space
    ldr     x20, [x19, #gc_gen1 + gen_start]
    ldr     x21, [x19, #gc_gen1 + gen_top]
    
1:  cmp     x20, x21
    b.ge    2f                      // Done
    
    // Check age
    ldrb    w22, [x20, #obj_age]
    cmp     w22, #AGE_THRESHOLD
    b.lt    3f                      // Not old enough
    
    // Promote to tenured
    bl      promote_to_tenured
    
3:  // Next object
    ldr     x22, [x20, #obj_size]
    add     x20, x20, x22
    b       1b
    
2:  ret

// Promote Object to Tenured Space  
promote_to_tenured:
    // Get tenured allocation point
    ldr     x22, [x19, #gc_gen2 + gen_top]
    
    // Copy object
    ldr     x23, [x20, #obj_size]
    mov     x0, x22
    mov     x1, x20
    mov     x2, x23
    bl      memcpy
    
    // Update generation
    mov     w23, #2                 // Tenured generation
    strb    w23, [x22, #obj_gen]
    
    // Update tenured top
    add     x22, x22, x23
    str     x22, [x19, #gc_gen2 + gen_top]
    
    ret