### Loops

```arm
.text            
.global _start

    /********************
     * Syscall format
     * x8 - syscall number
     * If there are args:
     * x0 - first argument
     * x1 - second argument
     * x2 - third argument
     * and so on
     *********************/
_start:
main:
    /**************************
     * Program overview:
     * Demonstration of a 
     * simple loop construct.
     * The loop will be used to 
     * print a message 10 times
     **************************/
    

 /***************
  * Declare a label .begin: that
  * we will jump to. Labels can generally
  * be used to jump anywhere in the program
  * For example, let's use a jump to avoid
  * this call to exit:
  ***************/
    b .begin

    /********************
     * Exit syscall
     *********************/
    mov x8, #93
    svc 0
    
 /***************
  * Thanks to b .begin
  * our program arrives here
  * and avoids the premature
  * exit
  ****************/
  
.begin:
    /******************
     * Our counter will be
     * kept in x4
     ******************/
     mov x4, #10
.loop:
    /********************
     * Write syscall
     *********************/ 
    mov x0, #1
    ldr x1, =message
    ldr x2, =len
    mov x8, 0x40
    svc 0
    /********************
     * Decrement x4
     *********************/
    sub  x4, x4, 1
    /********************
     * Compare and Branch on Zero compares a given 
     * register to zero. If the register does not equal
     * zero, the jump is taken to the given label
     *********************/    
    cbnz x4, .loop
    /********************
     * Exit syscall
     *********************/
    mov x8, #93
    svc 0  

.data
message: .asciz "hello world\n"
len = . - message
```

```go
// Vectorizable
for i := 0; i < 4; i++ {
    result[i] = a[i] + b[i]
}

// Not vectorizable (irregular access)
for i := 0; i < 4; i++ {
    result[i] = a[i*2] + b[i]
}

// Not vectorizable (dependency)
for i := 1; i < 4; i++ {
    result[i] = result[i-1] + a[i]
}
```

`for(i=0; i<=10; i++)` increment and check condition before every loop
```arm
        MOV     r0, #0      @ i
        MOV     r1, #0      @ total
for     CMP     r0, #10
        BGE     fordone     @ first branch
        ADD     r1, r1, r0
        ADD     r0, r0, #1
        B       for         @ second branch
fordone
```
`for(i=10; i>0; i--)` decrement with using only one branch at the end`
```arm
        MOV     r0, #10     @ i
        MOV     r1, #0      @ total
ford    ADD     r1, r1, r0  
        SUBS    r0, r0, #1
        BNE     ford        @ first branch
forddone
```
`while(i>0)` initial branch only run once
```arm
        MOV     r0, #10     @ i
        MOV     r1, #0      @ total
        B       wtest       @ initial branch
while   ADD     r1, r1, r0
        SUB     r0, r0, #1
wtest   CMP     r0, #0
        BGE     while       @ first branch
```

 `do..while(i>0)` same as while but without initial branch
 
```arm
        MOV     r0, #10     @ i
        MOV     r1, #0      @ total
dwhile  ADD     r1, r1, r0
        SUB     r0, r0, #1
dwtest  CMP     r0, #0
        BGE     while       @ first branch
        
        END
```

### NEON 

process 4 elements per iteration
```arm
basic_loop:
        VLD1.32     {Q0}, [R0]!           @ Load 4 elements from array a
        VLD1.32     {Q1}, [R1]!           @ Load 4 elements from array b
        VADD.I32    Q0, Q0, Q1            @ Add vectors
        VST1.32     {Q0}, [R2]!           @ Store result
        SUBS        R4, R4, #4            @ Decrement counter
        BGE         basic_loop
```

### Unrolled 2x 
process 8 elements per iteration
```arm
unrolled_2x_loop:
        @ Load 8 elements (2 quad words) from each array
        VLD1.32     {Q0,Q1}, [R0]!        @ Load 8 elements from array a
        VLD1.32     {Q2,Q3}, [R1]!        @ Load 8 elements from array b
        
        @ Perform addition
        VADD.I32    Q0, Q0, Q2            @ Add first 4 elements
        VADD.I32    Q1, Q1, Q3            @ Add second 4 elements
        
        @ Store results
        VST1.32     {Q0,Q1}, [R2]!        @ Store 8 results
        SUBS        R4, R4, #8            @ Decrement counter
        BGE         unrolled_2x_loop
```

### Unrolled 4x with prefetch 
process 16 elements per iteration
```arm
unrolled_4x_loop:
        @ Prefetch next iterations
        PLD         [R0, #64]             @ Prefetch array a
        PLD         [R1, #64]             @ Prefetch array b
        
        @ Load 16 elements (4 quad words) from each array
        VLD1.32     {Q0,Q1}, [R0]!        @ Load first 8 elements from a
        VLD1.32     {Q2,Q3}, [R0]!        @ Load second 8 elements from a
        VLD1.32     {Q8,Q9}, [R1]!        @ Load first 8 elements from b
        VLD1.32     {Q10,Q11}, [R1]!      @ Load second 8 elements from b
        
        @ Perform additions
        VADD.I32    Q0, Q0, Q8            @ Add first 4 elements
        VADD.I32    Q1, Q1, Q9            @ Add second 4 elements
        VADD.I32    Q2, Q2, Q10           @ Add third 4 elements
        VADD.I32    Q3, Q3, Q11           @ Add fourth 4 elements
        
        @ Store results
        VST1.32     {Q0,Q1}, [R2]!        @ Store first 8 results
        VST1.32     {Q2,Q3}, [R2]!        @ Store second 8 results
        SUBS        R4, R4, #16           @ Decrement counter
        BGE         unrolled_4x_loop
```

### Interleaved loads/compute with dual accumulators
```arm
interleaved_loop:
        @ Load first set while processing previous data
        VLD1.32     {Q0,Q1}, [R0]!        @ Load from array a
        VADD.I32    Q8, Q4, Q6            @ Process previous data (first set)
        VLD1.32     {Q2,Q3}, [R1]!        @ Load from array b
        VADD.I32    Q9, Q5, Q7            @ Process previous data (second set)
        
        @ Store previous results while loading new data
        VST1.32     {Q8,Q9}, [R2]!        @ Store previous results
        VLD1.32     {Q4,Q5}, [R0]!        @ Load next set from array a
        VLD1.32     {Q6,Q7}, [R1]!        @ Load next set from array b
        
        @ Continue with current data
        VADD.I32    Q10, Q0, Q2           @ Process current data (first set)
        VADD.I32    Q11, Q1, Q3           @ Process current data (second set)
        VST1.32     {Q10,Q11}, [R2]!      @ Store current results
        
        SUBS        R4, R4, #16           @ Decrement counter
        BGE         interleaved_loop
```

###  Pipelining
NOTE: Needs init and epilogue code
```arm
pipelined_loop:
        @ Load
        VLD1.32     {Q0,Q1}, [R0]!        @ Load new data from array a
        
        @ Load + Process previous
        VLD1.32     {Q2,Q3}, [R1]!        @ Load new data from array b
        VADD.I32    Q8, Q4, Q6            @ Process previous iteration
        
        @ Process current + Store previous
        VADD.I32    Q9, Q0, Q2            @ Process current data
        VST1.32     {Q8}, [R2]!           @ Store previous result
        
        @ Update pipeline registers
        VMOV        Q4, Q0                @ Move current to previous
        VMOV        Q6, Q2                @ Move current to previous
        
        SUBS        R4, R4, #8            @ Decrement counter
        BGE         pipelined_loop
```

### Advanced SIMD with multiple accumulators and prefetch
Concurrency optimized for high-end ARM processors with multiple NEON units
```arm
multi_acc_loop:
        @ Prefetch next cache lines
        PLD         [R0, #128]            @ Prefetch array a
        PLD         [R1, #128]            @ Prefetch array b
        
        @ Load first block while processing previous block
        VLD1.32     {Q0,Q1}, [R0]!        @ Load block 1 from array a
        VADD.I32    Q12, Q8, Q10          @ Process previous block 1
        VLD1.32     {Q2,Q3}, [R1]!        @ Load block 1 from array b
        VADD.I32    Q13, Q9, Q11          @ Process previous block 2
        
        @ Load second block while storing previous results
        VLD1.32     {Q4,Q5}, [R0]!        @ Load block 2 from array a
        VST1.32     {Q12,Q13}, [R2]!      @ Store previous results
        VLD1.32     {Q6,Q7}, [R1]!        @ Load block 2 from array b
        
        @ Process current blocks
        VADD.I32    Q8, Q0, Q2            @ Process current block 1
        VADD.I32    Q9, Q1, Q3            @ Process current block 1
        VADD.I32    Q10, Q4, Q6           @ Process current block 2
        VADD.I32    Q11, Q5, Q7           @ Process current block 2
        
        @ Store current results
        VST1.32     {Q8,Q9,Q10,Q11}, [R2]!  @ Store all current results
        
        SUBS        R4, R4, #32           @ Decrement counter
        BGE         multi_acc_loop
```

### IRPC (repeat character)

```asm
.irpc   round, 123456789101112
    // Code block that will be repeated 12 times
.endr
```
 equivalent to manually writing out 12 nearly identical rounds, but with less code. Each time the loop:

- round will take on the values: '1', '2', '3', ..., '12'
- code block inside will be repeated 12 times
- each repetition can use the current round value if needed

```asm
.irpc   x, 123
    // will expand to three separate blocks
    mov r0, #&x   // First iteration: mov r0, #1
                  // Second iteration: mov r0, #2
                  // Third iteration: mov r0, #3
.endr
```

#### AES example

```arm
.irpc   round, 123456789101112
    // Parallel SubBytes
    aese.8  v0.16b, v2.16b
    
    // AddRoundKey
    eor     v0.16b, v0.16b, v2.16b
    
    // ShiftRows and MixColumns using crypto extensions
    aesmc.8 v0.16b, v0.16b
.endr
```
expands to 12 sets of these three instructions, effectively unrolling the AES encryption rounds at compile time.

[NEXT -> math](11_math.md)

<div align="center">
    <img src="../img/argo-mascot.jpg" alt="Logo">
</div>
<p align="center">
    <img src="https://raw.githubusercontent.com/bornmay/bornmay/Update/svg/Bottom.svg" alt="Github Stats" />
</p>
<p align="right">(<a href="#top">back to top</a>)</p>