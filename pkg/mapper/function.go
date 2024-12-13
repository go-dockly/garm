package mapper

import (
	"fmt"

	"github.com/algoboyz/garm/pkg/alloc"
	"github.com/algoboyz/garm/pkg/ir"
	"golang.org/x/tools/go/ssa"
)

func (m *SSAMapper) MapFunction(fn *ssa.Function) (irFunc *ir.Function, err error) {
	if fn.Blocks == nil {
		fmt.Println("Skipping external function:", fn.Name())
		return nil, nil // Skip external functions
	}
	// Reset mapper state
	// m.reset(ssaFunc)
	m.currentFunc = fn

	// Generate labels for basic blocks
	// m.generateBlockLabels()
	params, err := m.processParams(fn.Params)
	if err != nil {
		return nil, fmt.Errorf("processing parameters: %w", err)
	}
	m.currentIR = ir.NewFunction(fn.Name(), m.debug)
	m.currentIR.Params = params
	// Allocate registers for parameters and return values
	// err := m.alloc.AllocateFn(irFunc)
	// if err != nil {
	// 	return nil, fmt.Errorf("allocating registers: %w", err)
	// }

	// Create function prologue
	m.currentIR.Blocks = append(m.currentIR.Blocks, m.prologue()...)

	// Iterate through SSA instructions
	for _, block := range fn.Blocks {
		m.currentBlock = block
		if err = m.MapBlock(block); err != nil {
			return nil, fmt.Errorf("mapping block %d: %w", block.Index, err)
		}
	}
	// Create function epilogue
	m.currentIR.Blocks = append(m.currentIR.Blocks, m.epilogue()...)

	return m.currentIR, nil
}

func (m *SSAMapper) processParams(params []*ssa.Parameter) (map[string]alloc.Location, error) {
	irParams := make(map[string]alloc.Location)
	for _, param := range params {
		val := param.Object()
		prim := alloc.NewPrimitive(val.Type().String())
		typ := alloc.NewType(param.Name(), val.String(), prim, 64)
		paramAlloc, err := m.alloc.AllocateRegister(typ)
		if err != nil {
			return nil, fmt.Errorf("allocating parameter: %w", err)
		}
		irParams[param.Name()] = paramAlloc

	}
	return irParams, nil
}

// When creating parameter types in the mapper:
// func (m *SSAMapper) MapParams(params []*ssa.Parameter) (map[string]*alloc.ResourceAllocation, error) {
// 	irParams := make(map[string]*alloc.ResourceAllocation)
// 	for _, param := range params {
// 		val := param.Object()
// 		typ := &BasicType{
// 			name:    param.Name(),
// 			typeStr: val.Type().String(),
// 			kind:    inferKindFromSSAType(val.Type()),
// 			size:    64, // Or infer from type
// 			flags:   IsParameterFlag,
// 		}
// 		allocation, err := m.alloc.Allocate(typ)
// 		if err != nil {
// 			return nil, fmt.Errorf("allocating parameter: %w", err)
// 		}
// 		irParams[param.Name()] = allocation
// 	}
// 	return irParams, nil
// }

func (m *SSAMapper) prologue() (instructions []ir.Instruction) {
	if m.currentIR.Label == "main" {
		instructions = ir.PrologueMain()
	} else {
		instructions = ir.FuncPrologue(m.currentIR.Label)
	}
	// Create a new frame
	// frame := m.currentIR.Frames.PushFrame(true) // true means we need a frame pointer

	// // Save some registers
	// offset, _ := m.currentIR.Frames.SaveRegister(reg.LR)
	// fmt.Printf("LR saved at sp+%d\n", offset)

	// // Allocate local variables
	// localOffset, _ := m.currentIR.Frames.AllocateStackSlot(8, 8) // 8-byte variable, 8-byte aligned

	// spew.Dump(frame)
	// spew.Dump(localOffset)

	// // Generate prologue
	// instructions = m.currentIR.Frames.GenerateFrameSetup()
	// for _, inst := range instructions {
	// 	fmt.Println(inst)
	// }

	return instructions
}

func (m *SSAMapper) epilogue() (instructions []ir.Instruction) {
	// Generate epilogue
	if m.currentIR.Label == "main" {
		instructions = ir.EpilogueMain()
	} else {
		instructions = ir.FuncEpilogue()
	}
	// instructions = m.currentIR.Frames.GenerateFrameTeardown()
	// for _, inst := range instructions {
	// 	fmt.Println(inst)
	// }
	// // Clean up
	// m.currentIR.Frames.PopFrame()
	return instructions
}
