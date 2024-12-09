package mapper

import (
	"fmt"
	"go/types"

	"github.com/algoboyz/garm/pkg/alloc"
)

func (m *SSAMapper) MapBasicType(name string, typ types.Type) (p alloc.Primitive, err error) {
	switch t := typ.(type) {
	case *types.Basic:
		return m.mapBasicType(name, t)
	case *types.Pointer:
		return alloc.Pointer, nil
	case *types.Named:
		return m.MapBasicType(name, t.Underlying())
	default:
		return p, fmt.Errorf("unsupported type: %T", typ)
	}
}

func (m *SSAMapper) mapBasicType(name string, typ *types.Basic) (p alloc.Primitive, err error) {
	switch typ.Kind() {
	case types.Int8:
		return alloc.Int8, nil
	case types.Int16:
		return alloc.Int16, nil
	case types.Int32:
		return alloc.Int32, nil
	case types.Int64, types.Int:
		return alloc.Int64, nil
	case types.Float32:
		return alloc.Float32, nil
	case types.Float64:
		return alloc.Float64, nil
	case types.Bool:
		return alloc.Bool, nil
	case types.String:
		return alloc.String, nil
	default:
		return p, fmt.Errorf("unsupported basic type: %v", typ)
	}
}
