## 

```asm
// upper.s (ARM64 version)
.global upper

upper:
    // Save frame pointer and link register
    STP X29, X30, [SP, #-16]!
    MOV X29, SP

    // Save additional registers we'll use
    STP X19, X20, [SP, #-16]!

    // Save input parameters
    MOV X19, X1      // Save output string pointer
    MOV X20, X1      // Keep a copy for length calculation

loop:
    // Load byte and increment input pointer
    LDRB W4, [X0], #1

    // Check if character is lowercase
    CMP W4, #'a'
    BLO cont          // If below 'a', not lowercase
    CMP W4, #'z'
    BHI cont          // If above 'z', not lowercase

    // Convert to uppercase by subtracting 32
    SUB W4, W4, #('a' - 'A')

cont:
    // Store converted character
    STRB W4, [X1], #1

    // Check for null terminator
    CBZ W4, done

    B loop

done:
    // Calculate string length
    SUB X0, X1, X20

    // Restore registers
    LDP X19, X20, [SP], #16
    LDP X29, X30, [SP], #16
    RET
```

```asm
// main.s (ARM64 version)
.data
input_str:  .asciz "Hello World!"    // Input string to convert
output_str: .skip 100                // Buffer for output string

.text
.global main
.extern upper                        // Declare upper as external function

main:
    // Save frame pointer and link register
    STP X29, X30, [SP, #-16]!
    MOV X29, SP

    // Prepare parameters for upper function
    ADR X0, input_str    // Load address of input string
    ADR X1, output_str   // Load address of output buffer
    BL upper             // Call upper function

    // At this point:
    // X0 contains the length of the string
    // output_str contains the converted string

    // Restore frame pointer and link register
    LDP X29, X30, [SP], #16
    RET
```

Compile and link assembly files:

```bash
# Assemble the files
$ as -o upper.o upper.s
$ as -o build/main.o main.asm
# Link the object files
$ ld -o program build/main.o build/upper.o
```

[NEXT -> instructions](3_instruction.md)

<div align="center">
  <img src="../img/argo-mascot.jpg" alt="Logo">
</div>
<p align="center">
	<img src="https://raw.githubusercontent.com/bornmay/bornmay/Update/svg/Bottom.svg" alt="Github Stats" />
</p>
<p align="right">(<a href="#top">back to top</a>)</p>
