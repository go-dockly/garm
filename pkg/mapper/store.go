package mapper

import (
	"fmt"

	"github.com/algoboyz/garm/pkg/alloc"
	"github.com/algoboyz/garm/pkg/ir"
	"github.com/algoboyz/garm/pkg/op"
	"github.com/algoboyz/garm/pkg/reg"
	"github.com/davecgh/go-spew/spew"
	"go.uber.org/zap"
	"golang.org/x/tools/go/ssa"
)

func (m *SSAMapper) MapStore(v *ssa.Store) error {
	spew.Dump(v)
	// Special handling for init guard
	if isInitGuard(v.Addr) {
		return m.mapInitGuardStore(v)
	}
	// Get the address (destination)
	addr, ok := v.Addr.(*ssa.Alloc)
	if !ok {
		return fmt.Errorf("unsupported store address type: %T", v.Addr)
	}

	// Get the value to store
	val, err := m.MapLiteral(v.Val.Name(), v.Val.Type())
	if err != nil {
		return fmt.Errorf("mapping store value: %w", err)
	}

	// Get allocation info
	dstAlloc, ok := m.currentIR.Locals[addr.Name()]
	if !ok {
		return fmt.Errorf("undefined variable: %s", addr.Name())
	}

	// Allocate a temporary register for the value if needed
	srcReg, err := m.alloc.AllocateRegister(val)
	if err != nil {
		return fmt.Errorf("allocating source register: %w", err)
	}
	defer m.alloc.Free(srcReg)

	// Generate store instruction
	m.currentIR.Blocks = append(m.currentIR.Blocks, ir.Instruction{
		Op:  op.STR, // Or STRB/STRH depending on size
		Dst: srcReg.GetRegister(),
		Src: []reg.Operand{
			reg.NewMemOperand(reg.FP, dstAlloc.GetMemory().Offset), // FP-relative addressing
		},
		Comment: fmt.Sprintf("store %s -> [%s]", v.Val.Name(), addr.Name()),
	})

	return nil
}

// isInitGuard checks if the address operand is an init guard variable
func isInitGuard(addr ssa.Value) bool {
	// Check if it's a global variable
	global, ok := addr.(*ssa.Global)
	if !ok {
		return false
	}

	// Check if it's an init guard by name
	return global.Name() == "init$guard"
}

func (m *SSAMapper) mapInitGuardStore(v *ssa.Store) error {
	name := v.Val.Name() // "*init$guard"
	// For init guards, we typically want to store a boolean true value
	trueVal := alloc.NewType(name, "bool", alloc.Bool, 1)
	srcReg, err := m.alloc.AllocateRegister(trueVal)
	if err != nil {
		return fmt.Errorf("allocating init guard register: %w", err)
	}
	defer m.alloc.Free(srcReg)

	// Get or create a global section for the init guard
	guardSection := m.getOrCreateInitGuardSection()
	spew.Dump(guardSection)
	// Need a temporary register for the address
	ptr := alloc.NewType(name, "bool", alloc.Pointer, 16)
	addrReg, err := m.alloc.AllocateRegister(ptr)
	if err != nil {
		return fmt.Errorf("allocating address register: %w", err)
	}
	defer m.alloc.Free(addrReg)

	// Generate instructions for loading global address and storing
	m.currentIR.Blocks = append(m.currentIR.Blocks,
		// Load page address of global variable
		// adrp x0, init$guard
		ir.Instruction{
			Op:  op.ADRP,
			Dst: addrReg.GetRegister(),
			Src: []reg.Operand{
				reg.NewRegOperand("x0"), // Use x0 as a temporary register
				reg.NewImmediateOperand(name),
			},
			Comment: "load page address of " + name,
		},
		// Add low 12 bits of offset
		ir.Instruction{
			Op:  op.ADD,
			Dst: addrReg.GetRegister(),
			Src: []reg.Operand{
				reg.NewRegOperand(addrReg.GetRegister().String()),
				reg.NewImmediateOperand(name), // todo add x0, x0, :lo12:init$guard // Add low 12 bits
			},
			Comment: "add low bits of init$guard offset",
		},
		// Store the value
		ir.Instruction{
			Op: op.STRB,
			Src: []reg.Operand{
				reg.NewRegOperand(srcReg.GetRegister().String()),
				reg.NewMemOperand(addrReg.GetRegister(), 0), // Base register + 0 offset
			},
			Comment: "store true -> [init$guard]",
		},
	)

	return nil
}

// mapInitGuardStore handles storing to init guard variables
// func (m *SSAMapper) mapInitGuardStore(v *ssa.Store) error {
// 	// For init guards, we typically want to store a boolean true value

// 	// Allocate a temporary register for the true value
// 	trueVal := reg.NewImmediateOperand(1) // 1 represents true
// 	srcReg, err := m.alloc.AllocateRegister(trueVal)
// 	if err != nil {
// 		return fmt.Errorf("allocating init guard register: %w", err)
// 	}
// 	defer m.alloc.Free(srcReg)

// 	// Get or create a global section for the init guard
// 	guardSection := m.getOrCreateInitGuardSection()

// 	// Generate store instruction for init guard
// 	m.currentIR.Blocks = append(m.currentIR.Blocks, ir.Instruction{
// 		Op:  op.STRB, // Use byte store since it's a boolean
// 		Dst: srcReg.GetRegister(),
// 		Src: []reg.Operand{
// 			reg.NewMemOperand(reg.GP, guardSection.GetMemory().Offset),
// 		},
// 		Comment: "store true -> [init$guard]",
// 	})

// 	return nil
// }

// getOrCreateInitGuardSection returns the memory section for init guard
func (m *SSAMapper) getOrCreateInitGuardSection() alloc.Location {
	// Check if we already have an init guard section
	if section, exists := m.currentIR.Globals["init$guard"]; exists {
		return section
	}

	// Create new section for init guard
	section := alloc.MemoryLocation{
		Name:   "init$guard",
		Size:   1, // boolean size
		Offset: m.currentIR.AllocateGlobalOffset(1),
	}

	location, err := m.alloc.AllocateStack(section) // Allocate memory for the section
	if err != nil {
		m.debug.Fatal("allocating init guard section", zap.Error(err))
	}
	// Store in globals map
	m.currentIR.Globals["init$guard"] = location

	return location
}
