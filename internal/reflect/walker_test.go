package reflect_test

import (
	"testing"

	"github.com/tkrop/go-config/internal/reflect"
	"github.com/tkrop/go-testing/mock"
	"github.com/tkrop/go-testing/test"
)

//revive:disable:line-length-limit // go:generate line length.

//go:generate mockgen -package=reflect_test -destination=mock_callback_test.go -source=walker_test.go Callback

//revive:enable:line-length-limit

// Callback is a mock interface for testing.
type Callback interface {
	Call(path string, value any)
}

// Call calls the Call method of the given mocks.
func Call(path string, value any) mock.SetupFunc {
	return func(mocks *mock.Mocks) any {
		return mock.Get(mocks, NewMockCallback).EXPECT().Call(path, value).
			DoAndReturn(mocks.Do(Callback.Call))
	}
}

// tagWalkerParam contains a value and the expected tags.
type tagWalkerParam struct {
	value  any
	key    string
	zero   bool
	expect mock.SetupFunc
}

//revive:disable:nested-structs // simplifies test cases a lot.

// tagWalkerTestCases contains test cases for TagWalker.Walk.
var tagWalkerTestCases = map[string]tagWalkerParam{
	// Test build-in values.
	"nil": {
		value: nil,
	},
	"bool": {
		value:  true,
		expect: Call("", true),
	},
	"int": {
		value:  int(1),
		expect: Call("", 1),
	},
	"uint": {
		value:  uint(1),
		expect: Call("", uint(1)),
	},
	"float": {
		value:  float64(1.0),
		expect: Call("", float64(1.0)),
	},
	"complex": {
		value:  complex128(1.0),
		expect: Call("", complex128(1.0)),
	},
	"string": {
		value:  string("test"),
		expect: Call("", "test"),
	},
	"byte": {
		value:  byte('a'),
		expect: Call("", byte('a')),
	},
	"rune": {
		value:  rune('a'),
		expect: Call("", rune('a')),
	},
	"any": {
		value:  any(1),
		expect: Call("", any(1)),
	},

	// Test build-in pointer values.
	"ptr-bool": {
		value: new(bool),
	},
	"ptr-int": {
		value: new(int),
	},
	"ptr-uint": {
		value: new(uint),
	},
	"ptr-float": {
		value: new(float64),
	},
	"ptr-complex": {
		value: new(complex128),
	},
	"ptr-string": {
		value: new(string),
	},
	"ptr-byte": {
		value: new(byte),
	},
	"ptr-rune": {
		value: new(rune),
	},
	"ptr-any": {
		value: new(any),
	},
	"ptr-slice": {
		value: new([]any),
	},

	// Test build-in slice values.
	"slice-bool": {
		value: []bool{true, false},
		zero:  true,
		expect: mock.Chain(
			Call("0", true),
			Call("1", false),
		),
	},
	"slice-int": {
		value: []int{1, 0},
		zero:  true,
		expect: mock.Chain(
			Call("0", 1),
			Call("1", 0),
		),
	},
	"slice-uint": {
		value: []uint{1, 0},
		zero:  true,
		expect: mock.Chain(
			Call("0", uint(1)),
			Call("1", uint(0)),
		),
	},
	"slice-float": {
		value: []float64{1.0, 0.0},
		zero:  true,
		expect: mock.Chain(
			Call("0", 1.0),
			Call("1", 0.0),
		),
	},
	"slice-complex": {
		value: []complex128{1.0, 0.0},
		zero:  true,
		expect: mock.Chain(
			Call("0", complex128(1.0)),
			Call("1", complex128(0.0)),
		),
	},
	"slice-string": {
		zero:  true,
		value: []string{"test", ""},
		expect: mock.Chain(
			Call("0", "test"),
			Call("1", ""),
		),
	},
	"slice-byte": {
		value: []byte{'a', 'b'},
		expect: mock.Chain(
			Call("0", byte('a')),
			Call("1", byte('b')),
		),
	},
	"slice-rune": {
		value: []rune{'a', 'b'},
		expect: mock.Chain(
			Call("0", rune('a')),
			Call("1", rune('b')),
		),
	},
	"slice-any": {
		value: []any{0, "test"},
		expect: mock.Chain(
			Call("0", 0),
			Call("1", "test"),
		),
	},

	// Test struct with field tags.
	"struct-bool-tags": {
		value: struct {
			hidden  bool
			Visible bool `tag:"visible"`
		}{},
		expect: Call("visible", "visible"),
	},
	"struct-ints-tags": {
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
		expect: mock.Chain(
			Call("i", "int"),
			Call("pi", "*int"),
			Call("si", "[]int"),
			Call("psi", "*[]int"),
			Call("i8", "int8"),
			Call("i16", "int16"),
			Call("i32", "int32"),
			Call("i64", "int64"),
		),
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
		expect: mock.Chain(
			Call("ui", "uint"),
			Call("pui", "*uint"),
			Call("sui", "[]uint"),
			Call("psui", "*[]uint"),
			Call("ui8", "uint8"),
			Call("ui16", "uint16"),
			Call("ui32", "uint32"),
			Call("ui64", "uint64"),
		),
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
		expect: mock.Chain(
			Call("f32", "float32"),
			Call("f64", "float64"),
			Call("pf32", "*float32"),
			Call("pf64", "*float64"),
			Call("sf32", "[]float32"),
			Call("sf64", "[]float64"),
			Call("psf32", "*[]float32"),
			Call("psf64", "*[]float64"),
		),
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
		expect: mock.Chain(
			Call("f32", "complex64"),
			Call("f64", "complex128"),
			Call("pf32", "*complex64"),
			Call("pf64", "*complex128"),
			Call("sf32", "[]complex64"),
			Call("sf64", "[]complex128"),
			Call("psf32", "*[]complex64"),
			Call("psf64", "*[]complex128"),
		),
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
		expect: mock.Chain(
			Call("s", "string"),
			Call("ps", "*string"),
			Call("b", "uint8"),
			Call("pb", "*uint8"),
			Call("sb", "[]uint8"),
			Call("psb", "*[]uint8"),
			Call("r", "int32"),
			Call("pr", "*int32"),
			Call("sr", "[]int32"),
			Call("psr", "*[]int32"),
		),
	},

	// Test structs with field values.
	"struct-all-values": {
		value: struct {
			Bool   bool    `map:"bool" default:"false"`
			Int    int     `map:"int" default:"-2"`
			Uint   uint    `map:"uint" default:"2"`
			Float  float64 `map:"float" default:"3.0"`
			String string  `map:"string" default:"STRING"`
			Byte   byte    `map:"byte" default:"A"`
			Rune   rune    `map:"rune" default:"B"`
			Any    any     `map:"any" default:"ANY"`
		}{
			Bool:   true,
			Int:    int(-1),
			Uint:   uint(1),
			Float:  float64(2.0),
			String: "string",
			Byte:   byte('a'),
			Rune:   rune('b'),
			Any:    map[string]any{"key": "value"},
		},
		expect: mock.Chain(
			Call("bool", true),
			Call("int", -1),
			Call("uint", uint(1)),
			Call("float", float64(2.0)),
			Call("string", "string"),
			Call("byte", byte('a')),
			Call("rune", rune('b')),
			Call("any", map[string]any{"key": "value"}),
		),
	},

	// Test struct with nested structs.
	"struct-empty": {
		value: struct{}{},
	},
	"struct-any": {
		value: struct {
			A any `tag:"any"`
		}{},
		expect: Call("a", "any"),
	},
	"struct-struct": {
		value: struct {
			S struct {
				A any `tag:"any"`
			} `tag:"struct{any}"`
		}{},
		expect: mock.Chain(
			Call("s.a", "any"),
		),
	},
	"struct-ptr-struct": {
		value: struct {
			S *struct {
				A any `tag:"any"`
			} `tag:"*struct{any}"`
		}{},
		expect: mock.Chain(
			Call("s.a", "any"),
		),
	},

	// Test pointer struct with nested structs.
	"ptr-struct-empty": {
		value: &struct{}{},
	},
	"ptr-struct-any": {
		value: &struct {
			A any `tag:"any"`
		}{},
		expect: Call("a", "any"),
	},
	"ptr-struct-struct": {
		value: &struct {
			S struct {
				A any `tag:"any"`
			} `tag:"struct{any}"`
		}{},
		expect: mock.Chain(
			Call("s.a", "any"),
		),
	},
	"ptr-struct-ptr-struct": {
		value: &struct {
			S *struct {
				A any `tag:"any"`
			} `tag:"*struct{any}"`
		}{},
		expect: mock.Chain(
			Call("s.a", "any"),
		),
	},

	// Test struct with nested slices and tags.
	"struct-slice-tag": {
		value: struct {
			S []any `tag:"[]any"`
		}{},
		expect: Call("s", "[]any"),
	},
	"struct-slice-struct-tag": {
		value: struct {
			S []struct {
				A any `tag:"any"`
			} `tag:"[]struct{any}"`
		}{},
		expect: Call("s", "[]struct{any}"),
	},
	"struct-slice-ptr-struct-tag": {
		value: struct {
			S []*struct {
				A any `tag:"any"`
			} `tag:"[]*struct{any}"`
		}{},
		expect: Call("s", "[]*struct{any}"),
	},
	"struct-ptr-slice-ptr-struct-tag": {
		value: struct {
			S *[]*struct {
				A any `tag:"any"`
			} `tag:"*[]*struct{any}"`
		}{},
		expect: Call("s", "*[]*struct{any}"),
	},

	// Test struct with nested slices and values.
	"struct-slice-value": {
		value: struct {
			S []any `tag:"[]any"`
		}{S: []any{1, 2}},
		expect: mock.Chain(
			Call("s.0", 1),
			Call("s.1", 2),
		),
	},
	"struct-slice-struct-value": {
		value: struct {
			S []struct {
				A any `tag:"any"`
			} `tag:"[]struct{any}"`
		}{S: []struct {
			A any `tag:"any"`
		}{{A: 1}, {A: 2}}},
		expect: mock.Chain(
			Call("s.0.a", 1),
			Call("s.1.a", 2),
		),
	},
	"struct-slice-ptr-struct-value": {
		value: struct {
			S []*struct {
				A any `tag:"any"`
			} `tag:"[]*struct{any}"`
		}{S: []*struct {
			A any `tag:"any"`
		}{{A: 1}, {A: 2}}},
		expect: mock.Chain(
			Call("s.0.a", 1),
			Call("s.1.a", 2),
		),
	},
	"struct-ptr-slice-ptr-struct-value": {
		value: struct {
			S *[]*struct {
				A any `tag:"any"`
			} `tag:"*[]*struct{any}"`
		}{S: &[]*struct {
			A any `tag:"any"`
		}{{A: 1}, {A: 2}}},
		expect: mock.Chain(
			Call("s.0.a", 1),
			Call("s.1.a", 2),
		),
	},

	// Test struct with nested maps.
	"struct-map-tag": {
		value: struct {
			M map[string]any `tag:"map[string]any"`
		}{},
		expect: Call("m", "map[string]any"),
	},
	"struct-ptr-map-tag": {
		value: struct {
			M *map[string]any `tag:"*map[string]any"`
		}{},
		expect: Call("m", "*map[string]any"),
	},
	"struct-map-struct-tag": {
		value: &struct {
			M map[string]struct {
				A any `tag:"any"`
			} `tag:"map[string]struct{any}"`
		}{},
		expect: Call("m", "map[string]struct{any}"),
	},
	"struct-ptr-map-struct-tag": {
		value: &struct {
			M map[string]struct {
				A any `tag:"any"`
			} `tag:"*map[string]struct{any}"`
		}{},
		expect: Call("m", "*map[string]struct{any}"),
	},
	"struct-map-ptr-struct-tag": {
		value: &struct {
			M map[string]*struct {
				A any `tag:"any"`
			} `tag:"map[string]*struct{any}"`
		}{},
		expect: Call("m", "map[string]*struct{any}"),
	},

	// Test struct with nested maps.
	"struct-map-value": {
		value: struct {
			M map[string]any `tag:"map[string]any"`
		}{M: map[string]any{"key": "value"}},
		expect: Call("m.key", "value"),
	},
	"struct-ptr-map-value": {
		value: struct {
			M *map[string]any `tag:"*map[string]any"`
		}{M: &map[string]any{"key": "value"}},
		expect: Call("m.key", "value"),
	},
	"struct-map-struct-value": {
		value: struct {
			M map[string]struct {
				A any `tag:"any"`
			} `tag:"map[string]struct{any}"`
		}{M: map[string]struct {
			A any `tag:"any"`
		}{"key-0": {A: 1}, "key-1": {A: 2}}},
		expect: mock.Setup(
			Call("m.key-0.a", 1),
			Call("m.key-1.a", 2),
		),
	},
	"struct-ptr-map-struct-value": {
		value: struct {
			M *map[string]struct {
				A any `tag:"any"`
			} `tag:"*map[string]struct{any}"`
		}{M: &map[string]struct {
			A any `tag:"any"`
		}{"key-0": {A: 1}, "key-1": {A: 2}}},
		expect: mock.Setup(
			Call("m.key-0.a", 1),
			Call("m.key-1.a", 2),
		),
	},
	"struct-ptr-map-ptr-struct-value": {
		value: struct {
			M *map[string]*struct {
				A any `tag:"any"`
			} `tag:"map[string]*struct{any}"`
		}{M: &map[string]*struct {
			A any `tag:"any"`
		}{"key-0": {A: 1}, "key-1": {A: 2}}},
		expect: mock.Setup(
			Call("m.key-0.a", 1),
			Call("m.key-1.a", 2),
		),
	},

	// Test map structure tags.
	"map-name": {
		value: &struct {
			M any `map:"X" tag:"any"`
		}{},
		expect: Call("x", "any"),
	},
	"map-squash": {
		value: &struct {
			S struct {
				A *any `map:"X" tag:"*any"`
			} `map:",squash" tag:"struct{*any}"`
		}{},
		expect: mock.Chain(
			Call("x", "*any"),
		),
	},
	"map-empty": {
		value: &struct {
			S struct {
				A *any `map:",omitempty" tag:"*any"`
			} `map:",squash" tag:"struct{*any}"`
		}{},
		expect: mock.Chain(
			Call("a", "*any"),
		),
	},
	"map-remain": {
		value: &struct {
			S struct {
				A *any           `map:",omitempty" tag:"*any"`
				R map[string]any `map:",remain" tag:"map[string]any"`
			} `map:",squash" tag:"struct{*any}"`
		}{},
		expect: mock.Chain(
			Call("a", "*any"),
			Call("r", "map[string]any"),
		),
	},
	"map-comma": {
		value: &struct {
			S struct {
				A *any `map:"," tag:"*any"`
			} `map:",squash" tag:"struct{*any}"`
		}{},
		expect: mock.Chain(
			Call("a", "*any"),
		),
	},
}

// TestTagWalker_Walk tests TagWalker.Walk.
func TestTagWalker_Walk(t *testing.T) {
	test.Map(t, tagWalkerTestCases).
		Run(func(t test.Test, param tagWalkerParam) {
			// Given
			mocks := mock.NewMocks(t).Expect(param.expect)
			walker := reflect.NewTagWalker("tag", "map", param.zero)

			// When
			walker.Walk(param.key, param.value,
				mock.Get(mocks, NewMockCallback).Call)

			// Then
		})
}
