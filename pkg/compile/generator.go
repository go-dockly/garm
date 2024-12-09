package compile

import (
	"fmt"
	"strings"

	"github.com/algoboyz/garm/pkg/dbg"
	"github.com/algoboyz/garm/pkg/ir"
	"github.com/charmbracelet/glamour"
)

// CodeGenerator handles the final assembly generation with peephole optimization
type Generator struct {
	instructions []ir.Instruction
	dbg          *dbg.Debugger
}

func NewCodeGenerator(debug *dbg.Debugger) *Generator {
	return &Generator{
		instructions: make([]ir.Instruction, 0),
		dbg:          debug,
	}
}

// Generate produces the final ARM64 assembly
func (g *Generator) Generate(program Program) string {
	var sb strings.Builder
	sb.WriteString("\t.global main\n")
	sb.WriteString("\t.text\n")
	for _, f := range program.Functions {
		if f.Public {
			sb.WriteString(fmt.Sprintf(".global %s:\n", f.Label))
		}
		for _, inst := range f.Blocks {
			sb.WriteString(inst.String(g.dbg.ModeDebug))
		}
	}
	return sb.String()
}

// Render assembly outout as markdown to ANSI-styled terminal output
func (g *Generator) Glamour(program Program) (string, error) {
	content := fmt.Sprintf("```asm\n%s```\n", g.Generate(program))
	out, err := glamour.Render(content, "dark")
	if err != nil {
		return "", fmt.Errorf("Glamour: %w", err)
	}
	// todo use https://github.com/charmbracelet/glow to render output file
	fmt.Print(out)
	return content, nil
}
