// Data section
.data
    // Rectangle data
    width:  .double 3.0
    height: .double 4.0
    // Circle data
    radius: .double 2.0
    pi:     .double 3.14159265359
    fmt:    .string "Area: %f\n"

// Define geometry interface vtable structure
.section .rodata
geometry_vtable_type:
    .quad 0     // area method pointer

// Rectangle implementation of geometry interface
rect_vtable:
    .quad rect_area    // Implements geometry's area method

// Circle implementation of geometry interface
circle_vtable:
    .quad circle_area  // Implements geometry's area method

// BSS section
.bss
// geometry interface instances
geometry1:
    .skip 8     // Space for pointer to first shape
geometry2:
    .skip 8     // Space for pointer to second shape

// Shape instances
rectangle:
    .skip 24    // vtable pointer (8) + width (8) + height (8)
circle:
    .skip 16    // vtable pointer (8) + radius (8)

// Text section
.text
.global _start
.align 2

// Rectangle's implementation of area
rect_area:
    // Input: x0 = pointer to rectangle struct
    // Output: d0 = area result
    ldr     d0, [x0, #8]     // Load width
    ldr     d1, [x0, #16]    // Load height
    fmul    d0, d0, d1       // Calculate area
    ret

// Circle's implementation of area
circle_area:
    // Input: x0 = pointer to circle struct
    // Output: d0 = area result
    ldr     d0, [x0, #8]     // Load radius
    fmul    d0, d0, d0       // radius squared
    adrp    x1, pi
    add     x1, x1, :lo12:pi
    ldr     d1, [x1]         // Load pi
    fmul    d0, d0, d1       // pi * r^2
    ret

// Initialize a rectangle
init_rectangle:
    // Input: x0 = pointer to rectangle struct
    //        d0 = width
    //        d1 = height
    adrp    x1, rect_vtable
    add     x1, x1, :lo12:rect_vtable
    str     x1, [x0]         // Store vtable pointer
    str     d0, [x0, #8]     // Store width
    str     d1, [x0, #16]    // Store height
    ret

// Initialize a circle
init_circle:
    // Input: x0 = pointer to circle struct
    //        d0 = radius
    adrp    x1, circle_vtable
    add     x1, x1, :lo12:circle_vtable
    str     x1, [x0]         // Store vtable pointer
    str     d0, [x0, #8]     // Store radius
    ret

// Calculate and print area through geometry interface
print_area:
    // Input: x0 = pointer to geometry interface
    stp     x29, x30, [sp, #-16]!  // Save registers
    
    ldr     x0, [x0]         // Load implementing type pointer
    ldr     x1, [x0]         // Load vtable
    ldr     x1, [x1]         // Load area implementation
    blr     x1              // Call area through interface

    // Print result (area is in d0)
    adrp    x0, fmt
    add     x0, x0, :lo12:fmt
    bl      printf
    
    ldp     x29, x30, [sp], #16   // Restore registers
    ret

_start:
    stp     x29, x30, [sp, #-16]!

    // Initialize rectangle
    adrp    x0, rectangle
    add     x0, x0, :lo12:rectangle
    adrp    x1, width
    add     x1, x1, :lo12:width
    ldr     d0, [x1]         // Load width
    adrp    x1, height
    add     x1, x1, :lo12:height
    ldr     d1, [x1]         // Load height
    bl      init_rectangle

    // Initialize circle
    adrp    x0, circle
    add     x0, x0, :lo12:circle
    adrp    x1, radius
    add     x1, x1, :lo12:radius
    ldr     d0, [x1]         // Load radius
    bl      init_circle

    // Set up first geometry interface to point to rectangle
    adrp    x0, geometry1
    add     x0, x0, :lo12:geometry1
    adrp    x1, rectangle
    add     x1, x1, :lo12:rectangle
    str     x1, [x0]

    // Set up second geometry interface to point to circle
    adrp    x0, geometry2
    add     x0, x0, :lo12:geometry2
    adrp    x1, circle
    add     x1, x1, :lo12:circle
    str     x1, [x0]

    // Print rectangle area through geometry interface
    adrp    x0, geometry1
    add     x0, x0, :lo12:geometry1
    bl      print_area

    // Print circle area through geometry interface
    adrp    x0, geometry2
    add     x0, x0, :lo12:geometry2
    bl      print_area

    // Exit program
    ldp     x29, x30, [sp], #16
    mov     x0, #0
    mov     x8, #93
    svc     #0

.size _start, .-_start