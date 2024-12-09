// Input:
//   X0 = interface value
//   X1 = method index
// Output:
//   X0 = method address
//   All other registers preserved
interface_dispatch:
    // Load itab pointer from interface
    LDR     X0, [X0]
    
    // Check if itab is nil
    CBZ     X0, .Lnil_panic
    
    // Load method address from vtable
    ADD     X0, X0, #offset_vtable
    LDR     X0, [X0, X1, LSL #3]    // Scale index by 8 (pointer size)
    
    // Return method address
    RET

.Lnil_panic:
    // Handle nil interface panic
    // This will be replaced with a call to the arGo runtime
    B       runtime.nilPanic