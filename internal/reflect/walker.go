// Package reflect provides function based on reflection.
package reflect

import (
	"reflect"
	"slices"
	"strconv"
	"strings"
)

// TagWalker provides a way to walk through a struct and apply a function to
// each field that is settable.
type TagWalker struct {
	dtag, mtag string
	zero       bool
}

// NewTagWalker creates a new TagWalker with the given default tag name and
// given map tag name. The walker integrates the [go-defaults][go-defaults] and
// [mapstructure][mapstructure] packages to setup default values in the config
// reader. However, the implementation is not dependent on these packages and
// can be used without them or similar packages.
//
// [go-defaults]: <https://github.com/mcuadros/go-defaults>
// [mapstructure]: <https://github.com/go-viper/mapstructure>
func NewTagWalker(dtag, mtag string, zero bool) *TagWalker {
	return &TagWalker{dtag: dtag, mtag: mtag, zero: zero}
}

// Walk walks through the fields of the given value and calls the given
// function with the path and tag of each field that has a tag.
func (w *TagWalker) Walk(
	key string, value any,
	call func(path string, value any),
) {
	w.walk(strings.ToLower(key), reflect.ValueOf(value), call)
}

// walk is the internal walker function that is called recursively for each
// element of the given value. The function calls the given function for each
// value to apply the path and tag of the field to ensure that all paths can be
// provided via environment variables to the config reader.
func (w *TagWalker) walk(
	key string, value reflect.Value,
	call func(path string, value any),
) {
	switch value.Kind() {
	case reflect.Ptr:
		// TODO: Find test case for this code!
		// if value.IsZero() {
		// 	value = reflect.New(value.Type().Elem())
		// }
		w.walk(key, value.Elem(), call)
	case reflect.Slice, reflect.Array:
		for index := 0; index < value.Len(); index++ {
			nkey := w.key(key, strconv.Itoa(index))
			w.walk(nkey, value.Index(index), call)
		}
	case reflect.Map:
		for _, fkey := range value.MapKeys() {
			nkey := w.key(key, fkey.String())
			w.walk(nkey, value.MapIndex(fkey), call)
		}
	case reflect.Struct:
		w.walkStruct(key, value, call)
	default:
		if value.IsValid() && (!value.IsZero() || w.zero) {
			call(key, value.Interface())
		}
	}
}

// walkStruct walks through the fields of the given struct value and calls the
// given function with the path and tag of each field that has a tag. On each
// field it also calls recursively the `walk` function depth-first.
func (w *TagWalker) walkStruct(
	key string, value reflect.Value,
	call func(path string, value any),
) {
	vtype := value.Type()
	num := value.NumField()
	for index := 0; index < num; index++ {
		field := vtype.Field(index)
		if field.IsExported() {
			w.walkField(w.field(key, field),
				value.Field(index), field, call)
		}
	}
}

// walkField walks through the given field value and calls the given function
// with the path and tag of the field. If the field is a struct, the function
// calls the `walkStruct` function to walk through the struct fields. If the
// field is a pointer, slice, array, or map, the function calls the `walk`
// function to walk through the field elements.
func (w *TagWalker) walkField(
	key string, value reflect.Value,
	field reflect.StructField,
	call func(path string, value any),
) {
	switch value.Kind() {
	case reflect.Struct:
		w.walkStruct(key, value, call)
	case reflect.Ptr:
		if value.IsZero() {
			value = reflect.New(value.Type().Elem())
		}
		w.walkField(key, value.Elem(), field, call)
	case reflect.Slice, reflect.Array, reflect.Map:
		if value.Len() == 0 {
			call(key, field.Tag.Get(w.dtag))
		} else {
			w.walk(key, value, call)
		}
	default:
		if value.IsValid() && !value.IsZero() {
			call(key, value.Interface())
		} else {
			call(key, field.Tag.Get(w.dtag))
		}
	}
}

// field returns the field key for the given field and whether it is squashed.
// If the field has a tag, the tag is used as terminal field name. If the tag
// is empty, the field name is used as terminal field name. If the tag contains
// a `squash` option, the key is not extended with the field name.
func (w *TagWalker) field(
	key string, field reflect.StructField,
) string {
	mtag := field.Tag.Get(w.mtag)
	if mtag == "" {
		return w.key(key, field.Name)
	}

	args := strings.Split(mtag, ",")
	if isStruct(field) && slices.Contains(args[1:], "squash") {
		return key
	} else if args[0] != "" {
		return w.key(key, args[0])
	}
	return w.key(key, field.Name)
}

// isStruct evaluates whether the given field is a struct or a pointer to a
// struct.
func isStruct(field reflect.StructField) bool {
	return (field.Type.Kind() == reflect.Struct ||
		field.Type.Kind() == reflect.Ptr &&
			field.Type.Elem().Kind() == reflect.Struct)
}

// key is the default key building function. It concatenates the current key
// with the field name separated by a dot `.`. If the key is empty, the field
// name is used as base key.
func (w *TagWalker) key(key, name string) string {
	if key != "" {
		return key + "." + strings.ToLower(name)
	}
	return strings.ToLower(name)
}
