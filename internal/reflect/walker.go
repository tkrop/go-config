// Package reflect provides function based on reflection.
package reflect

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// ErrTagWalker is a common error to indicate a tag walker error.
var ErrTagWalker = errors.New("tag walker")

// NewErrTagWalker is a convenience method to create a new tag walker error
// with the given path value context wrapping the original error.
func NewErrTagWalker(message, path string, value any, err error) error {
	return fmt.Errorf("%w - %s [%s=%#v]: %w",
		ErrTagWalker, message, path, value, err)
}

// TagWalker provides a way to walk through a struct and apply a function to
// each field that is settable.
type TagWalker struct {
	dtag, mtag string
	zero       bool
	call       func(path string, value any)
	errors     []error
}

// NewTagWalker creates a new TagWalker with the given default tag name, map
// tag name, a flag to determine whether to signal zero values, and a callback
// function to be called for each field with a non-zero default value. The
// walker integrates the [go-defaults][defaults] and [mapstructure][mapstr]
// packages to signal paths and default values. However, the implementation is
// not dependent on these packages and can be used with different tag names and
// packages.
//
// The `TagWalker` is used to setup default values in the [reader][reader]
// using a template config with `mapstructure` and `default` tags.
//
// [defaults]: <https://github.com/mcuadros/go-defaults>
// [mapstr]: <https://github.com/go-viper/mapstructure>
// [reader]: <https://github.com/tkrop/go-config>
func NewTagWalker(
	dtag, mtag string, zero bool,
	call func(path string, value any),
) *TagWalker {
	return &TagWalker{
		dtag: dtag, mtag: mtag,
		zero: zero, call: call,
	}
}

// Walk walks through the fields of the given value and calls the callback
// function with the path and tag of each field that has a tag.
func (w *TagWalker) Walk(path string, value any) error {
	w.walk(strings.ToLower(path), reflect.ValueOf(value))
	return errors.Join(w.errors...)
}

// walk is the internal walker function that is called recursively for each
// element of the given value. The function calls the callback function for each
// value to apply the path and tag of the field to ensure that all paths can be
// provided via environment variables to the config reader.
func (w *TagWalker) walk(path string, value reflect.Value) {
	switch value.Kind() {
	case reflect.Ptr:
		if value.IsZero() {
			value = reflect.New(value.Type().Elem())
		}
		w.walk(path, value.Elem())
	case reflect.Slice, reflect.Array:
		for index := range value.Len() {
			npath := w.path(path, strconv.Itoa(index))
			w.walk(npath, value.Index(index))
		}
	case reflect.Map:
		for _, fpath := range value.MapKeys() {
			npath := w.path(path, fpath.String())
			w.walk(npath, value.MapIndex(fpath))
		}
	case reflect.Struct:
		w.walkStruct(path, value)
	default:
		if value.IsValid() && (!value.IsZero() || w.zero) {
			w.call(path, value.Interface())
		}
	}
}

// walkStruct walks through the fields of the given struct value and calls the
// callback function with the path and tag of each field that has a tag. On
// each field it also calls recursively the `walk` function depth-first.
func (w *TagWalker) walkStruct(path string, value reflect.Value) {
	vtype := value.Type()
	for index := range value.NumField() {
		field := vtype.Field(index)
		if field.IsExported() {
			w.walkField(w.field(path, field),
				field, value.Field(index))
		}
	}
}

// walkField walks through the given field value and calls the callback
// function with the path and tag of the field. If the field is a struct, the
// function calls the `walkStruct` function to walk through the struct fields.
// If the field is a pointer, slice, array, or map, the function calls the
// `walk` function to walk through the field elements.
func (w *TagWalker) walkField(
	path string, field reflect.StructField, value reflect.Value,
) {
	switch value.Kind() {
	case reflect.Struct:
		if field.Tag.Get(w.dtag) != "" {
			w.callField(path, field)
		} else {
			w.walkStruct(path, value)
		}
	case reflect.Ptr:
		if value.IsZero() {
			value = reflect.New(value.Type().Elem())
		}
		w.walkField(path, field, value.Elem())
	case reflect.Slice, reflect.Array, reflect.Map:
		if value.Len() != 0 {
			w.walk(path, value)
		}
		w.callField(path, field)
	default:
		if value.IsValid() && !value.IsZero() {
			w.call(path, value.Interface())
		} else if field.Tag.Get(w.dtag) != "" || w.zero {
			w.callField(path, field)
		}
	}
}

// callField is the generic callback wrapper function that creates a structured
// value for the callback from the attached `default` tag by parsing it as yaml.
// For complex numbers, which YAML doesn't support natively, we first parse the
// value as string/[]string, then convert to the actual complex type.
func (w *TagWalker) callField(path string, field reflect.StructField) {
	if value := field.Tag.Get(w.dtag); value != "" {
		fieldType := field.Type
		parseType := parseType(fieldType)
		ptr := reflect.New(parseType)

		if err := yaml.Unmarshal([]byte(value), ptr.Interface()); err != nil {
			w.errors = append(w.errors,
				NewErrTagWalker("yaml parsing", path, value, err))
			w.call(path, value)
		} else if parseType != fieldType {
			w.callComplex(path, ptr.Elem().Interface(), fieldType)
		} else {
			w.call(path, ptr.Elem().Interface())
		}
	}
}

// callComplex is a specialized callback wrapper function to handle complex
// number types. It converts the parsed string/[]string value into the actual
// complex type and calls the generic callback function. If conversion fails,
// it appends an error to the walker errors and calls the generic callback
// function with the original value.
func (w *TagWalker) callComplex(path string, value any, fieldType reflect.Type) {
	if number, err := toComplex(value, fieldType); err != nil {
		w.errors = append(w.errors,
			NewErrTagWalker("complex parsing", path, value, err))
		w.call(path, value)
	} else {
		w.call(path, number)
	}
}

// path is the default path building function. It concatenates the current path
// with the field name separated by a dot `.`. If the path is empty, the field
// name is used as base path.
func (*TagWalker) path(path, name string) string {
	if path != "" {
		return path + "." + strings.ToLower(name)
	}
	return strings.ToLower(name)
}

// field returns the field path for the given field and whether it is squashed.
// If the field has a tag, the tag is used as terminal field name. If the tag
// is empty, the field name is used as terminal field name. If the tag contains
// a `squash` option, the path is not extended with the field name.
func (w *TagWalker) field(path string, field reflect.StructField) string {
	mtag := field.Tag.Get(w.mtag)
	if mtag == "" {
		return w.path(path, field.Name)
	}

	args := strings.Split(mtag, ",")
	if w.isStruct(field) && slices.Contains(args[1:], "squash") {
		return path
	} else if args[0] != "" {
		return w.path(path, args[0])
	}
	return w.path(path, field.Name)
}

// pointer wraps the given result into a pointer type and returns the pointer.
func pointer(result any) any {
	ptr := reflect.New(reflect.TypeOf(result))
	ptr.Elem().Set(reflect.ValueOf(result))
	return ptr.Interface()
}

// parseType returns the appropriate type for parsing the default tag value.
// For complex numbers, which YAML doesn't support natively, we return string
// or []string types so that we first can parse the value in yaml.
func parseType(fieldType reflect.Type) reflect.Type {
	switch fieldType.Kind() {
	case reflect.Complex64, reflect.Complex128:
		return reflect.TypeOf("")
	case reflect.Slice, reflect.Array:
		if isComplex(fieldType.Elem()) {
			return reflect.TypeOf([]string{})
		}
	case reflect.Ptr:
		ftype := parseType(fieldType.Elem())
		if ftype != fieldType.Elem() {
			return ftype
		}
	}
	return fieldType
}

// isStruct evaluates whether the given field is a struct or a pointer to a
// struct.
func (*TagWalker) isStruct(field reflect.StructField) bool {
	return (field.Type.Kind() == reflect.Struct ||
		field.Type.Kind() == reflect.Ptr &&
			field.Type.Elem().Kind() == reflect.Struct)
}

// isComplex evaluates whether the given type is a complex number type.
func isComplex(elemType reflect.Type) bool {
	return elemType.Kind() == reflect.Complex64 ||
		elemType.Kind() == reflect.Complex128
}

// Complex bit sizes for parsing.
const (
	// Complex64BitSize is the bit size for parsing complex64 numbers.
	Complex64BitSize = 64
	// Complex128BitSize is the bit size for parsing complex128 numbers.
	Complex128BitSize = 128
)

// toComplex converts the parsed string/[]string value into the actual complex
// type based on the target field type. It supports complex64 and complex128
// types, as well as pointers to these types.
func toComplex(parsed any, fieldType reflect.Type) (any, error) {
	switch fieldType.Kind() {
	case reflect.Complex64:
		return parseValue[complex64](parsed.(string), Complex64BitSize)
	case reflect.Complex128:
		return parseValue[complex128](parsed.(string), Complex128BitSize)

	case reflect.Slice, reflect.Array:
		switch fieldType.Elem().Kind() {
		case reflect.Complex64:
			return parseSlice[complex64](parsed.([]string), Complex64BitSize)
		case reflect.Complex128:
			return parseSlice[complex128](parsed.([]string), Complex128BitSize)
		}

	case reflect.Ptr:
		if result, err := toComplex(parsed, fieldType.Elem()); err == nil {
			return pointer(result), nil
		} else {
			return nil, err
		}
	}
	panic(NewErrTagWalker("unsupported type", "<unknown>", fieldType, nil))
}

// parseValue parses a single complex number string into the specified
// complex type (complex64 or complex128) based on the target bitSize.
func parseValue[T complex64 | complex128](
	str string, bitSize int,
) (T, error) {
	if number, err := strconv.ParseComplex(
		strings.TrimSpace(str), bitSize); err != nil {
		return T(0), err
	} else {
		return T(number), nil
	}
}

// parseSlice parses a slice of complex number strings into a slice
// of the specified complex type based on the target bitSize.
func parseSlice[T complex64 | complex128](
	strs []string, bitSize int,
) ([]T, error) {
	result := make([]T, len(strs))
	for index, str := range strs {
		if number, err := parseValue[T](str, bitSize); err != nil {
			return nil, err
		} else {
			result[index] = number
		}
	}
	return result, nil
}
