// Package log provides common log handling based on [logrus][logrus] for
// services, jobs, and commands with integrated configuration loading.
package log

import (
	log "github.com/sirupsen/logrus"
)

// Config common configuration for logging.
type Config struct {
	// Level is defining the logger level (default `info`).
	Level string `default:"info"`
	// TImeFormat is defining the logger time format.
	TimeFormat string `default:"2006-01-02T15:04:05.999999"`
	// Caller is defining whether the caller is logged (default `false`).
	Caller bool `default:"false"`
	// File is defining the file name for logger output.
	File string `default:"/dev/stderr"`
}

// Exported log types to be used in the application.
type (
	// Logger is the logrus logger.
	Logger = log.Logger
	// Entry is the logrus entry.
	Entry = log.Entry
	// Fields is the logrus fields.
	Fields = log.Fields
	// Level is the logrus level.
	Level = log.Level
	// Hook is the logrus hook.
	Hook = log.Hook

	// Formatter is the logrus formatter.
	Formatter = log.Formatter
	// TextFormatter is the logrus text formatter.
	TextFormatter = log.TextFormatter
	// JSONFormatter is the logrus JSON formatter.
	JSONFormatter = log.JSONFormatter
)

//revive:disable:max-public-structs // export log types

// Exported log functions to be used in the application.
var (
	// New creates a new logger.
	New = log.New
	// StandardLogger returns the standard logger.
	StandardLogger = log.StandardLogger
	// ParseLevel parses a log level.
	ParseLevel = log.ParseLevel
	// GetLevel returns the current log level.
	GetLevel = log.GetLevel
	// SetLevel sets the log level of the logger.
	SetLevel = log.SetLevel
	// IsLevelEnabled checks if the log level is enabled.
	IsLevelEnabled = log.IsLevelEnabled

	// SetOutput sets the output of the logger.
	SetOutput = log.SetOutput
	// SetFormatter sets the formatter of the logger.
	SetFormatter = log.SetFormatter
	// SetReportCaller sets the report caller flag of the logger.
	SetReportCaller = log.SetReportCaller
	// AddHook adds a hook to the logger.
	AddHook = log.AddHook

	// WithTime adds the current time to the entry.
	WithTime = log.WithTime
	// WithContext adds the context to the entry.
	WithContext = log.WithContext
	// WithError adds the error to the entry.
	WithError = log.WithError
	// WithField adds a field to the entry.
	WithField = log.WithField
	// WithFields adds fields to the entry.
	WithFields = log.WithFields

	// Tracef logs a message at level Trace.
	Tracef = log.Tracef
	// Debugf logs a message at level Debug.
	Debugf = log.Debugf
	// Infof logs a message at level Info.
	Infof = log.Infof
	// Printf logs a message at level Info.
	Printf = log.Printf
	// Warnf logs a message at level Warn.
	Warnf = log.Warnf
	// Warningf logs a message at level Warn.
	Warningf = log.Warningf
	// Errorf logs a message at level Error.
	Errorf = log.Errorf
	// Fatalf logs a message at level Fatal.
	Fatalf = log.Fatalf

	// Traceln logs a message at level Trace.
	Traceln = log.Traceln
	// Debugln logs a message at level Debug.
	Debugln = log.Debugln
	// Infoln logs a message at level Info.
	Infoln = log.Infoln
	// Println logs a message at level Info.
	Println = log.Println
	// Warnln logs a message at level Warn.
	Warnln = log.Warnln
	// Warningln logs a message at level Warn.
	Warningln = log.Warningln
	// Errorln logs a message at level Error.
	Errorln = log.Errorln
	// Panicln logs a message at level Panic.
	Panicln = log.Panicln
	// Fatalln logs a message at level Fatal.
	Fatalln = log.Fatalln

	// Trace logs a message at level Trace.
	Trace = log.Trace
	// Debug logs a message at level Debug.
	Debug = log.Debug
	// Info logs a message at level Info.
	Info = log.Info
	// Print logs a message at level Info.
	Print = log.Print
	// Warn logs a message at level Warn.
	Warn = log.Warn
	// Warning logs a message at level Warn.
	Warning = log.Warning
	// Error logs a message at level Error.
	Error = log.Error
	// Panic logs a message at level Panic.
	Panic = log.Panic
	// Fatal logs a message at level Fatal.
	Fatal = log.Fatal

	// TraceLevel is the log level Trace.
	TraceLevel = log.TraceLevel
	// DebugLevel is the log level Debug.
	DebugLevel = log.DebugLevel
	// InfoLevel is the log level Info.
	InfoLevel = log.InfoLevel
	// WarnLevel is the log level Warn.
	WarnLevel = log.WarnLevel
	// ErrorLevel is the log level Error.
	ErrorLevel = log.ErrorLevel
	// PanicLevel is the log level Panic.
	PanicLevel = log.PanicLevel
	// FatalLevel is the log level Fatal.
	FatalLevel = log.FatalLevel
)

// Setup is setting up the given logger using. It sets the formatter, the log
// level, and the report caller flag. If no logger is given, the standard
// logger is used.
func (c *Config) Setup(logger *Logger) *Logger {
	// Uses the standard logger if no logger is given.
	if logger == nil {
		logger = StandardLogger()
	}

	// Sets up the text formatter with the given time format.
	if c.TimeFormat != "" {
		logger.SetFormatter(&TextFormatter{
			TimestampFormat: c.TimeFormat,
			ForceColors:     true,
		})
	}

	// Sets up the log level if given.
	logLevel, err := ParseLevel(c.Level)
	if err != nil {
		logger.WithError(err).WithFields(Fields{
			"config": c.Level,
		}).Info("failed setting log level")
	} else {
		logger.SetLevel(logLevel)
		WithFields(Fields{
			"level": c.Level,
		}).Info("setting up log level")
	}

	// Helpful setting in certain debug situations.
	logger.SetReportCaller(c.Caller)

	return logger
}
