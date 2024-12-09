
## Concurrency

minimal bare-metal assembly example that illustrates [what multicore assembly looks like](https://stackoverflow.com/questions/980999/what-does-multicore-assembly-language-look-like/33651438#33651438)

- [ARM syncronization primitives](https://developer.arm.com/documentation/dht0008/a/arm-synchronization-primitives/practical-uses/power-saving-features?lang=en)
- [Memory barriers](https://preshing.com/20120710/memory-barriers-are-like-source-control-operations/)
- [Lock free programming](https://www.boost.org/doc/libs/1_63_0/doc/html/lockfree.html)
- [futex](https://www.collabora.com/news-and-blog/blog/2022/02/08/landing-a-new-syscall-part-what-is-futex/)

### Semaphores
Semaphores are synchronization primitives used to control access to a common resource in concurrent programming.

```arm
.macro CHAN_SIGNAL_SEMAPHORE sem_ptr
    mov     w1, #1
    str     w1, [\sem_ptr]         // Set semaphore
    sev                            // Signal event
.endm
.macro CHAN_WAIT_SEMAPHORE sem_ptr
    mov     x0, \sem_ptr
    mov     x1, #0
    mov     x2, #0
    mov     x8, #98
    svc     #0
.endm
```
CHAN_WAIT_SEMAPHORE takes argument, sem_ptr, expected to be a pointer to a semaphore. 

`mov x0, \sem_ptr` 
- `\` before sem_ptr is macro syntax to reference the argument passed to it.

`mov x1, #0`
- sets register x1 to 0 or false indicating non-blocking.

`mov x8, #98`
- 98 in Arm Linux means futex (Fast Userspace muTEX - kernel system call implementing locking and synchronization primitives)

`svc #0`
- triggers a supervisor call, which switches the processor into supervisor mode to execute the sysCall in x8.

Fast Path (CHAN_SIGNAL_SEMAPHORE):

- Direct memory write for signaling
- `sev` for immediate notification
- No syscall overhead
- Good for when contention is low

Slow Path (CHAN_WAIT_SEMAPHORE):

- Linux Kernel handles thread sleeping/waking through
- futex syscall (fast userspace operations with kernel fallback) for waiting

```arm
// Thread 1 - Waiter
check_sem:
    ldr     w0, [sem_ptr]        // Check semaphore value
    cbz     w0, do_wait          // If zero, need to wait
    b       got_sem              // Otherwise, we got it

do_wait:
    CHAN_WAIT_SEMAPHORE sem_ptr  // Sleep until signaled
    b       check_sem            // Check again after wakeup

// Thread 2 - Signaler
CHAN_SIGNAL_SEMAPHORE sem_ptr    // Signal waiting thread
```

### Wait For Event (WFE)
conceptually equivalent to

```c
while (!event_has_occurred) /*do nothing*/;
```

except it turns the CPU off instead of running in a loop.

Several things that can interrupt WFE
- explicit wake up event from another CPU
- an interrupt 

If an interrupt happens during WFE
- processor switches to IRQ or FIQ mode
- jumps to the IRQ or FIQ handler
- address of WFE instruction (plus offset) is placed in the link register
- if the CPU is triggered by a WAKE_UP event
- execution proceeds with the next instruction after WFE.

ARM assembly WFE concept
```arm
wait_loop:
    ldr     w0, [sem_ptr]        // Load current value
    cbnz    w0, exit_wait        // Exit if value is non-zero
    wfe                          // Sleep until event
    b       wait_loop           // Check again after wake
exit_wait:
```
handle wake-up scenarios and potential race conditions

```arm
.macro WAIT_FOR_EVENT sem_ptr
    // Save state if needed
    stp     x29, x30, [sp, #-16]!
    
wait_loop:
    // Load-Exclusive for atomic operation
    ldaxr   w0, [\sem_ptr]       // Load with acquire semantics
    cbnz    w0, exit_wait        // Check if signaled
    
    // Ensure store atomicity
    dmb     ish                  // Data Memory Barrier
    
    // Wait for event
    wfe                          // Power-efficient wait
    
    // Optional: Check for interrupt
    mrs     x1, DAIF            // Get interrupt status
    tst     x1, #(1 << 7)       // Check IRQ bit
    b.ne    handle_interrupt
    
    b       wait_loop           // Check again if spurious wake-up

exit_wait:
    // Clear the event
    mov     w1, #0
    stlxr   w2, w1, [\sem_ptr]   // Store with release semantics
    cbnz    w2, exit_wait        // Retry if store failed
    
    // Memory barrier after operation
    dmb     ish
    
    // Restore state
    ldp     x29, x30, [sp], #16
    ret

handle_interrupt:
    // Handle interrupt case
    // IRQ handler will return to next instruction
    b       wait_loop
.endm

.macro SIGNAL_EVENT sem_ptr
    // Ensure atomicity
    dmb     ish                  // Memory barrier before signal
    
    // Set event
    mov     w1, #1
    str     w1, [\sem_ptr]       // Store new value
    
    // Ensure visibility
    dmb     ish                  // Memory barrier after store
    
    sev                         // Send event to all cores
.endm
```

Memory Ordering:

```arm
dmb     ish                  // Data Memory Barrier
```
ensure memory operations complete in order
to prevent race conditions between cores


Atomic Operations:
```arm
ldaxr   w0, [\sem_ptr]       // Load-Exclusive
stlxr   w2, w1, [\sem_ptr]   // Store-Exclusive
```
use exclusive monitors for atomic operations


Event Handling:
```arm
wfe                         // Wait For Event
sev                         // Send Event
```
WFE puts processor in low-power state
SEV wakes up other processors

Interrupt Awareness:
```arm
mrs     x1, DAIF            // Get interrupt status
tst     x1, #(1 << 7)       // Check IRQ bit
```
Check for interrupt occurrence handle IRQ/FIQ

Usage example:
```arm
// Initialize event
    mov     w0, #0
    str     w0, [event_ptr]

// Thread 1: Waiter
    WAIT_FOR_EVENT event_ptr

// Thread 2: Signaler
    SIGNAL_EVENT event_ptr
```

### Common pitfalls

Memory Barriers:

- always use appropriate DMB instructions to ensure visibility across cores

Spurious Wake-ups:

- always re-check condition after WFE don't assume wake means event occurred

Race Conditions:

- [Data Race Patterns in Go](https://www.uber.com/en-SE/blog/data-race-patterns-in-go/)
- use exclusive operations for atomic access

Interrupt Handling:

- save state then return from interrupt

Power Consumption:

- WFE is more efficient than spinning but don't WFE in critical timing sections