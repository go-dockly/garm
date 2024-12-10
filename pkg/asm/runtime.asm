// Garbage Collector Constants and Structures
.equ COLOR_WHITE,    0    // Unreachable/not visited
.equ COLOR_GREY,     1    // Reachable but not scanned
.equ COLOR_BLACK,    2    // Reachable and scanned
.equ PTR_BITS,       48   // Useful pointer bits
.equ PAGE_SIZE,      4096
.equ CARD_SIZE,      512  // Write barrier card table granularity
.equ MIN_HEAP_SIZE,  64 * 1024 * 1024  // 64MB

// Object Header (16 bytes)
.struct 0
obj_size:    .space 8     // Object size including header
obj_type:    .space 4     // Type descriptor index
obj_color:   .space 1     // GC color
obj_flags:   .space 3     // Flags (finalizer, etc)
.end

// Heap Page Header (64 bytes)
.struct 0
page_next:   .space 8     // Next page in list
page_prev:   .space 8     // Previous page in list
page_start:  .space 8     // Start of allocatable space
page_end:    .space 8     // End of allocatable space
page_free:   .space 8     // Free list head
page_used:   .space 8     // Bytes allocated
page_mark:   .space 8     // Mark bits
page_spans:  .space 8     // Span information
.end

// Write Barrier Card Table
.data
card_table:  .space (1 << 20)  // 1MB card table for 512GB address space

// GC State
.struct 0
gc_enabled:  .space 1     // GC enabled flag
gc_phase:    .space 1     // Current GC phase
gc_workers:  .space 2     // Number of GC workers
gc_trigger:  .space 8     // Next GC trigger point
gc_heap:     .space 8     // Heap start
gc_limit:    .space 8     // Heap limit
gc_alloc:    .space 8     // Allocation pointer
gc_roots:    .space 8     // Root set pointer
.end

.text
// Initialize Garbage Collector
.macro GC_INIT heap_size=MIN_HEAP_SIZE
    // Allocate heap space
    mov     x0, #0
    ldr     x1, =\heap_size
    mov     x2, #(PROT_READ | PROT_WRITE)
    mov     x3, #(MAP_PRIVATE | MAP_ANONYMOUS)
    mov     x8, #SYS_mmap
    svc     #0
    
    // Initialize GC state
    adr     x19, gc_state
    str     x0, [x19, #gc_heap]
    add     x1, x0, \heap_size
    str     x1, [x19, #gc_limit]
    str     x0, [x19, #gc_alloc]
    
    // Initialize card table
    adr     x0, card_table
    mov     x1, #0
    ldr     x2, =(1 << 20)
    bl      memset
    
    // Enable GC
    mov     w0, #1
    strb    w0, [x19, #gc_enabled]
.endm

// Write Barrier Implementation
.macro WRITE_BARRIER ptr, val
    // Get card table index
    lsr     x15, \ptr, #CARD_SIZE_SHIFT
    adr     x16, card_table
    mov     w17, #1
    strb    w17, [x16, x15]      // Mark card as dirty
    
    // Perform write
    str     \val, [\ptr]
    
    // Check if we're in GC phase
    adr     x15, gc_state
    ldrb    w16, [x15, #gc_phase]
    cbz     w16, 1f              // Skip if not in GC
    
    // Handle concurrent write
    mov     x0, \val
    bl      shade_grey           // Make value grey if needed
    
1:  // Continue
.endm

// Mark Phase Implementation
.macro MARK_PHASE
    // Save registers
    INIT_FRAME 6
    
    // Start marking from roots
    adr     x19, gc_state
    ldr     x20, [x19, #gc_roots]
    
    // Process root set
1:  cbz     x20, 2f              // No more roots
    ldr     x0, [x20]            // Load root pointer
    bl      mark_grey            // Mark as grey
    ldr     x20, [x20, #8]       // Next root
    b       1b
    
2:  // Process grey objects
    bl      process_grey_objects
    
    SAFE_STACK_RET 6
.endm

// Process Grey Objects
process_grey_objects:
    INIT_FRAME 4
    
    // While grey set not empty
1:  bl      get_grey_object
    cbz     x0, 2f               // No more grey objects
    
    // Process object
    mov     x19, x0              // Save object pointer
    ldr     w20, [x19, #obj_type]
    
    // Get type descriptor
    adr     x21, type_table
    ldr     x21, [x21, x20, lsl #3]
    
    // Scan object fields
    mov     x0, x19
    mov     x1, x21
    bl      scan_object_fields
    
    // Mark object black
    mov     w0, #COLOR_BLACK
    strb    w0, [x19, #obj_color]
    
    b       1b                   // Continue processing
    
2:  SAFE_STACK_RET 4

// Scan Object Fields
scan_object_fields:
    INIT_FRAME 4
    
    // Get field information
    ldr     x19, [x1, #type_fields]
    ldr     w20, [x1, #type_field_count]
    
    // For each field
1:  cbz     w20, 2f              // No more fields
    ldr     x21, [x19]           // Load field offset
    ldr     x22, [x0, x21]       // Load field value
    
    // If pointer field
    ldr     w23, [x19, #8]       // Load field type
    tbz     w23, #0, 3f          // Skip if not pointer
    
    // Mark field grey
    mov     x0, x22
    bl      mark_grey
    
3:  add     x19, x19, #16        // Next field
    sub     w20, w20, #1
    b       1b
    
2:  SAFE_STACK_RET 4

// Sweep Phase Implementation
.macro SWEEP_PHASE
    INIT_FRAME 4
    
    // Get heap bounds
    adr     x19, gc_state
    ldr     x20, [x19, #gc_heap]
    ldr     x21, [x19, #gc_limit]
    
    // For each page
1:  cmp     x20, x21
    b.ge    2f                   // Done
    
    // Get page header
    mov     x0, x20
    bl      sweep_page
    
    // Next page
    ldr     x22, [x20, #page_next]
    mov     x20, x22
    b       1b
    
2:  SAFE_STACK_RET 4
.endm

// Sweep Single Page
sweep_page:
    INIT_FRAME 4
    
    // Get page bounds
    ldr     x19, [x0, #page_start]
    ldr     x20, [x0, #page_end]
    
    // For each object in page
1:  cmp     x19, x20
    b.ge    2f                   // Done with page
    
    // Check color
    ldrb    w21, [x19, #obj_color]
    cmp     w21, #COLOR_WHITE
    b.ne    3f                   // Skip if not white
    
    // Free white object
    mov     x0, x19
    bl      free_object
    
3:  // Next object
    ldr     x22, [x19, #obj_size]
    add     x19, x19, x22
    b       1b
    
2:  SAFE_STACK_RET 4

// Concurrent GC Worker
.macro GC_WORKER
    INIT_FRAME 2
    
1:  // Wait for work
    bl      wait_for_gc_work
    
    // Check GC phase
    adr     x19, gc_state
    ldrb    w20, [x19, #gc_phase]
    
    // Handle different phases
    cmp     w20, #GC_MARK_PHASE
    b.eq    2f
    cmp     w20, #GC_SWEEP_PHASE
    b.eq    3f
    b       1b
    
2:  // Mark phase work
    bl      process_grey_objects
    b       1b
    
3:  // Sweep phase work
    bl      sweep_worker
    b       1b
    
    SAFE_STACK_RET 2
.endm

// Allocator with GC Integration
.macro GC_ALLOC size, type
    // Try fast path allocation
    adr     x19, gc_state
    ldr     x20, [x19, #gc_alloc]
    add     x21, x20, \size
    ldr     x22, [x19, #gc_trigger]
    cmp     x21, x22
    b.gt    1f                   // Need GC
    
    // Fast path
    str     x21, [x19, #gc_alloc]
    mov     x0, x20
    
    // Initialize object header
    str     \size, [x0, #obj_size]
    str     \type, [x0, #obj_type]
    mov     w1, #COLOR_BLACK
    strb    w1, [x0, #obj_color]
    b       2f
    
1:  // Trigger GC
    bl      trigger_gc
    
    // Retry allocation
    GC_ALLOC \size, \type
    
2:  // Done
.endm

// Root Registration
.macro REGISTER_ROOT ptr
    adr     x19, gc_state
    ldr     x20, [x19, #gc_roots]
    
    // Allocate root node
    mov     x0, #16
    bl      malloc
    
    // Initialize node
    str     \ptr, [x0]
    str     x20, [x0, #8]
    
    // Update root list
    str     x0, [x19, #gc_roots]
.endm

// Finalizer Support
.macro REGISTER_FINALIZER obj, func
    // Add finalizer to object
    ldr     w0, [x19, #obj_flags]
    orr     w0, w0, #FLAG_FINALIZER
    str     w0, [x19, #obj_flags]
    
    // Register finalizer function
    adr     x0, finalizer_table
    str     \func, [x0, \obj, lsl #3]
.endm