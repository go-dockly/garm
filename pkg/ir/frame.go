package ir

// import (
// 	"fmt"
// 	"sync"

// 	"github.com/algoboyz/garm/pkg/alloc"
// 	"github.com/algoboyz/garm/pkg/op"
// 	"github.com/algoboyz/garm/pkg/reg"
// )

// // FrameManager handles stack frame operations
// type FrameManager struct {
// 	mu sync.RWMutex

// 	currentFrame *StackFrame
// 	frames       []*StackFrame
// }

// // StackFrame represents an ARM64 stack frame
// type StackFrame struct {
// 	Size      int
// 	Offset    int
// 	HasFP     bool
// 	SavedRegs map[*reg.Register]int // Map of saved registers to their stack offsets
// 	LocalSize int                   // Size of local variables
// 	SpillSize int                   // Size of spilled registers
// 	ArgSize   int                   // Size of arguments to other functions
// }

// // NewFrameManager creates a new frame manager
// func NewFrameManager() *FrameManager {
// 	return &FrameManager{
// 		frames: make([]*StackFrame, 0),
// 		mu:     sync.RWMutex{},
// 	}
// }

// // PushFrame creates a new stack frame
// func (fm *FrameManager) PushFrame(needsFP bool) *StackFrame {
// 	fm.mu.Lock()
// 	defer fm.mu.Unlock()

// 	frame := &StackFrame{
// 		HasFP:     needsFP,
// 		SavedRegs: make(map[*reg.Register]int),
// 	}

// 	if len(fm.frames) > 0 {
// 		frame.Offset = fm.frames[len(fm.frames)-1].Offset + fm.frames[len(fm.frames)-1].Size
// 	}

// 	fm.frames = append(fm.frames, frame)
// 	fm.currentFrame = frame
// 	return frame
// }

// // PopFrame removes the current stack frame
// func (fm *FrameManager) PopFrame() {
// 	fm.mu.Lock()
// 	defer fm.mu.Unlock()

// 	if len(fm.frames) > 0 {
// 		fm.frames = fm.frames[:len(fm.frames)-1]
// 		if len(fm.frames) > 0 {
// 			fm.currentFrame = fm.frames[len(fm.frames)-1]
// 		} else {
// 			fm.currentFrame = nil
// 		}
// 	}
// }

// // AllocateStackSlot allocates space in the current frame
// func (fm *FrameManager) AllocateStackSlot(size, align int) (int, error) {
// 	fm.mu.Lock()
// 	defer fm.mu.Unlock()

// 	if fm.currentFrame == nil {
// 		return 0, fmt.Errorf("no active stack frame")
// 	}

// 	// Align the current size
// 	offset := alloc.AlignSize(fm.currentFrame.Size, align)
// 	fm.currentFrame.Size = offset + size
// 	fm.currentFrame.LocalSize += size

// 	return offset, nil
// }

// // SaveRegister allocates space for and records a register that needs to be saved
// func (fm *FrameManager) SaveRegister(reg *reg.Register) (int, error) {
// 	fm.mu.Lock()
// 	defer fm.mu.Unlock()

// 	if fm.currentFrame == nil {
// 		return 0, fmt.Errorf("no active stack frame")
// 	}

// 	// Check if register is already saved
// 	if offset, exists := fm.currentFrame.SavedRegs[reg]; exists {
// 		return offset, nil
// 	}

// 	// Allocate space for the register (all registers are 8 bytes on ARM64)
// 	offset := alloc.AlignSize(fm.currentFrame.Size, 8)
// 	fm.currentFrame.Size = offset + 8
// 	fm.currentFrame.SpillSize += 8
// 	fm.currentFrame.SavedRegs[reg] = offset

// 	return offset, nil
// }

// // GenerateFrameSetup generates the assembly instructions for frame setup
// func (fm *FrameManager) GenerateFrameSetup() (instructions []Instruction) {
// 	if fm.currentFrame == nil {
// 		return nil
// 	}

// 	// Save LR if needed
// 	// if _, needsSave := fm.currentFrame.SavedRegs[reg.LR]; needsSave {
// 	// 	// alignedSize := alloc.AlignSize(m.alloc.stackSize+15) & ^15
// 	// 	instructions = append(instructions, Instruction{
// 	// 		Op:  op.STP,
// 	// 		Dst: reg.FP,
// 	// 		Src: []reg.Operand{
// 	// 			reg.NewRegOperand(reg.LR.String()),
// 	// 			reg.NewMemOperand(reg.SP, -16, true),
// 	// 		},
// 	// 		Comment: "save LR",
// 	// 	})
// 	// }

// 	// Setup frame pointer if needed
// 	if fm.currentFrame.HasFP {
// 		if len(instructions) == 0 {
// 			instructions = append(instructions, Instruction{
// 				Op:  op.MOV,
// 				Dst: reg.SP,
// 				Src: []reg.Operand{
// 					reg.NewRegOperand(reg.FP.String()),
// 					reg.NewRegOperand(reg.SP.String()),
// 				},
// 				Comment: "Load SP into FP",
// 			})
// 		}
// 	}

// 	// Allocate stack space
// 	if fm.currentFrame.Size > 0 {
// 		size := alloc.AlignSize(fm.currentFrame.Size, 16) // ARM64 requires 16-byte stack alignment
// 		instructions = append(instructions, Instruction{
// 			Op:  op.SUB,
// 			Dst: reg.SP,
// 			Src: []reg.Operand{
// 				reg.NewRegOperand(reg.SP.String()),
// 				reg.NewImmediateOperand(fmt.Sprintf("%d", size)),
// 			},
// 			Comment: "subtract offset from stack pointer",
// 		})
// 	}

// 	return instructions
// }

// // GenerateFrameTeardown generates the assembly instructions for frame teardown
// func (fm *FrameManager) GenerateFrameTeardown() (instructions []Instruction) {
// 	if fm.currentFrame == nil {
// 		return nil
// 	}

// 	// Restore stack pointer
// 	if fm.currentFrame.Size > 0 {
// 		size := alloc.AlignSize(fm.currentFrame.Size, 16)
// 		instructions = append(instructions, Instruction{
// 			Op:  op.ADD,
// 			Dst: reg.SP,
// 			Src: []reg.Operand{
// 				reg.NewRegOperand(reg.SP.String()),
// 				reg.NewImmediateOperand(fmt.Sprintf("%d", size)),
// 			},
// 			Comment: "add offset to stack pointer",
// 		})
// 	}

// 	// Restore FP/LR if they were saved
// 	if _, needsRestore := fm.currentFrame.SavedRegs[reg.LR]; needsRestore {
// 		instructions = append(instructions, Instruction{
// 			Op:  op.LDP,
// 			Dst: reg.FP,
// 			Src: []reg.Operand{
// 				reg.NewRegOperand(reg.LR.String()),
// 				reg.NewMemOperand(reg.SP, 16, false, true),
// 			},
// 			Comment: "restore LR",
// 		})
// 	}

// 	return instructions
// }

// // GetStackOffset returns the offset from SP for a given local variable
// func (fm *FrameManager) GetStackOffset(index int) int {
// 	if fm.currentFrame == nil {
// 		return 0
// 	}
// 	return fm.currentFrame.Offset + index
// }

// // GetFrameSize returns the total size of the current frame
// func (fm *FrameManager) GetFrameSize() int {
// 	if fm.currentFrame == nil {
// 		return 0
// 	}
// 	return alloc.AlignSize(fm.currentFrame.Size, 16)
// }
