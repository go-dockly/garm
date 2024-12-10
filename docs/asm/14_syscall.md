## SysCalls

All registers except those required to return values are preserved. 
System calls return results in x0 while everything else remains the same, including the conditional flags. 
The main system instruction for shellcodes on linux is supervisor call (SVC)

```asm
   // read condition flags
  .equ OVERFLOW_FLAG, 1 << 28
  .equ CARRY_FLAG,    1 << 29
  .equ ZERO_FLAG,     1 << 30
  .equ NEGATIVE_FLAG, 1 << 31

  mrs    x0, nzcv

   // set C (arry) flag
  mov    w0, CARRY_FLAG
  msr    nzcv, x0
```

### Executing a shell

```asm
 // 40 bytes

    .arch armv8-a

    .include "macro.inc"

    .global _start
    .text

_start:
     // execve("/bin/sh", NULL, NULL);
    mov    x8, SYS_execve
    mov    x2, xzr            // NULL
    mov    x1, xzr            // NULL
    movq   x3, BINSH          // "/bin/sh"
    str    x3, [sp, -16]!     // stores string on stack
    mov    x0, sp
    svc    0
```

### Executing a command

```asm
    .arch armv8-a
    .align 4

    .include "macro.inc"

    .global _start
    .text

_start:
     // execve("/bin/sh", {"/bin/sh", "-c", cmd, NULL}, NULL);
    movq   x0, BINSH              // x0 = "/bin/sh\0"
    str    x0, [sp, -64]!
    mov    x0, sp
    mov    x1, 0x632D             // x1 = "-c"
    str    x1, [sp, 16]
    add    x1, sp, 16
    adr    x2, cmd                // x2 = cmd
    stp    x0, x1,  [sp, 32]      // store "-c", "/bin/sh"
    stp    x2, xzr, [sp, 48]      // store cmd, NULL
    mov    x2, xzr                // penv = NULL
    add    x1, sp, 32             // x1 = argv
    mov    x8, SYS_execve
    svc    0
cmd:
    .asciz "echo Hello, World!"
```

### Check for syscall errors

linux sysCall to yield to another thread on the same core voluntarily
```asm
thread_yield:
    INIT_FRAME
    mov     x8, #93             // `sched_yield` linux syscall
    svc     #0                  // call linux
    cmn     x0, #4095           // check if return value indicates error
    b.hi    error_handler       // if so branch
    RESTORE_STACK
```

[NEXT -> network](15_network.md)

<div align="center">
  <img src="../img/argo-mascot.jpg" alt="Logo">
</div>
<p align="center">
    <img src="https://raw.githubusercontent.com/bornmay/bornmay/Update/svg/Bottom.svg" alt="Github Stats" />
</p>
<p align="right">(<a href="#top">back to top</a>)</p>

