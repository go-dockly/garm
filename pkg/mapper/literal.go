package mapper

import (
	"fmt"
	"go/types"

	"github.com/algoboyz/garm/pkg/alloc"
)

func (m *SSAMapper) MapLiteral(name string, lit types.Type) (alloc.ARM64Type, error) {
	var typ alloc.Primitive
	size := 64

	switch lit.String() {
	case "int":
		typ = alloc.Int64
	case "float":
		typ = alloc.Float64
	case "string":
		typ = alloc.String
		size = alloc.AlignSize(len(lit.String())+1, alloc.WordSize) // +1 for null terminator
	default:
		return nil, fmt.Errorf("unsupported literal type: %s", lit)
	}

	return alloc.NewType(name, lit.String(), typ, size), nil
}
