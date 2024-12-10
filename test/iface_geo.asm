// Data section
.data
    width:  .double 3.0
    height: .double 4.0
    fmt:    .string "Area: %f\n"

// vtable for rectangle
.section .rodata
rect_vtable:
    .quad rect_area    // Function pointer to area implementation

// BSS section for our rectangle instance
.bss
rectangle:
    .skip 24   // Space for struct (16 bytes) + vtable pointer (8 bytes)

// Text section
.text
.global _start
.align 2

// Rectangle's area implementation
rect_area:
    // Input: x0 = pointer to rectangle struct
    // Output: d0 = area result
    
    // Load width and height from struct
    ldr     d0, [x0, #0]     // Load width
    ldr     d1, [x0, #8]     // Load height
    fmul    d0, d0, d1       // Calculate area
    ret

// Initialize a rectangle (constructor-like function)
init_rectangle:
    // Input: x0 = pointer to rectangle struct
    //        d0 = width
    //        d1 = height
    
    // Store the vtable pointer
    adrp    x1, rect_vtable
    add     x1, x1, :lo12:rect_vtable
    str     x1, [x0]         // Store vtable pointer at start of struct
    
    // Store width and height
    str     d0, [x0, #8]     // Store width after vtable pointer
    str     d1, [x0, #16]    // Store height
    ret

_start:
    // Save link register
    stp     x29, x30, [sp, #-16]!
    mov     x29, sp

    // Initialize rectangle
    adrp    x0, rectangle
    add     x0, x0, :lo12:rectangle
    
    // Load initial width and height
    adrp    x1, width
    add     x1, x1, :lo12:width
    ldr     d0, [x1]         // Load width
    adrp    x1, height
    add     x1, x1, :lo12:height
    ldr     d1, [x1]         // Load height
    
    bl      init_rectangle

    // Call area method through vtable
    adrp    x0, rectangle
    add     x0, x0, :lo12:rectangle
    ldr     x1, [x0]         // Load vtable pointer
    ldr     x1, [x1]         // Load area function pointer
    blr     x1              // Call area function

    // Print result (area is already in d0)
    adrp    x0, fmt
    add     x0, x0, :lo12:fmt
    bl      printf

    // Cleanup and exit
    ldp     x29, x30, [sp], #16
    mov     x0, #0
    mov     x8, #93
    svc     #0

.size _start, .-_start