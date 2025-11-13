package reflect_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tkrop/go-config/internal/reflect"
	"github.com/tkrop/go-testing/mock"
	"github.com/tkrop/go-testing/test"
	"gopkg.in/yaml.v3"
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

// TagWalkerParams contains a value and the expected tags.
type TagWalkerParams struct {
	value  any
	key    string
	zero   bool
	expect mock.SetupFunc
	error  error
}

//revive:disable:nested-structs // simplifies test cases a lot.

// tagWalkerTestCases contains test cases for TagWalker.Walk.
var tagWalkerTestCases = map[string]TagWalkerParams{
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
	"slice-bool-values": {
		value: []bool{true, false},
		zero:  true,
		expect: mock.Chain(
			Call("0", true),
			Call("1", false),
		),
	},
	"slice-int-values": {
		value: []int{1, 0},
		zero:  true,
		expect: mock.Chain(
			Call("0", 1),
			Call("1", 0),
		),
	},
	"slice-uint-values": {
		value: []uint{1, 0},
		zero:  true,
		expect: mock.Chain(
			Call("0", uint(1)),
			Call("1", uint(0)),
		),
	},
	"slice-float-values": {
		value: []float64{1.0, 0.0},
		zero:  true,
		expect: mock.Chain(
			Call("0", 1.0),
			Call("1", 0.0),
		),
	},
	"slice-complex-values": {
		value: []complex128{1 + 2i, 3 + 4i},
		zero:  true,
		expect: mock.Chain(
			Call("0", complex128(1+2i)),
			Call("1", complex128(3+4i)),
		),
	},
	"slice-string-values": {
		zero:  true,
		value: []string{"test", ""},
		expect: mock.Chain(
			Call("0", "test"),
			Call("1", ""),
		),
	},
	"slice-byte-values": {
		value: []byte{'a', 'b'},
		expect: mock.Chain(
			Call("0", byte('a')),
			Call("1", byte('b')),
		),
	},
	"slice-rune-values": {
		value: []rune{'a', 'b'},
		expect: mock.Chain(
			Call("0", rune('a')),
			Call("1", rune('b')),
		),
	},
	"slice-any-values": {
		value: []any{0, "test"},
		expect: mock.Chain(
			Call("0", 0),
			Call("1", "test"),
		),
	},
	"slice-nil-ptr-values": {
		value: []*struct {
			A any `tag:"any"`
		}{nil},
		expect: Call("0.a", "any"),
	},

	// Test struct with field tags.
	"struct-bool-tags": {
		value: struct {
			hidden  bool
			Visible bool `tag:"true"`
		}{},
		expect: Call("visible", true),
	},

	"struct-ints-tags": {
		value: struct {
			I   int    `tag:"1"`
			PI  *int   `tag:"2"`
			SI  []int  `tag:"[1,2,3]"`
			PSI *[]int `tag:"[1,2,3]"`
			I8  int8   `tag:"8"`
			I16 int16  `tag:"16"`
			I32 int32  `tag:"32"`
			I64 int64  `tag:"64"`
		}{},
		expect: mock.Chain(
			Call("i", int(1)),
			Call("pi", test.Ptr(int(2))),
			Call("si", []int{1, 2, 3}),
			Call("psi", test.Ptr([]int{1, 2, 3})),
			Call("i8", int8(8)),
			Call("i16", int16(16)),
			Call("i32", int32(32)),
			Call("i64", int64(64)),
		),
	},

	"struct-uint-tags": {
		value: struct {
			UI   uint    `tag:"1"`
			PUI  *uint   `tag:"2"`
			SUI  []uint  `tag:"[1,2,3]"`
			PSUI *[]uint `tag:"[1,2,3]"`
			UI8  uint8   `tag:"8"`
			UI16 uint16  `tag:"16"`
			UI32 uint32  `tag:"32"`
			UI64 uint64  `tag:"64"`
		}{},
		expect: mock.Chain(
			Call("ui", uint(1)),
			Call("pui", test.Ptr(uint(2))),
			Call("sui", []uint{1, 2, 3}),
			Call("psui", test.Ptr([]uint{1, 2, 3})),
			Call("ui8", uint8(8)),
			Call("ui16", uint16(16)),
			Call("ui32", uint32(32)),
			Call("ui64", uint64(64)),
		),
	},

	"struct-float-tags": {
		value: struct {
			F32   float32    `tag:"32e-1"`
			F64   float64    `tag:"64e-1"`
			PF32  *float32   `tag:"32e-1"`
			PF64  *float64   `tag:"64e-1"`
			SF32  []float32  `tag:"[32e-1]"`
			SF64  []float64  `tag:"[64e-1]"`
			PSF32 *[]float32 `tag:"[32e-1]"`
			PSF64 *[]float64 `tag:"[64e-1]"`
		}{},
		expect: mock.Chain(
			Call("f32", float32(32e-1)),
			Call("f64", float64(64e-1)),
			Call("pf32", test.Ptr(float32(32e-1))),
			Call("pf64", test.Ptr(float64(64e-1))),
			Call("sf32", []float32{32e-1}),
			Call("sf64", []float64{64e-1}),
			Call("psf32", test.Ptr([]float32{32e-1})),
			Call("psf64", test.Ptr([]float64{64e-1})),
		),
	},

	"struct-complex-tags": {
		value: struct {
			C64    complex64     `tag:"64+2i"`
			C128   complex128    `tag:"128+2i"`
			PC64   *complex64    `tag:"64+2i"`
			PC128  *complex128   `tag:"128+2i"`
			SC64   []complex64   `tag:"[64+2i, 32+1i]"`
			SC128  []complex128  `tag:"[128+2i, 64+1i]"`
			PSC64  *[]complex64  `tag:"[64+2i, 32+1i]"`
			PSC128 *[]complex128 `tag:"[128+2i, 64+1i]"`
		}{},
		expect: mock.Chain(
			Call("c64", complex64(64+2i)),
			Call("c128", complex128(128+2i)),
			Call("pc64", test.Ptr(complex64(64+2i))),
			Call("pc128", test.Ptr(complex128(128+2i))),
			Call("sc64", []complex64{64 + 2i, 32 + 1i}),
			Call("sc128", []complex128{128 + 2i, 64 + 1i}),
			Call("psc64", test.Ptr([]complex64{64 + 2i, 32 + 1i})),
			Call("psc128", test.Ptr([]complex128{128 + 2i, 64 + 1i})),
		),
	},

	"struct-string-tags": {
		value: struct {
			S   string  `tag:"string"`
			PS  *string `tag:"string"`
			B   byte    `tag:"117"`
			PB  *byte   `tag:"42"`
			SB  []byte  `tag:"[117,105,110,116,56]"`
			PSB *[]byte `tag:"[42,117,105,110,116,56]"`
			R   rune    `tag:"105"`
			PR  *rune   `tag:"42"`
			SR  []rune  `tag:"[105,110,116,51,50]"`
			PSR *[]rune `tag:"[42,105,110,116,51,50]"`
		}{},
		expect: mock.Chain(
			Call("s", "string"),
			Call("ps", test.Ptr("string")),
			Call("b", byte('u')),
			Call("pb", test.Ptr(byte('*'))),
			Call("sb", []byte("uint8")),
			Call("psb", test.Ptr([]byte("*uint8"))),
			Call("r", rune('i')),
			Call("pr", test.Ptr(rune('*'))),
			Call("sr", []rune("int32")),
			Call("psr", test.Ptr([]rune("*int32"))),
		),
	},

	// Test structs with field values.
	"struct-all-values": {
		value: struct {
			Bool   bool    `map:"bool" tag:"false"`
			Int    int     `map:"int" tag:"-2"`
			Uint   uint    `map:"uint" tag:"2"`
			Float  float64 `map:"float" tag:"3.0"`
			String string  `map:"string" tag:"STRING"`
			Byte   byte    `map:"byte" tag:"A"`
			Rune   rune    `map:"rune" tag:"B"`
			Any    any     `map:"any" tag:"ANY"`
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
			}
		}{},
		expect: Call("s.a", "any"),
	},
	"struct-ptr-struct": {
		value: struct {
			S *struct {
				A any `tag:"any"`
			}
		}{},
		expect: Call("s.a", "any"),
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
			}
		}{},
		expect: Call("s.a", "any"),
	},
	"ptr-struct-ptr-struct-tags": {
		value: &struct {
			S *struct {
				A any `tag:"any"`
			} `tag:"{a: any}"`
		}{},
		expect: mock.Chain(
			Call("s", &struct {
				A any `tag:"any"`
			}{A: "any"}),
		),
	},

	// Test struct with nested slices and tags.
	"struct-slice-tags": {
		value: struct {
			S []any `tag:"[any,all]"`
		}{},
		expect: Call("s", []any{"any", "all"}),
	},
	"struct-slice-struct-tags": {
		value: struct {
			S []struct {
				A any `tag:"any"`
			} `tag:"[{a: any},{a: all}]"`
		}{},
		expect: Call("s", []struct {
			A any `tag:"any"`
		}{{A: "any"}, {A: "all"}}),
	},
	"struct-slice-ptr-struct-tags": {
		value: struct {
			S []*struct {
				A any `tag:"any"`
			} `tag:"[{a: any},{a: all}]"`
		}{},
		expect: Call("s", []*struct {
			A any `tag:"any"`
		}{{A: "any"}, {A: "all"}}),
	},
	"struct-ptr-slice-ptr-struct-tags": {
		value: struct {
			S *[]*struct {
				A any `tag:"any"`
			} `tag:"[{a: any},{a: all}]"`
		}{},
		expect: Call("s", &[]*struct {
			A any `tag:"any"`
		}{
			test.Ptr(struct {
				A any `tag:"any"`
			}{A: "any"}),
			test.Ptr(struct {
				A any `tag:"any"`
			}{A: "all"}),
		}),
	},

	// Test struct with nested slices and values.
	"struct-slice-values": {
		value: struct {
			S []any
		}{S: []any{1, 2}},
		expect: mock.Chain(
			Call("s.0", 1),
			Call("s.1", 2),
		),
	},
	"struct-slice-struct-values": {
		value: struct {
			S []struct {
				A any `tag:"any"`
			}
		}{S: []struct {
			A any `tag:"any"`
		}{{A: 1}, {A: 2}}},
		expect: mock.Chain(
			Call("s.0.a", 1),
			Call("s.1.a", 2),
		),
	},
	"struct-slice-ptr-struct-values": {
		value: struct {
			S []*struct {
				A any `tag:"any"`
			}
		}{S: []*struct {
			A any `tag:"any"`
		}{{A: 1}, {A: 2}}},
		expect: mock.Chain(
			Call("s.0.a", 1),
			Call("s.1.a", 2),
		),
	},
	"struct-ptr-slice-ptr-struct-values": {
		value: struct {
			S *[]*struct {
				A any `tag:"any"`
			}
		}{S: &[]*struct {
			A any `tag:"any"`
		}{{A: 1}, {A: 2}}},
		expect: mock.Chain(
			Call("s.0.a", 1),
			Call("s.1.a", 2),
		),
	},

	// Test struct with nested maps.
	"struct-map-tags": {
		value: struct {
			M map[string]any `tag:"{a: any, b: {a: all}}"`
		}{},
		expect: Call("m", map[string]any{
			"a": "any",
			"b": map[string]any{"a": "all"},
		}),
	},
	"struct-ptr-map-tags": {
		value: struct {
			M *map[string]any `tag:"{a: any, b: {a: all}}"`
		}{},
		expect: Call("m", &map[string]any{
			"a": "any",
			"b": map[string]any{"a": "all"},
		}),
	},
	"struct-map-struct-tags": {
		value: &struct {
			M map[string]struct {
				A any `tag:"any"`
			} `tag:"{a: {a: any},b: {a: all}}"`
		}{},
		expect: Call("m", map[string]struct {
			A any `tag:"any"`
		}{
			"a": {A: "any"},
			"b": {A: "all"},
		}),
	},
	"struct-ptr-map-struct-tags": {
		value: &struct {
			M map[string]struct {
				A any `tag:"any"`
			} `tag:"{a: {a: any},b: {a: all} }"`
		}{},
		expect: Call("m", map[string]struct {
			A any `tag:"any"`
		}{
			"a": {A: "any"},
			"b": {A: "all"},
		}),
	},
	"struct-map-ptr-struct-tags": {
		value: &struct {
			M map[string]*struct {
				A any `tag:"any"`
			} `tag:"{a: {a: any},b: {a: all}}"`
		}{},
		expect: Call("m", map[string]*struct {
			A any `tag:"any"`
		}{
			"a": {A: "any"},
			"b": {A: "all"},
		}),
	},

	// Test struct with nested maps.
	"struct-map-values": {
		value: struct {
			M map[string]any
		}{M: map[string]any{"key": "value"}},
		expect: Call("m.key", "value"),
	},
	"struct-ptr-map-values": {
		value: struct {
			M *map[string]any
		}{M: &map[string]any{"key": "value"}},
		expect: Call("m.key", "value"),
	},
	"struct-map-struct-values": {
		value: struct {
			M map[string]struct {
				A any `tag:"any"`
			}
		}{M: map[string]struct {
			A any `tag:"any"`
		}{"key-0": {A: 1}, "key-1": {A: 2}}},
		expect: mock.Setup(
			Call("m.key-0.a", 1),
			Call("m.key-1.a", 2),
		),
	},
	"struct-ptr-map-struct-values": {
		value: struct {
			M *map[string]struct {
				A any `tag:"any"`
			}
		}{M: &map[string]struct {
			A any `tag:"any"`
		}{"key-0": {A: 1}, "key-1": {A: 2}}},
		expect: mock.Setup(
			Call("m.key-0.a", 1),
			Call("m.key-1.a", 2),
		),
	},
	"struct-ptr-map-ptr-struct-values": {
		value: struct {
			M *map[string]*struct {
				A any `tag:"any"`
			}
		}{M: &map[string]*struct {
			A any `tag:"any"`
		}{"key-0": {A: 1}, "key-1": {A: 2}}},
		expect: mock.Setup(
			Call("m.key-0.a", 1),
			Call("m.key-1.a", 2),
		),
	},

	// Test with special tags.
	"tag-yaml-no-tags": {
		value: &struct {
			S string
		}{},
	},
	"tag-yaml-zero-no-tags": {
		value: &struct {
			S string
		}{},
		zero:   true,
		expect: Call("s", ""),
	},
	"tag-yaml-empty": {
		value: &struct {
			S string `tag:""`
		}{},
	},
	"tag-yaml-zero-empty": {
		value: &struct {
			S string `tag:""`
		}{},
		zero:   true,
		expect: Call("s", ""),
	},
	"tag-map-squash": {
		value: struct {
			S struct {
				A any `tag:"any"`
			} `map:",squash"`
		}{},
		expect: Call("a", "any"),
	},
	"tag-map-remain": {
		value: struct {
			Field any `map:",remain" tag:"any"`
		}{},
		expect: Call("field", "any"),
	},
	"tag-yaml-slice": {
		value: &struct {
			S []string `tag:"[a,b]"`
		}{},
		expect: Call("s", []string{"a", "b"}),
	},
	"tag-yaml-map": {
		value: &struct {
			S struct {
				A string `tag:"v"`
				B []string
				C string
			} `tag:"{a: a, b: [a,b], c: c}"`
		}{},
		expect: mock.Chain(
			Call("s", struct {
				A string `tag:"v"`
				B []string
				C string
			}{A: "a", B: []string{"a", "b"}, C: "c"}),
		),
	},
	"tag-yaml-error": {
		value: &struct {
			S []string `tag:"a,b"`
		}{},
		expect: mock.Setup(
			Call("s", "a,b"),
		),
		error: fmt.Errorf("%w - %s [%s=%s]: %w",
			reflect.ErrTagWalker, "yaml parsing", "s", "\"a,b\"",
			&yaml.TypeError{Errors: []string{
				"line 1: cannot unmarshal !!str `a,b` into []string",
			}}),
	},

	// Complex number parsing errors
	"tag-yaml-complex-invalid": {
		value: &struct {
			C complex64 `tag:"invalid"`
		}{},
		expect: mock.Setup(
			Call("c", "invalid"),
		),
		error: fmt.Errorf("%w - %s [%s=%s]: %w",
			reflect.ErrTagWalker, "complex parsing", "c", "\"invalid\"",
			&strconv.NumError{
				Func: "ParseComplex",
				Num:  "invalid",
				Err:  strconv.ErrSyntax,
			}),
	},
	"tag-yaml-complex-slice-invalid": {
		value: &struct {
			SC []complex64 `tag:"[invalid, 1+2i]"`
		}{},
		expect: mock.Setup(
			Call("sc", []string{"invalid", "1+2i"}),
		),
		error: fmt.Errorf("%w - %s [%s=%#v]: %w",
			reflect.ErrTagWalker, "complex parsing", "sc",
			[]string{"invalid", "1+2i"},
			&strconv.NumError{
				Func: "ParseComplex",
				Num:  "invalid",
				Err:  strconv.ErrSyntax,
			}),
	},
	"tag-yaml-complex-ptr-slice-invalid": {
		value: &struct {
			PSC *[]complex64 `tag:"[invalid, 1+2i]"`
		}{},
		expect: mock.Setup(
			Call("psc", []string{"invalid", "1+2i"}),
		),
		error: fmt.Errorf("%w - %s [%s=%#v]: %w",
			reflect.ErrTagWalker, "complex parsing", "psc",
			[]string{"invalid", "1+2i"},
			&strconv.NumError{
				Func: "ParseComplex",
				Num:  "invalid",
				Err:  strconv.ErrSyntax,
			}),
	},
}

func TestTagWalker(t *testing.T) {
	test.Map(t, tagWalkerTestCases).
		// Filter(test.Pattern[TagWalkerParams]("^tag-yaml-zero")).
		Run(func(t test.Test, param TagWalkerParams) {
			// Given
			mocks := mock.NewMocks(t).Expect(param.expect)
			walker := reflect.NewTagWalker("tag", "map", param.zero,
				mock.Get(mocks, NewMockCallback).Call)

			// When
			err := walker.Walk(param.key, param.value)

			// Then
			assert.Equal(t, param.error, err)
		})
}
