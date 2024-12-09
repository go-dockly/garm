package mapper

import (
	"fmt"

	"github.com/algoboyz/garm/pkg/ir"
	"github.com/algoboyz/garm/pkg/op"
	"golang.org/x/tools/go/ssa"
)

func (m *SSAMapper) MapCall(expr *ssa.Call) error {
	// m.debug.Println(fun)
	params := []string{}
	for _, arg := range expr.Call.Args {
		params = append(params, arg.Name())
	}
	block := ir.Instruction{
		Op:      op.BL,
		Labels:  []string{expr.Name()},
		Comment: fmt.Sprintf("Call %s with %s", expr.Name(), params),
	}
	m.currentIR.Blocks = append(m.currentIR.Blocks, block)
	return nil
}
