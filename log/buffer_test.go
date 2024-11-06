package log_test

import (
	"bytes"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/tkrop/go-config/log"
	"github.com/tkrop/go-testing/mock"
	"github.com/tkrop/go-testing/test"
)

//revive:disable:line-length-limit // go:generate line length

//go:generate mockgen -package=log_test -destination=mock_writer_test.go -source=buffer.go BufferWriter

//revive:enable:line-length-limit

// setupWriter sets up the writer for testing.
func setupWriter(
	mocks *mock.Mocks, expect mock.SetupFunc,
) log.BufferWriter {
	if expect != nil {
		return mock.Get(mocks, NewMockBufferWriter)
	}
	return &bytes.Buffer{}
}

type testBufferWriteParam struct {
	colorMode    log.ColorModeString
	error        error
	setup        func(*log.Buffer)
	expect       mock.SetupFunc
	expectError  error
	expectString string
}

var testBufferWriteParams = map[string]testBufferWriteParam{
	// Test write byte.
	"write byte error": {
		error: errAny,
		setup: func(buffer *log.Buffer) {
			buffer.WriteByte(' ')
		},
		expectError: errAny,
	},
	"write byte failure": {
		setup: func(buffer *log.Buffer) {
			buffer.WriteByte(' ')
		},
		expect: mock.Chain(func(mocks *mock.Mocks) any {
			return mock.Get(mocks, NewMockBufferWriter).EXPECT().WriteByte(uint8(' ')).
				DoAndReturn(mocks.Do(log.BufferWriter.WriteByte, errAny))
		}, func(mocks *mock.Mocks) any {
			return mock.Get(mocks, NewMockBufferWriter).EXPECT().Bytes().
				DoAndReturn(mocks.Do(log.BufferWriter.Bytes, []byte("")))
		}),
		expectError: errAny,
	},
	"write byte": {
		setup: func(buffer *log.Buffer) {
			buffer.WriteByte(' ')
		},
		expectString: " ",
	},

	// Test write string.
	"write string error": {
		error: errAny,
		setup: func(buffer *log.Buffer) {
			buffer.WriteString("string")
		},
		expectError: errAny,
	},
	"write string failure": {
		setup: func(buffer *log.Buffer) {
			buffer.WriteString("string")
		},
		expect: mock.Chain(func(mocks *mock.Mocks) any {
			return mock.Get(mocks, NewMockBufferWriter).EXPECT().WriteString("string").
				DoAndReturn(mocks.Do(log.BufferWriter.WriteString, 0, errAny))
		}, func(mocks *mock.Mocks) any {
			return mock.Get(mocks, NewMockBufferWriter).EXPECT().Bytes().
				DoAndReturn(mocks.Do(log.BufferWriter.Bytes, []byte("")))
		}),
		expectError: errAny,
	},
	"write string": {
		setup: func(buffer *log.Buffer) {
			buffer.WriteString("string")
		},
		expectString: "string",
	},

	// Test write colored.
	"write colored error": {
		error: errAny,
		setup: func(buffer *log.Buffer) {
			buffer.WriteColored(log.ColorField, "string")
		},
		expectError: errAny,
	},
	"write colored default": {
		setup: func(buffer *log.Buffer) {
			buffer.WriteColored(log.ColorField, "string")
		},
		expectString: fieldC("string"),
	},
	"write colored color-off": {
		colorMode: log.ColorModeOff,
		setup: func(buffer *log.Buffer) {
			buffer.WriteColored(log.ColorField, "string")
		},
		expectString: field("string"),
	},
	"write colored color-on": {
		colorMode: log.ColorModeOn,
		setup: func(buffer *log.Buffer) {
			buffer.WriteColored(log.ColorField, "string")
		},
		expectString: fieldC("string"),
	},

	// Test write level.
	"write level error": {
		error: errAny,
		setup: func(buffer *log.Buffer) {
			buffer.WriteLevel(log.PanicLevel)
		},
		expectError: errAny,
	},
	"write level default": {
		setup: func(buffer *log.Buffer) {
			buffer.WriteLevel(log.PanicLevel)
		},
		expectString: levelC(log.PanicLevel),
	},
	"write level color-on": {
		colorMode: log.ColorModeOn,
		setup: func(buffer *log.Buffer) {
			buffer.WriteLevel(log.PanicLevel)
		},
		expectString: levelC(log.PanicLevel),
	},
	"write level color-off": {
		colorMode: log.ColorModeOff,
		setup: func(buffer *log.Buffer) {
			buffer.WriteLevel(log.PanicLevel)
		},
		expectString: level(log.PanicLevel),
	},

	// Test write colored field.
	"write field error": {
		error: errAny,
		setup: func(buffer *log.Buffer) {
			buffer.WriteField(log.FieldLevel, "value")
		},
		expectError: errAny,
	},
	"write field default": {
		setup: func(buffer *log.Buffer) {
			buffer.WriteField(log.FieldLevel, "value")
		},
		expectString: fieldC("value"),
	},
	"write field color-on": {
		colorMode: log.ColorModeOn,
		setup: func(buffer *log.Buffer) {
			buffer.WriteField(log.FieldLevel, "value")
		},
		expectString: fieldC("value"),
	},
	"write field color-off": {
		colorMode: log.ColorModeOff,
		setup: func(buffer *log.Buffer) {
			buffer.WriteField(log.FieldLevel, "value")
		},
		expectString: field("value"),
	},

	// Test write caller.
	"write caller error": {
		error: errAny,
		setup: func(buffer *log.Buffer) {
			buffer.WriteCaller(anyFrame)
		},
		expectError: errAny,
	},
	"write caller on": {
		setup: func(buffer *log.Buffer) {
			buffer.WriteCaller(anyFrame)
		},
		expectString: " [file:123#function]",
	},
	"write caller off": {
		setup: func(buffer *log.Buffer) {
			buffer.WriteCaller(nil)
		},
		expectString: "",
	},

	// Test write value.
	"write value error": {
		error: errAny,
		setup: func(buffer *log.Buffer) {
			buffer.WriteValue("value")
		},
		expectError: errAny,
	},
	"write value string": {
		setup: func(buffer *log.Buffer) {
			buffer.WriteValue("value")
		},
		expectString: "\"value\"",
	},
	"write value int": {
		setup: func(buffer *log.Buffer) {
			buffer.WriteValue(123)
		},
		expectString: "123",
	},
	"write value float": {
		setup: func(buffer *log.Buffer) {
			buffer.WriteValue(123.456)
		},
		expectString: "123.456",
	},
	"write value complex": {
		setup: func(buffer *log.Buffer) {
			buffer.WriteValue(123.456 + 789i)
		},
		expectString: "(123.456+789i)",
	},
	"write value bool": {
		setup: func(buffer *log.Buffer) {
			buffer.WriteValue(true)
		},
		expectString: "true",
	},

	// Test write data.
	"write data error": {
		error: errAny,
		setup: func(buffer *log.Buffer) {
			buffer.WriteData("key", "value")
		},
		expectError: errAny,
	},
	"write data default": {
		setup: func(buffer *log.Buffer) {
			buffer.WriteData("key", "value")
		},
		expectString: dataC("key", "value"),
	},
	"write data color-on error": {
		setup: func(buffer *log.Buffer) {
			buffer.WriteData(logrus.ErrorKey, errAny)
		},
		expectString: dataC(logrus.ErrorKey, errAny.Error()),
	},
	"write data color-on": {
		colorMode: log.ColorModeOn,
		setup: func(buffer *log.Buffer) {
			buffer.WriteData("key", "value")
		},
		expectString: dataC("key", "value"),
	},
	"write data color-off": {
		colorMode: log.ColorModeOff,
		setup: func(buffer *log.Buffer) {
			buffer.WriteData("key", "value")
		},
		expectString: data("key", "value"),
	},
}

func TestBufferWrite(t *testing.T) {
	test.Map(t, testBufferWriteParams).
		Run(func(t test.Test, param testBufferWriteParam) {
			// Given
			mocks := mock.NewMocks(t).Expect(param.expect)
			pretty := &log.Setup{
				ColorMode:   param.colorMode.Parse(true),
				ErrorName:   log.DefaultErrorName,
				LevelNames:  log.DefaultLevelNames,
				LevelColors: log.DefaultLevelColors,
			}

			buffer := log.NewBuffer(pretty,
				setupWriter(mocks, param.expect))
			test.NewAccessor(buffer).Set("err", param.error)

			// When
			param.setup(buffer)
			result, err := buffer.Bytes()

			// Then
			assert.Equal(t, param.expectError, err)
			assert.Equal(t, param.expectString, string(result))
			if param.expect == nil {
				assert.Equal(t, param.expectString, buffer.String())
			}
		})
}
