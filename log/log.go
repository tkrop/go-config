package log

import (
	"io"
	"os"
	"regexp"
	"strings"

	"golang.org/x/term"
)

// Default values for the log configuration.
const (
	// DefaultLevel is the default log level.
	DefaultLevel = LevelInfo
	// DefaultCaller is the default flag state reporting caller information.
	DefaultCaller = false
	// DefaultTimeFormat is the default time format.
	DefaultTimeFormat = "2006-01-02 15:04:05.999999"
)

// Default values for the log formatter.
var (
	// DefaultLevelColors is the default color mapping for the log levels.
	DefaultLevelColors = []string{
		ColorPanic, ColorFatal, ColorError,
		ColorWarn, ColorInfo, ColorDebug, ColorTrace, ColorField,
	}

	// DefaultLevelNames is the default name mapping for the log levels.
	DefaultLevelNames = []string{
		"PANIC", "FATAL", "ERROR", "WARN",
		"INFO", "DEBUG", "TRACE", "-",
	}

	// DefaultErrorName is the default name used for marking errors.
	DefaultErrorName = "error"
)

// Log levels.
const (
	// LevelPanic is the panic log level.
	LevelPanic string = "panic"
	// LevelFatal is the fatal log level.
	LevelFatal string = "fatal"
	// LevelError is the error log level.
	LevelError string = "error"
	// LevelWarn is the warn log level.
	LevelWarn string = "warn"
	// LevelWarning is the warning log level (alternative name).
	LevelWarning string = "warning"
	// LevelInfo is the info log level.
	LevelInfo string = "info"
	// LevelDebug is the debug log level.
	LevelDebug string = "debug"
	// LevelTrace is the trace log level.
	LevelTrace string = "trace"
)

// Level is the log level used for logging.
type Level int

// Log levels.
const (
	// PanicLevel is the log level used for panics.
	PanicLevel Level = 0
	// FatalLevel is the log level used for fatal errors.
	FatalLevel Level = 1
	// ErrorLevel is the log level used for errors.
	ErrorLevel Level = 2
	// WarnLevel is the log level used for warnings.
	WarnLevel Level = 3
	// InfoLevel is the log level used for information.
	InfoLevel Level = 4
	// DebugLevel is the log level used for debugging.
	DebugLevel Level = 5
	// TraceLevel is the log level used for tracing.
	TraceLevel Level = 6
	// FieldLevel is and extra log level used for field names.
	FieldLevel Level = 7
)

// ParseLevel parses the log level string and returns the corresponding level.
func ParseLevel(level string) Level {
	switch strings.ToLower(level) {
	case LevelPanic:
		return PanicLevel
	case LevelFatal:
		return FatalLevel
	case LevelError:
		return ErrorLevel
	case LevelWarn, LevelWarning:
		return WarnLevel
	case LevelInfo:
		return InfoLevel
	case LevelDebug:
		return DebugLevel
	case LevelTrace:
		return TraceLevel
	default:
		return InfoLevel
	}
}

// Formatter is the formatter used for logging.
type Formatter string

// Formatters.
const (
	// FormatterPretty is the pretty formatter.
	FormatterPretty Formatter = "pretty"
	// FormatterText is the text formatter.
	FormatterText Formatter = "text"
	// FormatterJSON is the JSON formatter.
	FormatterJSON Formatter = "json"
)

// Color codes for the different log levels.
const (
	// ColorRed is the color code for red.
	ColorRed = "1;91"
	// ColorGreen is the color code for green.
	ColorGreen = "1;92"
	// ColorYellow is the color code for yellow.
	ColorYellow = "1;93"
	// ColorBlue is the color code for blue.
	ColorBlue = "1;94"
	// ColorMagenta is the color code for magenta.
	ColorMagenta = "1;95"
	// ColorCyan is the color code for cyan.
	ColorCyan = "1;96"
	// ColorGray is the color code for gray.
	ColorGray = "1;37"

	// ColorPanic is the color code for panic.
	ColorPanic = ColorRed
	// ColorFatal is the color code for fatal.
	ColorFatal = ColorRed
	// ColorError is the color code for error.
	ColorError = ColorRed
	// ColorWarn is the color code for warn.
	ColorWarn = ColorYellow
	// ColorInfo is the color code for info.
	ColorInfo = ColorCyan
	// ColorDebug is the color code for debug.
	ColorDebug = ColorBlue
	// ColorTrace is the color code for trace.
	ColorTrace = ColorMagenta
	// ColorField is the color code for fields.
	ColorField = ColorGray
)

// ColorModeString is the color mode used for logging.
type ColorModeString string

// Color mode strings.
const (
	// ColorModeOff disables the color mode.
	ColorModeOff ColorModeString = "off"
	// ColorModeOn enables the color mode.
	ColorModeOn ColorModeString = "on"
	// ColorModeAuto enables the automatic color mode.
	ColorModeAuto ColorModeString = "auto"
	// ColorModeLevels enables the color mode for log level.
	ColorModeLevels ColorModeString = "levels"
	// ColorModeFields enables the color mode for fields.
	ColorModeFields ColorModeString = "fields"
)

var splitRegex = regexp.MustCompile(`[|,:;]`)

// Parse parses the color mode.
func (m ColorModeString) Parse(colorized bool) ColorMode {
	mode := ColorUnset
	for _, m := range splitRegex.Split(string(m), -1) {
		switch ColorModeString(m) {
		case ColorModeOff:
			mode = ColorOff
		case ColorModeOn:
			mode = ColorOn
		case ColorModeLevels:
			mode |= ColorLevels
		case ColorModeFields:
			mode |= ColorFields
		case ColorModeAuto:
			fallthrough
		default:
			if colorized {
				mode = ColorOn
			} else {
				mode = ColorOff
			}
		}
	}
	return mode
}

// ColorMode is the color mode used for logging.
type ColorMode int

// Color modes.
const (
	// ColorDefault is the default color mode.
	ColorDefault = ColorOn
	// ColorUnset is the unset color mode (activates the default).
	ColorUnset ColorMode = 0
	// ColorOff disables coloring of logs for all outputs files.
	ColorOff ColorMode = 1
	// ColorOn enables coloring of logs for all outputs files.
	ColorOn ColorMode = ColorFields | ColorLevels
	// ColorLevels enables coloring for log levels entries only.
	ColorLevels ColorMode = 2
	// ColorFields enables coloring for fields names only.
	ColorFields ColorMode = 4
)

// CheckFlag checks if the given color mode flag is set.
func (m ColorMode) CheckFlag(flag ColorMode) bool {
	return m&flag == flag
}

// OrderModeString is the order mode used for logging.
type OrderModeString string

// Order modes.
const (
	// OrderModeOff disables the order mode.
	OrderModeOff OrderModeString = "off"
	// OrderModeOn enables the order mode.
	OrderModeOn OrderModeString = "on"
)

// Parse parses the order mode.
func (m OrderModeString) Parse() OrderMode {
	switch m {
	case OrderModeOff:
		return OrderOff
	case OrderModeOn:
		return OrderOn
	default:
		return OrderOff
	}
}

// OrderMode is the order mode used for logging.
type OrderMode int

// Order modes.
const (
	// OrderDefault is the default order mode.
	OrderDefault = OrderOn
	// OrderUnset is the unset order mode.
	OrderUnset OrderMode = 0
	// OrderOff disables the order mode.
	OrderOff OrderMode = 1
	// OrderOn enables the order mode.
	OrderOn OrderMode = 2
)

// CheckFlag checks if the given order mode flag is set.
func (m OrderMode) CheckFlag(flag OrderMode) bool {
	return m&flag == flag
}

// IsTerminal checks whether the given writer is a terminal.
func IsTerminal(writer io.Writer) bool {
	if file, ok := writer.(*os.File); ok {
		// #nosec G115 // is a safe conversion for files.
		return term.IsTerminal(int(file.Fd()))
	}
	return false
}

// Config common configuration for logging.
type Config struct {
	// Level is defining the logger level (default `info`).
	Level string `default:"info"`
	// TImeFormat is defining the time format for timestamps.
	TimeFormat string `default:"2006-01-02 15:04:05.999999"`
	// Caller is defining whether the caller is logged (default `false`).
	Caller bool `default:"false"`
	// File is defining the file name used for the log output.
	File string `default:"/dev/stderr"`
	// ColorMode is defining the color mode used for logging.
	ColorMode ColorModeString `default:"auto"`
	// OrderMode is defining the order mode used for logging.
	OrderMode OrderModeString `default:"on"`
	// Formatter is defining the formatter used for logging.
	Formatter Formatter `default:"pretty"`

	// logger is the logger instance defined by the config.
	logger any
}

// Setup is a data structure that contains all necessary setup information to
// format logs into a pretty format.
type Setup struct {
	// TimeFormat is defining the time format used for printing timestamps.
	TimeFormat string
	// ColorMode is defining the color mode (default = ColorAuto).
	ColorMode ColorMode
	// OrderMode is defining the order mode.
	OrderMode OrderMode
	// Caller is defining whether the caller is reported.
	Caller bool

	// ErrorName is defining the name used for marking errors.
	ErrorName string
	// LevelNames is defining the names used for marking the different log
	// levels.
	LevelNames []string
	// LevelColors is defining the colors used for marking the different log
	// levels.
	LevelColors []string
}

// Setup creates a new pretty formatter config.
func (c *Config) Setup(writer io.Writer) *Setup {
	return &Setup{
		TimeFormat:  c.TimeFormat,
		ColorMode:   c.ColorMode.Parse(IsTerminal(writer)),
		OrderMode:   c.OrderMode.Parse(),
		Caller:      c.Caller,
		ErrorName:   DefaultErrorName,
		LevelNames:  DefaultLevelNames,
		LevelColors: DefaultLevelColors,
	}
}
