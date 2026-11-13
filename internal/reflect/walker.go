package reflect

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

// Parse parses a tagged field value.
type Parse func(path, tag string, field reflect.Type) (any, error)

// Call processes a config path and its parsed value.
type Call func(path string, value any)

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
	parse      Parse
	call       Call
	errors     []error
}

// NewTagWalker creates a new TagWalker with the given default tag name, map
// tag name, a flag to determine whether to signal zero values, a callback
// function for parsing the tag values and a callback function to be called for
// each field with a non-zero default value.
//
// The walker integrates the [go-defaults][defaults] and [mapstructure][mapstr]
// packages to signal paths and default values. However, the implementation is
// not dependent on these packages and can be used with different tag names and
// packages.
//
// The `TagWalker` is used to setup `default` values in the [reader][reader]
// using a template config with `mapstructure` and `default` tags, as well as
// to inject complex environment variables into the configuration.
//
// [defaults]: <https://github.com/mcuadros/go-defaults>
// [mapstr]: <https://github.com/go-viper/mapstructure>
// [reader]: <https://github.com/tkrop/go-config>
func NewTagWalker(
	dtag, mtag string, zero bool, parse Parse, call Call,
) *TagWalker {
	return &TagWalker{
		dtag: dtag, mtag: mtag, zero: zero, parse: parse, call: call,
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
	case reflect.Pointer:
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
			w.callTag(path, value.Interface())
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
	case reflect.Pointer:
		if value.IsZero() {
			value = reflect.New(value.Type().Elem())
		}
		w.walkField(path, field, value.Elem())
	case reflect.Slice, reflect.Array, reflect.Map:
		if value.Len() != 0 {
			w.walk(path, value)
		} else if field.Tag.Get(w.dtag) != "" || w.zero {
			w.callField(path, field)
		}
	default:
		if value.IsValid() && !value.IsZero() {
			w.callTag(path, value.Interface())
		} else if field.Tag.Get(w.dtag) != "" || w.zero {
			w.callField(path, field)
		}
	}
}

// callField parses the field's tag and invokes the callback with the result.
func (w *TagWalker) callField(path string, field reflect.StructField) {
	value := field.Tag.Get(w.dtag)
	parsed, err := w.parse(path, value, field.Type)
	if err != nil {
		w.errors = append(w.errors, err)
		return
	}
	w.callTag(path, parsed)
}

// callTag invokes the callback with the given path and value.
func (w *TagWalker) callTag(path string, value any) {
	w.call(path, value)
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

// isStruct evaluates whether the given field is a struct or a pointer to a
// struct.
func (*TagWalker) isStruct(field reflect.StructField) bool {
	return (field.Type.Kind() == reflect.Struct ||
		field.Type.Kind() == reflect.Pointer &&
			field.Type.Elem().Kind() == reflect.Struct)
}
