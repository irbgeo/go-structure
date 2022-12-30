package structure

import (
	"reflect"
)

type Structure interface {
	// Struct returns struct.
	Struct() any
	// AddTags adds tags by getTag function to struct's fields.
	AddTags(getTag func(fieldName string) string)
	// SaveInto save data to dst.
	SaveInto(dst any) error
	// AssignFrom load data to struct.
	AssignFrom(src any) error
}

type structure struct {
	st any
}

// New returns Structure by i.
func New(i any) (Structure, error) {
	if reflect.TypeOf(i).Kind() != reflect.Ptr {
		return nil, NeedPrtTypeErr
	}
	s := &structure{
		st: i,
	}
	return s, nil
}

func (s *structure) Struct() any {
	return s.st
}

func (s *structure) AddTags(getTag func(fieldName string) string) {
	var fields []reflect.StructField
	values := reflect.ValueOf(s.st).Elem()
	for i := 0; i < values.NumField(); i++ {
		f := values.Type().Field(i)
		f.Tag = reflect.StructTag(getTag(f.Name))
		fields = append(fields, f)
	}

	st := reflect.StructOf(fields)
	s.st = reflect.New(st).Interface()
}

func (s *structure) SaveInto(dst any) error {
	srcValue := reflect.ValueOf(s.st).Elem()
	if dstMap, ok := dst.(map[string]any); ok {
		toMap(dstMap, srcValue)
		return nil
	}

	dstValue := reflect.Indirect(reflect.ValueOf(dst))
	if !dstValue.CanSet() {
		return immutableErr
	}

	copy(dstValue, srcValue)
	return nil
}

func (s *structure) AssignFrom(src any) error {
	dstValue := reflect.ValueOf(s.st).Elem()
	if !dstValue.CanSet() {
		return immutableErr
	}
	if srcMap, ok := src.(map[string]any); ok {
		fromMap(dstValue, srcMap)
		return nil
	}
	srcValue := reflect.Indirect(reflect.ValueOf(src))

	copy(dstValue, srcValue)
	return nil
}

// Merge two structures by field's names and types.
func Merge(dst, src any) error {
	srcValue := reflect.Indirect(reflect.ValueOf(src))
	if !srcValue.CanSet() {
		return immutableErr
	}

	dstValue := reflect.Indirect(reflect.ValueOf(dst))
	if !dstValue.CanSet() {
		return immutableErr
	}

	copy(dstValue, srcValue)
	return nil
}

func copy(dst, src reflect.Value) {
	for i := 0; i < dst.NumField(); i++ {
		name := dst.Type().Field(i).Name
		if sf, ok := src.Type().FieldByName(name); ok && sf.Type == dst.Type().Field(i).Type {
			dst.FieldByIndex(sf.Index).Set(src.Field(i))
		}
	}
}

// SaveStructToMap saves data into dst map from src struct.
// key of map equals struct field name.
func SaveStructToMap(dst map[string]any, src any) error {
	toMap(dst, reflect.ValueOf(src).Elem())
	return nil
}

func toMap(dst map[string]any, srcValue reflect.Value) {
	for i := 0; i < srcValue.NumField(); i++ {
		field := srcValue.Type().Field(i)
		dst[field.Name] = srcValue.Field(i).Interface()
	}
}

// AssignStructFromMap saves data into dst struct from src map.
// key of map equals struct field name.
func AssignStructFromMap(dst any, src map[string]any) error {
	dstValue := reflect.ValueOf(dst).Elem()
	if !dstValue.CanSet() {
		return immutableErr
	}
	fromMap(dstValue, src)
	return nil
}

func fromMap(dstValue reflect.Value, src map[string]any) {
	for k, v := range src {
		val := dstValue.FieldByName(k)
		if val.IsValid() {
			val.Set(reflect.ValueOf(v))
		}
	}
}
