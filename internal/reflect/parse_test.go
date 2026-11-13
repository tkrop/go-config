package reflect_test

import (
	stdjson "encoding/json"
	"io"
	stdreflect "reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tkrop/go-config/internal/reflect"
	"github.com/tkrop/go-testing/mock"
	"github.com/tkrop/go-testing/test"
	"gopkg.in/yaml.v3"
)

type parseTagEnvParams struct {
	path   string
	input  string
	field  stdreflect.Type
	value  any
	expect any
	error  error
}

var parseTagEnvTestCases = map[string]parseTagEnvParams{
	"json-object": {
		path:   "field",
		input:  "json",
		field:  stdreflect.TypeOf(jsonStruct{}),
		value:  `{"name":"value"}`,
		expect: jsonStruct{Name: "value"},
	},
	"json-pointer-object": {
		path:   "field",
		input:  "json",
		field:  stdreflect.TypeOf((*jsonStruct)(nil)),
		value:  `{"name":"value"}`,
		expect: test.Ptr(jsonStruct{Name: "value"}),
	},
	"already-injected-object": {
		path:   "field",
		input:  "json",
		field:  stdreflect.TypeOf(jsonStruct{}),
		value:  jsonStruct{Name: "value"},
		expect: jsonStruct{Name: "value"},
	},
	"already-injected-string": {
		path:   "field",
		input:  "json",
		field:  stdreflect.TypeOf(""),
		value:  "value",
		expect: "value",
	},

	// Primitive JSON values.
	"json-bool": {
		path:   "field",
		input:  "json",
		field:  stdreflect.TypeOf(false),
		value:  "true",
		expect: true,
	},
	"json-string": {
		path:   "field",
		input:  "json",
		field:  stdreflect.TypeOf(""),
		value:  `"value"`,
		expect: `"value"`,
	},
	"json-int": {
		path:   "field",
		input:  "json",
		field:  stdreflect.TypeOf(int(0)),
		value:  "42",
		expect: int(42),
	},
	"json-int8": {
		path:   "field",
		input:  "json",
		field:  stdreflect.TypeOf(int8(0)),
		value:  "42",
		expect: int8(42),
	},
	"json-int16": {
		path:   "field",
		input:  "json",
		field:  stdreflect.TypeOf(int16(0)),
		value:  "42",
		expect: int16(42),
	},
	"json-int32": {
		path:   "field",
		input:  "json",
		field:  stdreflect.TypeOf(int32(0)),
		value:  "42",
		expect: int32(42),
	},
	"json-int64": {
		path:   "field",
		input:  "json",
		field:  stdreflect.TypeOf(int64(0)),
		value:  "42",
		expect: int64(42),
	},
	"json-uint": {
		path:   "field",
		input:  "json",
		field:  stdreflect.TypeOf(uint(0)),
		value:  "42",
		expect: uint(42),
	},
	"json-uint8": {
		path:   "field",
		input:  "json",
		field:  stdreflect.TypeOf(uint8(0)),
		value:  "42",
		expect: uint8(42),
	},
	"json-uint16": {
		path:   "field",
		input:  "json",
		field:  stdreflect.TypeOf(uint16(0)),
		value:  "42",
		expect: uint16(42),
	},
	"json-uint32": {
		path:   "field",
		input:  "json",
		field:  stdreflect.TypeOf(uint32(0)),
		value:  "42",
		expect: uint32(42),
	},
	"json-uint64": {
		path:   "field",
		input:  "json",
		field:  stdreflect.TypeOf(uint64(0)),
		value:  "42",
		expect: uint64(42),
	},
	"json-uintptr": {
		path:   "field",
		input:  "json",
		field:  stdreflect.TypeOf(uintptr(0)),
		value:  "42",
		expect: uintptr(42),
	},
	"json-float32": {
		path:   "field",
		input:  "json",
		field:  stdreflect.TypeOf(float32(0)),
		value:  "1.25",
		expect: float32(1.25),
	},
	"json-float64": {
		path:   "field",
		input:  "json",
		field:  stdreflect.TypeOf(float64(0)),
		value:  "1.25",
		expect: float64(1.25),
	},
	"json-complex64-unsupported": {
		path:  "field",
		input: "json",
		field: stdreflect.TypeOf(complex64(0)),
		value: "1",
		error: reflect.NewErrTagWalker("json parsing", "field", "1",
			stdjson.Unmarshal([]byte("1"), new(complex64))),
	},
	"json-complex128-unsupported": {
		path:  "field",
		input: "json",
		field: stdreflect.TypeOf(complex128(0)),
		value: "1",
		error: reflect.NewErrTagWalker("json parsing", "field", "1",
			stdjson.Unmarshal([]byte("1"), new(complex128))),
	},
	"json-int-pointer": {
		path:   "field",
		input:  "json",
		field:  stdreflect.TypeOf((*int)(nil)),
		value:  "42",
		expect: test.Ptr(42),
	},

	// Collections and nested values.
	"json-int-slice": {
		path:   "field",
		input:  "json",
		field:  stdreflect.TypeOf([]int{}),
		value:  "[1,2,3]",
		expect: []int{1, 2, 3},
	},
	"json-int-array": {
		path:   "field",
		input:  "json",
		field:  stdreflect.TypeOf([3]int{}),
		value:  "[1,2,3]",
		expect: [3]int{1, 2, 3},
	},
	"json-byte-slice": {
		path:   "field",
		input:  "json",
		field:  stdreflect.TypeOf([]byte{}),
		value:  `"AQID"`,
		expect: []byte{1, 2, 3},
	},
	"json-string-int-map": {
		path:   "field",
		input:  "json",
		field:  stdreflect.TypeOf(map[string]int{}),
		value:  `{"one":1,"two":2}`,
		expect: map[string]int{"one": 1, "two": 2},
	},
	"json-interface-map": {
		path:  "field",
		input: "json",
		field: stdreflect.TypeOf(map[string]any{}),
		value: `{"name":"value","count":2}`,
		expect: map[string]any{
			"name":  "value",
			"count": float64(2),
		},
	},
	"json-nested-collections": {
		path:  "field",
		input: "json",
		field: stdreflect.TypeOf(map[string][][2]int{}),
		value: `{"values":[[1,2],[3,4]]}`,
		expect: map[string][][2]int{
			"values": {{1, 2}, {3, 4}},
		},
	},
	"json-array-shorter-than-field": {
		path:   "field",
		input:  "json",
		field:  stdreflect.TypeOf([3]int{}),
		value:  "[1,2]",
		expect: [3]int{1, 2, 0},
	},

	"object-map": {
		path:  "field",
		input: "json",
		field: stdreflect.TypeOf(jsonStruct{}),
		value: map[string]any{"name": "value"},
	},
	"field-map": {
		path:  "field",
		input: "json",
		field: stdreflect.TypeOf(map[string]any{}),
		value: map[string]any{"name": "value"},
	},
	"nil-map": {
		path:  "field",
		input: "json",
		field: stdreflect.TypeOf(jsonStruct{}),
		value: map[string]any(nil),
	},
	"missing-value": {
		path:  "field",
		input: "json",
		field: stdreflect.TypeOf(jsonStruct{}),
	},
	"empty-value": {
		path:  "field",
		input: "json",
		field: stdreflect.TypeOf(jsonStruct{}),
		error: reflect.NewErrTagWalker("json parsing", "field", "",
			stdjson.Unmarshal([]byte(""), &jsonStruct{})),
		value: "",
	},
	"unsupported-format": {
		path:  "field",
		input: "yaml",
		field: stdreflect.TypeOf(jsonStruct{}),
		error: reflect.NewErrTagWalker(
			"env format", "field", "yaml", nil,
		),
	},
	"unexpected-end": {
		path:  "field",
		input: "json",
		field: stdreflect.TypeOf(jsonStruct{}),
		value: `{"name":`,
		error: reflect.NewErrTagWalker("json parsing", "field", `{"name":`,
			stdjson.Unmarshal([]byte(`{"name":`), &jsonStruct{})),
	},
	"wrong-injected-type": {
		path:  "field",
		input: "json",
		field: stdreflect.TypeOf(jsonStruct{}),
		value: 42,
		error: reflect.NewErrTagWalker("value type", "field", 42, nil),
	},
}

func TestParseTagEnv(t *testing.T) {
	test.Map(t, parseTagEnvTestCases).
		Run(func(t test.Test, param parseTagEnvParams) {
			// Given

			// When
			value, err := reflect.NewEnvParser(
				func(path string) any {
					if path != param.path {
						return nil
					}
					return param.value
				}).Parse(param.path, param.input, param.field)

			// Then
			assert.Equal(t, param.expect, value)
			assert.Equal(t, param.error, err)
		})
}

type jsonStruct struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type parseTagJsonParams struct {
	path   string
	input  string
	field  stdreflect.Type
	expect any
	error  error
}

var parseTagJsonTestCases = map[string]parseTagJsonParams{
	"object": {
		path:   "field",
		input:  `{"name":"value"}`,
		field:  stdreflect.TypeOf(jsonStruct{}),
		expect: jsonStruct{Name: "value"},
	},
	"pointer-object": {
		path:   "field",
		input:  `{"name":"value"}`,
		field:  stdreflect.TypeOf((*jsonStruct)(nil)),
		expect: test.Ptr(jsonStruct{Name: "value"}),
	},
	"unexpected-end": {
		path:  "field",
		input: `{"name":`,
		field: stdreflect.TypeOf(jsonStruct{}),
		error: reflect.NewErrTagWalker("json parsing", "field", `{"name":`,
			stdjson.Unmarshal([]byte(`{"name":`), &jsonStruct{})),
	},
	"relaxed-json": {
		path:  "field",
		input: `{name: value}`,
		field: stdreflect.TypeOf(jsonStruct{}),
		error: reflect.NewErrTagWalker(
			"json parsing", "field", `{name: value}`,
			stdjson.Unmarshal(
				[]byte(`{name: value}`), &jsonStruct{}),
		),
	},
}

func TestParseTagJson(t *testing.T) {
	test.Map(t, parseTagJsonTestCases).
		Run(func(t test.Test, param parseTagJsonParams) {
			// Given

			// When
			value, err := reflect.ParseTagJson(
				param.path, param.input, param.field)

			// Then
			assert.Equal(t, param.expect, value)
			assert.Equal(t, param.error, err)
		})
}

type parseTagJson5ExtParams struct {
	path   string
	input  string
	field  stdreflect.Type
	expect any
	error  error
}

var parseTagJson5ExtTestCases = map[string]parseTagJson5ExtParams{
	"relaxed-object": {
		path:   "field",
		input:  `{name: value, count: 2,}`,
		field:  stdreflect.TypeOf(jsonStruct{}),
		expect: jsonStruct{Name: "value", Count: 2},
	},
	"unexpected-end": {
		path:  "field",
		input: `{name: value, count: 2,`,
		field: stdreflect.TypeOf(jsonStruct{}),
		error: reflect.NewErrTagWalker(
			"json parsing", "field",
			`{name: value, count: 2,`, io.ErrUnexpectedEOF,
		),
	},
}

func TestParseTagJson5Ext(t *testing.T) {
	test.Map(t, parseTagJson5ExtTestCases).
		// TODO: disable me!
		Filter(test.Not(test.Pattern[parseTagJson5ExtParams](".*"))).
		Run(func(t test.Test, param parseTagJson5ExtParams) {
			// Given

			// When
			value, err := reflect.ParseTagJson5Ext(
				param.path, param.input, param.field)

			// Then
			assert.Equal(t, param.expect, value)
			assert.Equal(t, param.error, err)
		})
}

type yamlStruct struct {
	Name  string
	Count int
}

type parseTagYamlParams struct {
	path   string
	input  string
	field  stdreflect.Type
	expect any
	error  error
}

var parseTagYamlTestCases = map[string]parseTagYamlParams{
	"string": {
		path:   "field",
		input:  "value",
		field:  stdreflect.TypeOf(""),
		expect: "value",
	},
	"struct": {
		path:   "field",
		input:  "{name: value, count: 2}",
		field:  stdreflect.TypeOf(yamlStruct{}),
		expect: yamlStruct{Name: "value", Count: 2},
	},
	"complex": {
		path:   "field",
		input:  "1+2i",
		field:  stdreflect.TypeOf(complex64(0)),
		expect: complex64(1 + 2i),
	},
	"pointer-complex": {
		path:   "field",
		input:  "1+2i",
		field:  stdreflect.TypeOf((*complex128)(nil)),
		expect: test.Ptr(complex128(1 + 2i)),
	},
	"complex-slice": {
		path:   "field",
		input:  "[1+2i, 3+4i]",
		field:  stdreflect.TypeOf([]complex64{}),
		expect: []complex64{1 + 2i, 3 + 4i},
	},
	"pointer-complex-slice": {
		path:   "field",
		input:  "[1+2i, 3+4i]",
		field:  stdreflect.TypeOf((*[]complex128)(nil)),
		expect: test.Ptr([]complex128{1 + 2i, 3 + 4i}),
	},
	"yaml-parse-error": {
		path:  "field",
		input: "a,b",
		field: stdreflect.TypeOf([]string{}),
		error: reflect.NewErrTagWalker("yaml parsing", "field", "a,b",
			&yaml.TypeError{Errors: []string{
				"line 1: cannot unmarshal !!str `a,b` into []string",
			}}),
	},
	"complex-parse-error": {
		path:  "field",
		input: "invalid",
		field: stdreflect.TypeOf(complex64(0)),
		error: reflect.NewErrTagWalker("complex parsing", "field", "invalid",
			&strconv.NumError{
				Func: "ParseComplex",
				Num:  "invalid",
				Err:  strconv.ErrSyntax,
			}),
	},
}

func TestParseTagYaml(t *testing.T) {
	test.Map(t, parseTagYamlTestCases).
		Run(func(t test.Test, param parseTagYamlParams) {
			// Given

			// When
			value, err := reflect.ParseTagYaml(
				param.path, param.input, param.field)

			// Then
			assert.Equal(t, param.expect, value)
			assert.Equal(t, param.error, err)
		})
}

type toComplexParams struct {
	setup  mock.SetupFunc
	parsed any
	field  stdreflect.Type
	expect any
	error  error
}

var toComplexTestCases = map[string]toComplexParams{
	"complex64": {
		parsed: "1+2i",
		field:  stdreflect.TypeOf(complex64(0)),
		expect: complex64(1 + 2i),
	},
	"complex128": {
		parsed: "1+2i",
		field:  stdreflect.TypeOf(complex128(0)),
		expect: complex128(1 + 2i),
	},
	"complex64-error": {
		parsed: "invalid",
		field:  stdreflect.TypeOf(complex64(0)),
		expect: complex64(0),
		error: &strconv.NumError{
			Func: "ParseComplex",
			Num:  "invalid",
			Err:  strconv.ErrSyntax,
		},
	},
	"complex128-error": {
		parsed: "invalid",
		field:  stdreflect.TypeOf(complex128(0)),
		expect: complex128(0),
		error: &strconv.NumError{
			Func: "ParseComplex",
			Num:  "invalid",
			Err:  strconv.ErrSyntax,
		},
	},
	"complex64-slice": {
		parsed: []string{"1+2i", "3+4i"},
		field:  stdreflect.TypeOf([]complex64{}),
		expect: []complex64{1 + 2i, 3 + 4i},
	},
	"complex128-slice": {
		parsed: []string{"1+2i", "3+4i"},
		field:  stdreflect.TypeOf([]complex128{}),
		expect: []complex128{1 + 2i, 3 + 4i},
	},
	"complex64-slice-error": {
		parsed: []string{"1+2i", "invalid"},
		field:  stdreflect.TypeOf([]complex64{}),
		expect: []complex64(nil),
		error: &strconv.NumError{
			Func: "ParseComplex",
			Num:  "invalid",
			Err:  strconv.ErrSyntax,
		},
	},
	"complex128-slice-error": {
		parsed: []string{"1+2i", "invalid"},
		field:  stdreflect.TypeOf([]complex128{}),
		expect: []complex128(nil),
		error: &strconv.NumError{
			Func: "ParseComplex",
			Num:  "invalid",
			Err:  strconv.ErrSyntax,
		},
	},
	"complex64-array": {
		parsed: []string{"1+2i", "3+4i"},
		field:  stdreflect.TypeOf([2]complex64{}),
		expect: []complex64{1 + 2i, 3 + 4i},
	},
	"complex128-array": {
		parsed: []string{"1+2i", "3+4i"},
		field:  stdreflect.TypeOf([2]complex128{}),
		expect: []complex128{1 + 2i, 3 + 4i},
	},
	"complex64-pointer": {
		parsed: "1+2i",
		field:  stdreflect.TypeOf((*complex64)(nil)),
		expect: test.Ptr(complex64(1 + 2i)),
	},
	"complex128-pointer": {
		parsed: "1+2i",
		field:  stdreflect.TypeOf((*complex128)(nil)),
		expect: test.Ptr(complex128(1 + 2i)),
	},
	"pointer-error": {
		parsed: "invalid",
		field:  stdreflect.TypeOf((*complex64)(nil)),
		error: &strconv.NumError{
			Func: "ParseComplex",
			Num:  "invalid",
			Err:  strconv.ErrSyntax,
		},
	},
	"unsupported-type": {
		setup: test.Panic(reflect.NewErrTagWalker(
			"unsupported type", "<unknown>",
			stdreflect.TypeOf(""), nil,
		).Error()),
		field: stdreflect.TypeOf(""),
	},
	"unsupported-slice-element": {
		setup: test.Panic(reflect.NewErrTagWalker(
			"unsupported type", "<unknown>",
			stdreflect.TypeOf([]int{}), nil,
		).Error()),
		field: stdreflect.TypeOf([]int{}),
	},
}

func TestToComplex(t *testing.T) {
	test.Map(t, toComplexTestCases).
		RunSeq(func(t test.Test, param toComplexParams) {
			// Given
			mock.NewMocks(t).Expect(param.setup)

			// When
			value, err := reflect.ToComplex(param.parsed, param.field)

			// Then
			assert.Equal(t,
				[2]any{param.expect, param.error}, [2]any{value, err})
		})
}
