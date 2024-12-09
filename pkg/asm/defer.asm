.macro DEFER func, args:vararg
    // Save current stack pointer
    mov     x19, sp

    // Allocate defer record (16 bytes for function + args + 16 bytes for link/alignment)
    sub     sp, sp, #32

    // Store function pointer
    adr     x0, \func
    str     x0, [sp]

    // Store arguments if any
    .ifnb \args
        stp     \args, [sp, #16]    // Example pair of values.
    .endif

    // Link to previous defer iF any
    ldr     x0, [x29, #-8]        // Adjust offset if alignment is diff?
    str     x0, [sp, #8]          // store prev linked.
    str     sp, [x29, #-8]        // Update current defer pos.
.endm