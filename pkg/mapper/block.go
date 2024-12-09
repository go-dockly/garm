package mapper

import (
	"fmt"

	"golang.org/x/tools/go/ssa"
)

// processBlock converts an SSA basic block to ARM64 IR
func (m *SSAMapper) MapBlock(block *ssa.BasicBlock) (err error) {
	// Add block label
	// 	m.addBlockLabel(block)
	for _, instr := range block.Instrs {
		if err = m.MapInstruction(instr); err != nil {
			return fmt.Errorf("processing instruction %v: %w", instr, err)
		}
	}
	// Handle block terminator
	// if err = m.processTerminator(block.Instrs[len(block.Instrs)-1], irFunc); err != nil {
	// 	return nil, fmt.Errorf("processing terminator: %w", err)
	// }
	return nil
}
