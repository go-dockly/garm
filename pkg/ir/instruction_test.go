package ir

import (
	"testing"

	"github.com/algoboyz/garm/pkg/op"
	"github.com/algoboyz/garm/pkg/reg"
	"github.com/stretchr/testify/assert"
)

// todo missing tests

func TestInstructionString(t *testing.T) {
	tests := []struct {
		name     string
		instr    Instruction
		expected string
	}{
		// {
		// 	name: "Defer cleanup function with argument 42",
		// 	instr: Instruction{
		// 		Macro: &asm.DEFER,
		// 		Src: []reg.Operand{
		// 			reg.NewLabelOperand("cleanup"),
		// 			reg.NewRegOperand("X0"),
		// 		},
		// 		Comment: "Call defer for cleanup(42)",
		// 	},
		// 	expected: "DEFER X0, X1 // Move X1 to X0",
		// },
		{
			name: "Set up frame pointer",
			instr: Instruction{
				Op: op.STP,
				Dst: &reg.Register{
					ID:    29,
					Class: reg.RegisterClassGPR,
				},
				Src: []reg.Operand{
					reg.NewRegOperand("X30"),
					reg.NewMemOperand(reg.SP, -16, true),
				},
				Comment: "Set up frame pointer",
			},
			expected: "\tSTP X29, X30, [SP], #-16!\t// Set up frame pointer\n",
		},
		{
			name: "Instruction with flags",
			instr: Instruction{
				Op: op.ADD,
				Dst: &reg.Register{
					ID:    0,
					Class: reg.RegisterClassGPR,
				},
				Src: []reg.Operand{
					reg.NewRegOperand("X1"),
					reg.NewRegOperand("X2"),
				},
				Flags:   Wide,
				Comment: "Add X1 and X2, store in X0",
			},
			expected: "\tADD X0, X1, X2.W\t// Add X1 and X2, store in X0\n",
		},
		{
			name: "Instruction with predicates",
			instr: Instruction{
				Op: op.SUB,
				Dst: &reg.Register{
					ID:    0,
					Class: reg.RegisterClassGPR,
				},
				Src: []reg.Operand{
					reg.NewRegOperand("X1"),
					reg.NewRegOperand("X2"),
				},
				Pred: []op.Predicate{
					{Condition: op.Equal},
				},
				Comment: "Subtract X2 from X1 if equal",
			},
			expected: "\tSUB.EQ X0, X1, X2\t// Subtract X2 from X1 if equal\n",
		},
		{
			name: "Branch to Label Instruction",
			instr: Instruction{
				Op: op.B,
				Labels: []string{
					"label1",
				},
			},
			expected: "\tB label1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.instr.String(true)
			assert.Equal(t, tt.expected, result, "expected %q, got %q", tt.expected, result)
		})
	}
}
