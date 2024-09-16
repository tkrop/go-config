// Package reflect provides function based on reflection.
package reflect

import (
	"reflect"
	"strings"
)

// Exported reflect types.
type Value = reflect.Value

// TagWalker provides a way to walk through a struct and apply a function to
// each field that is settable.
type TagWalker struct {
	tag  string
	path func(path, name string) string
}

// NewTagWalker creates a new TagWalker with the given tag name and given
// separator for constructing paths.
func NewTagWalker(
	tag string,
	path func(path, name string) string,
) *TagWalker {
	if path == nil {
		// Default path function. Concatenates the current path with the field
		// name separated by a dot `.`. If the path is empty, the field name is
		// used as base path.
		path = func(path, name string) string {
			if path != "" {
				return path + "." + strings.ToLower(name)
			}
			return strings.ToLower(name)
		}
	}
	return &TagWalker{tag: tag, path: path}
}

// Walk walks through the fields of the given value and calls the given
// function with the path and tag of each field that has a tag.
func (w *TagWalker) Walk(
	value any, path string,
	call func(value reflect.Value, path, tag string),
) {
	w.walkTags(reflect.ValueOf(value), path, call)
}

func (w *TagWalker) walkTags(
	value reflect.Value, path string,
	call func(kind reflect.Value, path, tag string),
) {
	if value.Kind() != reflect.Struct {
		switch value.Kind() {
		case reflect.Ptr, reflect.Slice, reflect.Array:
			v := reflect.New(value.Type().Elem())
			w.walkTags(v.Elem(), path, call)
		}
		return // ignore non-struct values.
	}

	vtype := value.Type()
	num := value.NumField()
	for index := 0; index < num; index++ {
		field := vtype.Field(index)
		fvalue := value.Field(index)
		npath := w.path(path, field.Name)
		w.walkTags(fvalue, npath, call)

		tag := field.Tag.Get(w.tag)
		if tag != "" {
			call(fvalue, npath, tag)
		}
	}
}
