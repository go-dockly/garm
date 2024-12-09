### Functions 

Functions in ARM assembly are similar to functions in high-level languages. They are blocks of code that perform a specific task and can be called from other parts of the program. Functions in ARM assembly can take arguments, return values, and modify registers.

arguments are passed through registers (R0-R3) for integer values, and (S0-S7) for floating-point values. The return value is also passed through a register (R0 for integer values, and S0 for floating-point values).

```arm
// Input: x0 = channel ptr, x1 = buffer size, w2 = direction
// Output: x0 = channel ptr or NULL on failure
// Clobbers: x3-x6
chan_init:
    stp     x29, x30, [sp, #-48]!      // Save more registers for safety
    stp     x19, x20, [sp, #16]
    stp     x21, x22, [sp, #32]
    mov     x29, sp
    // ...
    ldp     x21, x22, [sp, #32]
    ldp     x19, x20, [sp, #16]
    ldp     x29, x30, [sp], #48
    ret
```
A typical ARM function frame setup that follows the ABI (Application Binary Interface). 
The function preserves the registers it needs to use (x19-x22) as required by the calling convention, and properly maintains the frame pointer (x29) and link register (x30) for function call chain maintenance and debugging.

The actual function body `// ...` would implement whatever channel logic is needed having registers x3-x6 plus any of the saved registers (x19-x22) free to use.

x0: Contains a pointer to channel structure
x1: Contains buffer size
x2: Contains direction val

The function returns channel pointer in x0, or NULL if it fails

Function prologue:
```arm
stp     x29, x30, [sp, #-48]!
stp     x19, x20, [sp, #16]
stp     x21, x22, [sp, #32]
mov     x29, sp
```
stp stands for "Store Pair" storing two registers at once
The ! at the end of the first instruction means pre-indexed 
first decrementing SP by 48 bytes befoe storing the arguments
in six registers total:

x29 (Frame Pointer) and x30 (Link Register)
x19 and x20 (preserved registers)
x21 and x22 (preserved registers)

Sets up x29 as the frame pointer pointing to the stack

Matching epilogue at the end
```arm
ldp     x21, x22, [sp, #32]
ldp     x19, x20, [sp, #16]
ldp     x29, x30, [sp], #48
ret
```
Restores all saved registers in reverse order using ldp (Load Pair)
ldp is post-indexed with `]`incrementing SP by 48 after loading
ret jumps to the address in x30 thereby returning from the function

### Labels
Best Practices

- Use descriptive, unique labels for major points in your code
- Use local labels (starting with . or numbers) for small, local jumps
- Document your labels with comments
- Keep a consistent naming convention
- Consider using a prefix for different sections of code
```arm
@ Option 1: Use unique names
loop_outer:
    // ...
loop_outer_done:
    // ...

@ Option 2: Use local labels (starting with number or dot)
.Loop:  @ Local label
    // ...
    
@ Option 3: .L\@ in the label names. The .L means local and the \@ is a special asm symbol that gets replaced with a unique number for each macro expansion. This ensures that each instance of the macro gets its own unique set of labels
```arm
// First use
.L123_oop:
    // ...
.L123_oop_done:

// Second use
.L124_oop:
    // ...
```
NOTE: dot prefix is reserved for:

- Local labels (.Lloop)
- Assembler directives (.global, .text, .data)

Valid Symbol Names:

- Can contain letters, digits, underscores
- Must start with a letter or underscore
- Cannot contain dots/periods within the name

Incorrect:
```arm
.global fmt.Println    @ Not allowed - contains dot
.global .test          @ Not allowed - starts with dot
```

Correct:
```arm
.global fmt_println   @ Valid - uses underscore
.global FmtPrintln    @ Valid - camelCase
.global FMT_PRINTLN   @ Valid - uppercase with underscore
```

[NEXT -> branching](branch.md)

<div align="center">
  <img src="../img/argo-mascot.jpg" alt="Logo">
</div>
<p align="center">
		<img src="https://raw.githubusercontent.com/bornmay/bornmay/Update/svg/Bottom.svg" alt="Github Stats" />
</p>
<p align="right">(<a href="#top">back to top</a>)</p>