.include "defer.s"
.section .text
.global defer_test
defer_test:
    // Prologue: set up frame pointer (assuming x29 and x30)
    stp     x29, x30, [sp, #-16]!    // Save frame pointer and link register
    mov     x29, sp                  // Set frame pointer

    sub     sp, sp, #16              // Allocate space for local variables

    // Defer cleanup function with argument 42
    ldr     x0, =42                  // Load argument into x0
    DEFER cleanup, x0                // Call defer for cleanup(42)

    // Simulate some work
    bl      foo             // Call some function

    // Epilogue: defer execution point
    // Normally this `cleanup` or DEFER unwind would trigger just before return

    ldr     x0, [x29, #-8]           // Restore Defer Handler
    blr     x0                       // Executes cleanup(42)

    mov     sp, x29
    ldp     x29, x30, [sp], #16
    ret

cleanup:
    // Input argument in x0
    ret

foo:
    // Placeholder for some operation
    ret