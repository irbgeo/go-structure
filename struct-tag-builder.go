package structtagbuilder

import (
	"fmt"
	"reflect"
)

type StructWithTags interface {
	Writable() interface{}
	SaveInto(i interface{}) error
	AssignFrom(i interface{}) error
}

type StructWithTagsBuilder interface {
	Build(getTag func(fieldName string) string) StructWithTags
}

type rowBuilderImpl struct {
	i any
}

func NewTagBuilder(i interface{}) (StructWithTagsBuilder, error) {
	ret := &rowBuilderImpl{
		i: i,
	}
	if reflect.TypeOf(i).Kind() != reflect.Ptr {
		return nil, fmt.Errorf("needs '*struct' type for sample")
	}
	return ret, nil
}

type rowImpl struct {
	row          interface{}
	originalType reflect.Type
}

func (rb rowBuilderImpl) Build(getTag func(fieldName string) string) StructWithTags {
	var fields []reflect.StructField
	values := reflect.ValueOf(rb.i).Elem()
	for i := 1; i < values.NumField(); i++ {
		f := values.Type().Field(i)
		f.Tag = reflect.StructTag(getTag(f.Name))
		fields = append(fields, f)
	}

	t := reflect.StructOf(fields)
	ret := &rowImpl{
		originalType: reflect.TypeOf(rb.i).Elem(),
		row:          reflect.New(t).Interface(),
	}
	return ret
}

func (r *rowImpl) Writable() interface{} {
	return r.row
}

func (r *rowImpl) SaveInto(dst interface{}) error {
	dstValue := reflect.Indirect(reflect.ValueOf(dst))
	if !dstValue.CanSet() {
		return fmt.Errorf("the type of dst is not acceptable")
	}
	if dstValue.Type() != r.originalType {
		return fmt.Errorf("the type of dst is not acceptable")
	}

	srcValues := reflect.ValueOf(r.row).Elem()
	for i := 0; i < srcValues.NumField(); i++ {
		name := srcValues.Type().Field(i).Name
		if sf, ok := dstValue.Type().FieldByName(name); ok {
			dstValue.FieldByIndex(sf.Index).Set(srcValues.Field(i))
		}
	}
	return nil
}

func (r *rowImpl) AssignFrom(dst interface{}) error {
	srcValue := reflect.Indirect(reflect.ValueOf(dst))
	if !srcValue.CanSet() {
		return fmt.Errorf("the type of dst is not acceptable")
	}
	if srcValue.Type() != r.originalType {
		return fmt.Errorf("the type of dst is not acceptable")
	}

	dstValues := reflect.ValueOf(r.row).Elem()
	for i := 0; i < dstValues.NumField(); i++ {
		name := dstValues.Type().Field(i).Name
		if sf, ok := srcValue.Type().FieldByName(name); ok {
			dstValues.FieldByIndex(sf.Index).Set(srcValue.Field(i))
		}
	}
	return nil
}
