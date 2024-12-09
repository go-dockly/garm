package asm

import (
	"embed"
	"fmt"
)

// Macro instructions
type Macro string

//go:embed *.asm
var Macros embed.FS

func LoadMacro(name string) (string, error) {
	path := fmt.Sprintf("macros/%s.asm", name)
	data, err := Macros.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
