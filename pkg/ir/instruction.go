package ir

import (
	"fmt"
	"strings"

	"github.com/algoboyz/garm/pkg/asm"
	"github.com/algoboyz/garm/pkg/op"
	"github.com/algoboyz/garm/pkg/reg"
)

// CompoundOp represents a compound operation that needs multiple instructions eg +=, -=, *=
type CompoundInstruction struct {
	LoadOp   []Instruction // Instructions to load value
	ModifyOp []Instruction // Instructions to modify value
	StoreOp  []Instruction // Instructions to store value
	Atomic   bool          // Whether operation needs to be atomic
}

// Instruction represents an ARM64 assembly instruction
type Instruction struct {
	Op      op.Op
	Pred    []op.Predicate // For conditional execution
	Dst     *reg.Register
	Src     []reg.Operand
	Labels  []string   // for better control flow
	Flags   InstrFlags // Instruction flags (e.g., S suffix)
	Macro   *asm.Macro // Macro instruction
	Comment string     // Optional comment
}

type InstrFlags uint8

const (
	UpdateFlags InstrFlags = 1 << iota // Sets condition flags
	Wide                               // Use wide (64-bit) variant
	Saturating                         // Use saturating arithmetic
)

func (i *Instruction) String(debug bool) string {
	var sb strings.Builder
	if !i.Op.IsBranch() {
		for _, label := range i.Labels {
			sb.WriteString(fmt.Sprintf("\n%s:\n", label))
		}
	}
	// Start building the instruction
	sb.WriteString("\t") // Indent for assembly format

	// Write operation with predicates
	if len(i.Pred) > 0 {
		sb.WriteString(i.Op.String())
		for _, p := range i.Pred {
			sb.WriteString("." + p.String())
		}
	} else {
		sb.WriteString(i.Op.String())
	}

	// Handle branch instructions specially
	if i.Op.IsBranch() && len(i.Labels) > 0 {
		sb.WriteString(" " + i.Labels[0])
		// Add comment if in debug mode
		if debug && i.Comment != "" {
			sb.WriteString("\t\t// " + i.Comment)
		}
		sb.WriteString("\n")
		return sb.String()
	}

	switch {
	case i.Macro != nil:
		sb.WriteString(fmt.Sprintf(" %s", *i.Macro))
	case i.Dst != nil:
		sb.WriteString(fmt.Sprintf(" %s", i.Dst.String()))
	}

	// Format operands
	var formattedOps []string
	for _, op := range i.Src {
		formattedOps = append(formattedOps, op.String())
	}

	if len(formattedOps) > 0 {
		sb.WriteString(", " + strings.Join(formattedOps, ", "))
	}

	// Add instruction flags
	if i.Flags&Saturating != 0 {
		sb.WriteString("S")
	}
	if i.Flags&Wide != 0 {
		sb.WriteString(".W")
	}

	// Add comment if in debug mode
	if debug && i.Comment != "" {
		sb.WriteString("\t// " + i.Comment)
	}

	sb.WriteString("\n")
	return sb.String()
}
