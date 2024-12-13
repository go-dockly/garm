package ir

import (
	"github.com/algoboyz/garm/pkg/op"
	"github.com/algoboyz/garm/pkg/reg"
)

// todo needs to regate registers for parameters
func FuncPrologue(label string) (instructions []Instruction) {
	return append(instructions, Instruction{
		Labels: []string{label},
	}, Instruction{
		Op: op.STP,
		Dst: &reg.Register{
			ID:    29,
			Class: reg.RegisterClassGPR,
		},
		Src: []reg.Operand{
			reg.NewRegOperand("X30"),
			reg.NewMemOperand(reg.SP, -16),
		},
		Comment: "Set up frame pointer",
	})
}

// todo needs to regate registers for return values
func FuncEpilogue() (instructions []Instruction) {
	return append(instructions, Instruction{
		Op: op.LDP,
		Dst: &reg.Register{
			ID:    29,
			Class: reg.RegisterClassGPR,
		},
		Src: []reg.Operand{
			reg.NewRegOperand("X30"),
			reg.NewMemOperand(reg.SP, 16),
		},
		Comment: "Restore frame pointer",
	}, Instruction{
		Op:      op.RET,
		Comment: "Function returns",
	})
}

func PrologueMain() (instructions []Instruction) {
	return append(instructions, Instruction{
		Labels: []string{"main"},
	}, Instruction{
		Op: op.STP,
		Dst: &reg.Register{
			ID:    29,
			Class: reg.RegisterClassGPR,
		},
		Src: []reg.Operand{
			reg.NewRegOperand("X30"),
			reg.NewMemOperand(reg.SP, -16),
		},
		Comment: "Set up frame pointer",
	}, Instruction{
		Op: op.MOV,
		Dst: &reg.Register{
			ID:    29,
			Class: reg.RegisterClassGPR,
		},
		Src: []reg.Operand{
			reg.NewRegOperand("SP"),
		},
		Comment: "Load SP into FP",
	})
}

func EpilogueMain() (instructions []Instruction) {
	return append(instructions, Instruction{
		Op: op.LDP,
		Dst: &reg.Register{
			ID:    29,
			Class: reg.RegisterClassGPR,
		},
		Src: []reg.Operand{
			reg.NewRegOperand("X30"),
			reg.NewMemOperand(reg.SP, 16, true),
		},
		Comment: "Restore frame pointer",
	}, Instruction{
		Op: op.MOV,
		Dst: &reg.Register{
			ID:    0,
			Class: reg.RegisterClassGPR,
		},
		Src: []reg.Operand{
			reg.NewImmediateOperand("0"),
		},
		Comment: "Cleanup and exit",
	}, Instruction{
		Op: op.MOV,
		Dst: &reg.Register{
			ID:    8,
			Class: reg.RegisterClassGPR,
		},
		Src: []reg.Operand{
			reg.NewImmediateOperand("93"),
		},
		Comment: "Prepare for syscall",
	}, Instruction{
		Op: op.SVC,
		Src: []reg.Operand{
			reg.NewImmediateOperand("0"),
		},
		Comment: "Call supervisor",
	})
}
