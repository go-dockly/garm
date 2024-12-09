package reg

import (
	"testing"
)

func TestMemoryOperandString(t *testing.T) {
	tests := []struct {
		name     string
		operand  MemoryOperand
		expected string
	}{
		{
			name: "Base register only",
			operand: MemoryOperand{
				BaseRegister: &Register{ID: 0, Class: RegisterClassGPR},
			},
			expected: "[X0]",
		},
		{
			name: "Base register with offset",
			operand: MemoryOperand{
				BaseRegister: &Register{ID: 1, Class: RegisterClassGPR},
				Offset:       "4",
			},
			expected: "[X1, #4]",
		},
		{

			name: "Base register with offset and index",
			operand: MemoryOperand{
				BaseRegister: &Register{ID: 2, Class: RegisterClassGPR},
				Offset:       "8",
				Index:        "X3",
			},
			expected: "[X2, #8, X3]",
		},
		{
			name: "Base register with write-back",
			operand: MemoryOperand{
				BaseRegister: &Register{ID: 4, Class: RegisterClassGPR},
				WriteBack:    true,
			},
			expected: "[X4]!",
		},
		{
			name: "Base register with offset and write-back",
			operand: MemoryOperand{
				BaseRegister: &Register{ID: 5, Class: RegisterClassGPR},
				Offset:       "#12",
				WriteBack:    true,
			},
			expected: "[X5, #12]!",
		},
		{
			name: "Load value of SP then advance SP by 16",
			operand: MemoryOperand{
				BaseRegister: SP,
				Offset:       "#16",
				Post:         true,
			},
			expected: "[SP], #16",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.operand.String()
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}
func TestOperandString(t *testing.T) {
	tests := []struct {
		name     string
		operand  Operand
		expected string
	}{
		{
			name: "Immediate operand",
			operand: Operand{
				Type: OperandImmediate,
				Var:  "10",
			},
			expected: "#10",
		},
		{
			name: "Shift operand",
			operand: Operand{
				Type: OperandShift,
				Var:  "X1",
				Shift: &Shift{
					Type:  "LSL",
					Value: "2",
				},
			},
			expected: "X1, LSL 2",
		},
		{
			name: "Post-index operand",
			operand: Operand{
				Type: OperandPostIndex,
				Var:  "X2",
			},
			expected: "X2!",
		},
		{
			name: "Memory operand",
			operand: Operand{
				Type: OperandMemory,
				Memory: &MemoryOperand{
					BaseRegister: &Register{ID: 3, Class: RegisterClassGPR},
					Offset:       "4",
				},
			},
			expected: "[X3, #4]",
		},
		{
			name: "Register operand",
			operand: Operand{
				Type: OperandRegister,
				Var:  "X4",
			},
			expected: "X4",
		},
		{
			name: "Label operand",
			operand: Operand{
				Type: OperandLabel,
				Var:  ".L1",
			},
			expected: ".L1",
		},
		{
			name: "Shifted register operand",
			operand: Operand{
				Type: ShiftedRegister,
				Var:  "X5",
			},
			expected: "X5",
		},
		{
			name: "Register list operand",
			operand: Operand{
				Type: RegisterList,
				Var:  "{X0, X1, X2}",
			},
			expected: "{X0, X1, X2}",
		},
		{
			name: "Register range operand",
			operand: Operand{
				Type: RegisterRange,
				Var:  "X0-X7",
			},
			expected: "X0-X7",
		},
		{
			name: "RegType operand",
			operand: Operand{
				Type: RegType,
				Var:  "X6",
			},
			expected: "X6",
		},
		{
			name: "ImmType operand",
			operand: Operand{
				Type: ImmType,
				Var:  "#15",
			},
			expected: "#15",
		},
		{
			name: "ShiftType operand",
			operand: Operand{
				Type: ShiftType,
				Var:  "LSL",
			},
			expected: "LSL",
		},
		{
			name: "MemType operand",
			operand: Operand{
				Type: MemType,
				Var:  "[X7]",
			},
			expected: "[X7]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.operand.String()
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}
