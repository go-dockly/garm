// Example Application Using Generational GC
.text
.global main

// Initialize Generational GC
init_gen_gc:
    // Allocate generation spaces
    mov     x0, #0
    ldr     x1, =GEN0_SIZE
    mov     x2, #(PROT_READ | PROT_WRITE)
    mov     x3, #(MAP_PRIVATE | MAP_ANONYMOUS)
    mov     x8, #SYS_mmap
    svc     #0
    
    // Initialize eden space (gen0)
    adr     x19, gc_state
    str     x0, [x19, #gc_gen0 + gen_start]
    add     x1, x0, #GEN0_SIZE
    str     x1, [x19, #gc_gen0 + gen_end]
    str     x0, [x19, #gc_gen0 + gen_top]
    
    // Allocate and init survivor space (gen1)
    add     x0, x1, #16            // Align to 16 bytes
    ldr     x1, =GEN1_SIZE
    mov     x8, #SYS_mmap
    svc     #0
    
    str     x0, [x19, #gc_gen1 + gen_start]
    add     x1, x0, #GEN1_SIZE
    str     x1, [x19, #gc_gen1 + gen_end]
    str     x0, [x19, #gc_gen1 + gen_top]
    
    // Allocate and init tenured space (gen2)
    add     x0, x1, #16
    ldr     x1, =GEN2_SIZE
    mov     x8, #SYS_mmap
    svc     #0
    
    str     x0, [x19, #gc_gen2 + gen_start]
    add     x1, x0, #GEN2_SIZE
    str     x1, [x19, #gc_gen2 + gen_end]
    str     x0, [x19, #gc_gen2 + gen_top]
    ret

// Example: Allocating Different Types of Objects
.macro ALLOC_STRING size
    // Calculate total size including header
    add     x0, \size, #16         // Add header size
    mov     x1, #TYPE_STRING       // String type ID
    GEN_ALLOC x0, x1              // Allocate in eden
.endm

.macro ALLOC_ARRAY length, elem_size
    // Calculate array size
    mov     x0, \length
    mul     x0, x0, \elem_size
    add     x0, x0, #24           // Header + length field
    mov     x1, #TYPE_ARRAY       // Array type ID
    GEN_ALLOC x0, x1             // Allocate in eden
    
    // Initialize length field
    str     \length, [x0, #16]    // Store length after header
.endm

// Example Usage in Application Code
example_allocation:
    // Save frame
    stp     x29, x30, [sp, #-16]!
    mov     x29, sp
    
    // Allocate a string
    ALLOC_STRING 32               // 32-byte string
    mov     x20, x0               // Save string pointer
    
    // Allocate an array
    ALLOC_ARRAY 100, 4           // Array of 100 integers
    mov     x21, x0               // Save array pointer
    
    // Create some garbage
    mov     x22, #0
1:  ALLOC_STRING 16              // Temporary strings
    add     x22, x22, #1
    cmp     x22, #1000
    b.lt    1b
    
    // Minor GC will happen automatically when eden fills
    
    // Access survivors (they'll be in gen1)
    ldr     x0, [x20]            // Load from string
    str     w0, [x21, #24]       // Store to array
    
    // Restore frame and return
    ldp     x29, x30, [sp], #16
    ret

// Example Root Registration
register_thread_roots:
    // Register stack roots
    mov     x0, sp               // Stack pointer
    mov     x1, x29              // Frame pointer
    bl      scan_stack_roots
    
    // Register global roots
    adr     x0, global_roots
    mov     x1, #GLOBAL_ROOT_COUNT
    bl      register_root_range
    ret

// Main Entry Point
main:
    // Initialize GC
    bl      init_gen_gc
    
    // Register roots
    bl      register_thread_roots
    
    // Run application code
    bl      example_allocation
    
    // Exit
    mov     x0, #0
    mov     x8, #SYS_exit
    svc     #0

// Data Section
.data
// Type descriptors
.align 3
type_table:
    .quad TYPE_STRING_DESC
    .quad TYPE_ARRAY_DESC
    
// Global roots
global_roots:
    .space 64                    // Space for global root pointers

// String type descriptor
TYPE_STRING_DESC:
    .quad 0                      // No fields to scan
    .word 0                      // No field count
    .word TYPE_STRING           // Type ID
    
// Array type descriptor    
TYPE_ARRAY_DESC:
    .quad array_elem_offsets    // Field offset table
    .word 1                     // One field (length)
    .word TYPE_ARRAY           // Type ID
    
array_elem_offsets:
    .quad 24                    // Offset to first element