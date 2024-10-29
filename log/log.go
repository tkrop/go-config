// Package log provides common log handling based on [logrus][logrus] for
// services, jobs, and commands with integrated configuration loading.
package log

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"golang.org/x/term"

	"github.com/tkrop/go-config/log/format"
)

// Config common configuration for logging.
type Config struct {
	// Level is defining the logger level (default `info`).
	Level string `default:"info"`
	// TImeFormat is defining the time format for timestamps.
	TimeFormat string `default:"2006-01-02T15:04:05.999999"`
	// Caller is defining whether the caller is logged (default `false`).
	Caller bool `default:"false"`
	// File is defining the file name used for the log output.
	File string `default:"/dev/stderr"`
	// ColorMode is defining the color mode used for logging.
	ColorMode format.ColorModeString `default:"auto"`
	// OrderMode is defining the order mode used for logging.
	OrderMode format.OrderModeString `default:"on"`
	// Formatter is defining the formatter used for logging.
	Formatter format.Formatter `default:"pretty"`
}

// ColorMode is the color mode used for logging.
type ColorModeString format.ColorModeString

// Color modes.
const (
	// ColorOff disables the color mode.
	ColorModeOff format.ColorModeString = format.ColorModeOff
	// ColorOn enables the color mode.
	ColorModeOn format.ColorModeString = format.ColorModeOn
	// ColorAuto enables the automatic color mode.
	ColorModeAuto format.ColorModeString = format.ColorModeAuto
	// ColorLevels enables the color mode for log level.
	ColorModeLevels format.ColorModeString = format.ColorModeLevels
	// ColorFields enables the color mode for fields.
	ColorModeFields format.ColorModeString = format.ColorModeFields
)

// OrderMode is the order mode used for logging.
type OrderModeString format.OrderModeString

// Order modes.
const (
	OrderModeAuto format.OrderModeString = ""
	// OrderOn enables the order mode.
	OrderModeOn format.OrderModeString = format.OrderModeOn
	// OrderOff disables the order mode.
	OrderModeOff format.OrderModeString = format.OrderModeOff
)

// Formatter is the formatter used for logging output.
type Formatter format.Formatter

// Supported formatters.
const (
	// Pretty is setting up a pretty formatter.
	FormatterPretty format.Formatter = format.FormatterPretty
	// Text is setting up a text formatter.
	FormatterText format.Formatter = format.FormatterText
	// JSON is setting up a JSON formatter.
	FormatterJSON format.Formatter = format.FormatterJSON
)

// IsTerminal checks whether the given writer is a terminal.
func IsTerminal(writer io.Writer) bool {
	if file, ok := writer.(*os.File); ok {
		// #nosec G115 // is a safe conversion for files.
		return term.IsTerminal(int(file.Fd()))
	}
	return false
}

// SetupRus is setting up the given logger using. It sets the formatter, the
// log level, and the report caller flag. If no logger is given, the standard
// logger is used.
func (c *Config) SetupRus(logger *logrus.Logger) *logrus.Logger {
	// Uses the standard logger if no logger is given.
	if logger == nil {
		logger = logrus.StandardLogger()
	}

	// Sets up the log output format.
	switch c.Formatter {
	case FormatterText:
		mode := c.ColorMode.Parse(IsTerminal(logger.Out))
		logger.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: c.TimeFormat,
			FullTimestamp:   true,
			ForceColors:     mode&format.ColorOn == format.ColorOn,
			DisableColors:   mode&format.ColorOff == format.ColorOff,
		})
	case FormatterJSON:
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: c.TimeFormat,
		})
	case FormatterPretty:
		fallthrough
	default:
		logger.SetFormatter(&format.Pretty{
			TimeFormat:  c.TimeFormat,
			ColorMode:   c.ColorMode.Parse(IsTerminal(logger.Out)),
			OrderMode:   c.OrderMode.Parse(),
			LevelNames:  format.DefaultLevelNames,
			LevelColors: format.DefaultLevelColors,
		})
	}

	// Helpful setting in certain debug situations.
	logger.SetReportCaller(c.Caller)

	// Sets up the log level if given.
	logLevel, err := logrus.ParseLevel(c.Level)
	if err != nil {
		logger.WithError(err).WithFields(logrus.Fields{
			"config": c.Level,
		}).Info("failed setting log level")
	} else {
		logger.SetLevel(logLevel)
		logger.WithFields(logrus.Fields{
			"level": c.Level,
		}).Info("setting up log level")
	}

	return logger
}
