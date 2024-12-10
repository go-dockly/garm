## Loops

```asm
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

### `for(i=0; i<=10; i++)` increment and check condition before every loop
```asm
        MOV     X0, #0      // i
        MOV     X1, #0      // total
for     CMP     X0, #10
        BGE     fordone     // first branch
        ADD     X1, X1, X0
        ADD     X0, X0, #1
        B       for         // second branch
fordone
```
### `for(i=10; i>0; i--)` decrement with using only one branch at the end`
```asm
        MOV     X0, #10     // i
        MOV     X1, #0      // total
ford    ADD     X1, X1, X0  
        SUBS    X0, X0, #1
        BNE     ford        // first branch
forddone
```
### `while(i>0)` initial branch only run once
```asm
        MOV     X0, #10     // i
        MOV     X1, #0      // total
        B       wtest       // initial branch
while   ADD     X1, X1, X0
        SUB     X0, X0, #1
wtest   CMP     X0, #0
        B.GE    while       // first branch
```

### `do..while(i>0)` same as while but without initial branch
 
```asm
        MOV     X0, #10     // i
        MOV     X1, #0      // total
dwhile  ADD     X1, X1, X0
        SUB     X0, X0, #1
dwtest  CMP     X0, #0
        B.GE    while       // first branch
        
        END
```

## NEON 

process 4 elements per iteration
```asm
basic_loop:
    LD1     {V0.4S}, [X0], #16      // Load 4 elements from array a
    LD1     {V1.4S}, [X1], #16      // Load 4 elements from array b
    ADD     V0.4S, V0.4S, V1.4S     // Add vectors
    ST1     {V0.4S}, [X2], #16      // Store result
    
    SUBS    X4, X4, #4              // Decrement counter
    B.GE    basic_loop              // Branch if greater or equal
```

### Unrolled 2x 
process 8 elements per iteration
```asm
unrolled_2x_loop:
    // Load 8 elements (2 quad words) from each array
    LD1     {V0.4S, V1.4S}, [X0], #32    // Load 8 elements from array a
    LD1     {V2.4S, V3.4S}, [X1], #32    // Load 8 elements from array b
    
    // Perform addition
    ADD     V0.4S, V0.4S, V2.4S     // Add first 4 elements
    ADD     V1.4S, V1.4S, V3.4S     // Add second 4 elements
    
    // Store results
    ST1     {V0.4S, V1.4S}, [X2], #32    // Store 8 results
    
    SUBS    X4, X4, #8              // Decrement counter
    B.GE    unrolled_2x_loop        // Branch if greater or equal
```

### Unrolled 4x with prefetch 
process 16 elements per iteration
```asm
unrolled_4x_loop:
    // Prefetch next iterations
    PRFM    PLDL1KEEP, [X0, #64]     // Prefetch array a
    PRFM    PLDL1KEEP, [X1, #64]     // Prefetch array b
    
    // Load 16 elements (4 quad words) from each array
    LD1     {V0.4S, V1.4S}, [X0], #32    // Load first 8 elements from a
    LD1     {V2.4S, V3.4S}, [X0], #32    // Load second 8 elements from a
    LD1     {V8.4S, V9.4S}, [X1], #32    // Load first 8 elements from b
    LD1     {V10.4S, V11.4S}, [X1], #32  // Load second 8 elements from b
    
    // Perform additions
    ADD     V0.4S, V0.4S, V8.4S     // Add first 4 elements
    ADD     V1.4S, V1.4S, V9.4S     // Add second 4 elements
    ADD     V2.4S, V2.4S, V10.4S    // Add third 4 elements
    ADD     V3.4S, V3.4S, V11.4S    // Add fourth 4 elements
    
    // Store results
    ST1     {V0.4S, V1.4S}, [X2], #32    // Store first 8 results
    ST1     {V2.4S, V3.4S}, [X2], #32    // Store second 8 results
    
    SUBS    X4, X4, #16             // Decrement counter
    B.GE    unrolled_4x_loop        // Branch if greater or equal
```

### Interleaved loads/compute with dual accumulators
```asm
interleaved_loop:
    // Load first set while processing previous data
    LD1     {V0.4S, V1.4S}, [X0], #32     // Load from array a
    ADD     V8.4S, V4.4S, V6.4S           // Process previous data (first set)
    LD1     {V2.4S, V3.4S}, [X1], #32     // Load from array b
    ADD     V9.4S, V5.4S, V7.4S           // Process previous data (second set)
    
    // Store previous results while loading new data
    ST1     {V8.4S, V9.4S}, [X2], #32     // Store previous results
    LD1     {V4.4S, V5.4S}, [X0], #32     // Load next set from array a
    LD1     {V6.4S, V7.4S}, [X1], #32     // Load next set from array b
    
    // Continue with current data
    ADD     V10.4S, V0.4S, V2.4S          // Process current data (first set)
    ADD     V11.4S, V1.4S, V3.4S          // Process current data (second set)
    ST1     {V10.4S, V11.4S}, [X2], #32   // Store current results
    
    SUBS    X4, X4, #16                   // Decrement counter
    B.GE    interleaved_loop              // Branch if greater or equal
```

##  Pipelining
NOTE: Needs init and epilogue code
```asm
pipelined_loop:
    // Load
    LD1     {V0.4S, V1.4S}, [X0], #32     // Load new data from array a
    
    // Load + Process previous
    LD1     {V2.4S, V3.4S}, [X1], #32     // Load new data from array b
    ADD     V8.4S, V4.4S, V6.4S           // Process previous iteration
    
    // Process current + Store previous
    ADD     V9.4S, V0.4S, V2.4S           // Process current data
    ST1     {V8.4S}, [X2], #16            // Store previous result
    
    // Update pipeline registers
    MOV     V4.16B, V0.16B                // Move current to previous
    MOV     V6.16B, V2.16B                // Move current to previous
    
    SUBS    X4, X4, #8                    // Decrement counter
    B.GE    pipelined_loop                // Branch if greater or equal
```

### Advanced SIMD with multiple accumulators and prefetch
Concurrency optimized for high-end ARM processors with multiple NEON units
```asm
multi_acc_loop:
    // Prefetch next cache lines
    PRFM    PLDL1KEEP, [X0, #128]     // Prefetch array a
    PRFM    PLDL1KEEP, [X1, #128]     // Prefetch array b
    
    // Load first block while processing previous block
    LD1     {V0.4S, V1.4S}, [X0], #32 // Load block 1 from array a
    ADD     V12.4S, V8.4S, V10.4S     // Process previous block 1
    LD1     {V2.4S, V3.4S}, [X1], #32 // Load block 1 from array b
    ADD     V13.4S, V9.4S, V11.4S     // Process previous block 2
    
    // Load second block while storing previous results
    LD1     {V4.4S, V5.4S}, [X0], #32 // Load block 2 from array a
    ST1     {V12.4S, V13.4S}, [X2], #32 // Store previous results
    LD1     {V6.4S, V7.4S}, [X1], #32 // Load block 2 from array b
    
    // Process current blocks
    ADD     V8.4S, V0.4S, V2.4S       // Process current block 1
    ADD     V9.4S, V1.4S, V3.4S       // Process current block 1
    ADD     V10.4S, V4.4S, V6.4S      // Process current block 2
    ADD     V11.4S, V5.4S, V7.4S      // Process current block 2
    
    // Store current results
    ST1     {V8.4S, V9.4S, V10.4S, V11.4S}, [X2], #64 // Store all current results
    
    SUBS    X4, X4, #32               // Decrement counter
    B.GE    multi_acc_loop            // Branch if greater or equal
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
    mov x0, #&x   // First iteration:  mov x0, #1
                  // Second iteration: mov x0, #2
                  // Third iteration:  mov x0, #3
.endr
```

#### AES example

```asm
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