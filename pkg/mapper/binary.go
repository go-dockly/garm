package mapper

import (
	"fmt"
	"strings"

	"github.com/algoboyz/garm/pkg/ir"
	"github.com/algoboyz/garm/pkg/op"
	"github.com/algoboyz/garm/pkg/reg"
	"golang.org/x/tools/go/ssa"
)

func (m *SSAMapper) MapBinaryOperation(expr *ssa.BinOp) error {
	var err error
	token, err := m.MapToken(expr.Op)
	if err != nil {
		return fmt.Errorf("mapping operator: %w", err)
	}
	x, err := m.MapLiteral(expr.X.Name(), expr.X.Type())
	if err != nil {
		return err
	}
	y, err := m.MapLiteral(expr.Y.Name(), expr.Y.Type())
	if err != nil {
		return err
	}
	result, err := m.MapLiteral(expr.Name(), expr.Type())
	if err != nil {
		return err
	}
	// Allocate registers
	dst, err := m.alloc.AllocateRegister(result)
	if err != nil {
		return err
	}
	m.currentIR.Locals[result.String()] = dst
	lhs, err := m.alloc.AllocateRegister(x)
	if err != nil {
		return err
	}
	rhs, err := m.alloc.AllocateRegister(y)
	if err != nil {
		return err
	}
	// Generate ARM64 instruction
	parts := strings.Split(expr.X.String(), ":")
	m.currentIR.Blocks = append(m.currentIR.Blocks, ir.Instruction{
		Op:  op.MOV,
		Dst: lhs.GetRegister(),
		Src: []reg.Operand{
			reg.NewImmediateOperand(parts[0]),
		},
		Comment: fmt.Sprintf("load: %s = %s", parts[1], parts[0]),
	})
	parts = strings.Split(expr.Y.String(), ":")
	m.currentIR.Blocks = append(m.currentIR.Blocks, ir.Instruction{
		Op:  op.MOV,
		Dst: rhs.GetRegister(),
		Src: []reg.Operand{
			reg.NewImmediateOperand(parts[0]),
		},
		Comment: fmt.Sprintf("load: %s = %s", parts[1], parts[0]),
	})
	m.currentIR.Blocks = append(m.currentIR.Blocks, ir.Instruction{
		Op:  token,
		Dst: dst.GetRegister(),
		Src: []reg.Operand{
			reg.NewRegOperand(lhs.GetRegister().String()),
			reg.NewRegOperand(rhs.GetRegister().String()),
		},
		Comment: fmt.Sprintf("%s = %s %s %s", expr.Name(), expr.X.Name(), expr.Op.String(), expr.Y.Name()),
	})
	return nil
}
