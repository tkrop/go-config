package format

import (
	"regexp"

	log "github.com/sirupsen/logrus"
)

// Formatter is the formatter used for logging.
type Formatter string

// Formatters.
const (
	// Pretty is the pretty formatter.
	FormatterPretty Formatter = "pretty"
	// Text is the text formatter.
	FormatterText Formatter = "text"
	// JSON is the JSON formatter.
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
	// ColorGray is the color code for gray.
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

	// FieldLevel is and extra log level used for field names.
	FieldLevel log.Level = 7

	// TImeFormat is defining default time format.
	DefaultTimeFormat = "2006-01-02 15:04:05.999999"
)

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
)

// ColorModeString is the color mode used for logging.
type ColorModeString string

// Color mode strings.
const (
	// ColorOff disables the color mode.
	ColorModeOff ColorModeString = "off"
	// ColorOn enables the color mode.
	ColorModeOn ColorModeString = "on"
	// ColorAuto enables the automatic color mode.
	ColorModeAuto ColorModeString = "auto"
	// ColorLevel enables the color mode for log level.
	ColorModeLevels ColorModeString = "levels"
	// ColorFields enables the color mode for fields.
	ColorModeFields ColorModeString = "fields"
)

var splitRegex = regexp.MustCompile(`[|,:;]`)

// Parse parses the color mode.
func (m ColorModeString) Parse() ColorMode {
	mode := ColorUnset
	for _, m := range splitRegex.Split(string(m), -1) {
		switch ColorModeString(m) {
		case ColorModeOff:
			mode = ColorOff
		case ColorModeOn:
			mode = ColorOn
		case ColorModeAuto:
			mode = ColorAuto
		case ColorModeLevels:
			mode |= ColorLevels
		case ColorModeFields:
			mode |= ColorFields
		default:
			mode = ColorDefault
		}
	}
	return mode
}

// ColorMode is the color mode used for logging.
type ColorMode uint

// Color modes.
const (
	// ColorDefault is the default color mode.
	ColorDefault = ColorAuto
	// ColorUnset is the unset color mode (activates the default).
	ColorUnset ColorMode = 0
	// ColorOff disables coloring of logs for all outputs files.
	ColorOff ColorMode = 1
	// ColorOn enables coloring of logs for all outputs files.
	ColorOn ColorMode = 2
	// ColorAuto enables the automatic coloring for tty outputs files.
	ColorAuto ColorMode = 4
	// ColorLevels enables coloring for log levels entries only.
	ColorLevels ColorMode = 8
	// ColorFields enables coloring for fields names only.
	ColorFields ColorMode = 16
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
		return OrderOn
	}
}

// OrderMode is the order mode used for logging.
type OrderMode uint

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
