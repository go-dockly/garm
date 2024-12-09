package ir

import (
	"fmt"

	"github.com/algoboyz/garm/pkg/alloc"
	"github.com/algoboyz/garm/pkg/dbg"
)

// Function represents a function in our IR
type Function struct {
	Label            string
	Public           bool
	StackSize        int
	nextGlobalOffset int
	Globals          map[string]alloc.Location
	Params           map[string]alloc.Location
	Locals           map[string]alloc.Location
	Blocks           []Instruction
	Returns          map[string]alloc.Location
	dbg              *dbg.Debugger
	Frames           *FrameManager
}

func NewFunction(label string, debug *dbg.Debugger) *Function {
	return &Function{
		Label:   label,
		Public:  false,
		Params:  make(map[string]alloc.Location),
		Locals:  make(map[string]alloc.Location),
		Globals: make(map[string]alloc.Location),
		Blocks:  make([]Instruction, 0),
		dbg:     debug,
		Frames:  NewFrameManager(),
	}
}

// AllocateGlobalOffset allocates space in the global section and returns the offset
func (f *Function) AllocateGlobalOffset(size int) int {
	// Align to 8-byte boundary for AArch64
	alignment := 8
	f.nextGlobalOffset = (f.nextGlobalOffset + alignment - 1) & ^(alignment - 1)

	offset := f.nextGlobalOffset
	f.nextGlobalOffset += size

	return offset
}

func (f *Function) StackFrame() int {
	// Calculate space needed for local vars
	var sum int
	for _, local := range f.Locals {
		size := local.GetMemory().Size
		if size%8 != 0 {
			size = (size + 7) & ^7
		}
		sum += size
	}
	// TODO Add space for spilled registers
	// spillSize := spillcount * 8
	spillSize := 8
	// Calculate total frame size (align to 16 bytes for ARM64)
	total := sum + spillSize
	if total%16 != 0 {
		total = (total + 15) & ^15
	}
	return total
}

func (f *Function) Has(variable string) (alloc.Location, error) {
	x, ok := f.Params[variable]
	if !ok {
		x, ok = f.Locals[variable]
		if !ok {
			return nil, fmt.Errorf("%s not found in scope", variable)
		}
	}
	return x, nil
}

func (f *Function) Debug() (md string) {
	md += "```yaml\n"
	md += "Function: " + f.Label + "\n"
	md += "Public: " + fmt.Sprintf("%t", f.Public) + "\n"
	md += "StackSize: " + fmt.Sprintf("%d", f.StackSize) + "\n"
	md += "Params:\n"
	// if len(f.Params) > 0 {
	// 	md += "```\n"
	// 	md += "```go\n"
	// 	for _, param := range f.Params {
	// 		md += "\t" + param.Value.Debug() + "\n"
	// 	}
	// 	md += "```\n"
	// 	md += "```yaml\n"
	// }
	// md += "Variables:\n"
	// if len(f.Locals) > 0 {
	// 	md += "```\n"
	// 	md += "```go\n"
	// 	for _, local := range f.Locals {
	// 		md += "\t" + local.Value.Debug() + "\n"
	// 	}
	// 	md += "```\n"
	// 	md += "```yaml\n"
	// }
	// md += "Returns:\n"
	// if len(f.Returns) > 0 {
	// 	md += "```\n"
	// 	md += "```go\n"
	// 	for _, ret := range f.Returns {
	// 		md += "\t" + ret.Value.Debug() + "\n"
	// 	}
	// 	md += "```\n"
	// 	md += "```yaml\n"
	// }
	md += "```\n"
	f.dbg.Render(md)
	return md
}
