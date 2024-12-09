package reg

import "fmt"

// RegisterClass represents different types of ARM64 registers
type RegisterClass uint8

const (
	RegisterClassGPR RegisterClass = iota // General Purpose Register
	RegisterClassFPR                      // Floating Point Register
	RegisterClassVec                      // Vector Register
	FramePointer
	LinkRegister
	StackPointer
)

var (
	FP = &Register{ID: 29, Class: FramePointer}
	LR = &Register{ID: 30, Class: LinkRegister}
	SP = &Register{Name: "sp", Class: StackPointer}
)

// Register represents an actual ARM64 register
type Register struct {
	ID    uint8
	Name  string
	Class RegisterClass
}

func (r *Register) String() string {
	if r.Name != "" {
		return r.Name
	}
	switch r.Class {
	case RegisterClassGPR:
		return fmt.Sprintf("x%d", r.ID)
	case RegisterClassFPR:
		return fmt.Sprintf("d%d", r.ID)
	case RegisterClassVec:
		return fmt.Sprintf("v%d", r.ID)
	case FramePointer:
		return "fp"
	case LinkRegister:
		return "lr"
	case StackPointer:
		return "sp"
	default:
		return fmt.Sprintf("unknown%d", r.ID)
	}
}
