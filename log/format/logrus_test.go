package format_test

import (
	"bytes"
	"errors"
	"runtime"
	"testing"
	"time"

	"github.com/mattn/go-tty"
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
	return key + "=" + value
}

// Helper functions for testing key-value data with color.
func dataC(key, value string) string {
	color := format.ColorField
	if key == log.ErrorKey {
		color = format.ColorError
	}
	return "\x1b[" + color + "m" + key + "\x1b[0m=" + value
}

type testPrettyFormatParam struct {
	timeFormat   string
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
			levelC(log.PanicLevel) + " panic message",
	},
	"level fatal default": {
		entry: &log.Entry{
			Level:   log.FatalLevel,
			Message: "fatal message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.FatalLevel) + " fatal message",
	},
	"level error default": {
		entry: &log.Entry{
			Level:   log.ErrorLevel,
			Message: "error message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.ErrorLevel) + " error message",
	},
	"level warn default": {
		entry: &log.Entry{
			Level:   log.WarnLevel,
			Message: "warn message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.WarnLevel) + " warn message",
	},
	"level info default": {
		entry: &log.Entry{
			Level:   log.InfoLevel,
			Message: "info message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.InfoLevel) + " info message",
	},
	"level debug default": {
		entry: &log.Entry{
			Level:   log.DebugLevel,
			Message: "debug message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.DebugLevel) + " debug message",
	},
	"level trace default": {
		entry: &log.Entry{
			Level:   log.TraceLevel,
			Message: "trace message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.TraceLevel) + " trace message",
	},

	// Test levels with color.
	"level panic color-on": {
		colorMode: format.ColorModeOn,
		entry: &log.Entry{
			Level:   log.PanicLevel,
			Message: "panic message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.PanicLevel) + " panic message",
	},
	"level fatal color-on": {
		colorMode: format.ColorModeOn,
		entry: &log.Entry{
			Level:   log.FatalLevel,
			Message: "fatal message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.FatalLevel) + " fatal message",
	},
	"level error color-on": {
		colorMode: format.ColorModeOn,
		entry: &log.Entry{
			Level:   log.ErrorLevel,
			Message: "error message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.ErrorLevel) + " error message",
	},
	"level warn color-on": {
		colorMode: format.ColorModeOn,
		entry: &log.Entry{
			Level:   log.WarnLevel,
			Message: "warn message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.WarnLevel) + " warn message",
	},
	"level info color-on": {
		colorMode: format.ColorModeOn,
		entry: &log.Entry{
			Level:   log.InfoLevel,
			Message: "info message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.InfoLevel) + " info message",
	},
	"level debug color-on": {
		colorMode: format.ColorModeOn,
		entry: &log.Entry{
			Level:   log.DebugLevel,
			Message: "debug message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.DebugLevel) + " debug message",
	},
	"level trace color-on": {
		colorMode: format.ColorModeOn,
		entry: &log.Entry{
			Level:   log.TraceLevel,
			Message: "trace message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.TraceLevel) + " trace message",
	},

	// Test levels with color.
	"level panic color-off": {
		colorMode: format.ColorModeOff,
		entry: &log.Entry{
			Level:   log.PanicLevel,
			Message: "panic message",
		},
		expectResult: otime[0:26] + " " +
			level(log.PanicLevel) + " panic message",
	},
	"level fatal color-off": {
		colorMode: format.ColorModeOff,
		entry: &log.Entry{
			Level:   log.FatalLevel,
			Message: "fatal message",
		},
		expectResult: otime[0:26] + " " +
			level(log.FatalLevel) + " fatal message",
	},
	"level error color-off": {
		colorMode: format.ColorModeOff,
		entry: &log.Entry{
			Level:   log.ErrorLevel,
			Message: "error message",
		},
		expectResult: otime[0:26] + " " +
			level(log.ErrorLevel) + " error message",
	},
	"level warn color-off": {
		colorMode: format.ColorModeOff,
		entry: &log.Entry{
			Level:   log.WarnLevel,
			Message: "warn message",
		},
		expectResult: otime[0:26] + " " +
			level(log.WarnLevel) + " warn message",
	},
	"level info color-off": {
		colorMode: format.ColorModeOff,
		entry: &log.Entry{
			Level:   log.InfoLevel,
			Message: "info message",
		},
		expectResult: otime[0:26] + " " +
			level(log.InfoLevel) + " info message",
	},
	"level debug color-off": {
		colorMode: format.ColorModeOff,
		entry: &log.Entry{
			Level:   log.DebugLevel,
			Message: "debug message",
		},
		expectResult: otime[0:26] + " " +
			level(log.DebugLevel) + " debug message",
	},
	"level trace color-off": {
		colorMode: format.ColorModeOff,
		entry: &log.Entry{
			Level:   log.TraceLevel,
			Message: "trace message",
		},
		expectResult: otime[0:26] + " " +
			level(log.TraceLevel) + " trace message",
	},

	// Test order key value data.
	"data default": {
		entry: &log.Entry{
			Message: "data message",
			Data:    anyData,
		},
		expectResult: otime[0:26] + " " +
			levelC(log.PanicLevel) + " data message " +
			dataC("key1", "value1") + " " + dataC("key2", "value2"),
	},
	"data ordered": {
		orderMode: format.OrderModeOn,
		entry: &log.Entry{
			Message: "data message",
			Data:    anyData,
		},
		expectResult: otime[0:26] + " " +
			levelC(log.PanicLevel) + " data message " +
			dataC("key1", "value1") + " " + dataC("key2", "value2"),
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
			data("key1", "value1") + " " + data("key2", "value2"),
	},
	"data color-on": {
		colorMode: format.ColorModeOn,
		entry: &log.Entry{
			Message: "data message",
			Data:    anyData,
		},
		expectResult: otime[0:26] + " " +
			levelC(log.PanicLevel) + " data message " +
			dataC("key1", "value1") + " " + dataC("key2", "value2"),
	},
	"data color-auto": {
		colorMode: format.ColorModeAuto,
		entry: &log.Entry{
			Message: "data message",
			Data:    anyData,
		},
		expectResult: otime[0:26] + " " +
			levelC(log.PanicLevel) + " data message " +
			dataC("key1", "value1") + " " + dataC("key2", "value2"),
	},
	"data color-auto no-tty": {
		colorMode: format.ColorModeAuto,
		entry: &log.Entry{
			Message: "data message",
			Data:    anyData,
			Logger:  &log.Logger{Out: nil},
		},
		expectResult: otime[0:26] + " " +
			level(log.PanicLevel) + " data message " +
			data("key1", "value1") + " " + data("key2", "value2"),
	},
	"data color-levels": {
		colorMode: format.ColorModeLevels,
		entry: &log.Entry{
			Message: "data message",
			Data:    anyData,
		},
		expectResult: otime[0:26] + " " +
			levelC(log.PanicLevel) + " data message " +
			data("key1", "value1") + " " + data("key2", "value2"),
	},
	"data color-fields": {
		colorMode: format.ColorModeFields,
		entry: &log.Entry{
			Message: "data message",
			Data:    anyData,
		},
		expectResult: otime[0:26] + " " +
			level(log.PanicLevel) + " data message " +
			dataC("key1", "value1") + " " + dataC("key2", "value2"),
	},
	"data color-levels+fields": {
		colorMode: format.ColorModeLevels + "|" + format.ColorModeFields,
		entry: &log.Entry{
			Message: "data message",
			Data:    anyData,
		},
		expectResult: otime[0:26] + " " +
			levelC(log.PanicLevel) + " data message " +
			dataC("key1", "value1") + " " + dataC("key2", "value2"),
	},

	// Time format.
	"time default": {
		entry: &log.Entry{
			Level:   log.PanicLevel,
			Message: "default time message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.PanicLevel) + " " +
			"default time message",
	},
	"time short": {
		timeFormat: "2006-01-02 15:04:05",
		entry: &log.Entry{
			Level:   log.PanicLevel,
			Message: "short time message",
		},
		expectResult: otime[0:19] + " " +
			levelC(log.PanicLevel) + " " +
			"short time message",
	},
	"time long": {
		timeFormat: "2006-01-02 15:04:05.000000000",
		entry: &log.Entry{
			Level:   log.PanicLevel,
			Message: "long time message",
		},
		expectResult: otime[0:29] + " " +
			levelC(log.PanicLevel) + " " +
			"long time message",
	},

	// Report caller.
	"caller only": {
		entry: &log.Entry{
			Message: "caller message",
			Caller:  anyFrame,
		},
		expectResult: otime[0:26] + " " +
			levelC(log.PanicLevel) + " " +
			"caller message",
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
			level(log.PanicLevel) + " " +
			"[file:123#function] caller report message",
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
			dataC("error", errAny.Error()),
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
			dataC("error", errAny.Error()),
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
			data("error", errAny.Error()),
	},
}

func TestPrettyFormat(t *testing.T) {
	tty, err := tty.Open()
	assert.NoError(t, err)

	test.Map(t, testPrettyFormatParams).
		Run(func(t test.Test, param testPrettyFormatParam) {
			// Given
			pretty := &format.Pretty{
				TimeFormat: param.timeFormat,
				ColorMode:  param.colorMode.Parse(),
				OrderMode:  param.orderMode.Parse(),
			}

			if param.entry.Time == (time.Time{}) {
				time, err := time.Parse(time.RFC3339Nano, itime)
				assert.NoError(t, err)
				param.entry.Time = time
			}
			if param.entry.Logger == nil {
				param.entry.Logger = log.New()
				param.entry.Logger.Out = tty.Output()
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
		}).
		Cleanup(func() {
			tty.Close()
		})
}

func setupWriter(mocks *mock.Mocks, expect mock.SetupFunc) format.BufferWriter {
	var writer format.BufferWriter
	if expect != nil {
		writer = mock.Get(mocks, NewMockBufferWriter)
	} else {
		writer = &bytes.Buffer{}
	}
	return writer
}

type testBufferWriteParam struct {
	pretty       *format.Pretty
	error        error
	setup        func(test.Test, *format.Buffer)
	expect       mock.SetupFunc
	expectError  error
	expectString string
}

var testBufferWriteParams = map[string]testBufferWriteParam{
	// Test write byte.
	"write byte error": {
		error: errAny,
		setup: func(t test.Test, buffer *format.Buffer) {
			buffer.WriteByte(' ')
		},
		expectError: errAny,
	},
	"write byte failure": {
		setup: func(t test.Test, buffer *format.Buffer) {
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
		setup: func(t test.Test, buffer *format.Buffer) {
			buffer.WriteByte(' ')
		},
		expectString: " ",
	},

	// Test write string.
	"write string error": {
		error: errAny,
		setup: func(t test.Test, buffer *format.Buffer) {
			buffer.WriteString("string")
		},
		expectError: errAny,
	},
	"write string failure": {
		setup: func(t test.Test, buffer *format.Buffer) {
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
		setup: func(t test.Test, buffer *format.Buffer) {
			buffer.WriteString("string")
		},
		expectString: "string",
	},

	// Test write colored.
	"write colored error": {
		error: errAny,
		setup: func(t test.Test, buffer *format.Buffer) {
			buffer.WriteColored(format.ColorField, "string")
		},
		expectError: errAny,
	},
	"write colored default": {
		setup: func(t test.Test, buffer *format.Buffer) {
			buffer.WriteColored(format.ColorField, "string")
		},
		expectString: fieldC("string"),
	},
	"write colored color-off": {
		pretty: &format.Pretty{
			ColorMode: format.ColorOff,
		},
		setup: func(t test.Test, buffer *format.Buffer) {
			buffer.WriteColored(format.ColorField, "string")
		},
		expectString: "string",
	},
	"write colored color-on": {
		pretty: &format.Pretty{
			ColorMode: format.ColorOn,
		},
		setup: func(t test.Test, buffer *format.Buffer) {
			buffer.WriteColored(format.ColorField, "string")
		},
		expectString: fieldC("string"),
	},

	// Test write level.
	"write level error": {
		error: errAny,
		setup: func(t test.Test, buffer *format.Buffer) {
			buffer.WriteLevel(log.PanicLevel)
		},
		expectError: errAny,
	},
	"write level default": {
		setup: func(t test.Test, buffer *format.Buffer) {
			buffer.WriteLevel(log.PanicLevel)
		},
		expectString: levelC(log.PanicLevel),
	},
	"write level color-on": {
		pretty: &format.Pretty{
			ColorMode: format.ColorOn,
		},
		setup: func(t test.Test, buffer *format.Buffer) {
			buffer.WriteLevel(log.PanicLevel)
		},
		expectString: levelC(log.PanicLevel),
	},
	"write level color-off": {
		pretty: &format.Pretty{
			ColorMode: format.ColorOff,
		},
		setup: func(t test.Test, buffer *format.Buffer) {
			buffer.WriteLevel(log.PanicLevel)
		},
		expectString: level(log.PanicLevel),
	},

	// Test write colored field.
	"write field error": {
		error: errAny,
		setup: func(t test.Test, buffer *format.Buffer) {
			buffer.WriteField(format.FieldLevel, "value")
		},
		expectError: errAny,
	},
	"write field default": {
		setup: func(t test.Test, buffer *format.Buffer) {
			buffer.WriteField(format.FieldLevel, "value")
		},
		expectString: fieldC("value"),
	},
	"write field color-on": {
		pretty: &format.Pretty{
			ColorMode: format.ColorOn,
		},
		setup: func(t test.Test, buffer *format.Buffer) {
			buffer.WriteField(format.FieldLevel, "value")
		},
		expectString: fieldC("value"),
	},
	"write field color-off": {
		pretty: &format.Pretty{
			ColorMode: format.ColorOff,
		},
		setup: func(t test.Test, buffer *format.Buffer) {
			buffer.WriteField(format.FieldLevel, "value")
		},
		expectString: "value",
	},

	// Test write caller.
	"write caller error": {
		error: errAny,
		setup: func(t test.Test, buffer *format.Buffer) {
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
		setup: func(t test.Test, buffer *format.Buffer) {
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
		setup: func(t test.Test, buffer *format.Buffer) {
			buffer.WriteCaller(&log.Entry{
				Logger: &log.Logger{
					ReportCaller: false,
				},
			})
		},
		expectString: "",
	},

	// Test write data.
	"write data error": {
		error: errAny,
		setup: func(t test.Test, buffer *format.Buffer) {
			buffer.WriteData("key", "value")
		},
		expectError: errAny,
	},
	"write data default": {
		setup: func(t test.Test, buffer *format.Buffer) {
			buffer.WriteData("key", "value")
		},
		expectString: dataC("key", "value"),
	},
	"write data color-on error": {
		setup: func(t test.Test, buffer *format.Buffer) {
			buffer.WriteData(log.ErrorKey, errAny)
		},
		expectString: dataC(log.ErrorKey, errAny.Error()),
	},
	"write data color-on": {
		pretty: &format.Pretty{
			ColorMode: format.ColorOn,
		},
		setup: func(t test.Test, buffer *format.Buffer) {
			buffer.WriteData("key", "value")
		},
		expectString: dataC("key", "value"),
	},
	"write data color-off": {
		pretty: &format.Pretty{
			ColorMode: format.ColorOff,
		},
		setup: func(t test.Test, buffer *format.Buffer) {
			buffer.WriteData("key", "value")
		},
		expectString: data("key", "value"),
	},
}

func TestBufferWrite(t *testing.T) {
	tty, err := tty.Open()
	assert.NoError(t, err)

	test.Map(t, testBufferWriteParams).
		Run(func(t test.Test, param testBufferWriteParam) {
			// Given
			mocks := mock.NewMocks(t).Expect(param.expect)
			if param.pretty == nil {
				param.pretty = &format.Pretty{}
			}
			param.pretty.Init(tty.Output())
			buffer := format.NewBuffer(param.pretty,
				setupWriter(mocks, param.expect))
			test.NewAccessor(buffer).Set("err", param.error)

			// When
			param.setup(t, buffer)
			result, err := buffer.Bytes()

			// Then
			assert.Equal(t, param.expectError, err)
			assert.Equal(t, param.expectString, string(result))
		}).
		Cleanup(func() {
			tty.Close()
		})
}
