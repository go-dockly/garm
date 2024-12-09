package mapper

import (
	"github.com/davecgh/go-spew/spew"
	"golang.org/x/tools/go/ssa"
)

func (m *SSAMapper) MapInstruction(instr ssa.Instruction) error {
	switch v := instr.(type) {
	case *ssa.Alloc:
		return m.MapAlloc(v)
	case *ssa.Store:
		return m.MapStore(v)
	case *ssa.BinOp:
		return m.MapBinaryOperation(v)
	case *ssa.Call:
		return m.MapCall(v)
	case *ssa.Convert:
		// m.MapTypeConversion(v)
	case *ssa.Jump:
		// m.MapJump(v)
	case *ssa.If:
		// m.MapConditional(v)
	case *ssa.Phi:
		// return m.MapPhi(v)
	case *ssa.Return:
		// return m.MapReturn(v)
	default:
		spew.Dump(instr)
		// return fmt.Errorf("unsupported instruction type: %T", instr)
	}
	return nil
}
