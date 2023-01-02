package structure

import (
	"fmt"
	"reflect"
)

type Structure interface {
	// Struct returns struct.
	Struct() any
	// ChangeTags changes tags by getNewTag function to struct's fields.
	ChangeTags(getNewTag func(fieldName, fieldTag string) string)
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
	kind := reflect.TypeOf(i).Kind()
	if kind != reflect.Ptr {
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

func (s *structure) ChangeTags(getNewTag func(fieldName, fieldTag string) string) {
	st := changeTags(reflect.ValueOf(s.st).Elem(), getNewTag)
	s.st = reflect.New(st).Interface()
}

func changeTags(v reflect.Value, getNewTag func(fieldName, fieldTag string) string) reflect.Type {
	var fields []reflect.StructField
	for i := 0; i < v.NumField(); i++ {
		f := v.Type().Field(i)
		f.Tag = reflect.StructTag(getNewTag(f.Name, string(f.Tag)))
		if v.Field(i).Kind() == reflect.Struct {
			f.Type = changeTags(v.Field(i), getNewTag)
		}
		fields = append(fields, f)
	}

	return reflect.StructOf(fields)
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
	kind := reflect.TypeOf(src).Kind()
	if kind != reflect.Ptr {
		return fmt.Errorf("src: %w", NeedPrtTypeErr)
	}

	kind = reflect.TypeOf(dst).Kind()
	if kind != reflect.Ptr {
		return fmt.Errorf("dst: %w", NeedPrtTypeErr)
	}

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
		dstField := dst.Type().Field(i)
		if sf, ok := src.Type().FieldByName(dstField.Name); ok && sf.Type == dst.Type().Field(i).Type {
			if dst.FieldByIndex(sf.Index).Kind() == reflect.Struct && src.Field(i).Kind() == reflect.Struct {
				copy(dst.FieldByIndex(sf.Index), src.Field(i))
			}
			dst.FieldByIndex(sf.Index).Set(src.Field(i))
		}
	}
}

// SaveStructToMap saves data into dst map from src struct.
// key of map equals struct field name.
func SaveStructToMap(dst map[string]any, src any) error {
	kind := reflect.TypeOf(src).Kind()
	if kind != reflect.Ptr {
		return fmt.Errorf("src: %w", NeedPrtTypeErr)
	}

	toMap(dst, reflect.ValueOf(src).Elem())
	return nil
}

func toMap(dst map[string]any, srcValue reflect.Value) {
	for i := 0; i < srcValue.NumField(); i++ {
		field := srcValue.Type().Field(i)
		if field.Type.Kind() == reflect.Struct {
			subDst := make(map[string]any)
			toMap(subDst, srcValue.Field(i))
			dst[field.Name] = subDst
			continue
		}
		dst[field.Name] = srcValue.Field(i).Interface()
	}
}

// AssignStructFromMap saves data into dst struct from src map.
// key of map equals struct field name.
func AssignStructFromMap(dst any, src map[string]any) error {
	kind := reflect.TypeOf(dst).Kind()
	if kind != reflect.Ptr {
		return fmt.Errorf("dst: %w", NeedPrtTypeErr)
	}

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
			if subSrc, ok := v.(map[string]any); ok {
				fromMap(val, subSrc)
			}
			val.Set(reflect.ValueOf(v))
		}
	}
}
