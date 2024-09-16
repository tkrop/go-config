package reflect_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tkrop/go-config/internal/reflect"
	"github.com/tkrop/go-testing/test"
)

// tagWalkerParam contains a value and the expected tags.
type tagWalkerParam struct {
	value       any
	path        string
	expectPaths []string
}

// testTagWalkerParams contains test cases for TagWalker.Walk.
var testTagWalkerParams = map[string]tagWalkerParam{
	// Ignore non-struct values.
	"bool": {
		value: true,
	},
	"pbool": {
		value: new(bool),
	},
	"sbool": {
		value: []bool{},
	},
	"int": {
		value: int(1),
	},
	"pint": {
		value: new(int),
	},
	"sint": {
		value: []int{},
	},
	"uint": {
		value: uint(1),
	},
	"puint": {
		value: new(uint),
	},
	"suint": {
		value: []uint{},
	},
	"float": {
		value: float64(1.0),
	},
	"complex": {
		value: complex128(1.0),
	},
	"pcomplex": {
		value: new(complex128),
	},
	"scomplex": {
		value: []complex128{},
	},
	"pfloat": {
		value: new(float64),
	},
	"sfloat": {
		value: []float64{},
	},
	"string": {
		value: string(""),
	},
	"pstring": {
		value: new(string),
	},
	"sstring": {
		value: []string{},
	},
	"byte": {
		value: byte(0),
	},
	"pbyte": {
		value: new(byte),
	},
	"sbyte": {
		value: []byte{},
	},
	"rune": {
		value: rune(0),
	},
	"prune": {
		value: new(rune),
	},
	"srune": {
		value: []rune{},
	},
	"any": {
		value: any(0),
	},
	"pany": {
		value: new(any),
	},
	"sany": {
		value: []any{},
	},

	// Test struct fields.
	"struct-ints": {
		value: struct {
			I   int    `tag:"int"`
			PI  *int   `tag:"*int"`
			SI  []int  `tag:"[]int"`
			PSI *[]int `tag:"*[]int"`
			I8  int8   `tag:"int8"`
			I16 int16  `tag:"int16"`
			I32 int32  `tag:"int32"`
			I64 int64  `tag:"int64"`
		}{},
		expectPaths: []string{
			"i", "pi", "si", "psi",
			"i8", "i16", "i32", "i64",
		},
	},
	"struct-uints": {
		value: struct {
			UI   uint    `tag:"uint"`
			PUI  *uint   `tag:"*uint"`
			SUI  []uint  `tag:"[]uint"`
			PSUI *[]uint `tag:"*[]uint"`
			UI8  uint8   `tag:"uint8"`
			UI16 uint16  `tag:"uint16"`
			UI32 uint32  `tag:"uint32"`
			UI64 uint64  `tag:"uint64"`
		}{},
		expectPaths: []string{
			"ui", "pui", "sui", "psui",
			"ui8", "ui16", "ui32", "ui65",
		},
	},
	"struct-floats": {
		value: struct {
			F32   float32    `tag:"float32"`
			F64   float64    `tag:"float64"`
			PF32  *float32   `tag:"*float32"`
			PF64  *float64   `tag:"*float64"`
			SF32  []float32  `tag:"[]float32"`
			SF64  []float64  `tag:"[]float64"`
			PSF32 *[]float32 `tag:"*[]float32"`
			PSF64 *[]float64 `tag:"*[]float64"`
		}{},
		expectPaths: []string{
			"f32", "f64", "pf32", "pf64",
			"sf32", "sf64", "psf32", "psf64",
		},
	},
	"struct-complex": {
		value: struct {
			F32   complex64     `tag:"complex64"`
			F64   complex128    `tag:"complex128"`
			PF32  *complex64    `tag:"*complex64"`
			PF64  *complex128   `tag:"*complex128"`
			SF32  []complex64   `tag:"[]complex64"`
			SF64  []complex128  `tag:"[]complex128"`
			PSF32 *[]complex64  `tag:"*[]complex64"`
			PSF64 *[]complex128 `tag:"*[]complex128"`
		}{},
		expectPaths: []string{
			"f32", "f64", "pf32", "pf64",
			"sf32", "sf64", "psf32", "psf64",
		},
	},
	"struct-strings": {
		value: struct {
			S   string  `tag:"string"`
			PS  *string `tag:"*string"`
			B   byte    `tag:"uint8"`
			PB  *byte   `tag:"*uint8"`
			SB  []byte  `tag:"[]uint8"`
			PSB *[]byte `tag:"*[]uint8"`
			R   rune    `tag:"int32"`
			PR  *rune   `tag:"*int32"`
			SR  []rune  `tag:"[]int32"`
			PSR *[]rune `tag:"*[]int32"`
		}{},
		expectPaths: []string{
			"s", "ps",
			"b", "pb", "sb", "psb",
			"r", "pr", "sr", "psr",
		},
	},

	// Test struct with nested structs.
	"struct": {
		value: struct {
			A any `tag:"interface {}"`
		}{},
		expectPaths: []string{"a"},
	},
	"struct-map": {
		value: struct {
			M map[string]any `tag:"map[string]interface {}"`
		}{},
		expectPaths: []string{"m"},
	},
	"struct-struct": {
		value: struct {
			S struct {
				A any `tag:"interface {}"`
			} `tag:"struct { A interface {} \"tag:\\\"interface {}\\\"\" }"`
		}{},
		expectPaths: []string{"s.a", "s"},
	},
	"struct-ptr-struct": {
		value: struct {
			S *struct {
				A any `tag:"interface {}"`
			} `tag:"*struct { A interface {} \"tag:\\\"interface {}\\\"\" }"`
		}{},
		expectPaths: []string{"s.a", "s"},
	},
	"struct-slice-struct": {
		value: struct {
			S []struct {
				A any `tag:"interface {}"`
			} `tag:"[]struct { A interface {} \"tag:\\\"interface {}\\\"\" }"`
		}{},
		expectPaths: []string{"s.a", "s"},
	},
	"struct-slice-ptr-struct": {
		value: struct {
			S []*struct {
				A any `tag:"interface {}"`
			} `tag:"[]*struct { A interface {} \"tag:\\\"interface {}\\\"\" }"`
		}{},
		expectPaths: []string{"s.a", "s"},
	},
	"struct-ptr-slice-ptr-struct": {
		value: struct {
			S *[]*struct {
				A any `tag:"interface {}"`
			} `tag:"*[]*struct { A interface {} \"tag:\\\"interface {}\\\"\" }"`
		}{},
		expectPaths: []string{"s.a", "s"},
	},

	// Test pointer struct with nested structs.
	"ptr-struct": {
		value: &struct {
			A any `tag:"interface {}"`
		}{},
		expectPaths: []string{"a"},
	},
	"ptr-struct-map": {
		value: &struct {
			M map[string]any `tag:"map[string]interface {}"`
		}{},
		expectPaths: []string{"m"},
	},
	"ptr-struct-struct": {
		value: &struct {
			S struct {
				A any `tag:"interface {}"`
			} `tag:"struct { A interface {} \"tag:\\\"interface {}\\\"\" }"`
		}{},
		expectPaths: []string{"s.a", "s"},
	},
	"ptr-struct-ptr-struct": {
		value: &struct {
			S *struct {
				A any `tag:"interface {}"`
			} `tag:"*struct { A interface {} \"tag:\\\"interface {}\\\"\" }"`
		}{},
		expectPaths: []string{"s.a", "s"},
	},
	"ptr-struct-slice-struct": {
		value: &struct {
			S []struct {
				A any `tag:"interface {}"`
			} `tag:"[]struct { A interface {} \"tag:\\\"interface {}\\\"\" }"`
		}{},
		expectPaths: []string{"s.a", "s"},
	},
	"ptr-struct-slice-ptr-struct": {
		value: &struct {
			S []*struct {
				A any `tag:"interface {}"`
			} `tag:"[]*struct { A interface {} \"tag:\\\"interface {}\\\"\" }"`
		}{},
		expectPaths: []string{"s.a", "s"},
	},
	"ptr-struct-ptr-slice-ptr-struct": {
		value: &struct {
			S *[]*struct {
				A any `tag:"interface {}"`
			} `tag:"*[]*struct { A interface {} \"tag:\\\"interface {}\\\"\" }"`
		}{},
		expectPaths: []string{"s.a", "s"},
	},
}

// TestTagWalker_Walk tests TagWalker.Walk.
func TestTagWalker_Walk(t *testing.T) {
	test.Map(t, testTagWalkerParams).
		Run(func(t test.Test, param tagWalkerParam) {
			// Given
			var paths []string
			walker := reflect.NewTagWalker("tag", nil)

			// When
			walker.Walk(param.value, param.path,
				func(value reflect.Value, path, tag string) {
					paths = append(paths, path)
					assert.Equal(t, tag, value.Type().String())
				})

			// Then
			assert.Equal(t, param.expectPaths, paths)
		})
}
