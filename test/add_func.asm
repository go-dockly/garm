.global main

// Function: add
// x0: first parameter (a)
// x1: second parameter (b)
// x0: return value (c)
add:
    // Set up frame pointer
    stp x29, x30, [sp, #-16]!
    mov x29, sp

    add x0, x0, x1      // c = a + b

    // Restore frame pointer
    ldp x29, x30, [sp], #16
    ret                 // return to caller

main:
    // Set up frame pointer
    stp x29, x30, [sp, #-16]!
    mov x29, sp

    // Load arguments directly into parameter registers
    mov x0, #1          // First parameter (a)
    mov x1, #2          // Second parameter (b)
    
    // Call add function
    bl add
    
    // now c is in x0

    // Clean up and exitx
    ldp x29, x30, [sp], #16
    mov x0, #0          // Return 0
    mov x8, #93         // Exit syscall number
    svc #0              // Make syscall to exit