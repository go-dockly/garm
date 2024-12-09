package alloc

import (
	"fmt"

	"github.com/algoboyz/garm/pkg/reg"
)

// ARM64Value represents a value that can be allocated to registers or memory
type ARM64Value struct {
	Type     ARM64Type
	Layout   MemoryLayout
	Info     TypeInfo
	location Location // Current location (register or memory)
}

// TypeInfo represents high-level type information from Go
type TypeInfo struct {
	Name     string
	Kind     TypeKind
	IsParam  bool
	Elements []TypeInfo // For composite types
}

// TypeKind represents the fundamental kinds of types we support
type TypeKind uint8

const (
	KindScalar TypeKind = iota
	KindVector
	KindComposite
	KindArray
)

// Primitive represents basic ARM64 types
type Primitive uint8

const (
	Invalid Primitive = iota
	Int8
	Int16
	Int32
	Int64
	Float32
	Float64
	String
	Bool
	Pointer
)

func NewPrimitive(v string) Primitive {
	switch v {
	case "int8":
		return Int8
	case "int16":
		return Int16
	case "int32":
		return Int32
	case "int64":
		return Int64
	case "float32":
		return Float32
	case "float64":
		return Float64
	case "string":
		return String
	case "bool":
		return Bool
	case "unsafe.Pointer":
		return Pointer
	default:
		return Invalid
	}
}

// baseARM64Type implements the ARM64Type interface
type baseARM64Type struct {
	name      string
	goType    string
	primitive Primitive
	size      int
	align     int
	reg       reg.RegisterClass
}

func (t *baseARM64Type) Size() int {
	return t.size
}

func (t *baseARM64Type) Alignment() int {
	return t.align
}

func (t *baseARM64Type) Register() reg.RegisterClass {
	return t.reg
}

func (t *baseARM64Type) String() string {
	return fmt.Sprintf("%s(%s)", t.name, t.goType)
}

// NewType creates a new ARM64Type with the given properties
func NewType(name, goType string, prim Primitive, size int) ARM64Type {
	t := &baseARM64Type{
		name:      name,
		goType:    goType,
		primitive: prim,
		size:      size,
	}

	// Set alignment and register class based on primitive type
	switch prim {
	case Int8, Int16, Int32, Int64, Bool, Pointer:
		t.reg = reg.RegisterClassGPR
		t.align = (size + 7) / 8 // Round up to nearest byte
	case Float32, Float64:
		t.reg = reg.RegisterClassFPR
		t.align = (size + 7) / 8
	case String:
		t.reg = reg.RegisterClassGPR
		t.align = 8 // Strings are 8-byte aligned in Go
	default:
		t.reg = reg.RegisterClassGPR
		t.align = 8
	}

	return t
}

// TypeMapper handles mapping between Go types and ARM64 types
type TypeMapper struct {
	types map[string]ARM64Type
}

func NewTypeMapper() *TypeMapper {
	return &TypeMapper{
		types: make(map[string]ARM64Type),
	}
}

// RegisterType adds a new type mapping
func (m *TypeMapper) RegisterType(goType string, armType ARM64Type) {
	m.types[goType] = armType
}

// GetType retrieves a mapped type
func (m *TypeMapper) GetType(goType string) (ARM64Type, bool) {
	t, ok := m.types[goType]
	return t, ok
}

// Standard type sizes for ARM64
const (
	PtrSize   = 8
	WordSize  = 8
	Int64Size = 8
	Int32Size = 4
	Int16Size = 2
	Int8Size  = 1
)

// TypeSet provides predefined ARM64 types
var TypeSet = struct {
	Int8    ARM64Type
	Int16   ARM64Type
	Int32   ARM64Type
	Int64   ARM64Type
	Float32 ARM64Type
	Float64 ARM64Type
	Bool    ARM64Type
	String  ARM64Type
	Pointer ARM64Type
}{
	Int8:    NewType("int8", "int8", Int8, Int8Size),
	Int16:   NewType("int16", "int16", Int16, Int16Size),
	Int32:   NewType("int32", "int32", Int32, Int32Size),
	Int64:   NewType("int64", "int64", Int64, Int64Size),
	Float32: NewType("float32", "float32", Float32, 4),
	Float64: NewType("float64", "float64", Float64, 8),
	Bool:    NewType("bool", "bool", Bool, 1),
	String:  NewType("string", "string", String, PtrSize*2), // ptr + len
	Pointer: NewType("ptr", "unsafe.Pointer", Pointer, PtrSize),
}

// AlignSize aligns the given size to the specified alignment
func AlignSize(size, align int) int {
	return (size + align - 1) & ^(align - 1)
}
