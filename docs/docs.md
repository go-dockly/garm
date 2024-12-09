
## `Σ(°△°)ꪱꪱ` Documentation ♨

generate instructions like so:
```go
// Vector add
instr := op.NewVectorOperation("VADD.F32")

// Compare
instr := op.NewBinaryOperation("CMP",
    CreateRegisterOperand("R0"),
    CreateImmediateOperand("#0"),
    nil,
    "Compare with zero")

// Branch
instr := op.NewBranch("B", ".L1", "Branch to loop start")
```

allocate registers like so:
```go
// Initialize allocator
allocator := NewAllocator()
...
// Assign a register for myVar
reg, err := allocator.GetRegister("myVar")
if err != nil {
    Handle allocation error
}
// Use the register
...
// Free the register when done
allocator.FreeRegister("myVar")
```