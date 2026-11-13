package reflect

import (
	stdjson "encoding/json"
	"reflect"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Lookup retrieves the environment-backed value for a config path.
type Lookup func(path string) any

// EnvParser is responsible for parsing environment-backed tag values.
type EnvParser struct {
	lookup func(path string) any
}

// NewEnvParser creates a new parser for environment-backed tags with the
// given lookup function.
func NewEnvParser(lookup func(path string) any) *EnvParser {
	return &EnvParser{lookup: lookup}
}

// Parse looks up and parses a tagged JSON environment object.
func (p *EnvParser) Parse(
	path, value string, field reflect.Type,
) (any, error) {
	if value != "json" {
		return nil, NewErrTagWalker("env format", path, value, nil)
	}

	env := p.lookup(path)
	if env == nil {
		return nil, nil //nolint:nilnil // no sentinel error possible.
	}

	etype := reflect.TypeOf(env)
	if etype.Kind() == reflect.Map {
		return nil, nil //nolint:nilnil // no sentinel error possible.
	} else if etype == field {
		return env, nil
	}

	json, ok := env.(string)
	if !ok {
		return nil, NewErrTagWalker("value type", path, env, nil)
	}

	return ParseTagJson(path, json, field)
}

// ParseTagJson parses a JSON tag value into the specified field type.
func ParseTagJson(
	path, value string, field reflect.Type,
) (any, error) {
	ptr := reflect.New(field)
	if err := stdjson.Unmarshal([]byte(value), ptr.Interface()); err != nil {
		return nil, NewErrTagWalker("json parsing", path, value, err)
	}

	return ptr.Elem().Interface(), nil
}

// ParseTagJson5Ext parses a JSON tag value into the specified field type.
// It supports complex numbers and pointers to complex numbers.
func ParseTagJson5Ext(
	_, _ string, field reflect.Type,
	// path, value string, field reflect.Type,
) (any, error) {
	ptr := reflect.New(field)
	// TODO: enable me!
	// dec := json.NewDecoder(strings.NewReader(value)).
	// 	Mode(json.Retain | json.Relaxed | json.Extended | json.CaseIgnore)
	// if err := dec.Decode(ptr.Interface()); err != nil {
	// 	return nil, NewErrTagWalker("json parsing", path, value, err)
	// }

	return ptr.Elem().Interface(), nil
}

// ParseTagYaml parses a YAML tag value into the specified field type.
// It supports complex numbers and pointers to complex numbers.
func ParseTagYaml(
	path, value string, field reflect.Type,
) (any, error) {
	ptype := parseType(field)
	ptr := reflect.New(ptype)
	if err := yaml.Unmarshal([]byte(value), ptr.Interface()); err != nil {
		return nil, NewErrTagWalker("yaml parsing", path, value, err)
	}
	parsed := ptr.Elem().Interface()
	if ptype == field {
		return parsed, nil
	}

	number, err := ToComplex(parsed, field)
	if err != nil {
		return nil, NewErrTagWalker("complex parsing", path, parsed, err)
	}
	return number, nil
}

// parseType returns the appropriate type for parsing the default tag value.
// For complex numbers, which YAML doesn't support natively, we return string
// or []string types so that we first can parse the value in yaml.
func parseType(field reflect.Type) reflect.Type {
	switch field.Kind() {
	case reflect.Complex64, reflect.Complex128:
		return reflect.TypeOf("")
	case reflect.Slice, reflect.Array:
		if isComplex(field.Elem()) {
			return reflect.TypeOf([]string{})
		}
	case reflect.Pointer:
		ftype := parseType(field.Elem())
		if ftype != field.Elem() {
			return ftype
		}
	}
	return field
}

// isComplex evaluates whether the given type is a complex number type.
func isComplex(field reflect.Type) bool {
	return field.Kind() == reflect.Complex64 ||
		field.Kind() == reflect.Complex128
}

// Complex bit sizes for parsing.
const (
	// Complex64BitSize is the bit size for parsing complex64 numbers.
	Complex64BitSize = 64
	// Complex128BitSize is the bit size for parsing complex128 numbers.
	Complex128BitSize = 128
)

// ToComplex converts the parsed string/[]string value into the actual complex
// type based on the target field type. It supports complex64 and complex128
// types, as well as pointers to these types.
func ToComplex(parsed any, field reflect.Type) (any, error) {
	switch field.Kind() {
	case reflect.Complex64:
		return parseValue[complex64](parsed.(string), Complex64BitSize)
	case reflect.Complex128:
		return parseValue[complex128](parsed.(string), Complex128BitSize)

	case reflect.Slice, reflect.Array:
		switch field.Elem().Kind() {
		case reflect.Complex64:
			return parseSlice[complex64](parsed.([]string), Complex64BitSize)
		case reflect.Complex128:
			return parseSlice[complex128](parsed.([]string), Complex128BitSize)
		}

	case reflect.Pointer:
		if result, err := ToComplex(parsed, field.Elem()); err == nil {
			return pointer(result), nil
		} else {
			return nil, err
		}
	}
	panic(NewErrTagWalker("unsupported type", "<unknown>", field, nil))
}

// parseValue parses a single complex number string into the specified
// complex type (complex64 or complex128) based on the target bitSize.
func parseValue[T complex64 | complex128](
	str string, bsize int,
) (T, error) {
	number, err := strconv.ParseComplex(strings.TrimSpace(str), bsize)
	if err != nil {
		return T(0), err
	}

	return T(number), nil
}

// parseSlice parses a slice of complex number strings into a slice
// of the specified complex type based on the target bitSize.
func parseSlice[T complex64 | complex128](
	strs []string, bsize int,
) ([]T, error) {
	result := make([]T, len(strs))
	for index, str := range strs {
		if number, err := parseValue[T](str, bsize); err != nil {
			return nil, err
		} else {
			result[index] = number
		}
	}
	return result, nil
}

// pointer wraps the given result into a pointer type and returns the pointer.
func pointer(result any) any {
	ptr := reflect.New(reflect.TypeOf(result))
	ptr.Elem().Set(reflect.ValueOf(result))
	return ptr.Interface()
}
