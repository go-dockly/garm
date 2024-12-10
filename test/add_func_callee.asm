.global start

add:
    // Set up frame pointer
    stp x29, x30, [sp, #-16]!
    mov x29, sp

    add x0, x0, x1      // c = a + b

    // Restore frame pointer
    ldp x29, x30, [sp], #16
    ret                 // return to caller

start:
    // Save all registers we'll need
    stp x29, x30, [sp, #-16]!
    mov x29, sp
    stp x19, x20, [sp, #-16]!   // Save callee-saved registers we'll use
    stp x21, x22, [sp, #-16]!   // Save more if needed

    // First calculation
    mov x0, #1
    mov x1, #2
    bl add
    mov x19, x0         // Store in callee-saved register

    // Second calculation
    mov x0, #3
    mov x1, #4
    bl add
    mov x20, x0         // Store in callee-saved register

    // More calculations...
    mov x0, x19
    mov x1, x20
    bl add
    mov x21, x0         // Store in another callee-saved register

    // Restore registers in reverse order
    ldp x21, x22, [sp], #16
    ldp x19, x20, [sp], #16
    ldp x29, x30, [sp], #16

    mov x8, #93         // Exit syscall number
    svc #0