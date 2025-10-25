package log_test

import (
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/tkrop/go-testing/test"

	"github.com/tkrop/go-config/config"
	"github.com/tkrop/go-config/log"
)

func TestSetupRus(t *testing.T) {
	test.Map(t, setupTestCases).
		Run(func(t test.Test, param setupParams) {
			// Given
			logger := logrus.New()
			config := config.NewReader[config.Config]("TEST", "test").
				SetDefaultConfig("log", param.config, false).
				GetConfig(t.Name())

			// When
			config.Log.SetupRus(os.Stderr, logger)

			// Then
			switch param.config.Formatter {
			case log.FormatterText:
				assert.IsType(t, &logrus.TextFormatter{}, logger.Formatter)
				format := logger.Formatter.(*logrus.TextFormatter)
				assert.Equal(t, param.expectTimeFormat, format.TimestampFormat)
				assert.Equal(t, param.expectColorMode.CheckFlag(log.ColorOn),
					format.ForceColors)
			case log.FormatterJSON:
				assert.IsType(t, &logrus.JSONFormatter{}, logger.Formatter)
				assert.Equal(t, param.expectTimeFormat,
					logger.Formatter.(*logrus.JSONFormatter).TimestampFormat)
			case log.FormatterPretty:
				assert.IsType(t, &log.LogRusPretty{}, logger.Formatter)
				pretty := logger.Formatter.(*log.LogRusPretty).Setup
				assert.Equal(t, param.expectTimeFormat, pretty.TimeFormat)
				assert.Equal(t, param.expectColorMode, pretty.ColorMode)
				assert.Equal(t, param.expectOrderMode, pretty.OrderMode)
			default:
				assert.IsType(t, &log.LogRusPretty{}, logger.Formatter)
				assert.Equal(t, param.expectTimeFormat,
					logger.Formatter.(*log.LogRusPretty).TimeFormat)
			}

			assert.Equal(t, log.ParseLevel(param.expectLogLevel),
				log.ParseLevel(logger.GetLevel().String()))
			assert.Equal(t, param.expectLogCaller, logger.ReportCaller)
			assert.Equal(t, os.Stderr, logger.Out)
		})
}

func TestSetupNil(t *testing.T) {
	// Given
	config := config.NewReader[config.Config]("TEST", "test").
		GetConfig(t.Name())

	// When
	logger := config.Log.SetupRus(os.Stderr, nil)

	// Then
	assert.True(t, true)
	assert.Equal(t, logrus.StandardLogger(), logger)
}

// Arbitrary data for testing.
var anyData = logrus.Fields{
	"key1": "value1",
	"key2": "value2",
}

type testPrettyLogRusParam struct {
	noTerminal   bool
	config       log.Config
	entry        *logrus.Entry
	expect       func(t test.Test, result string, err error)
	expectResult string
}

var prettyLogRusTestCases = map[string]testPrettyLogRusParam{
	// Test levels with default.
	"level panic default": {
		config: log.Config{Level: "panic"},
		entry: &logrus.Entry{
			Level:   logrus.PanicLevel,
			Message: "panic message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.PanicLevel) + " panic message\n",
	},
	"level fatal default": {
		config: log.Config{Level: "fatal"},
		entry: &logrus.Entry{
			Level:   logrus.FatalLevel,
			Message: "fatal message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.FatalLevel) + " fatal message\n",
	},
	"level error default": {
		config: log.Config{Level: "error"},
		entry: &logrus.Entry{
			Level:   logrus.ErrorLevel,
			Message: "error message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.ErrorLevel) + " error message\n",
	},
	"level warn default": {
		config: log.Config{Level: "warn"},
		entry: &logrus.Entry{
			Level:   logrus.WarnLevel,
			Message: "warn message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.WarnLevel) + " warn message\n",
	},
	"level info default": {
		config: log.Config{Level: "info"},
		entry: &logrus.Entry{
			Level:   logrus.InfoLevel,
			Message: "info message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.InfoLevel) + " info message\n",
	},
	"level debug default": {
		config: log.Config{Level: "debug"},
		entry: &logrus.Entry{
			Level:   logrus.DebugLevel,
			Message: "debug message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.DebugLevel) + " debug message\n",
	},
	"level trace default": {
		config: log.Config{Level: "trace"},
		entry: &logrus.Entry{
			Level:   logrus.TraceLevel,
			Message: "trace message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.TraceLevel) + " trace message\n",
	},

	// Test levels with color.
	"level panic color-on": {
		config: log.Config{Level: "panic", ColorMode: log.ColorModeOn},
		entry: &logrus.Entry{
			Level:   logrus.PanicLevel,
			Message: "panic message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.PanicLevel) + " panic message\n",
	},
	"level fatal color-on": {
		config: log.Config{Level: "fatal", ColorMode: log.ColorModeOn},
		entry: &logrus.Entry{
			Level:   logrus.FatalLevel,
			Message: "fatal message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.FatalLevel) + " fatal message\n",
	},
	"level error color-on": {
		config: log.Config{Level: "error", ColorMode: log.ColorModeOn},
		entry: &logrus.Entry{
			Level:   logrus.ErrorLevel,
			Message: "error message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.ErrorLevel) + " error message\n",
	},
	"level warn color-on": {
		config: log.Config{Level: "warn", ColorMode: log.ColorModeOn},
		entry: &logrus.Entry{
			Level:   logrus.WarnLevel,
			Message: "warn message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.WarnLevel) + " warn message\n",
	},
	"level info color-on": {
		config: log.Config{Level: "info", ColorMode: log.ColorModeOn},
		entry: &logrus.Entry{
			Level:   logrus.InfoLevel,
			Message: "info message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.InfoLevel) + " info message\n",
	},
	"level debug color-on": {
		config: log.Config{Level: "debug", ColorMode: log.ColorModeOn},
		entry: &logrus.Entry{
			Level:   logrus.DebugLevel,
			Message: "debug message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.DebugLevel) + " debug message\n",
	},
	"level trace color-on": {
		config: log.Config{Level: "trace", ColorMode: log.ColorModeOn},
		entry: &logrus.Entry{
			Level:   logrus.TraceLevel,
			Message: "trace message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.TraceLevel) + " trace message\n",
	},

	// Test levels with color.
	"level panic color-off": {
		config: log.Config{Level: "panic", ColorMode: log.ColorModeOff},
		entry: &logrus.Entry{
			Level:   logrus.PanicLevel,
			Message: "panic message",
		},
		expectResult: otime[0:26] + " " +
			level(log.PanicLevel) + " panic message\n",
	},
	"level fatal color-off": {
		config: log.Config{Level: "fatal", ColorMode: log.ColorModeOff},
		entry: &logrus.Entry{
			Level:   logrus.FatalLevel,
			Message: "fatal message",
		},
		expectResult: otime[0:26] + " " +
			level(log.FatalLevel) + " fatal message\n",
	},
	"level error color-off": {
		config: log.Config{Level: "error", ColorMode: log.ColorModeOff},
		entry: &logrus.Entry{
			Level:   logrus.ErrorLevel,
			Message: "error message",
		},
		expectResult: otime[0:26] + " " +
			level(log.ErrorLevel) + " error message\n",
	},
	"level warn color-off": {
		config: log.Config{Level: "warning", ColorMode: log.ColorModeOff},
		entry: &logrus.Entry{
			Level:   logrus.WarnLevel,
			Message: "warn message",
		},
		expectResult: otime[0:26] + " " +
			level(log.WarnLevel) + " warn message\n",
	},
	"level info color-off": {
		config: log.Config{Level: "info", ColorMode: log.ColorModeOff},
		entry: &logrus.Entry{
			Level:   logrus.InfoLevel,
			Message: "info message",
		},
		expectResult: otime[0:26] + " " +
			level(log.InfoLevel) + " info message\n",
	},
	"level debug color-off": {
		config: log.Config{Level: "debug", ColorMode: log.ColorModeOff},
		entry: &logrus.Entry{
			Level:   logrus.DebugLevel,
			Message: "debug message",
		},
		expectResult: otime[0:26] + " " +
			level(log.DebugLevel) + " debug message\n",
	},
	"level trace color-off": {
		config: log.Config{Level: "trace", ColorMode: log.ColorModeOff},
		entry: &logrus.Entry{
			Level:   logrus.TraceLevel,
			Message: "trace message",
		},
		expectResult: otime[0:26] + " " +
			level(log.TraceLevel) + " trace message\n",
	},

	// Test order key value data.
	"data default": {
		entry: &logrus.Entry{
			Message: "data message",
			Data:    anyData,
		},
		expectResult: otime[0:26] + " " +
			levelC(log.PanicLevel) + " data message " +
			dataC("key1", "value1") + " " +
			dataC("key2", "value2") + "\n",
	},
	"data ordered": {
		config: log.Config{OrderMode: log.OrderModeOn},
		entry: &logrus.Entry{
			Message: "data message",
			Data:    anyData,
		},
		expectResult: otime[0:26] + " " +
			levelC(log.PanicLevel) + " data message " +
			dataC("key1", "value1") + " " +
			dataC("key2", "value2") + "\n",
	},
	"data unordered": {
		config: log.Config{OrderMode: log.OrderModeOff},
		entry: &logrus.Entry{
			Message: "data message",
			Data:    anyData,
		},
		expect: func(t test.Test, result string, _ error) {
			assert.Contains(t, result, otime[0:26]+" "+
				levelC(log.PanicLevel)+" "+"data message")
			assert.Contains(t, result, dataC("key1", "value1"))
			assert.Contains(t, result, dataC("key2", "value2"))
		},
	},

	// Test color modes.
	"data color-off": {
		config: log.Config{ColorMode: log.ColorModeOff},
		entry: &logrus.Entry{
			Message: "data message",
			Data:    anyData,
		},
		expectResult: otime[0:26] + " " +
			level(log.PanicLevel) + " data message " +
			data("key1", "value1") + " " +
			data("key2", "value2") + "\n",
	},
	"data color-on": {
		config: log.Config{ColorMode: log.ColorModeOn},
		entry: &logrus.Entry{
			Message: "data message",
			Data:    anyData,
		},
		expectResult: otime[0:26] + " " +
			levelC(log.PanicLevel) + " data message " +
			dataC("key1", "value1") + " " +
			dataC("key2", "value2") + "\n",
	},
	"data color-auto colorized": {
		config: log.Config{ColorMode: log.ColorModeAuto},
		entry: &logrus.Entry{
			Message: "data message",
			Data:    anyData,
		},
		expectResult: otime[0:26] + " " +
			levelC(log.PanicLevel) + " data message " +
			dataC("key1", "value1") + " " +
			dataC("key2", "value2") + "\n",
	},
	"data color-auto not-colorized": {
		noTerminal: true,
		config:     log.Config{ColorMode: log.ColorModeAuto},
		entry: &logrus.Entry{
			Message: "data message",
			Data:    anyData,
			Logger:  &logrus.Logger{Out: nil},
		},
		expectResult: otime[0:26] + " " +
			level(log.PanicLevel) + " data message " +
			data("key1", "value1") + " " +
			data("key2", "value2") + "\n",
	},
	"data color-levels": {
		config: log.Config{ColorMode: log.ColorModeLevels},
		entry: &logrus.Entry{
			Message: "data message",
			Data:    anyData,
		},
		expectResult: otime[0:26] + " " +
			levelC(log.PanicLevel) + " data message " +
			data("key1", "value1") + " " +
			data("key2", "value2") + "\n",
	},
	"data color-fields": {
		config: log.Config{ColorMode: log.ColorModeFields},
		entry: &logrus.Entry{
			Message: "data message",
			Data:    anyData,
		},
		expectResult: otime[0:26] + " " +
			level(log.PanicLevel) + " data message " +
			dataC("key1", "value1") + " " +
			dataC("key2", "value2") + "\n",
	},
	"data color-levels+fields": {
		config: log.Config{
			ColorMode: log.ColorModeLevels + "|" + log.ColorModeFields,
		},
		entry: &logrus.Entry{
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
		entry: &logrus.Entry{
			Level:   logrus.PanicLevel,
			Message: "default time message",
		},
		expectResult: otime[0:26] + " " +
			levelC(log.PanicLevel) + " " +
			"default time message\n",
	},
	"time short": {
		config: log.Config{
			TimeFormat: "2006-01-02 15:04:05",
		},
		entry: &logrus.Entry{
			Level:   logrus.PanicLevel,
			Message: "short time message",
		},
		expectResult: otime[0:19] + " " +
			levelC(log.PanicLevel) + " " +
			"short time message\n",
	},
	"time long": {
		config: log.Config{
			TimeFormat: "2006-01-02 15:04:05.000000000",
		},
		entry: &logrus.Entry{
			Level:   logrus.PanicLevel,
			Message: "long time message",
		},
		expectResult: otime[0:29] + " " +
			levelC(log.PanicLevel) + " " +
			"long time message\n",
	},

	// Report caller.
	"caller only": {
		entry: &logrus.Entry{
			Message: "caller message",
			Caller:  anyFrame,
		},
		expectResult: otime[0:26] + " " +
			levelC(log.PanicLevel) + " " +
			"caller message\n",
	},
	"caller report": {
		entry: &logrus.Entry{
			Message: "caller report message",
			Caller:  anyFrame,
			Logger: &logrus.Logger{
				ReportCaller: true,
			},
		},
		expectResult: otime[0:26] + " " +
			levelC(log.PanicLevel) + " " +
			"[file:123#function] caller report message\n",
	},

	// Test error.
	"error output": {
		entry: &logrus.Entry{
			Level:   logrus.PanicLevel,
			Message: "error message",
			Data: logrus.Fields{
				logrus.ErrorKey: assert.AnError,
			},
		},
		expectResult: otime[0:26] + " " +
			levelC(log.PanicLevel) + " error message " +
			dataC("error", assert.AnError.Error()) + "\n",
	},
	"error output color-on": {
		config: log.Config{ColorMode: log.ColorModeOn},
		entry: &logrus.Entry{
			Level:   logrus.PanicLevel,
			Message: "error message",
			Data: logrus.Fields{
				logrus.ErrorKey: assert.AnError,
			},
		},
		expectResult: otime[0:26] + " " +
			levelC(log.PanicLevel) + " error message " +
			dataC("error", assert.AnError.Error()) + "\n",
	},
	"error output color-off": {
		config: log.Config{ColorMode: log.ColorModeOff},
		entry: &logrus.Entry{
			Level:   logrus.PanicLevel,
			Message: "error message",
			Data: logrus.Fields{
				logrus.ErrorKey: assert.AnError,
			},
		},
		expectResult: otime[0:26] + " " +
			level(log.PanicLevel) + " error message " +
			data("error", assert.AnError.Error()) + "\n",
	},
}

func TestPrettyLogRus(t *testing.T) {
	test.Map(t, prettyLogRusTestCases).
		Run(func(t test.Test, param testPrettyLogRusParam) {
			// Given
			config := config.NewReader[config.Config]("X", "app").
				SetDefaultConfig("log", param.config, true).
				SetDefaults(func(r *config.Reader[config.Config]) {
					r.SetDefault("log.level", "trace")
				}).GetConfig("logrus")
			pretty := config.Log.SetupRus(os.Stderr, logrus.New()).Formatter
			pretty.(*log.LogRusPretty).
				ColorMode = param.config.ColorMode.Parse(!param.noTerminal)

			if param.entry.Time.Equal(time.Time{}) {
				time, err := time.Parse(time.RFC3339Nano, itime)
				assert.NoError(t, err)
				param.entry.Time = time
			}

			// When
			result, err := pretty.Format(param.entry)

			// Then
			if param.expect == nil {
				assert.Equal(t, param.expectResult, string(result))
			} else {
				param.expect(t, string(result), err)
			}
		})
}
