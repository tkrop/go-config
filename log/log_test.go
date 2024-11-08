package log_test

import (
	"errors"
	"runtime"
	"strconv"
	"time"

	"github.com/tkrop/go-config/log"
)

var (
	// otime is a fixed output time string for testing.
	otime = "2024-10-01 23:07:13.891012345Z"
	// itime is a fixed input time string for testing.
	itime = "2024-10-01T23:07:13.891012345Z"
	// ttime is a fixed time stamp for testing.
	ttime, terr = time.Parse(time.RFC3339Nano, itime)
	// Arbitrary frame for testing.
	anyFrame = &runtime.Frame{
		File:     "file",
		Function: "function",
		Line:     123,
	}
	// Arbitrary error for testing.
	errAny = errors.New("any error")
)

// caller returns the file and line of the caller.
func caller(offset int) string {
	if _, file, line, ok := runtime.Caller(1); ok {
		return file + ":" + strconv.Itoa(line+offset)
	}
	return "unknown"
}

// Helper functions for testing log levels without color.
func level(level log.Level) string {
	return log.DefaultLevelNames[level]
}

// Helper functions for testing log levels with color.
func levelC(level log.Level) string {
	return "\x1b[" + log.DefaultLevelColors[level] +
		"m" + log.DefaultLevelNames[level] + "\x1b[0m"
}

// Helper functions for testing fields without color.
func field(value string) string {
	return value
}

// Helper functions for testing fields with color.
func fieldC(value string) string {
	return "\x1b[" + log.ColorField + "m" + value + "\x1b[0m"
}

// Helper functions for testing key data without color.
func key(key string) string {
	return key + "="
}

// Helper functions for testing key data with color.
func keyC(key string) string {
	color := log.ColorField
	if key == log.DefaultErrorName {
		color = log.ColorError
	}
	return "\x1b[" + color + "m" + key + "\x1b[0m="
}

// Helper functions for testing key-value data without color.
func data(key, value string) string {
	return key + "=\"" + value + "\""
}

// Helper functions for testing key-value data with color.
func dataC(key, value string) string {
	color := log.ColorField
	if key == log.DefaultErrorName {
		color = log.ColorError
	}
	return "\x1b[" + color + "m" + key + "\x1b[0m=\"" + value + "\""
}

type setupParams struct {
	config           *log.Config
	expectTimeFormat string
	expectLogLevel   string
	expectLogCaller  bool
	expectColorMode  log.ColorMode
	expectOrderMode  log.OrderMode
}

var testSetupParams = map[string]setupParams{
	"read default config": {
		config:           &log.Config{},
		expectLogLevel:   log.DefaultLevel,
		expectTimeFormat: log.DefaultTimeFormat,
		expectLogCaller:  log.DefaultCaller,
		expectColorMode:  log.ColorOff,
		expectOrderMode:  log.OrderOn,
	},

	"log level panic": {
		config: &log.Config{
			Level: log.LevelPanic,
		},
		expectLogLevel:   log.LevelPanic,
		expectTimeFormat: log.DefaultTimeFormat,
		expectLogCaller:  log.DefaultCaller,
		expectColorMode:  log.ColorOff,
		expectOrderMode:  log.OrderOn,
	},

	"log level fatal": {
		config: &log.Config{
			Level: log.LevelFatal,
		},
		expectLogLevel:   log.LevelFatal,
		expectTimeFormat: log.DefaultTimeFormat,
		expectLogCaller:  log.DefaultCaller,
		expectColorMode:  log.ColorOff,
		expectOrderMode:  log.OrderOn,
	},

	"log level error": {
		config: &log.Config{
			Level: log.LevelError,
		},
		expectLogLevel:   log.LevelError,
		expectTimeFormat: log.DefaultTimeFormat,
		expectLogCaller:  log.DefaultCaller,
		expectColorMode:  log.ColorOff,
		expectOrderMode:  log.OrderOn,
	},

	"log level warn": {
		config: &log.Config{
			Level: log.LevelWarn,
		},
		expectLogLevel:   log.LevelWarn,
		expectTimeFormat: log.DefaultTimeFormat,
		expectLogCaller:  log.DefaultCaller,
		expectColorMode:  log.ColorOff,
		expectOrderMode:  log.OrderOn,
	},

	"log level warning": {
		config: &log.Config{
			Level: log.LevelWarning,
		},
		expectLogLevel:   log.LevelWarning,
		expectTimeFormat: log.DefaultTimeFormat,
		expectLogCaller:  log.DefaultCaller,
		expectColorMode:  log.ColorOff,
		expectOrderMode:  log.OrderOn,
	},

	"log level info": {
		config: &log.Config{
			Level: log.LevelInfo,
		},
		expectLogLevel:   log.LevelInfo,
		expectTimeFormat: log.DefaultTimeFormat,
		expectLogCaller:  log.DefaultCaller,
		expectColorMode:  log.ColorOff,
		expectOrderMode:  log.OrderOn,
	},

	"log level debug": {
		config: &log.Config{
			Level: log.LevelDebug,
		},
		expectLogLevel:   log.LevelDebug,
		expectTimeFormat: log.DefaultTimeFormat,
		expectLogCaller:  log.DefaultCaller,
		expectColorMode:  log.ColorOff,
		expectOrderMode:  log.OrderOn,
	},

	"log level trace": {
		config: &log.Config{
			Level: log.LevelTrace,
		},
		expectLogLevel:   log.LevelTrace,
		expectTimeFormat: log.DefaultTimeFormat,
		expectLogCaller:  log.DefaultCaller,
		expectColorMode:  log.ColorOff,
		expectOrderMode:  log.OrderOn,
	},

	"log level invalid": {
		config: &log.Config{
			Level: "invalid",
		},
		expectLogLevel:   log.LevelInfo,
		expectTimeFormat: log.DefaultTimeFormat,
		expectLogCaller:  log.DefaultCaller,
		expectColorMode:  log.ColorOff,
		expectOrderMode:  log.OrderOn,
	},

	"time format date": {
		config: &log.Config{
			TimeFormat: "2024-12-31",
		},
		expectLogLevel:   log.DefaultLevel,
		expectTimeFormat: "2024-12-31",
		expectLogCaller:  log.DefaultCaller,
		expectColorMode:  log.ColorOff,
		expectOrderMode:  log.OrderOn,
	},

	"caller enabled": {
		config: &log.Config{
			Caller: true,
		},
		expectLogLevel:   log.DefaultLevel,
		expectTimeFormat: log.DefaultTimeFormat,
		expectLogCaller:  true,
		expectColorMode:  log.ColorOff,
		expectOrderMode:  log.OrderOn,
	},

	"formatter text": {
		config: &log.Config{
			Formatter: log.FormatterText,
		},
		expectLogLevel:   log.DefaultLevel,
		expectTimeFormat: log.DefaultTimeFormat,
		expectLogCaller:  log.DefaultCaller,
		expectColorMode:  log.ColorOff,
		expectOrderMode:  log.OrderOn,
	},

	"formatter json": {
		config: &log.Config{
			Formatter: log.FormatterJSON,
		},
		expectLogLevel:   log.DefaultLevel,
		expectTimeFormat: log.DefaultTimeFormat,
		expectLogCaller:  log.DefaultCaller,
		expectColorMode:  log.ColorOff,
		expectOrderMode:  log.OrderOn,
	},

	"formatter pretty default": {
		config: &log.Config{
			Formatter: log.FormatterPretty,
		},
		expectLogLevel:   log.DefaultLevel,
		expectTimeFormat: log.DefaultTimeFormat,
		expectLogCaller:  log.DefaultCaller,
		expectColorMode:  log.ColorOff,
		expectOrderMode:  log.OrderOn,
	},

	"formatter pretty color-on": {
		config: &log.Config{
			Formatter: log.FormatterPretty,
			ColorMode: log.ColorModeOn,
		},
		expectLogLevel:   log.DefaultLevel,
		expectTimeFormat: log.DefaultTimeFormat,
		expectColorMode:  log.ColorOn,
		expectOrderMode:  log.OrderOn,
		expectLogCaller:  log.DefaultCaller,
	},

	"formatter pretty color-off": {
		config: &log.Config{
			Formatter: log.FormatterPretty,
			ColorMode: log.ColorModeOff,
		},
		expectLogLevel:   log.DefaultLevel,
		expectTimeFormat: log.DefaultTimeFormat,
		expectColorMode:  log.ColorOff,
		expectOrderMode:  log.OrderOn,
		expectLogCaller:  log.DefaultCaller,
	},

	"formatter pretty color-levels": {
		config: &log.Config{
			Formatter: log.FormatterPretty,
			ColorMode: log.ColorModeLevels,
		},
		expectLogLevel:   log.DefaultLevel,
		expectTimeFormat: log.DefaultTimeFormat,
		expectColorMode:  log.ColorLevels,
		expectOrderMode:  log.OrderOn,
		expectLogCaller:  log.DefaultCaller,
	},

	"formatter pretty color-fields": {
		config: &log.Config{
			Formatter: log.FormatterPretty,
			ColorMode: log.ColorModeFields,
		},
		expectLogLevel:   log.DefaultLevel,
		expectTimeFormat: log.DefaultTimeFormat,
		expectColorMode:  log.ColorFields,
		expectOrderMode:  log.OrderOn,
		expectLogCaller:  log.DefaultCaller,
	},

	"formatter pretty color-any": {
		config: &log.Config{
			Formatter: log.FormatterPretty,
			ColorMode: "any",
		},
		expectLogLevel:   log.DefaultLevel,
		expectTimeFormat: log.DefaultTimeFormat,
		expectColorMode:  log.ColorOff,
		expectOrderMode:  log.OrderOn,
		expectLogCaller:  log.DefaultCaller,
	},

	"formatter pretty order-on": {
		config: &log.Config{
			Formatter: log.FormatterPretty,
			OrderMode: log.OrderModeOn,
		},
		expectLogLevel:   log.DefaultLevel,
		expectTimeFormat: log.DefaultTimeFormat,
		expectColorMode:  log.ColorOff,
		expectOrderMode:  log.OrderOn,
		expectLogCaller:  log.DefaultCaller,
	},

	"formatter pretty order-off": {
		config: &log.Config{
			Formatter: log.FormatterPretty,
			OrderMode: log.OrderModeOff,
		},
		expectLogLevel:   log.DefaultLevel,
		expectTimeFormat: log.DefaultTimeFormat,
		expectColorMode:  log.ColorOff,
		expectOrderMode:  log.OrderOff,
		expectLogCaller:  log.DefaultCaller,
	},

	"formatter pretty order-any": {
		config: &log.Config{
			Formatter: log.FormatterPretty,
			OrderMode: "any",
		},
		expectLogLevel:   log.DefaultLevel,
		expectTimeFormat: log.DefaultTimeFormat,
		expectColorMode:  log.ColorOff,
		expectOrderMode:  log.OrderOff,
		expectLogCaller:  log.DefaultCaller,
	},
}
