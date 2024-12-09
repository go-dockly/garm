package alloc

import (
	"errors"
	"sync"

	"github.com/algoboyz/garm/pkg/reg"
)

// ARM64Type represents the core type information needed for ARM64 code generation
type ARM64Type interface {
	// Core properties
	String() string              // var name
	Size() int                   // size in bytes
	Alignment() int              // required alignment
	Register() reg.RegisterClass // register class
}

// MemoryLayout handles memory allocation concerns
type MemoryLayout interface {
	Offset() int
	RequiredAlignment() int
	TotalSize() int
}

// Location represents where a value is stored
type Location interface {
	GetRegister() *reg.Register
	IsRegister() bool
	GetMemory() *MemoryLocation
	IsMemory() bool
	String() string
}

// Allocator handles resource allocation
type Allocator interface {
	AllocateRegister(ARM64Type) (Location, error)
	AllocateStack(MemoryLocation) (Location, error)
	Free(Location)
}

// SimpleAllocator provides a basic implementation of the Allocator interface
type SimpleAllocator struct {
	mu sync.Mutex

	// Register pools
	gprPool []reg.Register
	fprPool []reg.Register
	vecPool []reg.Register

	// Stack allocation
	stackSize   int
	stackOffset int
}

// NewAllocator creates a new allocator with predefined register pools
func NewAllocator() *SimpleAllocator {
	a := &SimpleAllocator{
		stackSize:   4096, // Default stack frame size
		stackOffset: 16,   // Initial offset after frame setup
	}

	// Initialize register pools
	// GPRs: x0-x18 (x19-x28 are callee-saved, x29-x31 special purpose)
	for i := uint8(0); i < 19; i++ {
		a.gprPool = append(a.gprPool, reg.Register{ID: i, Class: reg.RegisterClassGPR})
	}

	// FPRs: d0-d7 (d8-d15 are callee-saved)
	for i := uint8(0); i < 8; i++ {
		a.fprPool = append(a.fprPool, reg.Register{ID: i, Class: reg.RegisterClassFPR})
	}

	// Vector: v0-v7 (v8-v15 are callee-saved)
	for i := uint8(0); i < 8; i++ {
		a.vecPool = append(a.vecPool, reg.Register{ID: i, Class: reg.RegisterClassVec})
	}

	return a
}

func (a *SimpleAllocator) AllocateRegister(t ARM64Type) (Location, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	var dst reg.Register
	var pool *[]reg.Register

	switch t.Register() {
	case reg.RegisterClassGPR:
		if len(a.gprPool) == 0 {
			return nil, errors.New("no available GPR registers")
		}
		pool = &a.gprPool
	case reg.RegisterClassFPR:
		if len(a.fprPool) == 0 {
			return nil, errors.New("no available FPR registers")
		}
		pool = &a.fprPool
	case reg.RegisterClassVec:
		if len(a.vecPool) == 0 {
			return nil, errors.New("no available vector registers")
		}
		pool = &a.vecPool
	default:
		return nil, errors.New("unknown register class")
	}

	// Pop a register from the pool
	dst = (*pool)[len(*pool)-1]
	*pool = (*pool)[:len(*pool)-1]

	return &locationImpl{register: &dst}, nil
}

func (a *SimpleAllocator) AllocateStack(t MemoryLocation) (Location, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	size := t.Size
	alignment := t.Alignment

	// Align the stack offset
	if alignment > 0 {
		a.stackOffset = (a.stackOffset + alignment - 1) & ^(alignment - 1)
	}

	// Check if we have enough stack space
	if a.stackOffset+size > a.stackSize {
		return nil, errors.New("stack overflow")
	}

	mem := &MemoryLocation{
		Offset:    a.stackOffset,
		Size:      size,
		Alignment: alignment,
	}

	// Update stack offset
	a.stackOffset += size

	return &locationImpl{memory: mem}, nil
}

func (a *SimpleAllocator) Free(loc Location) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if l, ok := loc.(*locationImpl); ok && l.IsRegister() {
		switch l.register.Class {
		case reg.RegisterClassGPR:
			a.gprPool = append(a.gprPool, *l.register)
		case reg.RegisterClassFPR:
			a.fprPool = append(a.fprPool, *l.register)
		case reg.RegisterClassVec:
			a.vecPool = append(a.vecPool, *l.register)
		}
	}
	// Memory locations don't need to be freed explicitly
	// They'll be cleaned up when the stack frame is destroyed
}

// Helper function to get the current stack offset
func (a *SimpleAllocator) CurrentStackOffset() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.stackOffset
}
