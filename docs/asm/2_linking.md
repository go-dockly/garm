<p align="center"> 
  Visitor count<br>
  <img src="https://profile-counter.glitch.me/sagar-viradiya/count.svg" />
</p>

```arm
@ upper.asm
.global upper

upper:	
	PUSH {R4-R6}			@ save registers
	MOV	R4, R1

loop:						@ until byte pointed to by R1 is non-zero
	LDRB R5, [R0], #1		@ load character and increment pointer
							@ Want to know if 'a' <= R5 <= 'z'
							@ First subtract 'a'
	SUB	R6, R5, #'a'
							@ Now want to know if R6 <= 25
	CMP	R6, #25	    		@ chars are 0-25 after shift
	BHI	cont
							@ if we got here then the letter is lowercase
	SUB	R5, #('a'-'A') 		@ convert to uppercase by subtracting 32 from ascii value
cont:						@ end if
	STRB R5, [R1], #1		@ store character to output str
	CMP	R5, #0				@ exit when encounter null character
	BNE	loop				@ loop if character isn't null
	SUB	R0, R1, R4  		@ get the length by subtracting the pointers
	POP	{R4-R6}				@ Restore the register we use.
	BX	LR					@ Return to caller
```

```arm
@ main.asm
    .data
input_str:  .asciz "Hello World!"    @ Input string to convert
output_str: .skip 100                @ Buffer for output string

    .text
    .global main
    .extern upper                    @ Declare upper as external function

main:
    PUSH {LR}                        @ Save link register

    @ Prepare parameters for upper function
    LDR R0, =input_str              @ Load address of input string
    LDR R1, =output_str             @ Load address of output buffer
    BL upper                        @ Call upper function

    @ At this point:
    @ R0 contains the length of the string
    @ output_str contains the converted string

    @ You can add code here to print or use the result

    POP {LR}                        @ Restore link register
    BX LR                           @ Return from main
```

Compile and link assembly files:

```bash
$ as -o build/main.o main.asm
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
