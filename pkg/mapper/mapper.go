package mapper

import (
	"fmt"

	"github.com/algoboyz/garm/pkg/alloc"
	"github.com/algoboyz/garm/pkg/dbg"
	"github.com/algoboyz/garm/pkg/ir"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

type SSAMapper struct {
	// SSA-specific context tracking
	prog         *ssa.Program
	pkgs         []*ssa.Package
	currentFunc  *ssa.Function
	currentBlock *ssa.BasicBlock
	currentIR    *ir.Function
	labelMap     map[*ssa.BasicBlock]string
	alloc        alloc.Allocator
	debug        *dbg.Debugger
}

func NewSSAMapper(debug *dbg.Debugger) *SSAMapper {
	return &SSAMapper{alloc: alloc.NewAllocator(), debug: debug}
}

// LoadPackage loads and builds SSA for a Go package
func (m *SSAMapper) Load(path string) error {
	cfg := &packages.Config{
		Mode: packages.NeedFiles |
			packages.NeedSyntax |
			packages.NeedTypesInfo |
			packages.NeedTypes |
			packages.NeedTypesSizes |
			packages.NeedDeps,
	}

	pkgs, err := packages.Load(cfg, path)
	if err != nil {
		return fmt.Errorf("loading package: %w", err)
	}

	// Create SSA program
	m.prog, m.pkgs = ssautil.Packages(pkgs, ssa.BuilderMode(ssa.SanityCheckFunctions))
	m.prog.Build()

	return nil
}

// // MapPackage processes an entire SSA package
func (m *SSAMapper) MapPackage() (fns []*ir.Function, err error) {
	// Process all functions in the package
	for _, pkg := range m.pkgs {
		for _, member := range pkg.Members {
			if fn, ok := member.(*ssa.Function); ok {
				fun, err := m.MapFunction(fn)
				if err != nil {
					return nil, fmt.Errorf("mapping function %s: %w", fn.Name(), err)
				}
				fns = append(fns, fun)
			}
		}
	}
	return fns, nil
}
