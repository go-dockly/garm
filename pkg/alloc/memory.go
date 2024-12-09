package alloc

import (
	"fmt"

	"github.com/algoboyz/garm/pkg/reg"
)

// MemoryLocation represents a stack location
type MemoryLocation struct {
	Name      string
	Offset    int
	Size      int
	Alignment int
}

// locationImpl implements the Location interface
type locationImpl struct {
	register *reg.Register
	memory   *MemoryLocation
}

func (l *locationImpl) GetRegister() *reg.Register {
	return l.register
}

func (l *locationImpl) IsRegister() bool {
	return l.register != nil
}

func (l *locationImpl) GetMemory() *MemoryLocation {
	return l.memory
}

func (l *locationImpl) IsMemory() bool {
	return l.memory != nil
}

func (l *locationImpl) String() string {
	if l.IsRegister() {
		switch l.register.Class {
		case reg.RegisterClassGPR:
			return fmt.Sprintf("x%d", l.register.ID)
		case reg.RegisterClassFPR:
			return fmt.Sprintf("d%d", l.register.ID)
		case reg.RegisterClassVec:
			return fmt.Sprintf("v%d", l.register.ID)
		}
	}
	return fmt.Sprintf("[sp, #%d]", l.memory.Offset)
}
