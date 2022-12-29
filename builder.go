package structure

import (
	"reflect"
)

type StructureBuilder interface {
	// Add field to structure.
	AddField(name string, typ interface{}, tag string)
	// Build structure
	Build() (Structure, error)
}

type builder struct {
	fields []reflect.StructField
}

// NewBuilder return structure builder.
func NewBuilder() *builder {
	return &builder{}
}

func (b *builder) AddField(name string, typ interface{}, tag string) {
	b.fields = append(b.fields, reflect.StructField{
		Name: name,
		Type: reflect.TypeOf(typ),
		Tag:  reflect.StructTag(tag),
	})
}

func (b *builder) Build() Structure {
	s, _ := New(reflect.New(reflect.StructOf(b.fields)).Interface())
	return s
}
