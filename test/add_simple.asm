.global main

main:
    // Set up frame pointer
    stp x29, x30, [sp, #-16]!
    mov x29, sp

    // Initialize a = 1
    mov x0, #1          // a = 1

    // Initialize b = 2
    mov x1, #2          // b = 2

    // Perform addition: c = a + b
    add x2, x0, x1      // c = a + b

    // Clean up and exit
    ldp x29, x30, [sp], #16
    mov x0, #0          // Return 0
    mov x8, #93         // Exit syscall number
    svc #0              // Make syscall to exit