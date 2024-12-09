package reg

import (
	"fmt"
	"strings"

	"github.com/algoboyz/garm/pkg/op"
)

// Operand represents an operand in ARM assembly
type Operand struct {
	Type       OperandType
	Var        string
	Shift      *Shift         // For shifted register operations
	Memory     *MemoryOperand // For memory operands
	Register   *RegisterClass // For detailed register info
	Name       string
	IsArray    bool
	IsParam    bool
	Size       int
	Alignment  int  // for proper memory alignment
	IsVolatile bool // for volatile variables
	// Scope      ScopeType // track variable scope todo
}

type Shift struct {
	Type  op.Op  // LSL, LSR, ASR, ROR
	Value string // Can be immediate or register
}

func (o *Operand) String() string {
	switch o.Type {
	case OperandImmediate:
		return "#" + o.Var
	case OperandShift:
		if o.Shift != nil {
			return fmt.Sprintf("%s, %s %s", o.Var, o.Shift.Type, o.Shift.Value)
		}
	case OperandPostIndex:
		return fmt.Sprintf("%s!", o.Var)
	case OperandMemory:
		return o.Memory.String()
	case OperandRegister:
	case OperandLabel:
	case ShiftedRegister:
	case RegisterList:
	case RegisterRange:
	case RegType:
	case ImmType:
	case ShiftType:
	case MemType:
	default:
	}
	return o.Var
}

// OperandType represents different kinds of operands in ARM assembly
// Enhanced operand types for ARM-specific patterns
type OperandType int

const (
	OperandRegister  OperandType = iota
	OperandImmediate             // immediate value
	OperandMemory
	OperandShift
	OperandPostIndex // post-indexed addressing
	OperandLabel
	ShiftedRegister
	RegisterList  // for push/pop operations
	RegisterRange // for register ranges
	RegType
	ImmType
	ShiftType
	MemType // todo deduplicate
)

func NewLabelOperand(label string) Operand {
	return Operand{
		Type: OperandLabel,
		Var:  label,
	}
}

// Helper functions for instruction generation
func NewRegOperand(register string) Operand {
	return Operand{
		Type: OperandRegister,
		Var:  register,
	}
}

func NewImmediateOperand(val string) Operand {
	return Operand{
		Type: OperandImmediate,
		Var:  strings.TrimPrefix(val, "#"),
	}
}

func NewOperandWithUpdate(value string) Operand {
	return Operand{
		Type: OperandPostIndex,
		Var:  value,
	}
}

func NewMemOperand(reg *Register, offset int, postIdx ...bool) Operand {
	return Operand{
		Type: OperandMemory,
		Memory: &MemoryOperand{
			BaseRegister: reg,
			Offset:       fmt.Sprint(offset),
			WriteBack:    true,
			Indirect:     reg == nil,
			Post:         len(postIdx) == 1 && postIdx[0],
		},
	}
}

type MemoryOperand struct {
	BaseRegister *Register
	Offset       string
	Index        string
	WriteBack    bool // Write-back for auto-increment
	Pre          bool // Pre-indexed addressing
	Post         bool // Post-indexed addressing
	Indirect     bool // For indirect memory access means no base register is used
}

// Helper function to format memory operands
func (m *MemoryOperand) String() string {
	var result strings.Builder
	result.WriteString("[")
	result.WriteString(m.BaseRegister.String())
	if m.Post {
		result.WriteString("]")
	}
	if m.Offset != "" {
		if strings.HasPrefix(m.Offset, "#") {
			result.WriteString(", " + m.Offset)
		} else {
			result.WriteString(", #" + m.Offset)
		}
	}
	if m.Index != "" {
		result.WriteString(", " + m.Index)
	}
	if !m.Post {
		result.WriteString("]")
	}
	if m.WriteBack {
		result.WriteString("!")
	}
	return result.String()
}
