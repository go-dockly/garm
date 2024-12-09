package compile

import (
	"fmt"
	"go/token"

	"github.com/algoboyz/garm/pkg/alloc"
	"github.com/algoboyz/garm/pkg/dbg"
	"github.com/algoboyz/garm/pkg/ir"
	"github.com/algoboyz/garm/pkg/mapper"
	"golang.org/x/tools/go/ssa"
)

// Compiler handles the conversion from Go AST to ARM64 assembly
type Compiler struct {
	prog   Program
	fset   *token.FileSet
	dbg    *dbg.Debugger
	mapper *mapper.SSAMapper
	gen    *Generator
}

// IRProgram represents the entire program
type Program struct {
	Functions []*ir.Function
	Globals   []alloc.ARM64Type
	Constants []alloc.ARM64Type // for constant handling
	Imports   []string          // to handle external dependencies
}

func New(debug *dbg.Debugger) *Compiler {
	compiler := &Compiler{
		prog: Program{
			Functions: make([]*ir.Function, 0),
			Globals:   make([]alloc.ARM64Type, 0),
			Constants: make([]alloc.ARM64Type, 0),
			Imports:   make([]string, 0),
		},
		mapper: mapper.NewSSAMapper(debug),
		fset:   token.NewFileSet(),
		gen:    NewCodeGenerator(debug),
		dbg:    debug,
	}
	return compiler
}

func (c *Compiler) Parse(target string, debug bool) (*ssa.Function, error) {
	if err := c.mapper.Load(target); err != nil {
		return nil, fmt.Errorf("loading package: %w", err)
	}
	fns, err := c.mapper.MapPackage()
	if err != nil {
		return nil, fmt.Errorf("loading package: %w", err)
	}

	c.prog.Functions = fns

	// f, err := parser.ParseFile(c.fset, target, nil, parser.ParseComments)
	// if err != nil {
	// 	return nil, fmt.Errorf("parsing AST: %w", err)
	// }
	// pkg, _, err := ssautil.BuildPackage(
	// 	&types.Config{Importer: importer.Default()},
	// 	c.fset,
	// 	types.NewPackage("main", ""),
	// 	[]*ast.File{f},
	// 	ssa.SanityCheckFunctions)
	// if err != nil {
	// 	return nil, fmt.Errorf("building SSA: %w", err)
	// }
	// fn := pkg.Func("main")
	// if fn == nil {
	// 	return nil, fmt.Errorf("main not found in %s", target)
	// }
	return nil, nil
}

func (c *Compiler) Map(fn *ssa.Function) error {
	fun, err := c.mapper.MapFunction(fn)
	if err != nil {
		return fmt.Errorf("mapping function %s: %w", fn.Name(), err)
	}
	c.prog.Functions = append(c.prog.Functions, fun)
	return nil
}

// Generate produces the final ARM64 assembly with optional optimization
func (c *Compiler) Generate() (string, error) {
	return c.gen.Glamour(c.prog)

	// if c.dbg.ModeDebug {
	// 	return c.gen.Glamour(c.prog)
	// }
	// return c.gen.Generate(c.prog), nil
}
