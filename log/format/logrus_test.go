package format_test

import (
	"bytes"
	"errors"
	"runtime"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/tkrop/go-config/log/format"
	"github.com/tkrop/go-testing/mock"
	"github.com/tkrop/go-testing/test"
)

//revive:disable:line-length-limit // go:generate line length

//go:generate mockgen -package=format_test -destination=mock_writer_test.go -source=logrus.go BufferWriter

//revive:enable:line-length-limit

var (
	// TestTime is a fixed time for testing.
	otime = "2024-10-01 23:07:13.891012345Z"
	// TestTime is a fixed time for testing.
	itime = "2024-10-01T23:07:13.891012345Z"

	// Arbitrary data for testing.
	anyData = log.Fields{
		"key1": "value1",
		"key2": "value2",
	}
	// Arbitrary frame for testing.
	anyFrame = &runtime.Frame{
		File:     "file",
		Function: "function",
		Line:     123,
	}
	// Arbitrary error for testing.
	errAny = errors.New("any error")
)

// setupTimeFormat sets up the time format for testing.
func setupTimeFormat(timeFormat string) string {
	if timeFormat == "" {
		return format.DefaultTimeFormat
	}
	return timeFormat
}

// setupWriter sets up the writer for testing.
func setupWriter(
	mocks *mock.Mocks, expect mock.SetupFunc,
) format.BufferWriter {
	if expect != nil {
		return mock.Get(mocks, NewMockBufferWriter)
	}
	return &bytes.Buffer{}
}

// Helper functions for testing log levels without color.
func level(level log.Level) string {
	return format.DefaultLevelNames[level]
}

// Helper functions for testing log levels with color.
func levelC(level log.Level) string {
	return "\x1b[" + format.DefaultLevelColors[level] +
		"m" + format.DefaultLevelNames[level] + "\x1b[0m"
}

// Helper functions for testing fields without color.
func fieldC(value string) string {
	return "\x1b[" + format.ColorField + "m" + value + "\x1b[0m"
}

// Helper functions for testing key-value data without color.
func data(key, value string) string {
	return key + "=\"" + value + "\""
}

// Helper functions for testing key-value data with color.
func dataC(key, value string) string {
	color := format.ColorField
	if key == log.ErrorKey {
		color = format.ColorError
	}
	return "\x1b[" + color + "m" + key + "\x1b[0m=\"" + value + "\""
}

type testPrettyFormatParam struct {
	timeFormat   string
	noTerminal   bool
	colorMode    format.ColorModeString
	orderMode    format.OrderModeString
	entry        *log.Entry
	expect       func(t test.Test, result string, err error)
	expectError  error
	expectResult string
}

var testPrettyFormatParams = map[string]testPrettyFormatParam{
	// Test levels with default.
	"level panic default": {
		entry: &log.Entry{
			Level:   log.PanicLevel,
			Message: "panic message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.PanicLevel) + " panic message\n",
	},
	"level fatal default": {
		entry: &log.Entry{
			Level:   log.FatalLevel,
			Message: "fatal message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.FatalLevel) + " fatal message\n",
	},
	"level error default": {
		entry: &log.Entry{
			Level:   log.ErrorLevel,
			Message: "error message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.ErrorLevel) + " error message\n",
	},
	"level warn default": {
		entry: &log.Entry{
			Level:   log.WarnLevel,
			Message: "warn message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.WarnLevel) + " warn message\n",
	},
	"level info default": {
		entry: &log.Entry{
			Level:   log.InfoLevel,
			Message: "info message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.InfoLevel) + " info message\n",
	},
	"level debug default": {
		entry: &log.Entry{
			Level:   log.DebugLevel,
			Message: "debug message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.DebugLevel) + " debug message\n",
	},
	"level trace default": {
		entry: &log.Entry{
			Level:   log.TraceLevel,
			Message: "trace message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.TraceLevel) + " trace message\n",
	},

	// Test levels with color.
	"level panic color-on": {
		colorMode: format.ColorModeOn,
		entry: &log.Entry{
			Level:   log.PanicLevel,
			Message: "panic message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.PanicLevel) + " panic message\n",
	},
	"level fatal color-on": {
		colorMode: format.ColorModeOn,
		entry: &log.Entry{
			Level:   log.FatalLevel,
			Message: "fatal message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.FatalLevel) + " fatal message\n",
	},
	"level error color-on": {
		colorMode: format.ColorModeOn,
		entry: &log.Entry{
			Level:   log.ErrorLevel,
			Message: "error message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.ErrorLevel) + " error message\n",
	},
	"level warn color-on": {
		colorMode: format.ColorModeOn,
		entry: &log.Entry{
			Level:   log.WarnLevel,
			Message: "warn message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.WarnLevel) + " warn message\n",
	},
	"level info color-on": {
		colorMode: format.ColorModeOn,
		entry: &log.Entry{
			Level:   log.InfoLevel,
			Message: "info message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.InfoLevel) + " info message\n",
	},
	"level debug color-on": {
		colorMode: format.ColorModeOn,
		entry: &log.Entry{
			Level:   log.DebugLevel,
			Message: "debug message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.DebugLevel) + " debug message\n",
	},
	"level trace color-on": {
		colorMode: format.ColorModeOn,
		entry: &log.Entry{
			Level:   log.TraceLevel,
			Message: "trace message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.TraceLevel) + " trace message\n",
	},

	// Test levels with color.
	"level panic color-off": {
		colorMode: format.ColorModeOff,
		entry: &log.Entry{
			Level:   log.PanicLevel,
			Message: "panic message",
		},
		expectResult: otime[0:26] + " " +
			level(log.PanicLevel) + " panic message\n",
	},
	"level fatal color-off": {
		colorMode: format.ColorModeOff,
		entry: &log.Entry{
			Level:   log.FatalLevel,
			Message: "fatal message",
		},
		expectResult: otime[0:26] + " " +
			level(log.FatalLevel) + " fatal message\n",
	},
	"level error color-off": {
		colorMode: format.ColorModeOff,
		entry: &log.Entry{
			Level:   log.ErrorLevel,
			Message: "error message",
		},
		expectResult: otime[0:26] + " " +
			level(log.ErrorLevel) + " error message\n",
	},
	"level warn color-off": {
		colorMode: format.ColorModeOff,
		entry: &log.Entry{
			Level:   log.WarnLevel,
			Message: "warn message",
		},
		expectResult: otime[0:26] + " " +
			level(log.WarnLevel) + " warn message\n",
	},
	"level info color-off": {
		colorMode: format.ColorModeOff,
		entry: &log.Entry{
			Level:   log.InfoLevel,
			Message: "info message",
		},
		expectResult: otime[0:26] + " " +
			level(log.InfoLevel) + " info message\n",
	},
	"level debug color-off": {
		colorMode: format.ColorModeOff,
		entry: &log.Entry{
			Level:   log.DebugLevel,
			Message: "debug message",
		},
		expectResult: otime[0:26] + " " +
			level(log.DebugLevel) + " debug message\n",
	},
	"level trace color-off": {
		colorMode: format.ColorModeOff,
		entry: &log.Entry{
			Level:   log.TraceLevel,
			Message: "trace message",
		},
		expectResult: otime[0:26] + " " +
			level(log.TraceLevel) + " trace message\n",
	},

	// Test order key value data.
	"data default": {
		entry: &log.Entry{
			Message: "data message",
			Data:    anyData,
		},
		expectResult: otime[0:26] + " " +
			levelC(log.PanicLevel) + " data message " +
			dataC("key1", "value1") + " " +
			dataC("key2", "value2") + "\n",
	},
	"data ordered": {
		orderMode: format.OrderModeOn,
		entry: &log.Entry{
			Message: "data message",
			Data:    anyData,
		},
		expectResult: otime[0:26] + " " +
			levelC(log.PanicLevel) + " data message " +
			dataC("key1", "value1") + " " +
			dataC("key2", "value2") + "\n",
	},
	"data unordered": {
		orderMode: format.OrderModeOff,
		entry: &log.Entry{
			Message: "data message",
			Data:    anyData,
		},
		expect: func(t test.Test, result string, err error) {
			assert.Contains(t, result, otime[0:26]+" "+
				levelC(log.PanicLevel)+" "+"data message")
			assert.Contains(t, result, dataC("key1", "value1"))
			assert.Contains(t, result, dataC("key2", "value2"))
		},
	},

	// Test color modes.
	"data color-off": {
		colorMode: format.ColorModeOff,
		entry: &log.Entry{
			Message: "data message",
			Data:    anyData,
		},
		expectResult: otime[0:26] + " " +
			level(log.PanicLevel) + " data message " +
			data("key1", "value1") + " " +
			data("key2", "value2") + "\n",
	},
	"data color-on": {
		colorMode: format.ColorModeOn,
		entry: &log.Entry{
			Message: "data message",
			Data:    anyData,
		},
		expectResult: otime[0:26] + " " +
			levelC(log.PanicLevel) + " data message " +
			dataC("key1", "value1") + " " +
			dataC("key2", "value2") + "\n",
	},
	"data color-auto colorized": {
		colorMode: format.ColorModeAuto,
		entry: &log.Entry{
			Message: "data message",
			Data:    anyData,
		},
		expectResult: otime[0:26] + " " +
			levelC(log.PanicLevel) + " data message " +
			dataC("key1", "value1") + " " +
			dataC("key2", "value2") + "\n",
	},
	"data color-auto not-colorized": {
		colorMode:  format.ColorModeAuto,
		noTerminal: true,
		entry: &log.Entry{
			Message: "data message",
			Data:    anyData,
			Logger:  &log.Logger{Out: nil},
		},
		expectResult: otime[0:26] + " " +
			level(log.PanicLevel) + " data message " +
			data("key1", "value1") + " " +
			data("key2", "value2") + "\n",
	},
	"data color-levels": {
		colorMode: format.ColorModeLevels,
		entry: &log.Entry{
			Message: "data message",
			Data:    anyData,
		},
		expectResult: otime[0:26] + " " +
			levelC(log.PanicLevel) + " data message " +
			data("key1", "value1") + " " +
			data("key2", "value2") + "\n",
	},
	"data color-fields": {
		colorMode: format.ColorModeFields,
		entry: &log.Entry{
			Message: "data message",
			Data:    anyData,
		},
		expectResult: otime[0:26] + " " +
			level(log.PanicLevel) + " data message " +
			dataC("key1", "value1") + " " +
			dataC("key2", "value2") + "\n",
	},
	"data color-levels+fields": {
		colorMode: format.ColorModeLevels + "|" + format.ColorModeFields,
		entry: &log.Entry{
			Message: "data message",
			Data:    anyData,
		},
		expectResult: otime[0:26] + " " +
			levelC(log.PanicLevel) + " data message " +
			dataC("key1", "value1") + " " +
			dataC("key2", "value2") + "\n",
	},

	// Time format.
	"time default": {
		entry: &log.Entry{
			Level:   log.PanicLevel,
			Message: "default time message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.PanicLevel) + " " +
			"default time message\n",
	},
	"time short": {
		timeFormat: "2006-01-02 15:04:05",
		entry: &log.Entry{
			Level:   log.PanicLevel,
			Message: "short time message",
		},
		expectResult: otime[0:19] + " " +
			levelC(log.PanicLevel) + " " +
			"short time message\n",
	},
	"time long": {
		timeFormat: "2006-01-02 15:04:05.000000000",
		entry: &log.Entry{
			Level:   log.PanicLevel,
			Message: "long time message",
		},
		expectResult: otime[0:29] + " " +
			levelC(log.PanicLevel) + " " +
			"long time message\n",
	},

	// Report caller.
	"caller only": {
		entry: &log.Entry{
			Message: "caller message",
			Caller:  anyFrame,
		},
		expectResult: otime[0:26] + " " +
			levelC(log.PanicLevel) + " " +
			"caller message\n",
	},
	"caller report": {
		entry: &log.Entry{
			Message: "caller report message",
			Caller:  anyFrame,
			Logger: &log.Logger{
				ReportCaller: true,
			},
		},
		expectResult: otime[0:26] + " " +
			levelC(log.PanicLevel) + " " +
			"[file:123#function] caller report message\n",
	},

	// Test error.
	"error output": {
		entry: &log.Entry{
			Level:   log.PanicLevel,
			Message: "error message",
			Data: log.Fields{
				log.ErrorKey: errAny,
			},
		},
		expectError: nil,
		expectResult: otime[0:26] + " " +
			levelC(log.PanicLevel) + " error message " +
			dataC("error", errAny.Error()) + "\n",
	},
	"error output color-on": {
		colorMode: format.ColorModeOn,
		entry: &log.Entry{
			Level:   log.PanicLevel,
			Message: "error message",
			Data: log.Fields{
				log.ErrorKey: errAny,
			},
		},
		expectError: nil,
		expectResult: otime[0:26] + " " +
			levelC(log.PanicLevel) + " error message " +
			dataC("error", errAny.Error()) + "\n",
	},
	"error output color-off": {
		colorMode: format.ColorModeOff,
		entry: &log.Entry{
			Level:   log.PanicLevel,
			Message: "error message",
			Data: log.Fields{
				log.ErrorKey: errAny,
			},
		},
		expectError: nil,
		expectResult: otime[0:26] + " " +
			level(log.PanicLevel) + " error message " +
			data("error", errAny.Error()) + "\n",
	},
}

func TestPrettyFormat(t *testing.T) {
	test.Map(t, testPrettyFormatParams).
		Run(func(t test.Test, param testPrettyFormatParam) {
			// Given
			pretty := &format.Pretty{
				TimeFormat:  setupTimeFormat(param.timeFormat),
				ColorMode:   param.colorMode.Parse(!param.noTerminal),
				OrderMode:   param.orderMode.Parse(),
				LevelNames:  format.DefaultLevelNames,
				LevelColors: format.DefaultLevelColors,
			}

			if param.entry.Time == (time.Time{}) {
				time, err := time.Parse(time.RFC3339Nano, itime)
				assert.NoError(t, err)
				param.entry.Time = time
			}

			// When
			result, err := pretty.Format(param.entry)

			// Then
			if param.expect == nil {
				assert.Equal(t, param.expectError, err)
				assert.Equal(t, param.expectResult, string(result))
			} else {
				param.expect(t, string(result), err)
			}
		})
}

type testBufferWriteParam struct {
	colorMode    format.ColorModeString
	error        error
	setup        func(*format.Buffer)
	expect       mock.SetupFunc
	expectError  error
	expectString string
}

var testBufferWriteParams = map[string]testBufferWriteParam{
	// Test write byte.
	"write byte error": {
		error: errAny,
		setup: func(buffer *format.Buffer) {
			buffer.WriteByte(' ')
		},
		expectError: errAny,
	},
	"write byte failure": {
		setup: func(buffer *format.Buffer) {
			buffer.WriteByte(' ')
		},
		expect: mock.Chain(func(mocks *mock.Mocks) any {
			return mock.Get(mocks, NewMockBufferWriter).EXPECT().WriteByte(uint8(' ')).
				DoAndReturn(mocks.Do(format.BufferWriter.WriteByte, errAny))
		}, func(mocks *mock.Mocks) any {
			return mock.Get(mocks, NewMockBufferWriter).EXPECT().Bytes().
				DoAndReturn(mocks.Do(format.BufferWriter.Bytes, []byte("")))
		}),
		expectError: errAny,
	},
	"write byte": {
		setup: func(buffer *format.Buffer) {
			buffer.WriteByte(' ')
		},
		expectString: " ",
	},

	// Test write string.
	"write string error": {
		error: errAny,
		setup: func(buffer *format.Buffer) {
			buffer.WriteString("string")
		},
		expectError: errAny,
	},
	"write string failure": {
		setup: func(buffer *format.Buffer) {
			buffer.WriteString("string")
		},
		expect: mock.Chain(func(mocks *mock.Mocks) any {
			return mock.Get(mocks, NewMockBufferWriter).EXPECT().WriteString("string").
				DoAndReturn(mocks.Do(format.BufferWriter.WriteString, 0, errAny))
		}, func(mocks *mock.Mocks) any {
			return mock.Get(mocks, NewMockBufferWriter).EXPECT().Bytes().
				DoAndReturn(mocks.Do(format.BufferWriter.Bytes, []byte("")))
		}),
		expectError: errAny,
	},
	"write string": {
		setup: func(buffer *format.Buffer) {
			buffer.WriteString("string")
		},
		expectString: "string",
	},

	// Test write colored.
	"write colored error": {
		error: errAny,
		setup: func(buffer *format.Buffer) {
			buffer.WriteColored(format.ColorField, "string")
		},
		expectError: errAny,
	},
	"write colored default": {
		setup: func(buffer *format.Buffer) {
			buffer.WriteColored(format.ColorField, "string")
		},
		expectString: fieldC("string"),
	},
	"write colored color-off": {
		colorMode: format.ColorModeOff,
		setup: func(buffer *format.Buffer) {
			buffer.WriteColored(format.ColorField, "string")
		},
		expectString: "string",
	},
	"write colored color-on": {
		colorMode: format.ColorModeOn,
		setup: func(buffer *format.Buffer) {
			buffer.WriteColored(format.ColorField, "string")
		},
		expectString: fieldC("string"),
	},

	// Test write level.
	"write level error": {
		error: errAny,
		setup: func(buffer *format.Buffer) {
			buffer.WriteLevel(log.PanicLevel)
		},
		expectError: errAny,
	},
	"write level default": {
		setup: func(buffer *format.Buffer) {
			buffer.WriteLevel(log.PanicLevel)
		},
		expectString: levelC(log.PanicLevel),
	},
	"write level color-on": {
		colorMode: format.ColorModeOn,
		setup: func(buffer *format.Buffer) {
			buffer.WriteLevel(log.PanicLevel)
		},
		expectString: levelC(log.PanicLevel),
	},
	"write level color-off": {
		colorMode: format.ColorModeOff,
		setup: func(buffer *format.Buffer) {
			buffer.WriteLevel(log.PanicLevel)
		},
		expectString: level(log.PanicLevel),
	},

	// Test write colored field.
	"write field error": {
		error: errAny,
		setup: func(buffer *format.Buffer) {
			buffer.WriteField(format.FieldLevel, "value")
		},
		expectError: errAny,
	},
	"write field default": {
		setup: func(buffer *format.Buffer) {
			buffer.WriteField(format.FieldLevel, "value")
		},
		expectString: fieldC("value"),
	},
	"write field color-on": {
		colorMode: format.ColorModeOn,
		setup: func(buffer *format.Buffer) {
			buffer.WriteField(format.FieldLevel, "value")
		},
		expectString: fieldC("value"),
	},
	"write field color-off": {
		colorMode: format.ColorModeOff,
		setup: func(buffer *format.Buffer) {
			buffer.WriteField(format.FieldLevel, "value")
		},
		expectString: "value",
	},

	// Test write caller.
	"write caller error": {
		error: errAny,
		setup: func(buffer *format.Buffer) {
			buffer.WriteCaller(&log.Entry{
				Caller: anyFrame,
				Logger: &log.Logger{
					ReportCaller: true,
				},
			})
		},
		expectError: errAny,
	},
	"write caller on": {
		setup: func(buffer *format.Buffer) {
			buffer.WriteCaller(&log.Entry{
				Caller: anyFrame,
				Logger: &log.Logger{
					ReportCaller: true,
				},
			})
		},
		expectString: " [file:123#function]",
	},
	"write caller off": {
		setup: func(buffer *format.Buffer) {
			buffer.WriteCaller(&log.Entry{
				Logger: &log.Logger{
					ReportCaller: false,
				},
			})
		},
		expectString: "",
	},

	// Test write value.
	"write value error": {
		error: errAny,
		setup: func(buffer *format.Buffer) {
			buffer.WriteValue("value")
		},
		expectError: errAny,
	},
	"write value string": {
		setup: func(buffer *format.Buffer) {
			buffer.WriteValue("value")
		},
		expectString: "\"value\"",
	},
	"write value int": {
		setup: func(buffer *format.Buffer) {
			buffer.WriteValue(123)
		},
		expectString: "123",
	},
	"write value float": {
		setup: func(buffer *format.Buffer) {
			buffer.WriteValue(123.456)
		},
		expectString: "123.456",
	},
	"write value complex": {
		setup: func(buffer *format.Buffer) {
			buffer.WriteValue(123.456 + 789i)
		},
		expectString: "(123.456+789i)",
	},
	"write value bool": {
		setup: func(buffer *format.Buffer) {
			buffer.WriteValue(true)
		},
		expectString: "true",
	},

	// Test write data.
	"write data error": {
		error: errAny,
		setup: func(buffer *format.Buffer) {
			buffer.WriteData("key", "value")
		},
		expectError: errAny,
	},
	"write data default": {
		setup: func(buffer *format.Buffer) {
			buffer.WriteData("key", "value")
		},
		expectString: dataC("key", "value"),
	},
	"write data color-on error": {
		setup: func(buffer *format.Buffer) {
			buffer.WriteData(log.ErrorKey, errAny)
		},
		expectString: dataC(log.ErrorKey, errAny.Error()),
	},
	"write data color-on": {
		colorMode: format.ColorModeOn,
		setup: func(buffer *format.Buffer) {
			buffer.WriteData("key", "value")
		},
		expectString: dataC("key", "value"),
	},
	"write data color-off": {
		colorMode: format.ColorModeOff,
		setup: func(buffer *format.Buffer) {
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
			pretty := &format.Pretty{
				ColorMode:   param.colorMode.Parse(true),
				LevelNames:  format.DefaultLevelNames,
				LevelColors: format.DefaultLevelColors,
			}

			buffer := format.NewBuffer(pretty,
				setupWriter(mocks, param.expect))
			test.NewAccessor(buffer).Set("err", param.error)

			// When
			param.setup(buffer)
			result, err := buffer.Bytes()

			// Then
			assert.Equal(t, param.expectError, err)
			assert.Equal(t, param.expectString, string(result))
		})
}
