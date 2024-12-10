.data
    // String constant "arGO彡" in UTF-8
    treasure:
        .ascii "arGO\xe5\xbd\xa1"
        .size treasure, 7

    .text
    .align 2
    .global main

main:
    // Function prologue
    stp     x29, x30, [sp, #-16]!    // Save frame pointer and link register
    mov     x29, sp                   // Set up frame pointer

    // In Go, strings are represented as a pointer and length pair
    // Here we're setting up the string but not using it (matching the _ = treasure)
    adrp    x0, treasure             // Load page address of string
    add     x0, x0, :lo12:treasure   // Add low 12 bits offset
    mov     x1, #7                   // Length of "arGO彡" in bytes

    // Function epilogue
    ldp     x29, x30, [sp], #16      // Restore frame pointer and link register
    ret
