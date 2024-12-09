package mapper

import (
	"fmt"

	"golang.org/x/tools/go/ssa"
)

func (m *SSAMapper) MapAlloc(v *ssa.Alloc) error {
	typ, err := m.MapLiteral(v.Name(), v.Type())
	if err != nil {
		return fmt.Errorf("mapping alloc type: %w", err)
	}

	allocation, err := m.alloc.AllocateRegister(typ)
	if err != nil {
		return fmt.Errorf("allocating %s: %w", v.Name(), err)
	}

	m.currentIR.Locals[v.Name()] = allocation
	return nil
}
